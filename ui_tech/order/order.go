package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func CreateOrder(client *http.Client, address string) {
	payload := map[string]string{
		"address": address,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/order", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("❌ Request error:", err)
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
		fmt.Println("✅ Order created")
	} else {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Println("❌ Error:", errResp["error"])
	}
}

func GetOrderById(client *http.Client, id string) {
	_, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	resp, err := client.Get("http://localhost:8080/api/v1/order/" + id)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var order map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&order)

	fmt.Println("✅ Order:", order)
}

func GetOrderItems(client *http.Client, id string) {
	_, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	resp, err := client.Get("http://localhost:8080/api/v1/order/" + id + "/items")
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var items []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&items)

	fmt.Println("✅ Order items:")
	for _, item := range items {
		fmt.Println(item)
	}
}

func GetFreeOrders(client *http.Client) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/order/freeorders", nil)
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
		fmt.Printf("Failed to get free orders. Status: %s\n", resp.Status)
		return
	}

	var orders []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		fmt.Println("Failed to decode response:", err)
		return
	}

	fmt.Println("Free Orders:")
	for _, order := range orders {
		fmt.Printf("- ID: %v, Status: %v, Address: %v\n", order["Id"], order["Status"], order["Address"])
	}
}

func ChangeOrderStatus(client *http.Client, id string, status string) {
	_, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	payload := map[string]string{"status": status}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("PUT", "http://localhost:8080/api/v1/order/"+id+"/status", bytes.NewBuffer(body))
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
		fmt.Println("✅ Status updated")
	} else {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Println("❌ Error:", errResp["error"])
	}
}

func DeleteOrder(client *http.Client, id string) {
	_, err := uuid.Parse(id)
	if err != nil {
		fmt.Println("❌ Invalid UUID:", err)
		return
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/order/"+id, nil)
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
		fmt.Println("✅ Order deleted")
	} else {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Println("❌ Error:", errResp["error"])
	}
}
