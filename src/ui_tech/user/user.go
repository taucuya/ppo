package client_user

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func GetUserByEmail(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	url := fmt.Sprintf("http://localhost:8080/api/v1/users?email=%s", url.QueryEscape(email))
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
		var user map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&user)
		fmt.Println("SUCCESS: User found by email")
		fmt.Printf("User ID: %v\n", user["Id"])
		fmt.Printf("Name: %v\n", user["Name"])
		fmt.Printf("Date of Birth: %v\n", user["Date_of_birth"])
		fmt.Printf("Email: %v\n", user["Mail"])
		fmt.Printf("Phone: %v\n", user["Phone"])
		fmt.Printf("Address: %v\n", user["Address"])
		fmt.Printf("Status: %v\n", user["Status"])
		fmt.Printf("Role: %v\n", user["Role"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetUserByPhone(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Phone: ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)

	url := fmt.Sprintf("http://localhost:8080/api/v1/users?phone=%s", url.QueryEscape(phone))
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
		var user map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&user)
		fmt.Println("SUCCESS: User found by phone")
		fmt.Printf("User ID: %v\n", user["Id"])
		fmt.Printf("Name: %v\n", user["Name"])
		fmt.Printf("Date of Birth: %v\n", user["Date_of_birth"])
		fmt.Printf("Email: %v\n", user["Mail"])
		fmt.Printf("Phone: %v\n", user["Phone"])
		fmt.Printf("Address: %v\n", user["Address"])
		fmt.Printf("Status: %v\n", user["Status"])
		fmt.Printf("Role: %v\n", user["Role"])
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}

func GetAllUsers(client *http.Client) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/users", nil)
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
		var users []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&users)
		fmt.Printf("SUCCESS: Found %d users\n", len(users))
		for i, user := range users {
			fmt.Printf("\n%d: User ID: %v, Name: %v, Email: %v, Role: %v\n",
				i+1, user["Id"], user["Name"], user["Mail"], user["Role"])
		}
	} else {
		var errorResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		fmt.Println("ERROR:", errorResponse["error"])
	}
}
