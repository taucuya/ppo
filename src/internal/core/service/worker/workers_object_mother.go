package worker

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

var errTest = errors.New("test error")

type WorkerMother struct{}

func NewWorkerMother() *WorkerMother {
	return &WorkerMother{}
}

func (m *WorkerMother) ValidWorker() structs.Worker {
	return structs.Worker{
		Id:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		IdUser:   uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		JobTitle: "Developer",
	}
}

func (m *WorkerMother) AnotherWorker() structs.Worker {
	return structs.Worker{
		Id:       uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		IdUser:   uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		JobTitle: "Manager",
	}
}

func (m *WorkerMother) WorkersList() []structs.Worker {
	return []structs.Worker{
		m.ValidWorker(),
		m.AnotherWorker(),
	}
}

func (m *WorkerMother) ValidOrder() structs.Order {
	return structs.Order{
		Id:     uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		IdUser: uuid.MustParse("66666666-6666-6666-6666-666666666666"),
		Status: "pending",
	}
}

func (m *WorkerMother) AnotherOrder() structs.Order {
	return structs.Order{
		Id:     uuid.MustParse("77777777-7777-7777-7777-777777777777"),
		IdUser: uuid.MustParse("88888888-8888-8888-8888-888888888888"),
		Status: "in_progress",
	}
}

func (m *WorkerMother) OrdersList() []structs.Order {
	return []structs.Order{
		m.ValidOrder(),
		m.AnotherOrder(),
	}
}

func (m *WorkerMother) WorkersOrder() structs.WorkersOrders {
	return structs.WorkersOrders{
		IdOrder:  uuid.MustParse("99999999-9999-9999-9999-999999999999"),
		IdWorker: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
	}
}

type TestFixture struct {
	t            *testing.T
	ctrl         *gomock.Controller
	ctx          context.Context
	workerMother *WorkerMother
}

func NewTestFixture(t *testing.T) *TestFixture {
	ctrl := gomock.NewController(t)
	return &TestFixture{
		t:            t,
		ctrl:         ctrl,
		ctx:          context.Background(),
		workerMother: NewWorkerMother(),
	}
}

func (f *TestFixture) Cleanup() {
	f.ctrl.Finish()
}

func (f *TestFixture) CreateServiceWithMocks() (*Service, *mock_structs.MockWorkerRepository) {
	mockRepo := mock_structs.NewMockWorkerRepository(f.ctrl)
	service := New(mockRepo)
	return service, mockRepo
}

func (f *TestFixture) AssertError(err error, expectedErr error) {
	if expectedErr != nil {
		if err == nil {
			f.t.Error("Expected error, got nil")
			return
		}
		if !errors.Is(err, expectedErr) && err.Error() != expectedErr.Error() {
			f.t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	} else if err != nil {
		f.t.Errorf("Unexpected error: %v", err)
	}
}
