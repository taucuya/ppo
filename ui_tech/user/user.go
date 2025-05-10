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

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	fmt.Printf("✅ User: %+v\n", data)
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

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	fmt.Printf("✅ User: %+v\n", data)
}
