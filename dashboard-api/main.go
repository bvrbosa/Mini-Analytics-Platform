package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Estrutura para retorno de métricas
type Metric struct {
	Date          string  `json:"date"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	TotalOrders   int     `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
}

func connectDB() (*sql.DB, error) {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=analytics sslmode=disable"
	return sql.Open("postgres", connStr)
}

// Endpoint de métricas
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Erro ao conectar no banco", 500)
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT date, status, payment_method, total_orders, total_revenue FROM aggregated.daily_metrics`)
	if err != nil {
		http.Error(w, "Erro ao consultar banco", 500)
		return
	}
	defer rows.Close()

	var metrics []Metric
	for rows.Next() {
		var m Metric
		err := rows.Scan(&m.Date, &m.Status, &m.PaymentMethod, &m.TotalOrders, &m.TotalRevenue)
		if err != nil {
			http.Error(w, "Erro ao ler dados", 500)
			return
		}
		metrics = append(metrics, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func main() {
	http.HandleFunc("/metrics", metricsHandler)
	log.Println("Dashboard API rodando na porta 5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
