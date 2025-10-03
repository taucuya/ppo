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
	url := "http://localhost:8080/api/v1/baskets/items"
	req, err := http.NewRequest("GET", url, nil)
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

	var items []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&items)

	fmt.Println("✅ Basket items:")
	for _, item := range items {
		fmt.Printf("ID: %v, Product ID: %v, Amount: %v\n", item["Id"], item["IdProduct"], item["Amount"])
	}
}

func GetBasket(client *http.Client) {
	url := "http://localhost:8080/api/v1/baskets"

	req, err := http.NewRequest("GET", url, nil)
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

	var basket map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&basket); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	fmt.Printf("✅ Basket ID: %v, User ID: %v, Date: %v\n", basket["Id"], basket["IdUser"], basket["Date"])
}

func AddToBasket(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID): ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	idProduct, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("❌ Invalid product ID:", err)
		return
	}

	fmt.Print("Enter Amount: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Println("❌ Invalid amount:", err)
		return
	}

	item := map[string]interface{}{
		"id_product": idProduct.String(),
		"amount":     amount,
	}

	body, err := json.Marshal(item)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/baskets/items", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("❌ Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Item added to basket!")
	} else {
		fmt.Printf("❌ Error: %v\n", respData["error"])
	}
}

func DeleteFromBasket(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID) to delete: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	idProduct, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("❌ Invalid product ID:", err)
		return
	}

	payload := map[string]interface{}{
		"product_id": idProduct,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/baskets/items", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("❌ Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Item deleted from basket!")
	} else {
		fmt.Println("❌ Error:", respData["error"])
	}
}

func UpdateItemAmount(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID) to update: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	idProduct, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("❌ Invalid product ID:", err)
		return
	}

	fmt.Print("Enter new amount: ")
	amountStr, _ := reader.ReadString('\n')
	amountStr = strings.TrimSpace(amountStr)

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Println("❌ Invalid amount:", err)
		return
	}

	payload := map[string]interface{}{
		"product_id": idProduct,
		"amount":     amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("PATCH", "http://localhost:8080/api/v1/baskets/items", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("❌ Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Item amount updated!")
	} else {
		fmt.Printf("❌ Error: %v\n", respData["error"])
	}
}
