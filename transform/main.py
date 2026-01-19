import psycopg2

DB_HOST = "db"
DB_NAME = "analytics"
DB_USER = "postgres"
DB_PASS = "postgres"

def get_conn():
    return psycopg2.connect(
        host=DB_HOST,
        database=DB_NAME,
        user=DB_USER,
        password=DB_PASS
    )

def run_transform():
    conn = get_conn()
    cur = conn.cursor()

    print("Iniciando transformação...")

    # Limpa tabela agregada antes de recalcular
    cur.execute("TRUNCATE TABLE aggregated.daily_metrics;")

    # Agregação principal
    cur.execute("""
        INSERT INTO aggregated.daily_metrics
        (date, status, payment_method, total_orders, total_revenue)
        SELECT
            DATE(created_at) AS date,
            status,
            payment_method,
            COUNT(*) AS total_orders,
            SUM(value) AS total_revenue
        FROM raw_data.orders
        GROUP BY 1, 2, 3
        ORDER BY 1;
    """)

    conn.commit()
    cur.close()
    conn.close()

    print("Transformação concluída com sucesso.")

if __name__ == "__main__":
    run_transform()
