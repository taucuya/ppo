package product_rep_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	structs "github.com/taucuya/ppo/internal/core/structs"
	basket_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/basket"
	brand_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/brand"
	order_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/order"
	product_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/product"
	review_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/review"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
	worker_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/worker"
)

func runSQLScripts(db *sqlx.DB, scripts []string) error {
	for _, path := range scripts {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute %s: %w", path, err)
		}
	}
	return nil
}

var repProd *product_rep.Repository
var repBrand *brand_rep.Repository
var db *sqlx.DB

func TestMain(m *testing.M) {
	var err error

	dsn := "postgres://test_user:test_password@test_db:5432/test_db?sslmode=disable"

	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	err = runSQLScripts(db, []string{
		"/app/internal/database/sql/01-create.sql",
		"/app/internal/database/sql/02-constraints.sql",
	})
	if err != nil {
		panic("failed to run SQL scripts: " + err.Error())
	}

	repProd = product_rep.New(db)
	repBrand = brand_rep.New(db)

	code := m.Run()

	_ = db.Close()
	os.Exit(code)
}

func TestCreateAndGetProduct(t *testing.T) {
	ctx := context.Background()

	brand := structs.Brand{
		Name:          "Dior",
		Description:   "Премиум косметика",
		PriceCategory: "lux",
	}

	var lastId uuid.UUID

	err := repBrand.Create(ctx, brand)

	if err != nil {
		t.Fatalf("failed to create brand: %v", err)
	}

	err = db.QueryRowContext(ctx, "SELECT id FROM brand WHERE name = $1", brand.Name).Scan(&lastId)
	if err != nil {
		t.Fatalf("failed to fetch brand ID: %v", err)
	}

	product := structs.Product{
		Name:        "Тональный крем",
		Description: "Матирующий эффект",
		Price:       899.99,
		Category:    "lux",
		Amount:      15,
		IdBrand:     lastId,
		PicLink:     "http://example.com/image.jpg",
		Articule:    "TON-123",
	}

	err = repProd.Create(ctx, product)
	if err != nil {
		t.Fatalf("failed to create product: %v", err)
	}

	err = db.QueryRowContext(ctx, "SELECT id FROM product WHERE art = $1", product.Articule).Scan(&lastId)
	if err != nil {
		t.Fatalf("failed to fetch product ID: %v", err)
	}

	fetched, err := repProd.GetById(ctx, lastId)
	if err != nil {
		t.Fatalf("failed to fetch product: %v", err)
	}

	if fetched.Name != product.Name || fetched.Articule != product.Articule {
		t.Errorf("expected %+v, got %+v", product, fetched)
	}

	_, _ = db.ExecContext(ctx, "DELETE FROM product WHERE id = $1", uuid.UUID(lastId))
}

func TestAddProductToBasket(t *testing.T) {
	ctx := context.Background()

	user := structs.User{
		Name:          "Алексей Петров",
		Date_of_birth: time.Now(),
		Mail:          "aleksei@example.com",
		Password:      "passpass",
		Phone:         "+79998887766",
	}
	userRep := user_rep.New(db)

	id, err := userRep.Create(ctx, user)
	if err != nil {
		t.Fatalf("не удалось создать пользователя: %v", err)
	}

	basketRep := basket_rep.New(db)
	basket := structs.Basket{
		IdUser: id,
		Date:   time.Now(),
	}
	if err := basketRep.Create(ctx, basket); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		t.Fatalf("не удалось создать корзину: %v", err)
	}

	var basketId uuid.UUID
	if err := db.QueryRowContext(ctx, `SELECT id FROM basket WHERE id_user = $1`, id).Scan(&basketId); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		t.Fatalf("не удалось получить id корзины: %v", err)
	}

	brand := structs.Brand{
		Name:          "MAC",
		Description:   "Макияж и косметика",
		PriceCategory: "mid",
	}
	if err := repBrand.Create(ctx, brand); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_, _ = db.ExecContext(ctx, "DELETE FROM basket WHERE id = $1", basketId)
		t.Fatalf("не удалось добавить бренд: %v", err)
	}
	var brandId uuid.UUID
	if err := db.QueryRowContext(ctx, `SELECT id FROM brand WHERE name = $1`, brand.Name).Scan(&brandId); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_, _ = db.ExecContext(ctx, "DELETE FROM basket WHERE id = $1", basketId)
		t.Fatalf("не удалось получить id бренда: %v", err)
	}

	product := structs.Product{
		Name:        "Карандаш для глаз",
		Description: "Черный",
		Price:       399.00,
		Category:    "mid",
		Amount:      25,
		IdBrand:     brandId,
		PicLink:     "http://example.com/eye-pencil.jpg",
		Articule:    "EYE-789",
	}
	if err := repProd.Create(ctx, product); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_, _ = db.ExecContext(ctx, "DELETE FROM basket WHERE id = $1", basketId)
		_ = repBrand.Delete(ctx, brandId)
		t.Fatalf("не удалось добавить продукт: %v", err)
	}

	var productId uuid.UUID
	if err := db.QueryRowContext(ctx, `SELECT id FROM product WHERE art = $1`, product.Articule).Scan(&productId); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_, _ = db.ExecContext(ctx, "DELETE FROM basket WHERE id = $1", basketId)
		_ = repBrand.Delete(ctx, brandId)
		t.Fatalf("не удалось получить id продукта: %v", err)
	}

	item := structs.BasketItem{
		IdProduct: productId,
		IdBasket:  basketId,
		Amount:    3,
	}
	if err := basketRep.AddItem(ctx, item); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_, _ = db.ExecContext(ctx, "DELETE FROM basket WHERE id = $1", basketId)
		_ = repBrand.Delete(ctx, brandId)
		_ = repProd.Delete(ctx, productId)
		t.Fatalf("не добавился элемент в корзину: %v", err)
	}
}

func TestAddReviewToProduct(t *testing.T) {
	ctx := context.Background()

	user := structs.User{
		Name:     "Мария Смирнова",
		Mail:     "maria@example.com",
		Password: "mypass123",
		Phone:    "+79990001122",
	}
	userRep := user_rep.New(db)
	id, err := userRep.Create(ctx, user)
	if err != nil {
		t.Fatalf("не удалось создать пользователя: %v", err)
	}

	brand := structs.Brand{
		Name:          "Lancome",
		Description:   "Французская косметика",
		PriceCategory: "premium",
	}
	if err := repBrand.Create(ctx, brand); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		t.Fatalf("не удалось создать бренд: %v", err)
	}

	var brandId uuid.UUID
	if err := db.QueryRowContext(ctx, `SELECT id FROM brand WHERE name = $1`, brand.Name).Scan(&brandId); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		t.Fatalf("не удалось получить id бренда: %v", err)
	}

	product := structs.Product{
		Name:        "Тушь для ресниц",
		Description: "Объемная",
		Price:       799.90,
		Category:    "premium",
		Amount:      10,
		IdBrand:     brandId,
		PicLink:     "http://example.com/mascara.jpg",
		Articule:    "MASC-001",
	}
	if err := repProd.Create(ctx, product); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_ = repBrand.Delete(ctx, brandId)
		t.Fatalf("не удалось создать продукт: %v", err)
	}

	var productId uuid.UUID
	if err := db.QueryRowContext(ctx, `SELECT id FROM product WHERE art = $1`, product.Articule).Scan(&productId); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_ = repBrand.Delete(ctx, brandId)
		t.Fatalf("не удалось получить id продукта: %v", err)
	}

	reviewRep := review_rep.New(db)
	review := structs.Review{
		IdUser:    id,
		IdProduct: productId,
		Rating:    5,
		Text:      "Очень понравилось!",
	}
	if err := reviewRep.Create(ctx, review); err != nil {
		_, _ = db.ExecContext(ctx, "DELETE FROM \"user\" WHERE id = $1", id)
		_ = repBrand.Delete(ctx, brandId)
		_ = repProd.Delete(ctx, productId)
		t.Fatalf("не удалось создать отзыв: %v", err)
	}
}

func TestFullFlow(t *testing.T) {
	ctx := context.Background()

	basket := basket_rep.New(db)
	brand := brand_rep.New(db)
	product := product_rep.New(db)
	review := review_rep.New(db)
	order := order_rep.New(db)
	user := user_rep.New(db)
	worker := worker_rep.New(db)

	id, err := user.Create(ctx, structs.User{
		Name:          "Alice",
		Mail:          "alice@example.com",
		Password:      "secure123",
		Phone:         "+71234567869",
		Address:       "Dream St. 42",
		Status:        "active",
		Role:          "customer",
		Date_of_birth: time.Now().AddDate(-30, 0, 0),
	})
	require.NoError(t, err)
	err = db.QueryRowContext(ctx, `SELECT id FROM "user" WHERE mail = $1`, "alice@example.com").Scan(&id)

	require.NoError(t, err)
	err = brand.Create(ctx, structs.Brand{
		Name:          "BrandX",
		Description:   "Quality Stuff",
		PriceCategory: "medium",
	})
	require.NoError(t, err)
	var brandId uuid.UUID
	err = db.QueryRowContext(ctx, `SELECT id FROM brand WHERE name = $1`, "BrandX").Scan(&brandId)
	require.NoError(t, err)

	err = product.Create(ctx, structs.Product{
		Name:        "Sneakers",
		Description: "Running shoes",
		Price:       99.99,
		Category:    "Footwear",
		Amount:      10,
		IdBrand:     brandId,
		PicLink:     "link/to/pic",
		Articule:    "ART123",
	})
	require.NoError(t, err)
	var productId uuid.UUID
	err = db.QueryRowContext(ctx, `SELECT id FROM product WHERE art = $1`, "ART123").Scan(&productId)
	require.NoError(t, err)
	err = basket.Create(ctx, structs.Basket{
		IdUser: id,
		Date:   time.Now(),
	})
	require.NoError(t, err)
	var basketId uuid.UUID
	db.QueryRowContext(ctx, `SELECT id FROM basket WHERE id_user = $1`, id).Scan(&basketId)
	require.NoError(t, err)
	err = basket.AddItem(ctx, structs.BasketItem{
		IdProduct: productId,
		IdBasket:  basketId,
		Amount:    1,
	})
	require.NoError(t, err)

	err = review.Create(ctx, structs.Review{
		IdProduct: productId,
		IdUser:    id,
		Rating:    5,
		Text:      "Excellent product!",
		Date:      time.Now(),
	})
	require.NoError(t, err)

	var reviewId uuid.UUID
	db.QueryRowContext(ctx, `SELECT id FROM review WHERE id_user = $1`, id).Scan(&reviewId)
	require.NoError(t, err)

	err = worker.Create(ctx, structs.Worker{
		IdUser:   id,
		JobTitle: "Courier",
	})
	require.NoError(t, err)

	var workerId uuid.UUID
	db.QueryRowContext(ctx, `SELECT id FROM "worker" WHERE id_user = $1`, id).Scan(&workerId)
	require.NoError(t, err)

	err = order.Create(ctx, structs.Order{
		Date:     time.Now(),
		IdUser:   id,
		Address:  "Dream St. 42",
		Status:   "pending",
		Price:    99.99,
		IdWorker: workerId,
	})
	require.NoError(t, err)

	var orderId uuid.UUID
	db.QueryRowContext(ctx, `SELECT id FROM "order" WHERE id_user = $1`, id).Scan(&orderId)
	require.NoError(t, err)

	gotWorker, err := worker.GetById(ctx, workerId)
	require.NoError(t, err)
	require.Equal(t, "Courier", gotWorker.JobTitle)
	require.Equal(t, id, gotWorker.IdUser)

	status, err := order.GetStatus(ctx, orderId)
	require.NoError(t, err)
	require.Equal(t, "pending", status)

	err = order.Delete(ctx, orderId)
	require.NoError(t, err)

	err = review.Delete(ctx, reviewId)
	require.NoError(t, err)

	err = product.Delete(ctx, productId)
	require.NoError(t, err)

	err = brand.Delete(ctx, brandId)
	require.NoError(t, err)

	err = worker.Delete(ctx, workerId)
	require.NoError(t, err)

}
