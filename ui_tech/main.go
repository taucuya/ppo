package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	client_auth "github.com/taucuya/ppo/ui_tech/auth"
	client_basket "github.com/taucuya/ppo/ui_tech/basket"
	client_brand "github.com/taucuya/ppo/ui_tech/brand"
	client_order "github.com/taucuya/ppo/ui_tech/order"
	client_product "github.com/taucuya/ppo/ui_tech/product"
	client_review "github.com/taucuya/ppo/ui_tech/review"
	client_user "github.com/taucuya/ppo/ui_tech/user"
	client_worker "github.com/taucuya/ppo/ui_tech/worker"
)

func mustCookieJar() http.CookieJar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("failed to create cookie jar: %v", err)
	}
	return jar
}

var client = &http.Client{
	Jar: mustCookieJar(),
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		cmdLine, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(cmdLine)

		switch {
		case cmd == "signup":
			client_auth.Signup(client, reader)

		case cmd == "login":
			client_auth.Login(client, reader)

		case cmd == "logout":
			client_auth.Logout(client)

		case cmd == "exit":
			return

		case strings.HasPrefix(cmd, "get-basket-items"):
			client_basket.GetBasketItems(client)

		case strings.HasPrefix(cmd, "get-basket"):
			client_basket.GetBasket(client)

		case strings.HasPrefix(cmd, "add-to-basket"):
			client_basket.AddToBasket(client, reader)

		case strings.HasPrefix(cmd, "delete-from-basket"):
			client_basket.DeleteFromBasket(client, reader)

		case strings.HasPrefix(cmd, "update-item-amount"):
			client_basket.UpdateItemAmount(client, reader)

		case strings.HasPrefix(cmd, "create-brand"):
			client_brand.CreateBrand(client, reader)

		case strings.HasPrefix(cmd, "get-brands-category"):
			client_brand.GetBrandsByCategory(client, reader)

		case strings.HasPrefix(cmd, "get-brand"):
			client_brand.GetBrandById(client, reader)

		case strings.HasPrefix(cmd, "delete-brand"):
			client_brand.DeleteBrand(client, reader)

		case strings.HasPrefix(cmd, "create-order"):
			client_order.CreateOrder(client, reader)

		case strings.HasPrefix(cmd, "get-order"):
			client_order.GetOrderById(client, reader)

		case strings.HasPrefix(cmd, "get-free-orders"):
			client_order.GetFreeOrders(client)

		case strings.HasPrefix(cmd, "get-order-items"):
			client_order.GetOrderItems(client, reader)

		case strings.HasPrefix(cmd, "change-order-status"):
			client_order.ChangeOrderStatus(client, reader)

		case strings.HasPrefix(cmd, "delete-order"):
			client_order.DeleteOrder(client, reader)

		case strings.HasPrefix(cmd, "create-product"):
			client_product.CreateProduct(client, reader)

		case strings.HasPrefix(cmd, "delete-product"):
			client_product.DeleteProduct(client, reader)

		case strings.HasPrefix(cmd, "get-products-brand"):
			client_product.GetProductsByBrand(client, reader)

		case strings.HasPrefix(cmd, "get-products-category"):
			client_product.GetProductsByCategory(client, reader)

		case strings.HasPrefix(cmd, "get-product"):
			client_product.GetProduct(client, reader)

		case strings.HasPrefix(cmd, "create-review"):
			client_review.CreateReview(client, reader)

		case strings.HasPrefix(cmd, "get-reviews-product"):
			client_review.GetReviewsByProductId(client, reader)

		case strings.HasPrefix(cmd, "get-review"):
			client_review.GetReviewById(client, reader)

		case strings.HasPrefix(cmd, "delete-review"):
			client_review.DeleteReview(client, reader)

		case strings.HasPrefix(cmd, "get-user-email"):
			client_user.GetUserByEmail(client, reader)

		case strings.HasPrefix(cmd, "get-user-phone"):
			client_user.GetUserByPhone(client, reader)

		case strings.HasPrefix(cmd, "create-worker"):
			client_worker.CreateWorker(client, reader)

		case strings.HasPrefix(cmd, "delete-worker"):
			client_worker.DeleteWorker(client, reader)

		case strings.HasPrefix(cmd, "get-worker-id"):
			client_worker.GetWorkerById(client, reader)

		case strings.HasPrefix(cmd, "get-workers"):
			client_worker.GetAllWorkers(client)

		case strings.HasPrefix(cmd, "accept-order"):
			client_worker.AcceptOrder(client, reader)

		case strings.HasPrefix(cmd, "get-my-order"):
			client_worker.GetWorkerOrders(client)

		default:
			fmt.Println("Unknown command")
		}
	}
}
