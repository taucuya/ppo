package favourites

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

func GetFavouritesItems(client *http.Client) {
	url := "http://localhost:8080/api/v1/users/me/favourite/items"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Failed to read response:", err)
		return
	}

	var result map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("ERROR: Failed to parse response:", err)
			return
		}
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var items []map[string]interface{}
		if err := json.Unmarshal(body, &items); err != nil {
			fmt.Println("ERROR: Failed to parse items:", err)
			return
		}
		fmt.Printf("SUCCESS: Found %d items in favourites\n", len(items))
		for i, item := range items {
			fmt.Printf("%d: ID: %v, Product ID: %v\n", i+1, item["Id"], item["IdProduct"])
		}
	case http.StatusBadRequest:
		fmt.Println("ERROR: Invalid ID format:", result["error"])
	case http.StatusUnauthorized:
		fmt.Println("ERROR: Unauthorized access:", result["error"])
	case http.StatusNotFound:
		fmt.Println("SUCCESS: No items in favourites")
	case http.StatusInternalServerError:
		fmt.Println("ERROR: Failed to get favourites:", result["error"])
	default:
		fmt.Printf("ERROR: Unexpected error (status %d): %v\n", resp.StatusCode, result["error"])
	}
}

func AddToFavourites(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Enter Product ID (UUID): ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	_, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("ERROR: Invalid product ID:", err)
		return
	}

	item := map[string]interface{}{
		"id_product": productID,
	}

	body, err := json.Marshal(item)
	if err != nil {
		fmt.Println("ERROR: Failed to encode JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/users/me/favourite/items", bytes.NewBuffer(body))
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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Failed to read response:", err)
		return
	}

	var result map[string]interface{}
	if len(responseBody) > 0 {
		if err := json.Unmarshal(responseBody, &result); err != nil {
			fmt.Println("ERROR: Failed to parse response:", err)
			return
		}
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		fmt.Println("SUCCESS: Item added to favourites")
	case http.StatusBadRequest:
		if result["error"] == "Item already in favourites" {
			fmt.Println("ERROR: Item already in favourites")
		} else {
			fmt.Println("ERROR: Invalid input data:", result["error"])
		}
	case http.StatusUnauthorized:
		fmt.Println("ERROR: Unauthorized access:", result["error"])
	case http.StatusInternalServerError:
		fmt.Println("ERROR: Failed to add item to favourites:", result["error"])
	default:
		fmt.Printf("ERROR: Unexpected error (status %d): %v\n", resp.StatusCode, result["error"])
	}
}

func DeleteFromFavourites(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Product ID to remove from favourites: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	_, err := uuid.Parse(productID)
	if err != nil {
		fmt.Println("ERROR: Invalid product ID:", err)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/users/me/favourite/items/%s", productID)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Failed to read response:", err)
		return
	}

	var result map[string]interface{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("ERROR: Failed to parse response:", err)
			return
		}
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println("SUCCESS: Item removed from favourites")
	case http.StatusBadRequest:
		fmt.Println("ERROR: Invalid product ID format:", result["error"])
	case http.StatusUnauthorized:
		fmt.Println("ERROR: Unauthorized access:", result["error"])
	case http.StatusNotFound:
		fmt.Println("ERROR: Item not found in favourites:", result["error"])
	case http.StatusInternalServerError:
		fmt.Println("ERROR: Failed to remove item from favourites:", result["error"])
	default:
		fmt.Printf("ERROR: Unexpected error (status %d): %v\n", resp.StatusCode, result["error"])
	}
}
