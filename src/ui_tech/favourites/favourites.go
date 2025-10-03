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
	fmt.Print("Item ID: ")
	itemID, _ := reader.ReadString('\n')
	itemID = strings.TrimSpace(itemID)

	id, err := uuid.Parse(itemID)
	if err != nil {
		fmt.Println("❌ Invalid item ID:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:8080/api/v1/favourites/items/%s", id.String())

	req, err := http.NewRequest("DELETE", url, nil)
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
