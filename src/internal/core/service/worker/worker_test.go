package worker

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taucuya/ppo/internal/core/mock_structs"
	"github.com/taucuya/ppo/internal/core/structs"
)

func TestCreate_AAA(t *testing.T) {
	fixture := NewTestFixture(t)
	testWorker := fixture.workerMother.ValidWorker()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockWorkerRepository, structs.Worker)
		expectedErr error
	}{
		{
			name: "successful creation",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, worker structs.Worker) {
				mockRepo.EXPECT().Create(fixture.ctx, worker).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, worker structs.Worker) {
				mockRepo.EXPECT().Create(fixture.ctx, worker).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testWorker)
			err := service.Create(fixture.ctx, testWorker)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetById_AAA(t *testing.T) {
	fixture := NewTestFixture(t)
	testWorker := fixture.workerMother.ValidWorker()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockWorkerRepository, structs.Worker)
		expectedRet structs.Worker
		expectedErr error
	}{
		{
			name: "successful get by id",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, worker structs.Worker) {
				mockRepo.EXPECT().GetById(fixture.ctx, worker.Id).Return(worker, nil)
			},
			expectedRet: testWorker,
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, worker structs.Worker) {
				mockRepo.EXPECT().GetById(fixture.ctx, worker.Id).Return(structs.Worker{}, errTest)
			},
			expectedRet: structs.Worker{},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testWorker)
			ret, err := service.GetById(fixture.ctx, testWorker.Id)
			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
		})
	}
	fixture.Cleanup()
}

func TestDelete_AAA(t *testing.T) {
	fixture := NewTestFixture(t)
	testWorker := fixture.workerMother.ValidWorker()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockWorkerRepository, uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful delete",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, id uuid.UUID) {
				mockRepo.EXPECT().Delete(fixture.ctx, id).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, id uuid.UUID) {
				mockRepo.EXPECT().Delete(fixture.ctx, id).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testWorker.Id)
			err := service.Delete(fixture.ctx, testWorker.Id)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}

func TestGetAllWorkers_AAA(t *testing.T) {
	fixture := NewTestFixture(t)
	testWorkers := fixture.workerMother.WorkersList()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockWorkerRepository, []structs.Worker)
		expectedRet []structs.Worker
		expectedErr error
	}{
		{
			name: "successful get all workers",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, workers []structs.Worker) {
				mockRepo.EXPECT().GetAllWorkers(fixture.ctx).Return(workers, nil)
			},
			expectedRet: testWorkers,
			expectedErr: nil,
		},
		{
			name: "empty list",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, workers []structs.Worker) {
				mockRepo.EXPECT().GetAllWorkers(fixture.ctx).Return([]structs.Worker{}, nil)
			},
			expectedRet: []structs.Worker{},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, workers []structs.Worker) {
				mockRepo.EXPECT().GetAllWorkers(fixture.ctx).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testWorkers)
			ret, err := service.GetAllWorkers(fixture.ctx)
			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
		})
	}
	fixture.Cleanup()
}

func TestGetOrders_AAA(t *testing.T) {
	fixture := NewTestFixture(t)
	testWorker := fixture.workerMother.ValidWorker()
	testOrders := fixture.workerMother.OrdersList()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockWorkerRepository, uuid.UUID, []structs.Order)
		expectedRet []structs.Order
		expectedErr error
	}{
		{
			name: "successful get orders",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, id uuid.UUID, orders []structs.Order) {
				mockRepo.EXPECT().GetOrders(fixture.ctx, id).Return(orders, nil)
			},
			expectedRet: testOrders,
			expectedErr: nil,
		},
		{
			name: "empty orders",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, id uuid.UUID, orders []structs.Order) {
				mockRepo.EXPECT().GetOrders(fixture.ctx, id).Return([]structs.Order{}, nil)
			},
			expectedRet: []structs.Order{},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, id uuid.UUID, orders []structs.Order) {
				mockRepo.EXPECT().GetOrders(fixture.ctx, id).Return(nil, errTest)
			},
			expectedRet: nil,
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testWorker.Id, testOrders)
			ret, err := service.GetOrders(fixture.ctx, testWorker.Id)
			if tt.expectedErr != nil {
				fixture.AssertError(err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRet, ret)
			}
		})
	}
	fixture.Cleanup()
}

func TestAcceptOrder_AAA(t *testing.T) {
	fixture := NewTestFixture(t)
	testWorker := fixture.workerMother.ValidWorker()
	testOrder := fixture.workerMother.ValidOrder()

	tests := []struct {
		name        string
		setupMocks  func(*mock_structs.MockWorkerRepository, uuid.UUID, uuid.UUID)
		expectedErr error
	}{
		{
			name: "successful accept order",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, orderID uuid.UUID, workerID uuid.UUID) {
				mockRepo.EXPECT().AcceptOrder(fixture.ctx, orderID, workerID).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",
			setupMocks: func(mockRepo *mock_structs.MockWorkerRepository, orderID uuid.UUID, workerID uuid.UUID) {
				mockRepo.EXPECT().AcceptOrder(fixture.ctx, orderID, workerID).Return(errTest)
			},
			expectedErr: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := fixture.CreateServiceWithMocks()
			tt.setupMocks(mockRepo, testOrder.Id, testWorker.Id)
			err := service.AcceptOrder(fixture.ctx, testOrder.Id, testWorker.Id)
			fixture.AssertError(err, tt.expectedErr)
		})
	}
	fixture.Cleanup()
}
