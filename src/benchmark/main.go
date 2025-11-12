package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

var (
	apiURL    = getEnv("API_URL", "http://api:8080")
	targetURL = fmt.Sprintf("%s/api/v1/auth/signup", apiURL)
	dot       int
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type SignUpRequest struct {
	Name          string `json:"name"`
	Date_of_birth string `json:"date_of_birth"`
	Mail          string `json:"email"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
}

type SignUpResult struct {
	ScenarioName       string            `json:"scenario_name"`
	StartTime          time.Time         `json:"start_time"`
	EndTime            time.Time         `json:"end_time"`
	Duration           time.Duration     `json:"duration"`
	TotalRequests      int               `json:"total_requests"`
	SuccessfulRequests int               `json:"successful_requests"`
	FailedRequests     int               `json:"failed_requests"`
	ErrorBreakdown     ErrorBreakdown    `json:"error_breakdown"`
	ResponseTimes      ResponseTimeStats `json:"response_times"`
	UserCreationRate   float64           `json:"user_creation_rate"`
	RequestsPerSecond  float64           `json:"requests_per_second"`
	DegradationPoint   int               `json:"degradation_point"`
}

type ErrorBreakdown struct {
	ValidationErrors int `json:"validation_errors"`
	DuplicateErrors  int `json:"duplicate_errors"`
	ServerErrors     int `json:"server_errors"`
	NetworkErrors    int `json:"network_errors"`
	OtherErrors      int `json:"other_errors"`
}

type ResponseTimeStats struct {
	Min    time.Duration `json:"min"`
	Max    time.Duration `json:"max"`
	Mean   time.Duration `json:"mean"`
	Median time.Duration `json:"median"`
	P75    time.Duration `json:"p75"`
	P90    time.Duration `json:"p90"`
	P95    time.Duration `json:"p95"`
	P99    time.Duration `json:"p99"`
}

type StatisticsCollector struct {
	measurements    []time.Duration
	successCount    int
	errorCount      int
	minResponseTime time.Duration
	maxResponseTime time.Duration
	mu              sync.RWMutex
}

func (s *StatisticsCollector) AddMeasurement(responseTime time.Duration, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.measurements = append(s.measurements, responseTime)

	if success {
		s.successCount++
	} else {
		s.errorCount++
	}

	if responseTime < s.minResponseTime || s.minResponseTime == 0 {
		s.minResponseTime = responseTime
	}
	if responseTime > s.maxResponseTime {
		s.maxResponseTime = responseTime
	}
}

func (s *StatisticsCollector) TotalRequests() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.measurements)
}

func (s *StatisticsCollector) GetPercentiles() map[string]time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.measurements) == 0 {
		return map[string]time.Duration{}
	}

	measurements := make([]time.Duration, len(s.measurements))
	copy(measurements, s.measurements)

	sort.Slice(measurements, func(i, j int) bool {
		return measurements[i] < measurements[j]
	})

	return map[string]time.Duration{
		"min":  measurements[0],
		"max":  measurements[len(measurements)-1],
		"p50":  calculatePercentile(measurements, 0.50),
		"p75":  calculatePercentile(measurements, 0.75),
		"p90":  calculatePercentile(measurements, 0.90),
		"p95":  calculatePercentile(measurements, 0.95),
		"p99":  calculatePercentile(measurements, 0.99),
		"p100": calculatePercentile(measurements, 1.00),
	}
}

func cleanupDatabase() error {
	dsn := "postgres://test_user:test_password@postgres:5432/test_db?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("truncate table \"user\";")
	if err != nil {
		return fmt.Errorf("failed to cleanup database: %v", err)
	}

	return nil
}

func calculatePercentile(measurements []time.Duration, percentile float64) time.Duration {
	if len(measurements) == 0 {
		return 0
	}
	index := int(float64(len(measurements)-1) * percentile)
	mean := 0.0
	for i := 0; i < index; i++ {
		mean += float64(measurements[i])
	}
	mean = mean / float64(index+1)
	return time.Duration(mean)
}

func calculateRPS(requests int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return float64(requests) / duration.Seconds()
}

func GenerateSignUpRequest(template SignUpRequest, iteration int, userIndex int) ([]byte, error) {
	uniqueID := uuid.New()
	uniqueSuffix := fmt.Sprintf("%d_%s", userIndex, uniqueID.String()[:8])

	phonePrefix := "89"
	randomNumbers := fmt.Sprintf("%09d", time.Now().UnixNano()%100000)
	if len(randomNumbers) > 9 {
		randomNumbers = randomNumbers[:9]
	}
	birthDate := time.Now().AddDate(-20, 0, 0)
	request := SignUpRequest{
		Name:          fmt.Sprintf("%s_%s", template.Name, uniqueSuffix),
		Date_of_birth: birthDate.Format("2006-01-02"),
		Mail:          fmt.Sprintf("%s_%s@test.com", template.Mail, uniqueSuffix),
		Password:      template.Password,
		Phone:         phonePrefix + randomNumbers,
		Address:       template.Address,
	}

	return json.Marshal(request)
}

const (
	criticalResponseTime = 3 * time.Second
)

func runSignUpGradualLoad() *SignUpResult {
	maxUsers := 1000
	step := 10

	collector := &StatisticsCollector{}
	var (
		successCount     int32
		validationErrors int32
		duplicateErrors  int32
		serverErrors     int32
		otherErrors      int32
	)

	startTime := time.Now()
	var degradationPoint int

	file, err := os.OpenFile("grad.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file_perc, err := os.OpenFile("percentiles.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file_perc.Close()

	for concurrentUsers := step; concurrentUsers <= maxUsers; concurrentUsers += step {
		var wg sync.WaitGroup
		var batchSuccess int32

		fmt.Printf("Testing with %d concurrent users\n", concurrentUsers)

		for i := 0; i < concurrentUsers; i++ {
			wg.Add(1)
			go func(userIndex int) {
				defer wg.Done()

				requestBody, err := GenerateSignUpRequest(
					SignUpRequest{
						Name:     "User",
						Mail:     "user",
						Password: "123456",
						Address:  "Test Address",
					},
					concurrentUsers,
					userIndex,
				)

				if err != nil {
					atomic.AddInt32(&validationErrors, 1)
					return
				}

				reqStartTime := time.Now()
				resp, err := http.Post("http://api:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(requestBody))
				responseTime := time.Since(reqStartTime)

				if err != nil {
					atomic.AddInt32(&otherErrors, 1)
					return
				}
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				var result map[string]interface{}
				json.Unmarshal(body, &result)

				success := false
				switch resp.StatusCode {
				case http.StatusCreated:
					atomic.AddInt32(&successCount, 1)
					atomic.AddInt32(&batchSuccess, 1)
					success = true
				case http.StatusBadRequest:
					atomic.AddInt32(&validationErrors, 1)
				case http.StatusConflict:
					atomic.AddInt32(&duplicateErrors, 1)
				case http.StatusInternalServerError:
					atomic.AddInt32(&serverErrors, 1)
				default:
					atomic.AddInt32(&otherErrors, 1)
				}

				collector.AddMeasurement(responseTime, success)
			}(i)
		}

		wg.Wait()

		percentiles := collector.GetPercentiles()
		if p100, exists := percentiles["p100"]; exists && p100 > criticalResponseTime {
			log.Printf("Performance degradation at %d users: P100 = %v", concurrentUsers, p100)
			degradationPoint = concurrentUsers
			dot = degradationPoint
			break
		}

		log.Println(percentiles["p100"], "for\t", concurrentUsers)
		_, err = fmt.Fprintln(file, percentiles["p95"], concurrentUsers)
		if err != nil {
			log.Fatal(err)
		}

		_, err = fmt.Fprintln(file_perc, percentiles["p50"], percentiles["p75"], percentiles["p90"], percentiles["p95"], percentiles["p99"], concurrentUsers)
		if err != nil {
			log.Fatal(err)
		}
		cleanupDatabase()
	}

	endTime := time.Now()
	totalRequests := collector.TotalRequests()
	percentiles := collector.GetPercentiles()

	return &SignUpResult{
		ScenarioName:       "gradual_load",
		StartTime:          startTime,
		EndTime:            endTime,
		Duration:           endTime.Sub(startTime),
		TotalRequests:      totalRequests,
		SuccessfulRequests: int(successCount),
		FailedRequests:     totalRequests - int(successCount),
		UserCreationRate:   float64(successCount) / float64(totalRequests) * 100,
		DegradationPoint:   degradationPoint,
		RequestsPerSecond:  calculateRPS(totalRequests, endTime.Sub(startTime)),
		ErrorBreakdown: ErrorBreakdown{
			ValidationErrors: int(validationErrors),
			DuplicateErrors:  int(duplicateErrors),
			ServerErrors:     int(serverErrors),
			OtherErrors:      int(otherErrors),
		},
		ResponseTimes: ResponseTimeStats{
			Min:    percentiles["min"],
			Max:    percentiles["max"],
			Median: percentiles["p50"],
			P75:    percentiles["p75"],
			P90:    percentiles["p90"],
			P95:    percentiles["p95"],
			P99:    percentiles["p99"],
		},
	}
}

func runSignUpConstantLoad() *SignUpResult {
	maxUsers := dot
	concurrentUsers := maxUsers * 80 / 100
	testDuration := 30 * time.Second

	collector := &StatisticsCollector{}
	var (
		successCount     int32
		validationErrors int32
		duplicateErrors  int32
		serverErrors     int32
		otherErrors      int32
	)

	fmt.Printf("Running constant load: %d users for %v\n", concurrentUsers, testDuration)

	var wg sync.WaitGroup
	startTime := time.Now()
	endTime := startTime.Add(testDuration)

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()

			for time.Now().Before(endTime) {
				requestBody, err := GenerateSignUpRequest(
					SignUpRequest{
						Name:     "User",
						Mail:     "user",
						Password: "123456",
						Address:  "Test Address",
					},
					concurrentUsers,
					userIndex,
				)
				if err != nil {
					atomic.AddInt32(&validationErrors, 1)
					continue
				}

				reqStartTime := time.Now()
				resp, err := http.Post("http://api:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(requestBody))
				responseTime := time.Since(reqStartTime)

				if err != nil {
					atomic.AddInt32(&otherErrors, 1)
					continue
				}
				defer resp.Body.Close()

				success := false
				switch resp.StatusCode {
				case http.StatusCreated:
					atomic.AddInt32(&successCount, 1)
					success = true
				case http.StatusBadRequest:
					atomic.AddInt32(&validationErrors, 1)
				case http.StatusConflict:
					atomic.AddInt32(&duplicateErrors, 1)
				case http.StatusInternalServerError:
					atomic.AddInt32(&serverErrors, 1)
				default:
					atomic.AddInt32(&otherErrors, 1)
				}

				collector.AddMeasurement(responseTime, success)
			}
		}(i)
	}
	wg.Wait()
	cleanupDatabase()
	finalEndTime := time.Now()
	totalRequests := collector.TotalRequests()
	percentiles := collector.GetPercentiles()

	return &SignUpResult{
		ScenarioName:       "constant_load",
		StartTime:          startTime,
		EndTime:            finalEndTime,
		Duration:           finalEndTime.Sub(startTime),
		TotalRequests:      totalRequests,
		SuccessfulRequests: int(successCount),
		FailedRequests:     totalRequests - int(successCount),
		UserCreationRate:   float64(successCount) / float64(totalRequests) * 100,
		RequestsPerSecond:  calculateRPS(totalRequests, finalEndTime.Sub(startTime)),
		ErrorBreakdown: ErrorBreakdown{
			ValidationErrors: int(validationErrors),
			DuplicateErrors:  int(duplicateErrors),
			ServerErrors:     int(serverErrors),
			OtherErrors:      int(otherErrors),
		},
		ResponseTimes: ResponseTimeStats{
			Min:    percentiles["min"],
			Max:    percentiles["max"],
			Mean:   percentiles["mean"],
			Median: percentiles["p50"],
			P75:    percentiles["p75"],
			P90:    percentiles["p90"],
			P95:    percentiles["p95"],
			P99:    percentiles["p99"],
		},
	}
}

func runSignUpOverloadThenReduce() *SignUpResult {
	maxUsers := dot
	overloadUsers := maxUsers * 120 / 100
	normalUsers := maxUsers / 2
	phaseDuration := 15 * time.Second

	collector := &StatisticsCollector{}
	var (
		successCount     int32
		validationErrors int32
		duplicateErrors  int32
		serverErrors     int32
		otherErrors      int32
	)

	startTime := time.Now()

	fmt.Printf("Phase 1: Overload with %d users for %v\n", overloadUsers, phaseDuration)
	runLoadPhase(overloadUsers, phaseDuration, collector, &successCount, &validationErrors, &duplicateErrors, &serverErrors, &otherErrors)

	fmt.Printf("Phase 2: Normal load with %d users for %v\n", normalUsers, phaseDuration)
	runLoadPhase(normalUsers, phaseDuration, collector, &successCount, &validationErrors, &duplicateErrors, &serverErrors, &otherErrors)

	endTime := time.Now()
	totalRequests := collector.TotalRequests()
	percentiles := collector.GetPercentiles()

	return &SignUpResult{
		ScenarioName:       "overload_then_reduce",
		StartTime:          startTime,
		EndTime:            endTime,
		Duration:           endTime.Sub(startTime),
		TotalRequests:      totalRequests,
		SuccessfulRequests: int(successCount),
		FailedRequests:     totalRequests - int(successCount),
		UserCreationRate:   float64(successCount) / float64(totalRequests) * 100,
		RequestsPerSecond:  calculateRPS(totalRequests, endTime.Sub(startTime)),
		ErrorBreakdown: ErrorBreakdown{
			ValidationErrors: int(validationErrors),
			DuplicateErrors:  int(duplicateErrors),
			ServerErrors:     int(serverErrors),
			OtherErrors:      int(otherErrors),
		},
		ResponseTimes: ResponseTimeStats{
			Min:    percentiles["min"],
			Max:    percentiles["max"],
			Mean:   percentiles["mean"],
			Median: percentiles["p50"],
			P75:    percentiles["p75"],
			P90:    percentiles["p90"],
			P95:    percentiles["p95"],
			P99:    percentiles["p99"],
		},
	}
}

func runLoadPhase(concurrentUsers int, duration time.Duration, collector *StatisticsCollector, successCount, validationErrors, duplicateErrors, serverErrors, otherErrors *int32) {
	var wg sync.WaitGroup
	endTime := time.Now().Add(duration)

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userIndex int) {
			defer wg.Done()

			for time.Now().Before(endTime) {

				requestBody, err := GenerateSignUpRequest(
					SignUpRequest{
						Name:     "User",
						Mail:     "user",
						Password: "123456",
						Address:  "Test Address",
					},
					concurrentUsers,
					userIndex,
				)
				if err != nil {
					atomic.AddInt32(validationErrors, 1)
					continue
				}
				reqStartTime := time.Now()
				resp, err := http.Post("http://api:8080/api/v1/auth/signup", "application/json", bytes.NewBuffer(requestBody))
				responseTime := time.Since(reqStartTime)

				if err != nil {
					atomic.AddInt32(otherErrors, 1)
					continue
				}
				defer resp.Body.Close()

				success := false
				switch resp.StatusCode {
				case http.StatusCreated:
					atomic.AddInt32(successCount, 1)
					success = true
				case http.StatusBadRequest:
					atomic.AddInt32(validationErrors, 1)
				case http.StatusConflict:
					atomic.AddInt32(duplicateErrors, 1)
				case http.StatusInternalServerError:
					atomic.AddInt32(serverErrors, 1)
				default:
					atomic.AddInt32(otherErrors, 1)
				}

				collector.AddMeasurement(responseTime, success)
			}
		}(i)
		cleanupDatabase()
	}

	wg.Wait()
	cleanupDatabase()
}

func main() {
	fmt.Println("Starting benchmark tests...")
	fmt.Printf("Target API: %s\n", targetURL)

	fmt.Println("\n=== СЦЕНАРИЙ 1: Постепенная нагрузка ===")
	result1 := runSignUpGradualLoad()
	printResult(result1)

	fmt.Println("\n=== СЦЕНАРИЙ 2: Постоянная нагрузка ===")
	result2 := runSignUpConstantLoad()
	printResult(result2)

	fmt.Println("\n=== СЦЕНАРИЙ 3: Перегрузка + снижение ===")
	result3 := runSignUpOverloadThenReduce()
	printResult(result3)
}

func printResult(result *SignUpResult) {
	fmt.Printf("Scenario: %s\n", result.ScenarioName)
	fmt.Printf("Duration: %v\n", result.Duration)
	fmt.Printf("Total Requests: %d\n", result.TotalRequests)
	fmt.Printf("Successful: %d (%.1f%%)\n", result.SuccessfulRequests, result.UserCreationRate)
	fmt.Printf("RPS: %.2f\n", result.RequestsPerSecond)
	fmt.Printf("Response Time P95: %v\n", result.ResponseTimes.P95)

	if result.DegradationPoint > 0 {
		fmt.Printf("Performance degradation at: %d users\n", result.DegradationPoint)
	} else {
		fmt.Printf("No performance degradation detected\n")
	}

	fmt.Printf("Errors: Validation=%d, Duplicate=%d, Server=%d, Other=%d\n",
		result.ErrorBreakdown.ValidationErrors,
		result.ErrorBreakdown.DuplicateErrors,
		result.ErrorBreakdown.ServerErrors,
		result.ErrorBreakdown.OtherErrors)
}
