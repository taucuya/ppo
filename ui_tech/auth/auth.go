package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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
		fmt.Println("Failed to encode JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("✅", result["message"])
	} else {
		fmt.Println("❌ Signup failed:", result["error"])
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
		fmt.Println("❌ Request failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("✅", "Logged in!")
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
		fmt.Println("❌ No refresh_token found")
		return
	}

	data := map[string]string{"refresh_token": rtoken}
	body, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/auth/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Logout failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("✅", "Logged out")
}
