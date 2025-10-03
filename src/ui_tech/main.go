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
	cyan := color.New(color.FgCyan).Add(color.Bold)
	yellow := color.New(color.FgYellow)
	client := color.New(color.FgGreen)
	worker := color.New(color.FgBlue)
	admin := color.New(color.FgRed)
	white := color.New(color.FgWhite)

	cyan.Println("\n╔════════════════════════════════════════════╗")
	cyan.Println("║           WELCOME TO VIRTUAL' CLI          ║")
	cyan.Println("╚════════════════════════════════════════════╝")

	white.Println("\nAvailable commands:")

	yellow.Println("\n  [ Auth ]")
	client.Printf("  %-20s", "signup")
	white.Println("- Register new account")
	client.Printf("  %-20s", "login")
	white.Println("- Log into your account")
	client.Printf("  %-20s", "logout")
	white.Println("- Log out from current account")

	yellow.Println("\n  [ User ]")
	admin.Printf("  %-20s", "get-user-email")
	white.Println("- Get user by email")
	admin.Printf("  %-20s", "get-user-phone")
	white.Println("- Get user by phone")
	admin.Printf("  %-20s", "get-all-user")
	white.Println("- Get all users")

	yellow.Println("\n  [ Basket ]")
	admin.Printf("  %-20s", "get-basket")
	white.Println("- View your basket")
	client.Printf("  %-20s", "get-basket-items")
	white.Println("- View items in your basket")
	client.Printf("  %-20s", "add-to-basket")
	white.Println("- Add product to basket")
	client.Printf("  %-20s", "delete-from-basket")
	white.Println("- Remove product from basket")
	client.Printf("  %-20s", "update-item-amount")
	white.Println("- Change product quantity in basket")

	yellow.Println("\n  [ Favourites ]")
	client.Printf("  %-20s", "get-favourites-items")
	white.Println("- View items in your favourites")
	client.Printf("  %-20s", "add-to-favourites")
	white.Println("- Add product to favourites")
	client.Printf("  %-20s", "delete-from-favourites")
	white.Println("- Remove product from favourites")

	yellow.Println("\n  [ Brand ]")
	admin.Printf("  %-20s", "create-brand")
	white.Println("- Add new brand")
	client.Printf("  %-20s", "get-brands-category")
	white.Println("- Get brands by category")
	admin.Printf("  %-20s", "get-brand")
	white.Println("- Get brand by its id")
	admin.Printf("  %-20s", "delete-brand")
	white.Println("- Delete brand from brand list")

	yellow.Println("\n  [ Products ]")
	admin.Printf("  %-20s", "create-product")
	white.Println("- Add new product (admin)")
	admin.Printf("  %-20s", "delete-product")
	white.Println("- Remove product (admin)")
	client.Printf("  %-20s", "get-product")
	white.Println("- View product details")
	client.Printf("  %-20s", "get-products-brand")
	white.Println("- List products by brand")
	client.Printf("  %-20s", "get-products-category")
	white.Println("- List products by category")

	yellow.Println("\n  [ Orders ]")
	client.Printf("  %-20s", "create-order")
	white.Println("- Create order from basket")
	admin.Printf("  %-20s", "get-order")
	white.Println("- View order details")
	client.Printf("  %-20s", "get-order-items")
	white.Println("- View items in order")
	worker.Printf("  %-20s", "get-free-orders")
	white.Println("- List available orders (worker)")
	worker.Printf("  %-20s", "get-user-orders")
	white.Println("- List done orders")
	worker.Printf("  %-20s", "change-order-status")
	white.Println("- Update order status (admin/worker)")
	admin.Printf("  %-20s", "delete-order")
	white.Println("- Cancel order")

	yellow.Println("\n  [ Reviews ]")
	client.Printf("  %-20s", "create-review")
	white.Println("- Add review for product")
	admin.Printf("  %-20s", "get-review")
	white.Println("- View review details")
	client.Printf("  %-20s", "get-reviews-product")
	white.Println("- List reviews for product")
	admin.Printf("  %-20s", "delete-review")
	white.Println("- Remove review")

	yellow.Println("\n  [ Workers ]")
	admin.Printf("  %-20s", "create-worker")
	white.Println("- Add worker (admin)")
	admin.Printf("  %-20s", "delete-worker")
	white.Println("- Remove worker (admin)")
	admin.Printf("  %-20s", "get-worker-id")
	white.Println("- View worker details")
	admin.Printf("  %-20s", "get-workers")
	white.Println("- List client workers")
	worker.Printf("  %-20s", "accept-order")
	white.Println("- Take order for delivery")
	worker.Printf("  %-20s", "get-my-order")
	white.Println("- View your assigned orders")

	yellow.Println("\n  [ System ]")
	client.Printf("  %-20s", "exit")
	white.Println("- Quit the application")

	color.New(color.FgHiBlack).Println("\nType a command and press Enter. For most commands you'll be prompted for additional input.")
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

		case cmd == "exit":
			color.New(color.FgHiMagenta).Println("Goodbye!")
			return

		case cmd == "help":
			printWelcome()

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

		case strings.HasPrefix(cmd, "get-favourites-items"):
			client_favourites.GetFavouritesItems(client)

		case strings.HasPrefix(cmd, "add-to-favourites"):
			client_favourites.AddToFavourites(client, reader)

		case strings.HasPrefix(cmd, "delete-from-favourites"):
			client_favourites.DeleteFromFavourites(client, reader)

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

		case strings.HasPrefix(cmd, "get-all-user"):
			client_user.GetAllUsers(client, reader)

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

		case strings.HasPrefix(cmd, "get-user-order"):
			client_order.GetOrdersByUser(client)

		default:
			color.New(color.FgRed).Println("Unknown command. Type 'help' to see available commands.")
		}
	}
}
