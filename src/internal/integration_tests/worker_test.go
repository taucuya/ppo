package integrationtests

import (
	"context"
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
}

func NewWorkerTestFixture(t *testing.T) *WorkerTestFixture {
	workerRepo := worker_rep.New(db)
	userRepo := user_rep.New(db)

	service := worker.New(workerRepo)

	return &WorkerTestFixture{
		t:          t,
		ctx:        context.Background(),
		service:    service,
		workerRepo: workerRepo,
		userRepo:   userRepo,
	}
}

func (f *WorkerTestFixture) createTestUser() structs.User {
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	return structs.User{
		Name:          "Test User",
		Date_of_birth: dob,
		Mail:          "test@example.com",
		Password:      "password123",
		Phone:         "89016475843",
		Address:       "123 Test St",
		Status:        "active",
		Role:          "обычный пользователь",
	}
}

func (f *WorkerTestFixture) createTestWorker(userID uuid.UUID) structs.Worker {
	return structs.Worker{
		IdUser:   userID,
		JobTitle: "работник склада",
	}
}

func (f *WorkerTestFixture) setupWorker() (uuid.UUID, uuid.UUID) {
	testUser := f.createTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)

	testWorker := f.createTestWorker(userID)
	err = f.workerRepo.Create(f.ctx, testWorker)
	require.NoError(f.t, err)

	var workers []struct {
		ID uuid.UUID `db:"id"`
	}
	err = db.SelectContext(f.ctx, &workers, "SELECT id FROM worker WHERE id_user = $1", userID)
	require.NoError(f.t, err)
	require.Len(f.t, workers, 1)

	return userID, workers[0].ID
}

func TestWorker_Create_AAA(t *testing.T) {
	fixture := NewWorkerTestFixture(t)

	tests := []struct {
		name        string
		setup       func() structs.Worker
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successful worker creation",
			setup: func() structs.Worker {
				truncateTables(t)
				testUser := fixture.createTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)
				return fixture.createTestWorker(userID)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to create worker for non-existent user",
			setup: func() structs.Worker {
				truncateTables(t)
				return fixture.createTestWorker(uuid.New())
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worker := tt.setup()
			defer tt.cleanup()

			err := fixture.service.Create(fixture.ctx, worker)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWorker_GetById_AAA(t *testing.T) {
	fixture := NewWorkerTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully get worker by id",
			setup: func() uuid.UUID {
				truncateTables(t)
				_, workerID := fixture.setupWorker()
				return workerID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to get non-existent worker by id",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workerID := tt.setup()
			defer tt.cleanup()

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
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully delete worker",
			setup: func() uuid.UUID {
				truncateTables(t)
				_, workerID := fixture.setupWorker()
				return workerID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent worker",
			setup: func() uuid.UUID {
				truncateTables(t)
				return uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workerID := tt.setup()
			defer tt.cleanup()

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

func TestWorker_GetAllWorkers_AAA(t *testing.T) {
	fixture := NewWorkerTestFixture(t)

	tests := []struct {
		name          string
		setup         func()
		cleanup       func()
		expectedCount int
		expectedErr   bool
	}{
		{
			name: "successfully get all workers",
			setup: func() {
				truncateTables(t)

				user1 := fixture.createTestUser()
				user1ID, err := fixture.userRepo.Create(fixture.ctx, user1)
				require.NoError(t, err)
				worker1 := fixture.createTestWorker(user1ID)
				err = fixture.workerRepo.Create(fixture.ctx, worker1)
				require.NoError(t, err)

				user2 := fixture.createTestUser()
				user2.Mail = "test2@example.com"
				user2.Phone = "89016475844"
				user2ID, err := fixture.userRepo.Create(fixture.ctx, user2)
				require.NoError(t, err)
				worker2 := fixture.createTestWorker(user2ID)
				err = fixture.workerRepo.Create(fixture.ctx, worker2)
				require.NoError(t, err)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name: "get empty workers list",
			setup: func() {
				truncateTables(t)
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedCount: 0,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.cleanup()

			workers, err := fixture.service.GetAllWorkers(fixture.ctx)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, workers, tt.expectedCount)
			}
		})
	}
}

func TestWorker_AcceptOrder_AAA(t *testing.T) {
	fixture := NewWorkerTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (uuid.UUID, uuid.UUID)
		cleanup     func()
		expectedErr bool
	}{
		{
			name: "successfully accept order",
			setup: func() (uuid.UUID, uuid.UUID) {
				truncateTables(t)
				userID, _ := fixture.setupWorker()

				basketID := uuid.New()
				_, err := db.Exec("INSERT INTO basket (id, id_user, date) VALUES ($1, $2, $3)",
					basketID, userID, time.Now())
				require.NoError(t, err)

				orderID := uuid.New()
				_, err = db.Exec(`INSERT INTO "order" (id, id_user, address, status) VALUES ($1, $2, $3, $4)`,
					orderID, userID, "123 Test St", "непринятый")
				require.NoError(t, err)

				return orderID, userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: false,
		},
		{
			name: "fail to accept non-existent order",
			setup: func() (uuid.UUID, uuid.UUID) {
				truncateTables(t)
				userID, _ := fixture.setupWorker()
				return uuid.New(), userID
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
		{
			name: "fail to accept order with non-existent worker",
			setup: func() (uuid.UUID, uuid.UUID) {
				truncateTables(t)
				testUser := fixture.createTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)

				basketID := uuid.New()
				_, err = db.Exec("INSERT INTO basket (id, id_user, date) VALUES ($1, $2, $3)",
					basketID, userID, time.Now())
				require.NoError(t, err)

				orderID := uuid.New()
				_, err = db.Exec(`INSERT INTO "order" (id, id_user, address, status) VALUES ($1, $2, $3, $4)`,
					orderID, userID, "123 Test St", "непринятый")
				require.NoError(t, err)

				return orderID, uuid.New()
			},
			cleanup: func() {
				truncateTables(t)
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID, userID := tt.setup()
			defer tt.cleanup()

			err := fixture.service.AcceptOrder(fixture.ctx, orderID, userID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
