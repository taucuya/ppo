package client_worker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("SUCCESS: Worker created successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func DeleteWorker(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Worker ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/v1/workers/"+id, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("SUCCESS: Worker deleted successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetWorkerById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Worker ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	url := fmt.Sprintf("http://localhost:8080/api/v1/workers/%s", id)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var worker map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&worker)
		fmt.Println("SUCCESS: Worker found")
		fmt.Printf("ID: %v\n", worker["Id"])
		fmt.Printf("User ID: %v\n", worker["IdUser"])
		fmt.Printf("Job Title: %v\n", worker["JobTitle"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetAllWorkers(client *http.Client) {
	resp, err := client.Get("http://localhost:8080/api/v1/workers")
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var workers []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&workers)
		fmt.Printf("SUCCESS: Found %d workers\n", len(workers))
		for i, worker := range workers {
			fmt.Printf("%d: ID: %v, User ID: %v, Job Title: %v\n",
				i+1, worker["Id"], worker["IdUser"], worker["JobTitle"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func AcceptOrder(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Order ID (UUID): ")
	orderID, _ := reader.ReadString('\n')
	orderID = strings.TrimSpace(orderID)

	url := fmt.Sprintf("http://localhost:8080/api/v1/workers/me/orders?order_id=%s", orderID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("ERROR: Failed to create request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("SUCCESS: Order accepted successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetWorkerOrders(client *http.Client) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/workers/me/orders", nil)
	if err != nil {
		fmt.Println("ERROR: Failed to create request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var orders []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&orders)
		fmt.Printf("SUCCESS: Found %d orders assigned to worker\n", len(orders))
		for i, order := range orders {
			fmt.Printf("%d: Order ID: %v, Status: %v, Address: %v\n",
				i+1, order["Id"], order["Status"], order["Address"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}
