package basket

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func GetBasketItems(client *http.Client) {
	url := "http://localhost:8080/api/v1/users/me/basket/items"
	req, err := http.NewRequest("GET", url, nil)
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
		var items []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&items)
		fmt.Printf("SUCCESS: Found %d items in basket\n", len(items))
		for i, item := range items {
			fmt.Printf("%d: ID: %v, Product ID: %v, Amount: %v\n", i+1, item["Id"], item["IdProduct"], item["Amount"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetBasket(client *http.Client) {
	url := "http://localhost:8080/api/v1/users/me/basket"

	req, err := http.NewRequest("GET", url, nil)
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
		var basket map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&basket)
		fmt.Println("SUCCESS: Basket retrieved")
		fmt.Printf("Basket ID: %v, User ID: %v, Date: %v\n", basket["Id"], basket["IdUser"], basket["Date"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func AddToBasket(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID): ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	_, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("ERROR: Invalid product ID:", err)
		return
	}

	fmt.Print("Enter Amount: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Println("ERROR: Invalid amount:", err)
		return
	}

	item := map[string]interface{}{
		"product_id": productID,
		"amount":     amount,
	}

	body, err := json.Marshal(item)
	if err != nil {
		fmt.Println("ERROR: Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/users/me/basket/items", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("ERROR: Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("SUCCESS: Item added to basket")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func DeleteFromBasket(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID) to delete: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	_, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("ERROR: Invalid product ID:", err)
		return
	}

	payload := map[string]interface{}{
		"product_id": productID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("ERROR: Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/users/me/basket/items", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("ERROR: Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("SUCCESS: Item deleted from basket")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func UpdateItemAmount(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID) to update: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	_, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("ERROR: Invalid product ID:", err)
		return
	}

	fmt.Print("Enter new amount: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Println("ERROR: Invalid amount:", err)
		return
	}

	payload := map[string]interface{}{
		"product_id": productID,
		"amount":     amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("ERROR: Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("PATCH", "http://localhost:8080/api/v1/users/me/basket/items", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("ERROR: Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("SUCCESS: Item amount updated")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}
