package client_worker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CreateWorker(client *http.Client, reader *bufio.Reader) {
	fmt.Print("User ID (UUID): ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	fmt.Print("Job Title: ")
	jobTitle, _ := reader.ReadString('\n')
	jobTitle = strings.TrimSpace(jobTitle)

	payload := map[string]string{
		"id_user":   userID,
		"job_title": jobTitle,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/workers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Worker created successfully!")
	} else {
		fmt.Println("❌ Failed to create worker:", result["error"])
	}
}

func DeleteWorker(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Worker ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/v1/workers/"+id, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Worker with ID %s successfully deleted!\n", id)
	} else {
		fmt.Printf("❌ Failed to delete worker. Error: %s\n", result["error"])
	}
}

func GetWorkerById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Worker ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	resp, err := client.Get("http://localhost:8080/api/v1/workers/" + id)
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to retrieve worker. Status: %s\n", resp.Status)
		return
	}

	var worker map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&worker); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	fmt.Println("✅ Worker Details:")
	fmt.Printf("ID: %v\n", worker["Id"])
	fmt.Printf("User ID: %v\n", worker["IdUser"])
	fmt.Printf("Job Title: %s\n", worker["JobTitle"])
}

func GetAllWorkers(client *http.Client) {
	resp, err := client.Get("http://localhost:8080/api/v1/workers/all")
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Failed to get workers. Status: %s\n", resp.Status)
		return
	}

	var workers []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&workers); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	if len(workers) == 0 {
		fmt.Println("❌ No workers found.")
		return
	}

	fmt.Println("✅ Workers List:")
	for _, worker := range workers {
		fmt.Printf("ID: %v\n", worker["Id"])
		fmt.Printf("User ID: %v\n", worker["IdUser"])
		fmt.Printf("Job Title: %s\n", worker["JobTitle"])
		fmt.Println("------------------------------")
	}
}

func AcceptOrder(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Order ID (UUID): ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/workers/accept?order_id="+orderID, nil)
	if err != nil {
		fmt.Println("❌ Failed to create request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("❌ Failed to decode JSON:", err)
		return
	}

	if msg, exists := response["message"]; exists {
		fmt.Printf("✅ %v\n", msg)
	} else {
		fmt.Println("❌ No success message found in response.")
	}
}

func GetWorkerOrders(client *http.Client) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/workers/orders", nil)
	if err != nil {
		fmt.Println("❌ Failed to create request:", err)
		return
	}

	resp, err := client.Do(req)
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

	var orders []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("No orders found for this worker.")
		return
	}

	fmt.Println("📦 Worker Orders:")
	for i, order := range orders {
		fmt.Printf("\n Order #%d\n", i+1)
		for key, value := range order {
			fmt.Printf("   %s: %v\n", key, value)
		}
	}
}
