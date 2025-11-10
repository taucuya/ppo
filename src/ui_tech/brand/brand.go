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
		fmt.Println("ERROR: Failed to encode brand JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/brands", bytes.NewBuffer(body))
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
		fmt.Println("SUCCESS: Brand created successfully")
	} else {
		body, _ := io.ReadAll(resp.Body)
		var errorResponse map[string]interface{}
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			fmt.Printf("ERROR: Failed to parse error response: %s\n", string(body))
		} else {
			fmt.Println("ERROR:", errorResponse["error"])
		}
	}
}

func GetBrandById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Brand ID: ")
	brandID, _ := reader.ReadString('\n')
	brandID = strings.TrimSpace(brandID)

	id, err := uuid.Parse(brandID)
	if err != nil {
		fmt.Println("ERROR: Invalid brand ID:", err)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/brands/%s", id.String())
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var brand map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&brand)
		fmt.Println("SUCCESS: Brand found")
		fmt.Printf("Brand ID: %v, Name: %v, Description: %v, Price category: %v\n",
			brand["Id"], brand["Name"], brand["Description"], brand["PriceCategory"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func DeleteBrand(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Brand ID: ")
	brandID, _ := reader.ReadString('\n')
	brandID = strings.TrimSpace(brandID)
	id, err := uuid.Parse(brandID)
	if err != nil {
		fmt.Println("ERROR: Invalid brand ID:", err)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/brands/%s", id.String())
	req, err := http.NewRequest("DELETE", url, nil)
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
		fmt.Println("SUCCESS: Brand deleted successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
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
		fmt.Println("ERROR: Request failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var brands []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&brands)
		fmt.Printf("SUCCESS: Found %d brands in category '%s'\n", len(brands), category)
		for i, brand := range brands {
			fmt.Printf("%d: %v (ID: %v)\n", i+1, brand["Name"], brand["Id"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}
