package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/taucuya/ppo/internal/docs"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	controller "github.com/taucuya/ppo/internal/controllers"
	"github.com/taucuya/ppo/internal/core/service/auth"
	"github.com/taucuya/ppo/internal/core/service/basket"
	"github.com/taucuya/ppo/internal/core/service/brand"
	"github.com/taucuya/ppo/internal/core/service/favourites"
	"github.com/taucuya/ppo/internal/core/service/order"
	"github.com/taucuya/ppo/internal/core/service/product"
	"github.com/taucuya/ppo/internal/core/service/review"
	"github.com/taucuya/ppo/internal/core/service/user"
	"github.com/taucuya/ppo/internal/core/service/worker"
	auth_prov "github.com/taucuya/ppo/internal/providers/jwt/auth"
	auth_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/auth"
	basket_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/basket"
	brand_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/brand"
	favourites_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/favourites"
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
	log.Println("SQL completed")
	return nil
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения")
	}
}

func main() {

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Cant open log file: %v", err)
	}

	log.SetOutput(logFile)

	log.Println("HELLO")

	loadEnv()
	dsn := os.Getenv("DB_DSN")
	log.Println("DSN,", dsn)
	key := []byte(os.Getenv("JWT_SECRET"))
	acstime, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_LIFETIME_MINUTES"))
	if err != nil {
		return
	}
	reftime, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_LIFETIME_DAYS"))
	if err != nil {
		return
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	_ = runSQLScripts(db, []string{
		"/home/taya/Desktop/ppo/src/internal/database/sql/delete.sql",
		"/home/taya/Desktop/ppo/src/internal/database/sql/01-create.sql",
		"/home/taya/Desktop/ppo/src/internal/database/sql/02-constraints.sql",
		"/home/taya/Desktop/ppo/src/internal/database/sql/03-inserts.sql",
		"/home/taya/Desktop/ppo/src/internal/database/sql/trigger_accept.sql",
		"/home/taya/Desktop/ppo/src/internal/database/sql/trigger_order.sql",
	})

	gin.DefaultWriter = logFile

	ar := auth_rep.New(db)
	ap := auth_prov.New(key, time.Duration(time.Duration(acstime)*time.Minute), time.Duration(time.Duration(reftime)*24*time.Hour))
	bar := basket_rep.New(db)
	brr := brand_rep.New(db)
	fr := favourites_rep.New(db)
	or := order_rep.New(db)
	pr := product_rep.New(db)
	rr := review_rep.New(db)
	ur := user_rep.New(db)
	wr := worker_rep.New(db)

	bas := basket.New(bar)
	fs := favourites.New(fr)
	us := user.New(ur, bas, fs)
	as := auth.New(ap, ar, us)
	brs := brand.New(brr)
	oss := order.New(or)
	ps := product.New(pr)
	rs := review.New(rr)
	ws := worker.New(wr)

	c := controller.Controller{
		BasketService:     *bas,
		UserService:       *us,
		AuthServise:       *as,
		BrandService:      *brs,
		FavouritesService: *fs,
		OrderService:      *oss,
		ProductService:    *ps,
		ReviewService:     *rs,
		WorkerService:     *ws,
	}

	router := gin.New()
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

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

		users := api.Group("/users")
		{
			users.GET("", c.GetUserByPrivatesHandler)

			me := users.Group("/me")
			{
				basket := me.Group("/basket")
				{
					basket.GET("", c.GetBasketByIdHandler)
					basketItems := basket.Group("/items")
					{
						basketItems.GET("", c.GetBasketItemsHandler)
						basketItems.POST("", c.AddBasketItemHandler)
						basketItems.DELETE("", c.DeleteBasketItemHandler)
						basketItems.PATCH("", c.UpdateBasketItemAmountHandler)
					}
				}

				favourite := me.Group("/favourite")
				{
					favouriteItems := favourite.Group("/items")
					{
						favouriteItems.GET("", c.GetFavouritesHandler)
						favouriteItems.POST("", c.AddFavouritesItemHandler)
						favouriteItems.DELETE("/:id_product", c.DeleteFavouritesItemHandler)
					}
				}

				orders := me.Group("/orders")
				{
					orders.GET("", c.GetOrdersHandler)
					orders.POST("", c.CreateOrderHandler)
					orders.GET("/:id", c.GetOrderByIdHandler)
					orders.PATCH("/:id", c.ChangeOrderStatusHandler)
					orders.DELETE("/:id", c.DeleteOrderHandler)
					orders.GET("/:id/items", c.GetOrderItemsHandler)
				}

				products := me.Group("/products")
				{
					products.POST("/:id_product/reviews", c.CreateReviewHandler)
				}
			}
		}

		brands := api.Group("/brands")
		{
			brands.GET("", c.GetAllBrandsInCategoryHander)
			brands.POST("", c.CreateBrandHandler)
			brands.GET("/:id", c.GetBrandByIdHandler)
			brands.DELETE("/:id", c.DeleteBrandHandler)
		}

		products := api.Group("/products")
		{
			products.GET("", c.GetProductsHandler)
			products.POST("", c.CreateProductHandler)
			products.DELETE("/:id", c.DeleteProductHandler)
			products.GET("/:id/reviews", c.GetReviewsForProductHandler)
			products.GET("/:id/reviews/:id", c.GetReviewByIdHandler)
			products.DELETE("/:id/reviews/:id", c.DeleteReviewHandler)
		}

		workers := api.Group("/workers")
		{
			workers.GET("", c.GetAllWorkersHandler)
			workers.POST("", c.CreateWorkerHandler)
			workers.GET("/:id", c.GetWorkerByIdHandler)
			workers.DELETE("/:id", c.DeleteWorkerHandler)

			me := workers.Group("/me")
			{
				orders := me.Group("/orders")
				{
					orders.GET("", c.GetWorkerOrders)
					orders.POST("", c.AcceptOrderHandler)
				}
			}
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
