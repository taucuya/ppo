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
	auth_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/auth"
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

func truncateTables(t *testing.T) {
	tables := []string{
		"review", "product", "brand", "order_item", "order_worker",
		"\"order\"", "basket_item", "basket", "worker", "\"user\"",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
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
		Phone:         "89998887766",
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
		Phone:    "89990001122",
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
		Phone:         "81234567869",
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
		Date:    time.Now(),
		IdUser:  id,
		Address: "Dream St. 42",
		Status:  "pending",
		Price:   99.99,
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

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	repo := user_rep.New(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Test User",
			Date_of_birth: dob,
			Mail:          "test@example.com",
			Password:      "password123",
			Phone:         "89016475843",
			Address:       "123 Test St",
			Status:        "active",
			Role:          "обычный покупатель",
		}

		id, err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		retrieved, err := repo.GetById(ctx, id)
		if err != nil {
			t.Fatalf("Failed to get user by ID: %v", err)
		}

		if retrieved.Name != user.Name {
			t.Errorf("Expected name %s, got %s", user.Name, retrieved.Name)
		}
		if retrieved.Mail != user.Mail {
			t.Errorf("Expected mail %s, got %s", user.Mail, retrieved.Mail)
		}
	})

	t.Run("GetByMail", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Mail Test User",
			Date_of_birth: dob,
			Mail:          "mailtest@example.com",
			Password:      "password123",
			Phone:         "89016475843",
			Address:       "123 Test St",
			Status:        "active",
			Role:          "обычный покупатель",
		}

		_, err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		retrieved, err := repo.GetByMail(ctx, user.Mail)
		if err != nil {
			t.Fatalf("Failed to get user by mail: %v", err)
		}

		if retrieved.Name != user.Name {
			t.Errorf("Expected name %s, got %s", user.Name, retrieved.Name)
		}
	})

	t.Run("GetByPhone", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Phone Test User",
			Date_of_birth: dob,
			Mail:          "phonetest@example.com",
			Password:      "password123",
			Phone:         "89016475843",
			Address:       "123 Test St",
			Status:        "active",
			Role:          "обычный покупатель",
		}

		_, err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		retrieved, err := repo.GetByPhone(ctx, user.Phone)
		if err != nil {
			t.Fatalf("Failed to get user by phone: %v", err)
		}

		if retrieved.Name != user.Name {
			t.Errorf("Expected name %s, got %s", user.Name, retrieved.Name)
		}
	})
}

func TestWorkerRepository(t *testing.T) {
	ctx := context.Background()
	userRepo := user_rep.New(db)
	workerRepo := worker_rep.New(db)

	t.Run("Create and GetByID", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Worker Test User",
			Date_of_birth: dob,
			Mail:          "worker@example.com",
			Password:      "password123",
			Phone:         "89016475843",
			Address:       "123 Worker St",
			Status:        "active",
			Role:          "обычный покупатель",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		worker := structs.Worker{
			IdUser:   userID,
			JobTitle: "Warehouse Worker",
		}

		err = workerRepo.Create(ctx, worker)
		if err != nil {
			t.Fatalf("Failed to create worker: %v", err)
		}
	})

	t.Run("AcceptOrder", func(t *testing.T) {
		truncateTables(t)
		brandRepo := brand_rep.New(db)
		productRepo := product_rep.New(db)
		basketRepo := basket_rep.New(db)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Order Worker",
			Date_of_birth: dob,
			Mail:          "orderworker@example.com",
			Password:      "password123",
			Phone:         "89016475843",
			Address:       "123 Worker St",
			Status:        "active",
			Role:          "обычный покупатель",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		worker := structs.Worker{
			IdUser:   userID,
			JobTitle: "Warehouse Worker",
		}
		err = workerRepo.Create(ctx, worker)
		if err != nil {
			t.Fatalf("Failed to create worker: %v", err)
		}

		brand := structs.Brand{
			Name:          "Test Brand",
			Description:   "Test Brand Description",
			PriceCategory: "premium",
		}
		err = brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id uuid.UUID
		if err = db.Get(&id, `select id from brand where name = $1`, brand.Name); err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		product := structs.Product{
			Name:        "Test Product",
			Description: "Test Product Description",
			Price:       10.99,
			Category:    "electronics",
			Amount:      100,
			IdBrand:     id,
			PicLink:     "http://example.com/product.jpg",
			Articule:    "TP123",
		}
		err = productRepo.Create(ctx, product)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}

		basket := structs.Basket{
			IdUser: userID,
			Date:   time.Now(),
		}
		err = basketRepo.Create(ctx, basket)
		if err != nil {
			t.Fatalf("Failed to create basket: %v", err)
		}

		var id_product uuid.UUID
		err = db.Get(&id_product, `select id from product where name = $1`, "Test Product")
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		var id_basket uuid.UUID
		err = db.Get(&id_basket, `select id from basket where id_user = $1`, userID)
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		basketItem := structs.BasketItem{
			IdProduct: id_product,
			IdBasket:  id_basket,
			Amount:    2,
		}
		err = basketRepo.AddItem(ctx, basketItem)
		if err != nil {
			t.Fatalf("Failed to add item to basket: %v", err)
		}

		orderRepo := order_rep.New(db)
		order := structs.Order{
			Date:    time.Now(),
			IdUser:  userID,
			Address: "123 Test St",
			Status:  "непринятый",
			Price:   0,
		}

		err = orderRepo.Create(ctx, order)
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		var orderID uuid.UUID
		err = db.GetContext(ctx, &orderID, `SELECT id FROM "order" WHERE id_user = $1 ORDER BY date DESC LIMIT 1`, userID)
		if err != nil {
			t.Fatalf("Failed to get order ID: %v", err)
		}

		err = workerRepo.AcceptOrder(ctx, orderID, userID)
		if err != nil {
			t.Fatalf("Failed to accept order: %v", err)
		}

		status, err := orderRepo.GetStatus(ctx, orderID)
		if err != nil {
			t.Fatalf("Failed to get order status: %v", err)
		}
		if status != "принятый" {
			t.Errorf("Expected order status 'принятый', got '%s'", status)
		}

		orders, err := workerRepo.GetOrders(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get worker's orders: %v", err)
		}
		if len(orders) != 1 {
			t.Errorf("Expected 1 order, got %d", len(orders))
		}
	})
}

func TestBasketRepository(t *testing.T) {
	ctx := context.Background()
	userRepo := user_rep.New(db)
	basketRepo := basket_rep.New(db)
	productRepo := product_rep.New(db)
	brandRepo := brand_rep.New(db)

	t.Run("Create and GetBasket", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Basket Test User",
			Date_of_birth: dob,
			Mail:          "basket@example.com",
			Password:      "password123",
			Phone:         "89016412843",
			Address:       "123 Basket St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		basket := structs.Basket{
			IdUser: userID,
			Date:   time.Now(),
		}

		err = basketRepo.Create(ctx, basket)
		if err != nil {
			t.Fatalf("Failed to create basket: %v", err)
		}

		var id_basket uuid.UUID
		err = db.Get(&id_basket, `select id from basket where id_user = $1`, userID)
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		retrieved, err := basketRepo.GetById(ctx, id_basket)
		if err != nil {
			t.Fatalf("Failed to get basket: %v", err)
		}

		if retrieved.IdUser != userID {
			t.Errorf("Expected user ID %s, got %s", userID, retrieved.IdUser)
		}
	})

	t.Run("Add and Get Basket Items", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Basket Items User",
			Date_of_birth: dob,
			Mail:          "basketitems@example.com",
			Password:      "password123",
			Phone:         "82226475843",
			Address:       "123 Basket St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		basket := structs.Basket{
			Id:     uuid.New(),
			IdUser: userID,
			Date:   time.Now(),
		}

		err = basketRepo.Create(ctx, basket)
		if err != nil {
			t.Fatalf("Failed to create basket: %v", err)
		}

		var id_basket uuid.UUID
		err = db.Get(&id_basket, `select id from basket where id_user = $1`, userID)
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		brand := structs.Brand{
			Name:          "Test Brand",
			Description:   "Test Brand Description",
			PriceCategory: "premium",
		}

		err = brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id_brand uuid.UUID
		err = db.Get(&id_brand, `select id from brand where name = $1`, "Test Brand")
		if err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		product := structs.Product{
			Name:        "Test Product",
			Description: "Test Product Description",
			Price:       10.99,
			Category:    "electronics",
			Amount:      100,
			IdBrand:     id_brand,
			PicLink:     "http://example.com/product.jpg",
			Articule:    "TP123",
		}

		err = productRepo.Create(ctx, product)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}

		var id_product uuid.UUID
		err = db.Get(&id_product, `select id from product where name = $1`, "Test Product")
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		basketItem := structs.BasketItem{
			IdProduct: id_product,
			IdBasket:  id_basket,
			Amount:    2,
		}

		err = basketRepo.AddItem(ctx, basketItem)
		if err != nil {
			t.Fatalf("Failed to add item to basket: %v", err)
		}

		items, err := basketRepo.GetItems(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get basket items: %v", err)
		}

		if len(items) != 1 {
			t.Errorf("Expected 1 item in basket, got %d", len(items))
		} else {
			if items[0].IdProduct != id_product {
				t.Errorf("Expected product ID %s, got %s", id_product, items[0].IdProduct)
			}
			if items[0].Amount != 2 {
				t.Errorf("Expected amount 2, got %d", items[0].Amount)
			}
		}
	})
}

func TestOrderRepository(t *testing.T) {
	ctx := context.Background()
	userRepo := user_rep.New(db)
	basketRepo := basket_rep.New(db)
	productRepo := product_rep.New(db)
	brandRepo := brand_rep.New(db)
	orderRepo := order_rep.New(db)

	t.Run("Create and Get Order", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Order Test User",
			Date_of_birth: dob,
			Mail:          "order@example.com",
			Password:      "password123",
			Phone:         "89012275843",
			Address:       "123 Order St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		basket := structs.Basket{
			IdUser: userID,
			Date:   time.Now(),
		}

		err = basketRepo.Create(ctx, basket)
		if err != nil {
			t.Fatalf("Failed to create basket: %v", err)
		}

		var id_basket uuid.UUID
		err = db.Get(&id_basket, `select id from basket where id_user = $1`, userID)
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		brand := structs.Brand{
			Name:          "Order Test Brand",
			Description:   "Order Test Brand Description",
			PriceCategory: "premium",
		}

		err = brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id_brand uuid.UUID
		err = db.Get(&id_brand, `select id from brand where name = $1`, "Order Test Brand")
		if err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		product := structs.Product{
			Name:        "Order Test Product",
			Description: "Order Test Product Description",
			Price:       15.99,
			Category:    "electronics",
			Amount:      50,
			IdBrand:     id_brand,
			PicLink:     "http://example.com/orderproduct.jpg",
			Articule:    "OTP123",
		}

		err = productRepo.Create(ctx, product)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}

		var id_product uuid.UUID
		err = db.Get(&id_product, `select id from product where name = $1`, "Order Test Product")
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		basketItem := structs.BasketItem{
			IdProduct: id_product,
			IdBasket:  id_basket,
			Amount:    3,
		}

		err = basketRepo.AddItem(ctx, basketItem)
		if err != nil {
			t.Fatalf("Failed to add item to basket: %v", err)
		}

		order := structs.Order{
			Date:    time.Now(),
			IdUser:  userID,
			Address: "456 Order St",
			Status:  "непринятый",
			Price:   0,
		}

		err = orderRepo.Create(ctx, order)
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		var orderID uuid.UUID
		err = db.GetContext(ctx, &orderID, `SELECT id FROM "order" WHERE id_user = $1 LIMIT 1`, userID)
		if err != nil {
			t.Fatalf("Failed to get order ID: %v", err)
		}

		retrieved, err := orderRepo.GetById(ctx, orderID)
		if err != nil {
			t.Fatalf("Failed to get order: %v", err)
		}
		if retrieved.IdUser != userID {
			t.Errorf("Expected user ID %s, got %s", userID, retrieved.IdUser)
		}
		expectedPrice := 15.99 * 3
		if retrieved.Price != expectedPrice {
			t.Errorf("Expected price %.2f, got %.2f", expectedPrice, retrieved.Price)
		}
		items, err := orderRepo.GetItems(ctx, orderID)
		if err != nil {
			t.Fatalf("Failed to get order items: %v", err)
		}

		if len(items) != 1 {
			t.Errorf("Expected 1 order item, got %d", len(items))
		} else {
			if items[0].IdProduct != id_product {
				t.Errorf("Expected product ID %s, got %s", id_product, items[0].IdProduct)
			}
			if items[0].Amount != 3 {
				t.Errorf("Expected amount 3, got %d", items[0].Amount)
			}
		}
	})

	t.Run("GetFreeOrders", func(t *testing.T) {
		truncateTables(t)
		workerRepo := worker_rep.New(db)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Free Orders User",
			Date_of_birth: dob,
			Mail:          "freeorders@example.com",
			Password:      "password123",
			Phone:         "89016475842",
			Address:       "123 Free Orders St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		brand := structs.Brand{
			Name:          "Test Brand",
			Description:   "Test Brand Description",
			PriceCategory: "premium",
		}
		err = brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id_brand uuid.UUID
		err = db.Get(&id_brand, `select id from brand where name = $1`, "Test Brand")
		if err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		product := structs.Product{
			Name:        "Test Product",
			Description: "Test Product Description",
			Price:       10.99,
			Category:    "electronics",
			Amount:      100,
			IdBrand:     id_brand,
			PicLink:     "http://example.com/product.jpg",
			Articule:    "TP123",
		}
		err = productRepo.Create(ctx, product)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}

		var id_product uuid.UUID
		err = db.Get(&id_product, `select id from product where name = $1`, "Test Product")
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		orderStatuses := []string{"непринятый", "принятый", "непринятый"}
		orderPrices := []float64{50.0, 75.0, 100.0}

		for i := 0; i < 3; i++ {

			basket := structs.Basket{
				Id:     uuid.New(),
				IdUser: userID,
				Date:   time.Now(),
			}
			err = basketRepo.Create(ctx, basket)
			if err != nil {
				t.Fatalf("Failed to create basket: %v", err)
			}

			var id_basket uuid.UUID
			err = db.Get(&id_basket, `select id from basket where id_user = $1`, userID)
			if err != nil {
				t.Fatalf("Failed to get product id: %v", err)
			}

			basketItem := structs.BasketItem{
				IdProduct: id_product,
				IdBasket:  id_basket,
				Amount:    i + 1,
			}
			err = basketRepo.AddItem(ctx, basketItem)
			if err != nil {
				t.Fatalf("Failed to add item to basket: %v", err)
			}

			order := structs.Order{
				Date:    time.Now(),
				IdUser:  userID,
				Address: "123 Free Orders St",
				Status:  orderStatuses[i],
				Price:   orderPrices[i],
			}

			err = orderRepo.Create(ctx, order)
			if err != nil {
				t.Fatalf("Failed to create order: %v", err)
			}

			if orderStatuses[i] == "принятый" {
				worker := structs.Worker{
					IdUser:   userID,
					JobTitle: "Test Worker",
				}
				err = workerRepo.Create(ctx, worker)
				if err != nil {
					t.Fatalf("Failed to create worker: %v", err)
				}

				var orderID uuid.UUID
				err = db.GetContext(ctx, &orderID,
					`SELECT id FROM "order" WHERE id_user = $1 ORDER BY date DESC LIMIT 1`,
					userID)
				if err != nil {
					t.Fatalf("Failed to get order ID: %v", err)
				}

				err = workerRepo.AcceptOrder(ctx, orderID, userID)
				if err != nil {
					t.Fatalf("Failed to accept order: %v", err)
				}
			}
		}

		freeOrders, err := orderRepo.GetFreeOrders(ctx)
		if err != nil {
			t.Fatalf("Failed to get free orders: %v", err)
		}

		if len(freeOrders) != 2 {
			t.Errorf("Expected 2 free orders, got %d", len(freeOrders))
		}

		for _, order := range freeOrders {
			if order.Status != "непринятый" {
				t.Errorf("Expected status 'непринятый', got '%s'", order.Status)
			}
		}
	})
}

func TestProductRepository(t *testing.T) {
	ctx := context.Background()
	productRepo := product_rep.New(db)
	brandRepo := brand_rep.New(db)

	t.Run("Create and Get Product", func(t *testing.T) {
		truncateTables(t)

		brand := structs.Brand{
			Name:          "Test Brand",
			Description:   "Test Brand Description",
			PriceCategory: "premium",
		}

		err := brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id_brand uuid.UUID
		err = db.Get(&id_brand, `select id from brand where name = $1`, "Test Brand")
		if err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		product := structs.Product{
			Id:          uuid.New(),
			Name:        "Test Product",
			Description: "Test Product Description",
			Price:       25.99,
			Category:    "electronics",
			Amount:      100,
			IdBrand:     id_brand,
			PicLink:     "http://example.com/testproduct.jpg",
			Articule:    "TP456",
		}

		err = productRepo.Create(ctx, product)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}

		var id_product uuid.UUID
		err = db.Get(&id_product, `select id from product where name = $1`, "Test Product")
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		retrieved, err := productRepo.GetById(ctx, id_product)
		if err != nil {
			t.Fatalf("Failed to get product by ID: %v", err)
		}

		if retrieved.Name != product.Name {
			t.Errorf("Expected name %s, got %s", product.Name, retrieved.Name)
		}
	})

	t.Run("GetByCategory", func(t *testing.T) {
		truncateTables(t)

		brand := structs.Brand{
			Name:          "Category Test Brand",
			Description:   "Category Test Brand Description",
			PriceCategory: "premium",
		}

		err := brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id_brand uuid.UUID
		err = db.Get(&id_brand, `select id from brand where name = $1`, "Category Test Brand")
		if err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		products := []structs.Product{
			{
				Name:        "Electronics Product 1",
				Description: "Electronics Description",
				Price:       99.99,
				Category:    "electronics",
				Amount:      50,
				IdBrand:     id_brand,
				PicLink:     "http://example.com/electronics1.jpg",
				Articule:    "ELEC1",
			},
			{
				Name:        "Clothing Product 1",
				Description: "Clothing Description",
				Price:       29.99,
				Category:    "clothing",
				Amount:      100,
				IdBrand:     id_brand,
				PicLink:     "http://example.com/clothing1.jpg",
				Articule:    "CLOTH1",
			},
			{
				Name:        "Electronics Product 2",
				Description: "Electronics Description",
				Price:       199.99,
				Category:    "electronics",
				Amount:      25,
				IdBrand:     id_brand,
				PicLink:     "http://example.com/electronics2.jpg",
				Articule:    "ELEC2",
			},
		}

		for _, product := range products {
			err = productRepo.Create(ctx, product)
			if err != nil {
				t.Fatalf("Failed to create product: %v", err)
			}
		}

		electronics, err := productRepo.GetByCategory(ctx, "electronics")
		if err != nil {
			t.Fatalf("Failed to get products by category: %v", err)
		}

		if len(electronics) != 2 {
			t.Errorf("Expected 2 electronics products, got %d", len(electronics))
		}

		for _, product := range electronics {
			if product.Category != "electronics" {
				t.Errorf("Expected category 'electronics', got '%s'", product.Category)
			}
		}
	})
}

func TestReviewRepository(t *testing.T) {
	ctx := context.Background()
	userRepo := user_rep.New(db)
	productRepo := product_rep.New(db)
	brandRepo := brand_rep.New(db)
	reviewRepo := review_rep.New(db)

	t.Run("Create and Get Review", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Review Test User",
			Date_of_birth: dob,
			Mail:          "review@example.com",
			Password:      "password123",
			Phone:         "89016475843",
			Address:       "123 Review St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		brand := structs.Brand{
			Name:          "Review Test Brand",
			Description:   "Review Test Brand Description",
			PriceCategory: "premium",
		}

		err = brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id_brand uuid.UUID
		err = db.Get(&id_brand, `select id from brand where name = $1`, "Review Test Brand")
		if err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		product := structs.Product{
			Name:        "Review Test Product",
			Description: "Review Test Product Description",
			Price:       49.99,
			Category:    "electronics",
			Amount:      75,
			IdBrand:     id_brand,
			PicLink:     "http://example.com/reviewproduct.jpg",
			Articule:    "RTP123",
		}

		err = productRepo.Create(ctx, product)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}

		var id_product uuid.UUID
		err = db.Get(&id_product, `select id from product where name = $1`, "Review Test Product")
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		review := structs.Review{
			IdProduct: id_product,
			IdUser:    userID,
			Rating:    5,
			Text:      "Great product!",
			Date:      time.Now(),
		}

		err = reviewRepo.Create(ctx, review)
		if err != nil {
			t.Fatalf("Failed to create review: %v", err)
		}

		var reviewID uuid.UUID
		err = db.GetContext(ctx, &reviewID, `SELECT id FROM review WHERE id_product = $1 AND id_user = $2`,
			id_product, userID)
		if err != nil {
			t.Fatalf("Failed to get review ID: %v", err)
		}

		retrieved, err := reviewRepo.GetById(ctx, reviewID)
		if err != nil {
			t.Fatalf("Failed to get review: %v", err)
		}

		if retrieved.IdProduct != id_product {
			t.Errorf("Expected product ID %s, got %s", id_product, retrieved.IdProduct)
		}
		if retrieved.Rating != 5 {
			t.Errorf("Expected rating 5, got %d", retrieved.Rating)
		}
	})

	t.Run("ReviewsForProduct", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Multi Review User",
			Date_of_birth: dob,
			Mail:          "multireview@example.com",
			Password:      "password123",
			Phone:         "89026475843",
			Address:       "123 Multi Review St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		brand := structs.Brand{
			Id:            uuid.New(),
			Name:          "Multi Review Brand",
			Description:   "Multi Review Brand Description",
			PriceCategory: "premium",
		}

		err = brandRepo.Create(ctx, brand)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		var id_brand uuid.UUID
		err = db.Get(&id_brand, `select id from brand where name = $1`, "Multi Review Brand")
		if err != nil {
			t.Fatalf("Failed to get brand id: %v", err)
		}

		product := structs.Product{
			Name:        "Multi Review Product",
			Description: "Multi Review Product Description",
			Price:       79.99,
			Category:    "electronics",
			Amount:      50,
			IdBrand:     id_brand,
			PicLink:     "http://example.com/multireview.jpg",
			Articule:    "MRP123",
		}

		err = productRepo.Create(ctx, product)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}

		var id_product uuid.UUID
		err = db.Get(&id_product, `select id from product where name = $1`, "Multi Review Product")
		if err != nil {
			t.Fatalf("Failed to get product id: %v", err)
		}

		reviews := []structs.Review{
			{
				IdProduct: id_product,
				IdUser:    userID,
				Rating:    5,
				Text:      "Excellent!",
				Date:      time.Now().Add(-24 * time.Hour),
			},
			{
				IdProduct: id_product,
				IdUser:    userID,
				Rating:    4,
				Text:      "Very good",
				Date:      time.Now(),
			},
		}

		for _, review := range reviews {
			err = reviewRepo.Create(ctx, review)
			if err != nil {
				t.Fatalf("Failed to create review: %v", err)
			}
		}

		productReviews, err := reviewRepo.ReviewsForProduct(ctx, id_product)
		if err != nil {
			t.Fatalf("Failed to get reviews for product: %v", err)
		}

		if len(productReviews) != 2 {
			t.Errorf("Expected 2 reviews, got %d", len(productReviews))
		} else {
			if productReviews[0].Rating != 4 {
				t.Errorf("Expected first review rating 4 (newest), got %d", productReviews[0].Rating)
			}
			if productReviews[1].Rating != 5 {
				t.Errorf("Expected second review rating 5 (oldest), got %d", productReviews[1].Rating)
			}
		}
	})
}

func TestAuthRepository(t *testing.T) {
	ctx := context.Background()
	authRepo := auth_rep.New(db)
	userRepo := user_rep.New(db)

	t.Run("Create and Verify Token", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Token Test User",
			Date_of_birth: dob,
			Mail:          "token@example.com",
			Password:      "password123",
			Phone:         "89016475899",
			Address:       "123 Token St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		testToken := "test_refresh_token_" + time.Now().Format(time.RFC3339Nano)
		err = authRepo.CreateToken(ctx, userID, testToken)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		tokenID, err := authRepo.VerifyToken(ctx, testToken)
		if err != nil {
			t.Fatalf("Failed to verify token: %v", err)
		}
		if tokenID == uuid.Nil {
			t.Error("Expected valid token ID, got nil")
		}
	})

	t.Run("Verify Non-Existent Token", func(t *testing.T) {
		truncateTables(t)

		_, err := authRepo.VerifyToken(ctx, "non_existing_token")
		if err == nil {
			t.Fatal("Expected error for non-existent token, got nil")
		}
	})

	t.Run("Delete Token", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Delete Token User",
			Date_of_birth: dob,
			Mail:          "delete@example.com",
			Password:      "password123",
			Phone:         "89016475888",
			Address:       "123 Delete St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		testToken := "test_token_to_delete_" + time.Now().Format(time.RFC3339Nano)
		err = authRepo.CreateToken(ctx, userID, testToken)
		if err != nil {
			t.Fatalf("Failed to create token: %v", err)
		}

		tokenID, err := authRepo.VerifyToken(ctx, testToken)
		if err != nil {
			t.Fatalf("Failed to verify token: %v", err)
		}

		err = authRepo.DeleteToken(ctx, tokenID)
		if err != nil {
			t.Fatalf("Failed to delete token: %v", err)
		}

		_, err = authRepo.VerifyToken(ctx, testToken)
		if err == nil {
			t.Fatal("Expected error after token deletion, got nil")
		}
	})

	t.Run("Check Admin Role", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Admin Test User",
			Date_of_birth: dob,
			Mail:          "admin@example.com",
			Password:      "password123",
			Phone:         "89016475877",
			Address:       "123 Admin St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		isAdmin := authRepo.CheckAdmin(ctx, userID)
		if isAdmin {
			t.Error("Expected user not to be admin")
		}

		_, err = db.Exec(`INSERT INTO worker (id_user, job_title) VALUES ($1, 'admin')`, userID)
		if err != nil {
			t.Fatalf("Failed to make user admin: %v", err)
		}

		isAdmin = authRepo.CheckAdmin(ctx, userID)
		if !isAdmin {
			t.Error("Expected user to be admin")
		}
	})

	t.Run("Check Worker Role", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Worker Test User",
			Date_of_birth: dob,
			Mail:          "worker@example.com",
			Password:      "password123",
			Phone:         "89016475866",
			Address:       "123 Worker St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		isWorker := authRepo.CheckWorker(ctx, userID)
		if isWorker {
			t.Error("Expected user not to be worker")
		}

		_, err = db.Exec(`INSERT INTO worker (id_user, job_title) VALUES ($1, 'worker')`, userID)
		if err != nil {
			t.Fatalf("Failed to make user worker: %v", err)
		}

		isWorker = authRepo.CheckWorker(ctx, userID)
		if !isWorker {
			t.Error("Expected user to be worker")
		}
	})

	t.Run("Check Worker with Different Role", func(t *testing.T) {
		truncateTables(t)

		dob, _ := time.Parse("2006-01-02", "1990-01-01")
		user := structs.User{
			Name:          "Manager Test User",
			Date_of_birth: dob,
			Mail:          "manager@example.com",
			Password:      "password123",
			Phone:         "89016475855",
			Address:       "123 Manager St",
			Status:        "active",
			Role:          "customer",
		}

		userID, err := userRepo.Create(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		_, err = db.Exec(`INSERT INTO worker (id_user, job_title) VALUES ($1, 'manager')`, userID)
		if err != nil {
			t.Fatalf("Failed to make user manager: %v", err)
		}

		isWorker := authRepo.CheckWorker(ctx, userID)
		if isWorker {
			t.Error("Expected manager not to be recognized as worker")
		}
	})
}
