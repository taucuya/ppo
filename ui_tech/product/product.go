package product

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func CreateProduct(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')

	fmt.Print("Articule: ")
	articule, _ := reader.ReadString('\n')

	fmt.Print("Category: ")
	category, _ := reader.ReadString('\n')

	fmt.Print("Brand: ")
	brand, _ := reader.ReadString('\n')

	fmt.Print("Price: ")
	price, _ := reader.ReadString('\n')

	fmt.Print("Stock: ")
	stock, _ := reader.ReadString('\n')
	stock = strings.TrimSpace(stock)
	amountInt, err := strconv.Atoi(stock)
	if err != nil {
		fmt.Println("❌ Invalid amount:", err)
		return
	}

	fmt.Print("Picture link: ")
	pic_link, _ := reader.ReadString('\n')

	payload := map[string]interface{}{
		"name":        strings.TrimSpace(name),
		"description": strings.TrimSpace(description),
		"articule":    strings.TrimSpace(articule),
		"category":    strings.TrimSpace(category),
		"id_brand":    strings.TrimSpace(brand),
		"price":       strings.TrimSpace(price),
		"amount":      amountInt,
		"pic_link":    strings.TrimSpace(pic_link),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/product", bytes.NewBuffer(body))
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

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Product created!")
	} else {
		fmt.Println("❌ Failed to create product:", result["error"])
	}
}

func DeleteProduct(client *http.Client, id string) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/product/"+id, nil)
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

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Product deleted!")
	} else {
		fmt.Println("❌ Error: Unable to delete product")
	}
}

func GetProduct(client *http.Client, query string) {
	url := "http://localhost:8080/api/v1/product/?" + query

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

	var product map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		fmt.Println("❌ Failed to decode product:", err)
		return
	}
	fmt.Printf("✅ Product: %+v\n", product)
}

func GetProductsByCategory(client *http.Client, category string) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/product/category/"+category, nil)
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

	var products []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		fmt.Println("❌ Failed to decode products:", err)
		return
	}
	fmt.Println("✅ Products in category:")
	for _, p := range products {
		fmt.Println(p)
	}
}

func GetProductsByBrand(client *http.Client, brand string) {
	url := fmt.Sprintf("http://localhost:8080/api/v1/product/brand/%s", brand)

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

	var products []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		fmt.Println("❌ Failed to decode products:", err)
		return
	}
	fmt.Println("✅ Products by brand:")
	for _, p := range products {
		fmt.Println(p)
	}
}
