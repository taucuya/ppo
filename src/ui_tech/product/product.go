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

	fmt.Print("Brand ID: ")
	brand, _ := reader.ReadString('\n')

	fmt.Print("Price: ")
	price, _ := reader.ReadString('\n')

	fmt.Print("Stock: ")
	stock, _ := reader.ReadString('\n')
	stock = strings.TrimSpace(stock)
	amountInt, err := strconv.Atoi(stock)
	if err != nil {
		fmt.Println("ERROR: Invalid amount:", err)
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
		fmt.Println("ERROR: Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/products", bytes.NewBuffer(body))
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
		fmt.Println("SUCCESS: Product created successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func DeleteProduct(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/products/"+id, nil)
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
		fmt.Println("SUCCESS: Product deleted successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetProductById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	url := "http://localhost:8080/api/v1/products?id=" + id
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
		var product map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&product)
		fmt.Println("SUCCESS: Product found")
		fmt.Printf("ID:          %v\n", product["Id"])
		fmt.Printf("Name:        %v\n", product["Name"])
		fmt.Printf("Description: %v\n", product["Description"])
		fmt.Printf("Price:       %.2f\n", product["Price"])
		fmt.Printf("Category:    %v\n", product["Category"])
		fmt.Printf("Amount:      %v\n", product["Amount"])
		fmt.Printf("Brand ID:    %v\n", product["IdBrand"])
		fmt.Printf("Articule:    %v\n", product["Articule"])
		fmt.Printf("Picture URL: %v\n", product["PicLink"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetProductByArticule(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product Articule: ")
	articule, _ := reader.ReadString('\n')
	articule = strings.TrimSpace(articule)

	url := "http://localhost:8080/api/v1/products?art=" + articule
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
		var product map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&product)
		fmt.Println("SUCCESS: Product found")
		fmt.Printf("ID:          %v\n", product["Id"])
		fmt.Printf("Name:        %v\n", product["Name"])
		fmt.Printf("Description: %v\n", product["Description"])
		fmt.Printf("Price:       %.2f\n", product["Price"])
		fmt.Printf("Category:    %v\n", product["Category"])
		fmt.Printf("Amount:      %v\n", product["Amount"])
		fmt.Printf("Brand ID:    %v\n", product["IdBrand"])
		fmt.Printf("Articule:    %v\n", product["Articule"])
		fmt.Printf("Picture URL: %v\n", product["PicLink"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetProductsByCategory(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Category: ")
	category, _ := reader.ReadString('\n')
	category = strings.TrimSpace(category)

	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/products?category="+category, nil)
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
		var products []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&products)
		fmt.Printf("SUCCESS: Found %d products in category '%s'\n", len(products), category)
		for i, p := range products {
			fmt.Printf("%d: %v - %.2f руб. (ID: %v)\n", i+1, p["Name"], p["Price"], p["Id"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetProductsByBrand(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Brand: ")
	brand, _ := reader.ReadString('\n')
	brand = strings.TrimSpace(brand)

	url := fmt.Sprintf("http://localhost:8080/api/v1/products?brand=%s", brand)

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
		var products []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&products)
		fmt.Printf("SUCCESS: Found %d products by brand '%s'\n", len(products), brand)
		for i, p := range products {
			fmt.Printf("%d: %v - %.2f руб. (ID: %v)\n", i+1, p["Name"], p["Price"], p["Id"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetReviewsForProduct(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	url := fmt.Sprintf("http://localhost:8080/api/v1/products/%s/reviews", productID)
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
		var reviews []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&reviews)
		fmt.Printf("SUCCESS: Found %d reviews for product\n", len(reviews))
		for i, review := range reviews {
			fmt.Printf("%d: Rating: %v, Text: %v, Date: %v\n",
				i+1, review["Rating"], review["Text"], review["Date"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}
