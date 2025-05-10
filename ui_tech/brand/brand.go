package brand

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

func GetBrandById(client *http.Client, brandID string) {
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

	if resp.StatusCode == http.StatusOK {
		var brand map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&brand)
		fmt.Println("✅ Brand details:")
		fmt.Printf("%+v\n", brand)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
	}
}

func DeleteBrand(client *http.Client, brandID string) {
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

func GetBrandsByCategory(client *http.Client, category string) {
	url := fmt.Sprintf("http://localhost:8080/api/v1/brand/category/%s", category)
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
			fmt.Printf("- %v\t%v\n", brand["Name"], brand["Id"])
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Error: %s\n", string(body))
	}
}
