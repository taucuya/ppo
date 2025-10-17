package e2e_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var baseURL string

func TestMain(m *testing.M) {
	baseURL = os.Getenv("APP_URL")
	if baseURL == "" {
		baseURL = "http://api:8080"
	}
	code := m.Run()
	os.Exit(code)
}

func TestE2E_UserOrderFlow(t *testing.T) {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	timestamp := time.Now().Unix()
	testEmail := fmt.Sprintf("testuser%d@example.com", timestamp)

	// 1) Регистрация пользователя
	registerReq := map[string]interface{}{
		"name":          "Test User",
		"date_of_birth": "1990-01-01",
		"email":         testEmail,
		"password":      "password123",
		"phone":         "89016475899",
		"address":       "123 Order St",
	}

	registerBody, _ := json.Marshal(registerReq)
	resp, err := client.Post(baseURL+"/api/v1/auth/signup", "application/json", bytes.NewBuffer(registerBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Signup failed: %s", getResponseBody(resp))

	// 2) Вход пользователя
	loginReq := map[string]string{
		"email":    testEmail,
		"password": "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	resp, err = client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(loginBody))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "Login failed: %s", getResponseBody(resp))

	// fmt.Println("Cookies after login:")
	// for _, cookie := range resp.Cookies() {
	// 	fmt.Printf("  %s: %s (Domain: %s, Path: %s)\n",
	// 		cookie.Name, cookie.Value, cookie.Domain, cookie.Path)
	// }

	// 3) Получить каталог товаров
	req, _ := http.NewRequest("GET", baseURL+"/api/v1/products?category=уход", nil)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var products []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&products)
	require.NotEmpty(t, products)

	productID, _ := products[0]["Id"].(string)

	// 4) Добавить товар в корзину
	addToBasketReq := map[string]interface{}{
		"product_id": productID,
		"amount":     2,
	}
	addToBasketBody, _ := json.Marshal(addToBasketReq)

	req, _ = http.NewRequest("POST", baseURL+"/api/v1/users/me/basket/items", bytes.NewBuffer(addToBasketBody))
	req.Header.Set("Content-Type", "application/json")

	// u, _ := url.Parse(baseURL)
	// cookies := client.Jar.Cookies(u)
	// fmt.Printf("🍪 Cookies before basket request: %d\n", len(cookies))
	// for _, cookie := range cookies {
	// 	fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	// }

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// 5) Посмотреть корзину
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/users/me/basket/items", nil)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// 6) Создать заказ
	createOrderReq := map[string]string{
		"address": "123 Delivery Address",
	}
	createOrderBody, _ := json.Marshal(createOrderReq)

	req, _ = http.NewRequest("POST", baseURL+"/api/v1/users/me/orders", bytes.NewBuffer(createOrderBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// 7) Проверить заказ
	req, _ = http.NewRequest("GET", baseURL+"/api/v1/users/me/orders/", nil)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	t.Logf("E2E тест завершен успешно! Пользователь %s создал заказ", testEmail)
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
