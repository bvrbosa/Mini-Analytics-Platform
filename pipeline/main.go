package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

func connectDB() (*sql.DB, error) {
	connStr := "host=db port=5432 user=postgres password=postgres dbname=analytics sslmode=disable"
	return sql.Open("postgres", connStr)
}

func fetchOrders() ([]Order, error) {
	resp, err := http.Get("http://flask_api:5002/orders")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var orders []Order
	err = json.NewDecoder(resp.Body).Decode(&orders)
	return orders, err
}

func insertOrders(db *sql.DB, orders []Order) error {
	query := `
		INSERT INTO raw_data.orders 
		(order_id, created_at, status, value, payment_method)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, o := range orders {
		createdAt, _ := time.Parse(time.RFC3339, o.CreatedAt)

		_, err := db.Exec(
			query,
			o.OrderID,
			createdAt,
			o.Status,
			o.Value,
			o.PaymentMethod,
		)

		if err != nil {
			log.Println("Erro ao inserir:", err)
		}
	}

	return nil
}

func main() {
	log.Println("Pipeline iniciado...")

	db, err := connectDB()
	if err != nil {
		log.Fatal("Erro ao conectar no banco:", err)
	}
	defer db.Close()

	orders, err := fetchOrders()
	if err != nil {
		log.Fatal("Erro ao buscar pedidos:", err)
	}

	err = insertOrders(db, orders)
	if err != nil {
		log.Fatal("Erro ao inserir pedidos:", err)
	}

	log.Printf("Pipeline finalizado. %d pedidos inseridos.\n", len(orders))

	if os.Getenv("RUN_ONCE") == "true" {
		return
	}

	select {}
}
