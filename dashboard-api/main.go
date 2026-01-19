package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Metric representa a estrutura de uma métrica do dashboard
type Metric struct {
	Date          string  `json:"date"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	TotalOrders   int     `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
}

// connectDB cria a conexão com o PostgreSQL
func connectDB() (*sql.DB, error) {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=analytics sslmode=disable"
	return sql.Open("postgres", connStr)
}

// middleware para permitir CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// metricsHandler retorna métricas agregadas filtradas
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Erro ao conectar no banco", 500)
		log.Println("DB connection error:", err)
		return
	}
	defer db.Close()

	// Captura filtros da query string
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	paymentMethod := r.URL.Query().Get("payment_method")

	query := `SELECT date, status, payment_method, total_orders, total_revenue
	          FROM aggregated.daily_metrics
	          WHERE 1=1`
	args := []interface{}{}
	argId := 1

	if start != "" {
		query += fmt.Sprintf(" AND date >= $%d", argId)
		args = append(args, start)
		argId++
	}
	if end != "" {
		query += fmt.Sprintf(" AND date <= $%d", argId)
		args = append(args, end)
		argId++
	}
	if paymentMethod != "" {
		query += fmt.Sprintf(" AND payment_method = $%d", argId)
		args = append(args, paymentMethod)
		argId++
	}

	query += " ORDER BY date"

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Erro ao consultar banco", 500)
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		err := rows.Scan(&m.Date, &m.Status, &m.PaymentMethod, &m.TotalOrders, &m.TotalRevenue)
		if err != nil {
			http.Error(w, "Erro ao ler dados", 500)
			log.Println("Scan error:", err)
			return
		}
		metrics = append(metrics, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metricsHandler)

	// Aplicando CORS
	handler := enableCORS(mux)

	log.Println("Dashboard API rodando na porta 5002")
	log.Fatal(http.ListenAndServe(":5002", handler))
}
