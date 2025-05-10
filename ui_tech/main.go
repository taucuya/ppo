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
			parts := strings.Fields(cmd)
			if len(parts) > 1 {
				fmt.Println("Usage: get-basket-items")
				continue
			}
			client_basket.GetBasketItems(client)
		case strings.HasPrefix(cmd, "get-basket"):
			parts := strings.Fields(cmd)
			if len(parts) > 1 {
				fmt.Println("Usage: get-basket")
				continue
			}
			client_basket.GetBasket(client)
		case strings.HasPrefix(cmd, "add-to-basket"):
			parts := strings.Fields(cmd)
			if len(parts) != 3 {
				fmt.Println("Usage: add-to-basket <product_id> <amount>")
				continue
			}
			client_basket.AddToBasket(client, parts[1], parts[2])
		case strings.HasPrefix(cmd, "delete-from-basket"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: delete-from-basket <product_id>")
				continue
			}
			client_basket.DeleteFromBasket(client, parts[1])
		case strings.HasPrefix(cmd, "update-item-amount"):
			parts := strings.Fields(cmd)
			if len(parts) != 3 {
				fmt.Println("Usage: update-item-amount <product_id> <amount>")
				continue
			}
			client_basket.UpdateItemAmount(client, parts[1], parts[2])
		case strings.HasPrefix(cmd, "create-brand"):
			parts := strings.Fields(cmd)
			if len(parts) != 1 {
				fmt.Println("Usage: create-brand")
				continue
			}
			client_brand.CreateBrand(client, reader)

		case strings.HasPrefix(cmd, "get-brands-category"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: get-brands-category <category>")
				continue
			}
			client_brand.GetBrandsByCategory(client, parts[1])

		case strings.HasPrefix(cmd, "get-brand"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: get-brand <brand_id>")
				continue
			}
			client_brand.GetBrandById(client, parts[1])

		case strings.HasPrefix(cmd, "delete-brand"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: delete-brand <brand_id>")
				continue
			}
			client_brand.DeleteBrand(client, parts[1])

		case strings.HasPrefix(cmd, "create-order"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: create-order <address>")
				continue
			}
			client_order.CreateOrder(client, parts[1])

		case strings.HasPrefix(cmd, "get-order"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: get-order <order_id>")
				continue
			}
			client_order.GetOrderById(client, parts[1])

		case strings.HasPrefix(cmd, "get-free-orders"):
			client_order.GetFreeOrders(client)

		case strings.HasPrefix(cmd, "get-order-items"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: get-order-items <order_id>")
				continue
			}
			client_order.GetOrderItems(client, parts[1])

		case strings.HasPrefix(cmd, "change-order-status"):
			parts := strings.Fields(cmd)
			if len(parts) != 3 {
				fmt.Println("Usage: change-order-status <order_id> <status>")
				continue
			}
			client_order.ChangeOrderStatus(client, parts[1], parts[2])

		case strings.HasPrefix(cmd, "delete-order"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: delete-order <order_id>")
				continue
			}
			client_order.DeleteOrder(client, parts[1])

		case strings.HasPrefix(cmd, "create-product"):
			client_product.CreateProduct(client, reader)

		case strings.HasPrefix(cmd, "delete-product"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: delete-product <product_id>")
				continue
			}
			client_product.DeleteProduct(client, parts[1])

		case strings.HasPrefix(cmd, "get-products-brand"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: get-products-brand <brand>")
				continue
			}
			client_product.GetProductsByBrand(client, parts[1])

		case strings.HasPrefix(cmd, "get-products-category"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: get-products-category <category>")
				continue
			}
			client_product.GetProductsByCategory(client, parts[1])

		case strings.HasPrefix(cmd, "get-product"):
			parts := strings.Fields(cmd)
			if len(parts) != 2 {
				fmt.Println("Usage: get-product <prod_pole>")
				continue
			}
			client_product.GetProduct(client, parts[1])

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
