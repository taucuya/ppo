package integrationtests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/worker"
	"github.com/taucuya/ppo/internal/core/structs"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
	worker_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/worker"
)

type WorkerTestFixture struct {
	t          *testing.T
	ctx        context.Context
	service    *worker.Service
	workerRepo *worker_rep.Repository
	userRepo   *user_rep.Repository
	testID     string
}

func NewWorkerTestFixture(t *testing.T) *WorkerTestFixture {
	workerRepo := worker_rep.New(db)
	userRepo := user_rep.New(db)

	service := worker.New(workerRepo)

	testID := uuid.New().String()[:8]

	return &WorkerTestFixture{
		t:          t,
		ctx:        context.Background(),
		service:    service,
		workerRepo: workerRepo,
		userRepo:   userRepo,
		testID:     testID,
	}
}

func (f *WorkerTestFixture) generateTestUser() structs.User {
	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("%s-%d", f.testID, timestamp)

	dob, _ := time.Parse("2006-01-02", "1990-01-01")

	phoneSuffix := fmt.Sprintf("%09d", timestamp%1000000000)
	phone := "89" + phoneSuffix
	if len(phone) > 11 {
		phone = phone[:11]
	}

	return structs.User{
		Name:          fmt.Sprintf("Test User %s", uniqueID),
		Date_of_birth: dob,
		Mail:          fmt.Sprintf("test%s@example.com", uniqueID),
		Password:      "password123",
		Phone:         phone,
		Address:       "123 Test St",
		Status:        "active",
		Role:          "обычный пользователь",
	}
}

func (f *WorkerTestFixture) generateTestWorker(userID uuid.UUID) structs.Worker {
	return structs.Worker{
		IdUser:   userID,
		JobTitle: "работник склада",
	}
}

func (f *WorkerTestFixture) cleanupWorkerData(workerID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM order_worker WHERE id_worker = $1", workerID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM worker WHERE id = $1", workerID)
}

func (f *WorkerTestFixture) cleanupUserData(userID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM worker WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"user\" WHERE id = $1", userID)
}

func (f *WorkerTestFixture) cleanupOrderData(orderID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM order_worker WHERE id_order = $1", orderID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM order_item WHERE id_order = $1", orderID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"order\" WHERE id = $1", orderID)
}

func (f *WorkerTestFixture) createWorkerForTest() (uuid.UUID, uuid.UUID, uuid.UUID) {
	testUser := f.generateTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)

	testWorker := f.generateTestWorker(userID)
	err = f.workerRepo.Create(f.ctx, testWorker)
	require.NoError(f.t, err)

	var workerID uuid.UUID
	err = db.GetContext(f.ctx, &workerID, "SELECT id FROM worker WHERE id_user = $1", userID)
	require.NoError(f.t, err)

	return userID, workerID, workerID
}

func TestWorker_Create_AAA(t *testing.T) {
	fixture := NewWorkerTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.Worker, []uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successful worker creation",
			setup: func() (structs.Worker, []uuid.UUID) {
				testUser := fixture.generateTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
				worker := fixture.generateTestWorker(userID)
				return worker, []uuid.UUID{userID}
			},
			expectedErr: false,
		},
		{
			name: "fail to create worker for non-existent user",
			setup: func() (structs.Worker, []uuid.UUID) {
				worker := fixture.generateTestWorker(uuid.New())
				return worker, []uuid.UUID{}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker, cleanupIDs := tt.setup()

			defer func() {
				for _, id := range cleanupIDs {
					if id != uuid.Nil {
						fixture.cleanupUserData(id)
					}
				}
			}()

			err := fixture.service.Create(fixture.ctx, worker)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var createdWorkerID uuid.UUID
				err = db.GetContext(fixture.ctx, &createdWorkerID,
					"SELECT id FROM worker WHERE id_user = $1", worker.IdUser)
				require.NoError(t, err)
				defer fixture.cleanupWorkerData(createdWorkerID)
			}
		})
	}
}

func TestWorker_GetById_AAA(t *testing.T) {
	fixture := NewWorkerTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectedErr bool
	}{
		{
			name: "successfully get worker by id",
			setup: func() uuid.UUID {
				_, workerID, _ := fixture.createWorkerForTest()
				return workerID
			},
			expectedErr: false,
		},
		{
			name: "fail to get non-existent worker by id",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workerID := tt.setup()
			if workerID != uuid.Nil {
				defer fixture.cleanupWorkerData(workerID)
			}

			result, err := fixture.service.GetById(fixture.ctx, workerID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, workerID, result.Id)
				require.Equal(t, "работник склада", result.JobTitle)
			}
		})
	}
}

func TestWorker_Delete_AAA(t *testing.T) {
	fixture := NewWorkerTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectedErr bool
	}{
		{
			name: "successfully delete worker",
			setup: func() uuid.UUID {
				_, workerID, _ := fixture.createWorkerForTest()
				return workerID
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent worker",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workerID := tt.setup()

			err := fixture.service.Delete(fixture.ctx, workerID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				_, err := fixture.service.GetById(fixture.ctx, workerID)
				require.Error(t, err)
			}
		})
	}
}
