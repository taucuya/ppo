package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	controller "github.com/taucuya/ppo/internal/controllers"
	"github.com/taucuya/ppo/internal/core/service/auth"
	"github.com/taucuya/ppo/internal/core/service/basket"
	"github.com/taucuya/ppo/internal/core/service/brand"
	"github.com/taucuya/ppo/internal/core/service/order"
	"github.com/taucuya/ppo/internal/core/service/product"
	"github.com/taucuya/ppo/internal/core/service/review"
	"github.com/taucuya/ppo/internal/core/service/user"
	"github.com/taucuya/ppo/internal/core/service/worker"
	auth_prov "github.com/taucuya/ppo/internal/providers/jwt/auth"
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

func main() {
	dsn := "postgres://test_user:test_password@localhost:5433/test_db?sslmode=disable"

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	_ = runSQLScripts(db, []string{
		"/home/taya/Desktop/ppoft/src/internal/database/sql/01-create.sql",
		"/home/taya/Desktop/ppoft/src/internal/database/sql/02-constraints.sql",
		"/home/taya/Desktop/ppoft/src/internal/database/sql/03-inserts.sql",
		"/home/taya/Desktop/ppoft/src/internal/database/sql/trigger_accept.sql",
		"/home/taya/Desktop/ppoft/src/internal/database/sql/trigger_order.sql",
	})

	key := []byte(uuid.New().String())

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Cant open log file: %v", err)
	}

	log.SetOutput(logFile)

	gin.DefaultWriter = logFile

	ar := auth_rep.New(db)
	ap := auth_prov.New(key, time.Duration(15*time.Minute), time.Duration(7*24*time.Hour))
	bar := basket_rep.New(db)
	brr := brand_rep.New(db)
	or := order_rep.New(db)
	pr := product_rep.New(db)
	rr := review_rep.New(db)
	ur := user_rep.New(db)
	wr := worker_rep.New(db)

	bas := basket.New(bar)
	us := user.New(ur, bas)
	as := auth.New(ap, ar, us)
	brs := brand.New(brr)
	oss := order.New(or)
	ps := product.New(pr)
	rs := review.New(rr)
	ws := worker.New(wr)

	c := controller.Controller{
		BasketService:  *bas,
		UserService:    *us,
		AuthServise:    *as,
		BrandService:   *brs,
		OrderService:   *oss,
		ProductService: *ps,
		ReviewService:  *rs,
		WorkerService:  *ws,
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", c.SignupHandler)
			auth.POST("/login", c.LoginHandler)
			auth.POST("/logout", c.LogoutHandler)
		}

		basket := api.Group("/basket")
		{
			basket.GET("/items", c.GetBasketItemsHandler)
			basket.POST("", c.AddBasketItemHandler)
			basket.GET("", c.GetBasketByIdHandler)
			basket.DELETE("", c.DeleteBasketItemHandler)
			basket.PUT("", c.UpdateBasketItemAmountHandler)
		}

		brand := api.Group("/brand")
		{
			brand.POST("", c.CreateBrandHandler)
			brand.DELETE("/:id", c.DeleteBrandHandler)
			brand.GET("/:id", c.GetBrandByIdHandler)
			brand.GET("/category/:cat", c.GetAllBrandsInCategoryHander)
		}

		order := api.Group("/order")
		{
			order.POST("", c.CreateOrderHandler)
			order.GET("/:id/items", c.GetOrderItemsHandler)
			order.GET("/freeorders", c.GetFreeOrdersHandler)
			order.GET("/:id", c.GetOrderByIdHandler)
			order.PUT("/:id/status", c.ChangeOrderStatusHandler)
			order.DELETE("/:id", c.DeleteOrderHandler)
			order.PUT("/:id", c.AcceptOrderHandler)
		}

		product := api.Group("/product")
		{
			product.POST("", c.CreateProductHandler)
			product.DELETE("/:id", c.DeleteProductHandler)
			product.GET("/", c.GetProductHandler)
			product.GET("/category/:category", c.GetProductsByCategoryHandler)
			product.GET("/brand/:brand", c.GetProductsByBrandHandler)
			product.GET("/reviews/:id", c.GetReviewsForProductHandler)
		}

		review := api.Group("/review")
		{
			review.POST("product/:id_product", c.CreateReviewHandler)
			review.GET("/:id", c.GetReviewByIdHandler)
			review.DELETE("/:id", c.DeleteReviewHandler)
		}

		user := api.Group("/user")
		{
			user.GET("/email", c.GetUserByEmailHandler)
			user.GET("/users", c.GetAllUsersHandler)
			user.GET("/phone", c.GetUserByPhoneHandler)
		}

		worker := api.Group("/worker")
		{
			worker.POST("", c.CreateWorkerHandler)
			worker.GET("", c.GetAllWorkersHandler)
			worker.POST("/accept", c.AcceptOrderHandler)
			worker.GET("/:id", c.GetWorkerByIdHandler)
			worker.DELETE("/:id", c.DeleteWorkerHandler)
			worker.GET("/orders", c.GetWorkerOrders)
		}
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}
