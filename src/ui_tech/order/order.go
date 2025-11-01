package order

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func CreateOrder(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter address for the order: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	payload := map[string]string{"address": address}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/users/me/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed")
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("ERROR: %s\n", extractErrorMessage(body))
		return
	}

	fmt.Println("SUCCESS: Order created")
}

func GetOrderById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID: ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	if _, err := uuid.Parse(orderID); err != nil {
		fmt.Println("ERROR: Invalid order ID")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/users/me/orders/%s", orderID)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("ERROR: Request failed")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: %s\n", extractErrorMessage(body))
		return
	}

	var order map[string]interface{}
	json.Unmarshal(body, &order)

	fmt.Println("SUCCESS: Order details")
	fmt.Printf("ID: %v\nStatus: %v\nAddress: %v\nPrice: %v\n",
		order["Id"], order["Status"], order["Address"], order["Price"])
}

func GetOrderItems(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID: ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	if _, err := uuid.Parse(orderID); err != nil {
		fmt.Println("ERROR: Invalid order ID")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/users/me/orders/%s/items", orderID)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("ERROR: Request failed")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: %s\n", extractErrorMessage(body))
		return
	}

	var items []map[string]interface{}
	json.Unmarshal(body, &items)

	fmt.Printf("SUCCESS: %d items found\n", len(items))
	for i, item := range items {
		fmt.Printf("%d. Product: %v, Amount: %v\n", i+1, item["IdProduct"], item["Amount"])
	}
}

func GetFreeOrders(client *http.Client) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/users/me/orders?status=непринятый", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: %s\n", extractErrorMessage(body))
		return
	}

	var orders []map[string]interface{}
	json.Unmarshal(body, &orders)

	fmt.Printf("SUCCESS: %d free orders\n", len(orders))
	for i, order := range orders {
		fmt.Printf("%d. ID: %v, Address: %v, Price: %v\n",
			i+1, order["Id"], order["Address"], order["Price"])
	}
}

func ChangeOrderStatus(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID: ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	if _, err := uuid.Parse(orderID); err != nil {
		fmt.Println("ERROR: Invalid order ID")
		return
	}

	fmt.Print("Enter new status: ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)

	url := fmt.Sprintf("http://localhost:8080/api/v1/users/me/orders/%s?status=%s", orderID, status)
	req, _ := http.NewRequest("PATCH", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: %s\n", extractErrorMessage(body))
		return
	}

	fmt.Println("SUCCESS: Status updated")
}

func DeleteOrder(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID: ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	if _, err := uuid.Parse(orderID); err != nil {
		fmt.Println("ERROR: Invalid order ID")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/users/me/orders/%s", orderID)
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: %s\n", extractErrorMessage(body))
		return
	}

	fmt.Println("SUCCESS: Order deleted")
}

func GetOrdersByUser(client *http.Client) {
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/users/me/orders", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: %s\n", extractErrorMessage(body))
		return
	}

	var orders []map[string]interface{}
	json.Unmarshal(body, &orders)

	fmt.Printf("SUCCESS: %d orders\n", len(orders))
	for i, order := range orders {
		fmt.Printf("%d. ID: %v, Status: %v, Price: %v\n",
			i+1, order["Id"], order["Status"], order["Price"])
	}
}

func extractErrorMessage(body []byte) string {
	if len(body) == 0 {
		return "Unknown error"
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return string(body)
	}

	if errorMsg, ok := result["error"].(string); ok {
		return errorMsg
	}

	return string(body)
}
