package product

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/products", bytes.NewBuffer(body))
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

func DeleteProduct(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/products/"+id, nil)
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

func GetProduct(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter search query (e.g. art=123456 or id=12345): ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	url := "http://localhost:8080/api/v1/products/?" + query
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Failed to get product. Status: %s\nDetails: %s\n", resp.Status, string(body))
		return
	}

	var product map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		fmt.Println("❌ Failed to decode product:", err)
		return
	}

	fmt.Println("✅ Product Information:")
	fmt.Printf("  ID:          %v\n", product["Id"])
	fmt.Printf("  Name:        %v\n", product["Name"])
	fmt.Printf("  Description: %v\n", product["Description"])
	fmt.Printf("  Price:       %.2f\n", product["Price"])
	fmt.Printf("  Category:    %v\n", product["Category"])
	fmt.Printf("  Amount:      %v\n", product["Amount"])
	fmt.Printf("  Brand ID:    %v\n", product["IdBrand"])
	fmt.Printf("  Articule:    %v\n", product["Articule"])
	fmt.Printf("  Picture URL: %v\n", product["PicLink"])
}

func GetProductsByCategory(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Category: ")
	category, _ := reader.ReadString('\n')
	category = strings.TrimSpace(category)

	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/products/?category="+category, nil)
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
		fmt.Printf("❌ Failed to get products. Status: %s\nDetails: %s\n", resp.Status, string(body))
		return
	}

	var products []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		fmt.Println("❌ Failed to decode products:", err)
		return
	}

	if len(products) == 0 {
		fmt.Println("No products found in this category.")
		return
	}
	fmt.Printf("✅ Products in category '%s':\n", category)
	for i, p := range products {
		fmt.Printf("%d. %v — %v руб. (ID: %v)\n", i+1, p["Name"], p["Price"], p["Id"])
	}
}

func GetProductsByBrand(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Brand: ")
	brand, _ := reader.ReadString('\n')
	brand = strings.TrimSpace(brand)

	url := fmt.Sprintf("http://localhost:8080/api/v1/products/?brand=%s", brand)

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
	for i, p := range products {
		fmt.Printf("%d. %v — %v руб. (ID: %v)\n", i+1, p["Name"], p["Price"], p["Id"])
	}
}
