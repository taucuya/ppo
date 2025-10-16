package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var baseURL string

func TestMain(m *testing.M) {
	baseURL = os.Getenv("APP_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081"
	}

	if err := waitForServer(baseURL); err != nil {
		fmt.Printf("Server is not available: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func waitForServer(url string) error {
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < 30; i++ {
		resp, err := client.Get(url + "/api/v1/auth/login")
		if err == nil && resp.StatusCode < 500 {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("server did not become ready in time")
}

func TestE2E_UserOrderFlow(t *testing.T) {
	client := &http.Client{}
	timestamp := time.Now().Unix()
	testEmail := fmt.Sprintf("testuser%d@example.com", timestamp)

	// 1) Регистрация пользователя
	registerReq := map[string]interface{}{
		"name":          "Test User",
		"date_of_birth": "1990-01-01",
		"mail":          testEmail,
		"password":      "password123",
		"phone":         "89016475899",
		"address":       "123 Order St",
	}

	registerBody, _ := json.Marshal(registerReq)
	resp, err := client.Post(baseURL+"/api/v1/auth/signup", "application/json", bytes.NewBuffer(registerBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Signup failed: %s", getResponseBody(resp))

	// 2) Вход пользователя
	loginReq := map[string]string{
		"mail":     testEmail,
		"password": "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	resp, err = client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(loginBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Login failed: %s", getResponseBody(resp))

	var loginResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NotEmpty(t, loginResp.AccessToken)

	// 3) Получить каталог товаров
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/products", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var products []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&products)
	require.NotEmpty(t, products)

	productID := products[0]["id"].(string)

	// 4) Добавить товар в корзину
	addToBasketReq := map[string]interface{}{
		"product_id": productID,
		"amount":     2,
	}
	addToBasketBody, _ := json.Marshal(addToBasketReq)

	req, _ = http.NewRequest("POST", baseURL+"/api/v1/users/me/basket/items", bytes.NewBuffer(addToBasketBody))
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// 5) Посмотреть корзину
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/users/me/basket/items", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// 6) Создать заказ
	createOrderReq := map[string]string{
		"address": "123 Delivery Address",
	}
	createOrderBody, _ := json.Marshal(createOrderReq)

	req, _ = http.NewRequest("POST", baseURL+"/api/v1/users/me/orders", bytes.NewBuffer(createOrderBody))
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var orderResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&orderResp)
	orderID := orderResp["id"].(string)
	require.NotEmpty(t, orderID)

	// 7) Проверить заказ
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/users/me/orders/"+orderID, nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	resp, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var order map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&order)
	require.Equal(t, "непринятый", order["status"])

	t.Logf("✅ E2E тест завершен успешно! Пользователь %s создал заказ %s", testEmail, orderID)
}

func getResponseBody(resp *http.Response) string {
	if resp.Body == nil {
		return ""
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String()
}
