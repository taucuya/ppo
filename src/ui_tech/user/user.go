package client_user

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetUserByEmail(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	resp, err := client.Get("http://localhost:8080/api/v1/user/email?email=" + email)
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Error:", resp.Status)
		return
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Println("❌ Failed to decode user data:", err)
		return
	}

	fmt.Println("✅ User Information:")
	fmt.Println("---------------------------")
	fmt.Printf("User ID: %v\n", user["Id"])
	fmt.Printf("Name: %v\n", user["Name"])
	fmt.Printf("Date of Birth: %v\n", user["Date_of_birth"])
	fmt.Printf("Email: %v\n", user["Mail"])
	fmt.Printf("Phone: %v\n", user["Phone"])
	fmt.Printf("Address: %v\n", user["Address"])
	fmt.Printf("Status: %v\n", user["Status"])
	fmt.Printf("Role: %v\n", user["Role"])
	fmt.Println("---------------------------")
}

func GetUserByPhone(client *http.Client, reader *bufio.Reader) {
	fmt.Print("Phone: ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)

	if strings.HasPrefix(phone, "+") {
		phone = strings.Replace(phone, "+", "%2B", 1)
	}

	resp, err := client.Get("http://localhost:8080/api/v1/user/phone?phone=" + phone)
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Error:", resp.Status)
		return
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Println("❌ Failed to decode user data:", err)
		return
	}

	fmt.Println("✅ User Information:")
	fmt.Println("---------------------------")
	fmt.Printf("User ID: %v\n", user["Id"])
	fmt.Printf("Name: %v\n", user["Name"])
	fmt.Printf("Date of Birth: %v\n", user["Date_of_birth"])
	fmt.Printf("Email: %v\n", user["Mail"])
	fmt.Printf("Phone: %v\n", user["Phone"])
	fmt.Printf("Address: %v\n", user["Address"])
	fmt.Printf("Status: %v\n", user["Status"])
	fmt.Printf("Role: %v\n", user["Role"])
	fmt.Println("---------------------------")
}

func GetAllUsers(client *http.Client, reader *bufio.Reader) {

	resp, err := client.Get("http://localhost:8080/api/v1/user/users")
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("❌ Error:", resp.Status)
		return
	}

	var user []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Println("❌ Failed to decode user data:", err)
		return
	}

	for _, v := range user {
		fmt.Println("✅ User Information:")
		fmt.Println("---------------------------")
		fmt.Printf("User ID: %v\n", v["Id"])
		fmt.Printf("Name: %v\n", v["Name"])
		fmt.Printf("Date of Birth: %v\n", v["Date_of_birth"])
		fmt.Printf("Email: %v\n", v["Mail"])
		fmt.Printf("Phone: %v\n", v["Phone"])
		fmt.Printf("Address: %v\n", v["Address"])
		fmt.Printf("Status: %v\n", v["Status"])
		fmt.Println("---------------------------")
	}

}
