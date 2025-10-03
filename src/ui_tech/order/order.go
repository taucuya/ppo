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

	payload := map[string]string{
		"address": address,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/orders", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("❌ Request creation failed:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Order created successfully.")
	} else {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Println("❌ Error:", errResp["error"])
	}
}

func GetOrderById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID (UUID): ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)
	_, err := uuid.Parse(orderID)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	resp, err := client.Get("http://localhost:8080/api/v1/orders/" + orderID)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var order map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&order)

	fmt.Println("✅ Order Info:")
	fmt.Println("ID:      ", order["Id"])
	fmt.Println("Date:    ", order["Date"])
	fmt.Println("User ID: ", order["IdUser"])
	fmt.Println("Address: ", order["Address"])
	fmt.Println("Status:  ", order["Status"])
	fmt.Println("Price:   ", order["Price"])
}

func GetOrderItems(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID (UUID) to fetch items: ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	_, err := uuid.Parse(orderID)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	resp, err := client.Get("http://localhost:8080/api/v1/order/items/" + orderID)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
		return
	}

	var items []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	if len(items) == 0 {
		fmt.Println("No items found in the order.")
		return
	}

	fmt.Println("✅ Order Items:")
	for i, item := range items {
		fmt.Printf("\nItem #%d:\n", i+1)
		fmt.Println("  ID:       ", item["Id"])
		fmt.Println("  Product:  ", item["IdProduct"])
		fmt.Println("  Amount:   ", item["Amount"])
	}
}

func GetFreeOrders(client *http.Client) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/orders/?status=непринятый", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to get free orders. Status: %s\n", resp.Status)
		return
	}

	var orders []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("ℹ️ No free orders available.")
		return
	}

	fmt.Println("✅ Free Orders:")
	for i, order := range orders {
		fmt.Printf("\nOrder #%d\n", i+1)
		fmt.Printf("  ID:      %v\n", order["Id"])
		fmt.Printf("  Status:  %v\n", order["Status"])
		fmt.Printf("  Address: %v\n", order["Address"])
		fmt.Printf("  Price:   %v\n", order["Price"])
		fmt.Printf("  Date:    %v\n", order["Date"])
	}
}

func ChangeOrderStatus(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID (UUID): ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	_, err := uuid.Parse(orderID)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	fmt.Print("Enter new status for the order: ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)

	payload := map[string]string{"status": status}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("PATCH", "http://localhost:8080/api/v1/orders/"+orderID+"?status="+status, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("❌ Request creation failed:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Status updated successfully.")
	} else {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Println("❌ Error:", errResp["error"])
	}
}

func DeleteOrder(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Order ID (UUID) to delete: ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	_, err := uuid.Parse(orderID)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/orders/"+orderID, nil)
	if err != nil {
		fmt.Println("❌ Request creation failed:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Order deleted successfully.")
	} else {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Println("❌ Error:", errResp["error"])
	}
}

func GetOrdersByUser(client *http.Client) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/orders/", nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to get orders. Status: %s\n", resp.Status)
		return
	}

	var orders []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("ℹ️ No orders available.")
		return
	}

	fmt.Println("✅ User Orders:")
	for i, order := range orders {
		fmt.Printf("\nOrder #%d\n", i+1)
		fmt.Printf("  ID:      %v\n", order["Id"])
		fmt.Printf("  Status:  %v\n", order["Status"])
		fmt.Printf("  Address: %v\n", order["Address"])
		fmt.Printf("  Price:   %v\n", order["Price"])
		fmt.Printf("  Date:    %v\n", order["Date"])
	}
}
