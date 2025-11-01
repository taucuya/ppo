package client_review

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func CreateReview(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Product ID: ")
	idProduct, _ := reader.ReadString('\n')
	idProduct = strings.TrimSpace(idProduct)

	fmt.Print("Rating (1-5): ")
	ratingStr, _ := reader.ReadString('\n')
	ratingStr = strings.TrimSpace(ratingStr)

	ratingInt, err := strconv.Atoi(ratingStr)
	if err != nil {
		fmt.Println("ERROR: Invalid rating. Please enter a number from 1 to 5.")
		return
	}

	fmt.Print("Text: ")
	text, _ := reader.ReadString('\n')

	payload := map[string]interface{}{
		"rating": ratingInt,
		"r_text": strings.TrimSpace(text),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("ERROR: Failed to encode JSON:", err)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/users/me/products/%s/reviews", idProduct)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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
		fmt.Println("SUCCESS: Review created successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetReviewById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Product ID: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	fmt.Print("Review ID: ")
	reviewID, _ := reader.ReadString('\n')
	reviewID = strings.TrimSpace(reviewID)

	url := fmt.Sprintf("http://localhost:8080/api/v1/products/%s/reviews/%s", productID, reviewID)
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
		var review map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&review)
		fmt.Println("SUCCESS: Review found")
		fmt.Printf("ID: %v\n", review["Id"])
		fmt.Printf("Product ID: %v\n", review["IdProduct"])
		fmt.Printf("User ID: %v\n", review["IdUser"])
		fmt.Printf("Rating: %v/5\n", review["Rating"])
		fmt.Printf("Text: %v\n", review["Text"])
		fmt.Printf("Date: %v\n", review["Date"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetReviewsByProductId(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Product ID: ")
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
		fmt.Printf("SUCCESS: Found %d reviews for product %s\n", len(reviews), productID)
		for i, review := range reviews {
			fmt.Printf("%d: Rating: %v/5, Text: %v, Date: %v\n",
				i+1, review["Rating"], review["Text"], review["Date"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func DeleteReview(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Product ID: ")
	productID, _ := reader.ReadString('\n')
	productID = strings.TrimSpace(productID)

	fmt.Print("Review ID: ")
	reviewID, _ := reader.ReadString('\n')
	reviewID = strings.TrimSpace(reviewID)

	url := fmt.Sprintf("http://localhost:8080/api/v1/products/%s/reviews/%s", productID, reviewID)
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
		fmt.Println("SUCCESS: Review deleted successfully")
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}
