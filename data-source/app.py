import csv
from flask import Flask, jsonify
from datetime import datetime

app = Flask(__name__)
CSV_FILE = "data.csv"

def parse_float(value_str):
    """Converte '199,90' para 199.90"""
    if value_str is None:
        return 0.0
    return float(value_str.replace(",", "."))

@app.route("/orders")
def get_orders():
    orders = []
    # Força o encoding e o delimitador ';'
    with open(CSV_FILE, newline="", encoding="utf-8-sig") as f:
        reader = csv.DictReader(f, delimiter=";")
        for row in reader:
            # Ignora linhas vazias
            if not row.get("order_id"):
                continue
            created_at = datetime.fromisoformat(row.get("created_at").replace("Z", "+00:00"))
            orders.append({
                "order_id": row.get("order_id"),
                "created_at": created_at.isoformat(),
                "status": row.get("status"),
                "value": parse_float(row.get("value")),
                "payment_method": row.get("payment_method")
            })
    return jsonify(orders)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5001, debug=False)