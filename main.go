package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

// Métrica agregada por status e método de pagamento
type Metric struct {
	Date          string  `json:"date"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	TotalOrders   int     `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
}

var db *sql.DB
var apiToken string

func init() {
	// Carrega variáveis de ambiente do .env
	_ = godotenv.Load()

	apiToken = os.Getenv("API_TOKEN")
	if apiToken == "" {
		apiToken = "secrettoken" // fallback
	}

	var err error
	connStr := "host=db port=5432 user=postgres password=postgres dbname=analytics sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar no DB:", err)
	}
}

// Middleware simples de autenticação via token
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != fmt.Sprintf("Bearer %s", apiToken) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Endpoint: métricas agregadas por período e filtro opcional
func getMetrics(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	query := `
		SELECT date, status, payment_method, total_orders, total_revenue
		FROM aggregated.daily_metrics
		WHERE ($1::date IS NULL OR date >= $1::date)
		  AND ($2::date IS NULL OR date <= $2::date)
		ORDER BY date ASC;
	`

	rows, err := db.Query(query, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.Date, &m.Status, &m.PaymentMethod, &m.TotalOrders, &m.TotalRevenue); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		metrics = append(metrics, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// Endpoint: séries temporais (exemplo simplificado)
func getTimeseries(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	query := `
		SELECT date, SUM(total_orders) as orders, SUM(total_revenue) as revenue
		FROM aggregated.daily_metrics
		WHERE ($1::date IS NULL OR date >= $1::date)
		  AND ($2::date IS NULL OR date <= $2::date)
		GROUP BY date
		ORDER BY date ASC;
	`

	type TimeSeries struct {
		Date    string  `json:"date"`
		Orders  int     `json:"orders"`
		Revenue float64 `json:"revenue"`
	}

	rows, err := db.Query(query, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var series []TimeSeries
	for rows.Next() {
		var s TimeSeries
		if err := rows.Scan(&s.Date, &s.Orders, &s.Revenue); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		series = append(series, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware)
	api.HandleFunc("/metrics", getMetrics).Methods("GET")
	api.HandleFunc("/timeseries", getTimeseries).Methods("GET")

	serverPort := "5002"
	if p := os.Getenv("PORT"); p != "" {
		serverPort = p
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + serverPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Dashboard API rodando na porta", serverPort)
	log.Fatal(srv.ListenAndServe())
}
