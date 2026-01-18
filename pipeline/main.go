package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Order struct {
	OrderID       string  `json:"order_id"`
	CreatedAt     string  `json:"created_at"`
	Status        string  `json:"status"`
	Value         float64 `json:"value"`
	PaymentMethod string  `json:"payment_method"`
}


func fetchOrders() ([]Order, error) {
	resp, err := http.Get("http://flask_api:5001/orders")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var orders []Order
	err = json.Unmarshal(body, &orders)
	return orders, err
}

func connectDB() (*sql.DB, error) {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=analytics sslmode=disable"
	return sql.Open("postgres", connStr)
}

func parseTimestamp(ts string) time.Time {
	if ts == "" {
		return time.Now().UTC()
	}

	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return time.Now().UTC()
	}

	return parsed
}

func insertOrders(db *sql.DB, orders []Order) error {
	query := `
	INSERT INTO raw_data.orders 
	(order_id, created_at, status, value, payment_method)
	VALUES ($1, $2, $3, $4, $5)
	`

	for _, o := range orders {
		createdAt := parseTimestamp(o.CreatedAt)

		_, err := db.Exec(
			query,
			o.OrderID,
			createdAt,
			o.Status,
			o.Value,
			o.PaymentMethod,
		)

		if err != nil {
			return err
		}
	}
	return nil
}

func aggregateDailyMetrics(db *sql.DB) error {
	query := `
	INSERT INTO aggregated.daily_metrics
	SELECT
		DATE(created_at) AS date,
		status,
		payment_method,
		COUNT(*) AS total_orders,
		SUM(value) AS total_revenue
	FROM raw_data.orders
	GROUP BY DATE(created_at), status, payment_method;
	`
	_, err := db.Exec(query)
	return err
}

func main() {
	fmt.Println("Iniciando pipeline...")

	orders, err := fetchOrders()
	if err != nil {
		log.Fatal("Erro ao buscar dados:", err)
	}

	db, err := connectDB()
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}
	defer db.Close()

	err = insertOrders(db, orders)
	if err != nil {
		log.Fatal("Erro ao inserir dados:", err)
	}

	err = aggregateDailyMetrics(db)
	if err != nil {
		log.Fatal("Erro na agregação:", err)
	}

	fmt.Println("Pipeline finalizado com sucesso.")
}