package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func Signup(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')

	fmt.Print("Date of Birth (YYYY-MM-DD): ")
	dob, _ := reader.ReadString('\n')

	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	fmt.Print("Phone: ")
	phone, _ := reader.ReadString('\n')

	fmt.Print("Address: ")
	address, _ := reader.ReadString('\n')

	payload := map[string]string{
		"name":          strings.TrimSpace(name),
		"date_of_birth": strings.TrimSpace(dob),
		"email":         strings.TrimSpace(email),
		"password":      strings.TrimSpace(password),
		"phone":         strings.TrimSpace(phone),
		"address":       strings.TrimSpace(address),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("ERROR: Failed to encode JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(body))
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
		fmt.Println("SUCCESS:", result["message"])
	case http.StatusBadRequest:
		fmt.Println("ERROR: Invalid input data:", result["error"])
	case http.StatusInternalServerError:
		fmt.Println("ERROR: Registration failed:", result["error"])
	default:
		fmt.Printf("ERROR: Unexpected error (status %d): %v\n", resp.StatusCode, result["error"])
	}
}

func Login(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	payload := map[string]string{
		"email":    strings.TrimSpace(email),
		"password": strings.TrimSpace(password),
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/auth/login", bytes.NewBuffer(body))
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
	case http.StatusOK:
		cookies := resp.Cookies()
		var gotAccessToken, gotRefreshToken bool
		for _, cookie := range cookies {
			if cookie.Name == "access_token" {
				gotAccessToken = true
			}
			if cookie.Name == "refresh_token" {
				gotRefreshToken = true
			}
		}

		if gotAccessToken && gotRefreshToken {
			fmt.Println("SUCCESS: Logged in successfully!")
		} else {
			fmt.Println("WARNING: Login successful but tokens not set in cookies")
		}
	case http.StatusBadRequest:
		fmt.Println("ERROR: Invalid input data:", result["error"])
	case http.StatusUnauthorized:
		fmt.Println("ERROR: Invalid credentials:", result["error"])
	default:
		fmt.Printf("ERROR: Login failed (status %d): %v\n", resp.StatusCode, result["error"])
	}
}

func Logout(client *http.Client) {
	u, _ := url.Parse("http://localhost")
	cookies := client.Jar.Cookies(u)
	var rtoken string
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			rtoken = c.Value
			break
		}
	}

	if rtoken == "" {
		fmt.Println("ERROR: No refresh token found")
		return
	}

	data := map[string]string{"refresh_token": rtoken}
	body, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/auth/logout", bytes.NewBuffer(body))
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
	case http.StatusOK:
		fmt.Println("SUCCESS: Logged out successfully")
	case http.StatusBadRequest:
		fmt.Println("ERROR: Invalid input data:", result["error"])
	case http.StatusInternalServerError:
		fmt.Println("ERROR: Logout failed:", result["error"])
	default:
		fmt.Printf("ERROR: Unexpected error (status %d): %v\n", resp.StatusCode, result["error"])
	}
}
