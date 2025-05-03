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
	"github.com/taucuya/ppo/internal/core/service/user"
	auth_prov "github.com/taucuya/ppo/internal/providers/jwt/auth"
	auth_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/auth"
	user_rep "github.com/taucuya/ppo/internal/repository/postgres/reps/user"
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

	err = runSQLScripts(db, []string{
		"/home/taya/Desktop/ppoft/src/internal/database/sql/01-create.sql",
		"/home/taya/Desktop/ppoft/src/internal/database/sql/02-constraints.sql",
		"/home/taya/Desktop/ppoft/src/internal/database/sql/03-inserts.sql",
	})

	fmt.Println(err)

	key := []byte(uuid.New().String())

	ar := auth_rep.New(db)
	ap := auth_prov.New(key, time.Duration(15*time.Minute), time.Duration(7*24*time.Hour))
	// bar := basket_rep.New(db)
	// brr := brand_rep.New(db)
	// or := order_rep.New(db)
	// pr := product_rep.New(db)
	// rr := review_rep.New(db)
	ur := user_rep.New(db)
	// wr := worker_rep.New(db)

	us := user.New(ur)
	as := auth.New(ap, ar, us)
	// bas := basket.New(bar)
	// brs := brand.New(brr)
	// os := order.New(or)
	// ps := product.New(pr)
	// rs := review.New(rr)
	// ws := worker.New(wr)

	c := controller.Controller{
		UserService: *us,
		AuthServise: *as,
	}

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", c.SignupHandler)
			auth.POST("/login", c.LoginHandler)
			auth.POST("/logout", c.LogoutHandler)
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
