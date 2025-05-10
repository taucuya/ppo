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
		fmt.Println("❌ Invalid rating. Please enter a number from 1 to 5.")
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
		fmt.Println("❌ Failed to encode JSON:", err)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/api/v1/review/product/%s", idProduct)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("❌ Failed to create request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YOUR_TOKEN")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅", result["message"])
	} else {
		fmt.Println("❌ Failed to create review:", result["error"])
	}
}

func GetReviewById(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Review ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	url := fmt.Sprintf("http://localhost:8080/api/v1/review/%s", id)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer YOUR_TOKEN")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var review map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&review); err != nil {
		fmt.Println("❌ Failed to decode review:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Review:", review)
	} else {
		fmt.Println("❌ Error:", review["error"])
	}
}

func GetReviewsByProductId(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Product ID: ")
	pid, _ := reader.ReadString('\n')
	pid = strings.TrimSpace(pid)

	url := fmt.Sprintf("http://localhost:8080/api/v1/product/reviews/%s", pid)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer YOUR_TOKEN")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var reviews []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&reviews); err != nil {
		fmt.Println("❌ Failed to decode reviews:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅ Reviews:")
		for _, r := range reviews {
			fmt.Println(r)
		}
	} else {
		fmt.Println("❌ Failed to get reviews")
	}
}

func DeleteReview(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Review ID to delete: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	url := fmt.Sprintf("http://localhost:8080/api/v1/review/%s", id)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer YOUR_TOKEN")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("✅", result["message"])
	} else {
		fmt.Println("❌ Failed to delete review:", result["error"])
	}
}
