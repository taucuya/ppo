package integrationtests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/taucuya/ppo/internal/core/service/review"
	"github.com/taucuya/ppo/internal/core/structs"
	review_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/review"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
)

type ReviewTestFixture struct {
	t          *testing.T
	ctx        context.Context
	service    *review.Service
	reviewRepo *review_rep.Repository
	userRepo   *user_rep.Repository
	testID     string
}

func NewReviewTestFixture(t *testing.T) *ReviewTestFixture {
	reviewRepo := review_rep.New(db)
	userRepo := user_rep.New(db)
	service := review.New(reviewRepo)

	testID := uuid.New().String()[:8]

	return &ReviewTestFixture{
		t:          t,
		ctx:        context.Background(),
		service:    service,
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
		testID:     testID,
	}
}

func (f *ReviewTestFixture) generateTestUser() structs.User {
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

func (f *ReviewTestFixture) generateTestProduct() uuid.UUID {
	brandName := fmt.Sprintf("Test Brand %d", time.Now().UnixNano())
	_, err := db.Exec("INSERT INTO brand (name, description, price_category) VALUES ($1, $2, $3)",
		brandName, "Test Description", "premium")
	require.NoError(f.t, err)

	var brandID uuid.UUID
	err = db.GetContext(f.ctx, &brandID, "SELECT id FROM brand WHERE name = $1", brandName)
	require.NoError(f.t, err)

	productID := uuid.New()
	uniqueArt := fmt.Sprintf("TEST-ART-%s", uuid.New().String()[:8])
	_, err = db.Exec("INSERT INTO product (id, name, description, price, id_brand, amount, art) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		productID, "Test Product", "Test Description", 1000, brandID, 10, uniqueArt)
	require.NoError(f.t, err)

	return productID
}

func (f *ReviewTestFixture) generateTestReview(userID, productID uuid.UUID) structs.Review {
	return structs.Review{
		IdProduct: productID,
		IdUser:    userID,
		Rating:    5,
		Text:      "Great product!",
		Date:      time.Now(),
	}
}

func (f *ReviewTestFixture) generateAnotherTestReview(userID, productID uuid.UUID) structs.Review {
	return structs.Review{
		IdProduct: productID,
		IdUser:    userID,
		Rating:    4,
		Text:      "Good product",
		Date:      time.Now().Add(-24 * time.Hour),
	}
}

func (f *ReviewTestFixture) cleanupReviewData(reviewID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM review WHERE id = $1", reviewID)
}

func (f *ReviewTestFixture) cleanupUserData(userID uuid.UUID) {
	_, _ = db.ExecContext(f.ctx, "DELETE FROM review WHERE id_user = $1", userID)
	_, _ = db.ExecContext(f.ctx, "DELETE FROM \"user\" WHERE id = $1", userID)
}

func (f *ReviewTestFixture) createReviewForTest() (uuid.UUID, uuid.UUID, uuid.UUID) {
	testUser := f.generateTestUser()
	userID, err := f.userRepo.Create(f.ctx, testUser)
	require.NoError(f.t, err)

	productID := f.generateTestProduct()

	testReview := f.generateTestReview(userID, productID)
	err = f.reviewRepo.Create(f.ctx, testReview)
	require.NoError(f.t, err)

	var reviewID uuid.UUID
	err = db.GetContext(f.ctx, &reviewID, "SELECT id FROM review WHERE id_user = $1 AND id_product = $2", userID, productID)
	require.NoError(f.t, err)

	return userID, productID, reviewID
}

func TestReview_Create_AAA(t *testing.T) {
	fixture := NewReviewTestFixture(t)

	tests := []struct {
		name        string
		setup       func() (structs.Review, []uuid.UUID)
		expectedErr bool
	}{
		{
			name: "successful review creation",
			setup: func() (structs.Review, []uuid.UUID) {
				testUser := fixture.generateTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)

				productID := fixture.generateTestProduct()
				review := fixture.generateTestReview(userID, productID)
				return review, []uuid.UUID{userID}
			},
			expectedErr: false,
		},
		{
			name: "fail to create review for non-existent user",
			setup: func() (structs.Review, []uuid.UUID) {
				productID := fixture.generateTestProduct()
				review := fixture.generateTestReview(uuid.New(), productID)
				return review, []uuid.UUID{}
			},
			expectedErr: true,
		},
		{
			name: "fail to create review for non-existent product",
			setup: func() (structs.Review, []uuid.UUID) {
				testUser := fixture.generateTestUser()
				userID, err := fixture.userRepo.Create(fixture.ctx, testUser)
				require.NoError(t, err)

				review := fixture.generateTestReview(userID, uuid.New())
				return review, []uuid.UUID{userID}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			review, cleanupIDs := tt.setup()

			defer func() {
				for _, id := range cleanupIDs {
					if id != uuid.Nil {
						fixture.cleanupUserData(id)
					}
				}
			}()

			err := fixture.service.Create(fixture.ctx, review)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var createdReviewID uuid.UUID
				err = db.GetContext(fixture.ctx, &createdReviewID,
					"SELECT id FROM review WHERE id_user = $1 AND id_product = $2",
					review.IdUser, review.IdProduct)
				if err == nil {
					defer fixture.cleanupReviewData(createdReviewID)
				}
			}
		})
	}
}

func TestReview_GetById_AAA(t *testing.T) {
	fixture := NewReviewTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectedErr bool
	}{
		{
			name: "successfully get review by id",
			setup: func() uuid.UUID {
				_, _, reviewID := fixture.createReviewForTest()
				return reviewID
			},
			expectedErr: false,
		},
		{
			name: "fail to get non-existent review by id",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewID := tt.setup()
			if reviewID != uuid.Nil {
				defer fixture.cleanupReviewData(reviewID)
			}

			result, err := fixture.service.GetById(fixture.ctx, reviewID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, 5, result.Rating)
				require.Equal(t, "Great product!", result.Text)
			}
		})
	}
}

func TestReview_Delete_AAA(t *testing.T) {
	fixture := NewReviewTestFixture(t)

	tests := []struct {
		name        string
		setup       func() uuid.UUID
		expectedErr bool
	}{
		{
			name: "successfully delete review",
			setup: func() uuid.UUID {
				_, _, reviewID := fixture.createReviewForTest()
				return reviewID
			},
			expectedErr: false,
		},
		{
			name: "fail to delete non-existent review",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewID := tt.setup()

			err := fixture.service.Delete(fixture.ctx, reviewID)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				_, err := fixture.service.GetById(fixture.ctx, reviewID)
				require.Error(t, err)
			}
		})
	}
}
