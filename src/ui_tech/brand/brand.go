package brand

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

func CreateBrand(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Brand name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Price category: ")
	priceCategory, _ := reader.ReadString('\n')
	priceCategory = strings.TrimSpace(priceCategory)

	brand := map[string]interface{}{
		"name":           name,
		"description":    description,
		"price_category": priceCategory,
	}

	body, err := json.Marshal(brand)
	if err != nil {
		fmt.Println("❌ Failed to encode brand JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/brand", bytes.NewBuffer(body))
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

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ Brand created successfully.")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
	}
}

func GetBrandById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Brand ID: ")
	brandID, _ := reader.ReadString('\n')
	brandID = strings.TrimSpace(brandID)

	id, err := uuid.Parse(brandID)
	if err != nil {
		fmt.Println("❌ Invalid brand ID:", err)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/brand/%s", id.String())
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var brand map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&brand); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	fmt.Printf("✅ Brand ID: %v, Name: %v, \n Description: %v, Price category: %v\n", brand["Id"],
		brand["Name"], brand["Description"], brand["PriceCategory"])
}

func DeleteBrand(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Brand ID: ")
	brandID, _ := reader.ReadString('\n')
	brandID = strings.TrimSpace(brandID)
	id, err := uuid.Parse(brandID)
	if err != nil {
		fmt.Println("❌ Invalid brand ID:", err)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/brand/%s", id.String())
	req, err := http.NewRequest("DELETE", url, nil)
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
		fmt.Println("✅ Brand deleted successfully.")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
	}
}

func GetBrandsByCategory(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Category: ")
	category, _ := reader.ReadString('\n')
	category = strings.TrimSpace(category)

	baseURL := "http://localhost:8080/api/v1/brands"
	params := url.Values{}
	params.Add("category", category)

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var brands []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&brands)
		fmt.Println("✅ Brands in category:", category)
		for _, brand := range brands {
			fmt.Printf("- %v ID: %v \n", brand["Name"], brand["Id"])
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
	}
}
