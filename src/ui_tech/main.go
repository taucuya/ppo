package main

import (
	"bufio"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/fatih/color"
	client_auth "github.com/taucuya/ppo/ui_tech/auth"
	client_basket "github.com/taucuya/ppo/ui_tech/basket"
	client_brand "github.com/taucuya/ppo/ui_tech/brand"
	client_favourites "github.com/taucuya/ppo/ui_tech/favourites"
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

func printWelcome() {
	yellow := color.New(color.FgYellow)
	client := color.New(color.FgGreen)
	worker := color.New(color.FgBlue)
	admin := color.New(color.FgRed)
	white := color.New(color.FgWhite)

	white.Println("\nAvailable commands:")

	yellow.Println("\n  [ Authentication ]")
	client.Printf("  %-25s", "signup")
	white.Println("- Register new user account")
	client.Printf("  %-25s", "login")
	white.Println("- Log into your account")
	client.Printf("  %-25s", "logout")
	white.Println("- Log out from current account")

	yellow.Println("\n  [ User Management ]")
	admin.Printf("  %-25s", "get-user-email")
	white.Println("- Get user by email (admin)")
	admin.Printf("  %-25s", "get-user-phone")
	white.Println("- Get user by phone (admin)")
	admin.Printf("  %-25s", "get-all-users")
	white.Println("- Get all users (admin)")

	yellow.Println("\n  [ Basket Operations ]")
	client.Printf("  %-25s", "get-basket")
	white.Println("- Get user's basket")
	client.Printf("  %-25s", "get-basket-items")
	white.Println("- Get items in basket")
	client.Printf("  %-25s", "add-to-basket")
	white.Println("- Add product to basket")
	client.Printf("  %-25s", "delete-from-basket")
	white.Println("- Remove product from basket")
	client.Printf("  %-25s", "update-item-amount")
	white.Println("- Update product quantity in basket")

	yellow.Println("\n  [ Favourites Management ]")
	client.Printf("  %-25s", "get-favourites")
	white.Println("- Get user's favourites")
	client.Printf("  %-25s", "add-to-favourites")
	white.Println("- Add product to favourites")
	client.Printf("  %-25s", "delete-from-favourites")
	white.Println("- Remove product from favourites")

	yellow.Println("\n  [ Brand Management ]")
	admin.Printf("  %-25s", "create-brand")
	white.Println("- Create new brand (admin)")
	admin.Printf("  %-25s", "get-brand-by-id")
	white.Println("- Get brand by ID (admin)")
	admin.Printf("  %-25s", "delete-brand")
	white.Println("- Delete brand (admin)")
	client.Printf("  %-25s", "get-brands-category")
	white.Println("- Get brands by price category")

	yellow.Println("\n  [ Product Management ]")
	admin.Printf("  %-25s", "create-product")
	white.Println("- Create new product (admin)")
	admin.Printf("  %-25s", "delete-product")
	white.Println("- Delete product (admin)")
	client.Printf("  %-25s", "get-product-by-id")
	white.Println("- Get product by ID")
	client.Printf("  %-25s", "get-product-by-art")
	white.Println("- Get product by articule")
	client.Printf("  %-25s", "get-products-brand")
	white.Println("- Get products by brand name")
	client.Printf("  %-25s", "get-products-category")
	white.Println("- Get products by category")
	client.Printf("  %-25s", "get-product-reviews")
	white.Println("- Get reviews for product")

	yellow.Println("\n  [ Order Management ]")
	client.Printf("  %-25s", "create-order")
	white.Println("- Create order from basket")
	client.Printf("  %-25s", "get-all-orders")
	white.Println("- Get orders")
	client.Printf("  %-25s", "get-user-orders")
	white.Println("- Get user's orders")
	// admin.Printf("  %-25s", "get-order-by-id")
	// white.Println("- Get order by ID (admin/worker)")
	client.Printf("  %-25s", "get-order-items")
	white.Println("- Get items in order")
	worker.Printf("  %-25s", "get-free-orders")
	white.Println("- Get available orders (worker)")
	worker.Printf("  %-25s", "change-order-status")
	white.Println("- Change order status (worker/admin)")
	admin.Printf("  %-25s", "delete-order")
	white.Println("- Delete order (admin)")

	yellow.Println("\n  [ Review Management ]")
	client.Printf("  %-25s", "create-review")
	white.Println("- Create review for product")
	client.Printf("  %-25s", "get-review-by-id")
	white.Println("- Get review by ID")
	client.Printf("  %-25s", "get-product-reviews")
	white.Println("- Get all reviews for product")
	admin.Printf("  %-25s", "delete-review")
	white.Println("- Delete review (admin)")

	yellow.Println("\n  [ Worker Management ]")
	admin.Printf("  %-25s", "create-worker")
	white.Println("- Create worker account (admin)")
	admin.Printf("  %-25s", "get-worker-by-id")
	white.Println("- Get worker by ID (admin)")
	admin.Printf("  %-25s", "get-all-workers")
	white.Println("- Get all workers (admin)")
	admin.Printf("  %-25s", "delete-worker")
	white.Println("- Delete worker (admin)")
	worker.Printf("  %-25s", "accept-order")
	white.Println("- Accept order for delivery (worker)")
	worker.Printf("  %-25s", "get-worker-orders")
	white.Println("- Get worker's assigned orders")

	yellow.Println("\n  [ System ]")
	client.Printf("  %-25s", "help")
	white.Println("- Show this help message")
	client.Printf("  %-25s", "exit")
	white.Println("- Quit the application")

	color.New(color.FgHiBlack).Println("\nLegend:")
	client.Print("  Client")
	white.Print(" - Available to all authenticated users | ")
	worker.Print("Worker")
	white.Print(" - Available to workers | ")
	admin.Print("Admin")
	white.Println(" - Available to administrators")

	color.New(color.FgHiBlack).Println("\nType a command and press Enter. You'll be prompted for additional input when needed.")
}

func main() {
	color.New(color.FgHiBlue).Println("Initializing PPO Shop CLI...")

	reader := bufio.NewReader(os.Stdin)
	printWelcome()

	for {
		color.New(color.FgCyan).Print("\n> ")
		cmdLine, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(cmdLine)

		switch {
		case cmd == "signup":
			client_auth.Signup(client, reader)
		case cmd == "login":
			client_auth.Login(client, reader)
		case cmd == "logout":
			client_auth.Logout(client)

		case cmd == "get-user-email":
			client_user.GetUserByEmail(client, reader)
		case cmd == "get-user-phone":
			client_user.GetUserByPhone(client, reader)
		case cmd == "get-all-users":
			client_user.GetAllUsers(client)

		case cmd == "get-basket":
			client_basket.GetBasket(client)
		case cmd == "get-basket-items":
			client_basket.GetBasketItems(client)
		case cmd == "add-to-basket":
			client_basket.AddToBasket(client, reader)
		case cmd == "delete-from-basket":
			client_basket.DeleteFromBasket(client, reader)
		case cmd == "update-item-amount":
			client_basket.UpdateItemAmount(client, reader)

		case cmd == "get-favourites":
			client_favourites.GetFavouritesItems(client)
		case cmd == "add-to-favourites":
			client_favourites.AddToFavourites(client, reader)
		case cmd == "delete-from-favourites":
			client_favourites.DeleteFromFavourites(client, reader)

		case cmd == "create-brand":
			client_brand.CreateBrand(client, reader)
		case cmd == "get-brand-by-id":
			client_brand.GetBrandById(client, reader)
		case cmd == "delete-brand":
			client_brand.DeleteBrand(client, reader)
		case cmd == "get-brands-category":
			client_brand.GetBrandsByCategory(client, reader)

		case cmd == "create-product":
			client_product.CreateProduct(client, reader)
		case cmd == "delete-product":
			client_product.DeleteProduct(client, reader)
		case cmd == "get-product-by-id":
			client_product.GetProductById(client, reader)
		case cmd == "get-product-by-art":
			client_product.GetProductByArticule(client, reader)
		case cmd == "get-products-brand":
			client_product.GetProductsByBrand(client, reader)
		case cmd == "get-products-category":
			client_product.GetProductsByCategory(client, reader)
		case cmd == "get-product-reviews":
			client_product.GetReviewsForProduct(client, reader)

		case cmd == "create-order":
			client_order.CreateOrder(client, reader)
		case cmd == "get-all-orders":
			client_order.GetAllOrders(client)
		case cmd == "get-user-orders":
			client_order.GetOrdersByUser(client)
		case cmd == "get-order-by-id":
			client_order.GetOrderById(client, reader)
		case cmd == "get-order-items":
			client_order.GetOrderItems(client, reader)
		case cmd == "get-free-orders":
			client_order.GetFreeOrders(client)
		case cmd == "change-order-status":
			client_order.ChangeOrderStatus(client, reader)
		case cmd == "delete-order":
			client_order.DeleteOrder(client, reader)

		case cmd == "create-review":
			client_review.CreateReview(client, reader)
		case cmd == "get-review-by-id":
			client_review.GetReviewById(client, reader)
		case cmd == "get-product-reviews":
			client_review.GetReviewsByProductId(client, reader)
		case cmd == "delete-review":
			client_review.DeleteReview(client, reader)

		case cmd == "create-worker":
			client_worker.CreateWorker(client, reader)
		case cmd == "get-worker-by-id":
			client_worker.GetWorkerById(client, reader)
		case cmd == "get-all-workers":
			client_worker.GetAllWorkers(client)
		case cmd == "delete-worker":
			client_worker.DeleteWorker(client, reader)
		case cmd == "accept-order":
			client_worker.AcceptOrder(client, reader)
		case cmd == "get-worker-orders":
			client_worker.GetWorkerOrders(client)

		case cmd == "help":
			printWelcome()
		case cmd == "exit":
			color.New(color.FgHiMagenta).Println("Thank you for using PPO Shop CLI. Goodbye!")
			return

		default:
			color.New(color.FgRed).Println("Unknown command. Type 'help' to see available commands.")
		}
	}
}
