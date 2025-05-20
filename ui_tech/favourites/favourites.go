package favourites

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func GetFavouritesItems(client *http.Client) {
	url := "http://localhost:8080/api/v1/favourites/items"
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

	fmt.Println("✅ Favourites items:")
	for _, item := range items {
		fmt.Printf("ID: %v, Product ID: %v", item["Id"], item["IdProduct"])
	}
}

func GetFavourites(client *http.Client) {
	url := "http://localhost:8080/api/v1/favourites"

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

	var favourites map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&favourites); err != nil {
		fmt.Println("❌ Failed to decode response:", err)
		return
	}

	fmt.Printf("✅ Favourites ID: %v, User ID: %v", favourites["Id"], favourites["IdUser"])
}

func AddToFavourites(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID): ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	idProduct, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("❌ Invalid product ID:", err)
		return
	}

	item := map[string]interface{}{
		"id_product": idProduct.String(),
	}

	body, err := json.Marshal(item)
	if err != nil {
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/favourites", bytes.NewBuffer(body))
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
		fmt.Println("✅ Item added to favourites!")
	} else {
		fmt.Printf("❌ Error: %v\n", respData["error"])
	}
}

func DeleteFromFavourites(client *http.Client, reader *bufio.Reader) {
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

	req, err := http.NewRequest("DELETE", "http://localhost:8080/api/v1/favourites", bytes.NewBuffer(body))
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
		fmt.Println("✅ Item deleted from favourites!")
	} else {
		fmt.Println("❌ Error:", respData["error"])
	}
}
