import React, { useState, useEffect, useCallback } from "react";
import { getMetrics, api } from "../services/api"; // <-- api importado corretamente
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  CartesianGrid,
  ResponsiveContainer,
} from "recharts";

function Dashboard() {
  const [metrics, setMetrics] = useState([]);
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [paymentMethod, setPaymentMethod] = useState("");
  const [loading, setLoading] = useState(false);

  const fetchMetrics = useCallback(async () => {
    setLoading(true);
    try {
      const data = await getMetrics({
        start: startDate || undefined,
        end: endDate || undefined,
        payment_method: paymentMethod || undefined,
      });
      setMetrics(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error("Erro ao buscar métricas:", error);
      setMetrics([]);
    } finally {
      setLoading(false);
    }
  }, [startDate, endDate, paymentMethod]);

  useEffect(() => {
    fetchMetrics();
  }, [fetchMetrics]);

  const summarize = (status) => {
    const filtered = metrics.filter((m) => m.status === status);
    return {
      revenue: filtered.reduce((sum, m) => sum + Number(m.total_revenue || 0), 0),
      orders: filtered.reduce((sum, m) => sum + Number(m.total_orders || 0), 0),
    };
  };

  const approved = summarize("approved");
  const pending = summarize("pending");
  const cancelled = summarize("cancelled");

  const chartData = metrics.map((m) => ({
    date: m.date,
    value: Number(m.total_revenue || 0),
    status: m.status,
  }));

  return (
    <div style={styles.container}>
      <h1 style={styles.title}>Dashboard</h1>

      <div style={styles.filters}>
        <input type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} />
        <input type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} />
        <select value={paymentMethod} onChange={(e) => setPaymentMethod(e.target.value)}>
          <option value="">Todos os métodos</option>
          <option value="pix">PIX</option>
          <option value="boleto">Boleto</option>
          <option value="credit_card">Cartão</option>
        </select>

        <button onClick={fetchMetrics} disabled={loading}>
          {loading ? "Carregando..." : "Atualizar"}
        </button>

        <button
          onClick={async () => {
            setLoading(true);
            try {
              await api.post("/sync");
              await fetchMetrics();
            } catch (err) {
              console.error("Erro ao sincronizar dados:", err);
            } finally {
              setLoading(false);
            }
          }}
          disabled={loading}
          style={{ marginLeft: 10, backgroundColor: "#38bdf8", color: "#fff" }}
        >
          {loading ? "Sincronizando..." : "Sincronizar dados"}
        </button>
      </div>

      <div style={styles.metricsContainer}>
        <div style={styles.card}>Receita Aprovada: {approved.revenue.toFixed(2)}</div>
        <div style={styles.card}>Receita Pendente: {pending.revenue.toFixed(2)}</div>
        <div style={styles.card}>Receita Cancelada: {cancelled.revenue.toFixed(2)}</div>
      </div>

      <div style={styles.metricsContainer}>
        <div style={styles.card}>Pedidos Aprovados: {approved.orders}</div>
        <div style={styles.card}>Pedidos Pendentes: {pending.orders}</div>
        <div style={styles.card}>Pedidos Cancelados: {cancelled.orders}</div>
      </div>

      <div style={{ width: "100%", height: 400 }}>
        <ResponsiveContainer>
          <LineChart data={chartData}>
            <CartesianGrid stroke="#eee" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="value" name="Receita" dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

const styles = {
  container: { padding: 20 },
  title: { marginBottom: 20 },
  filters: { display: "flex", gap: 10, marginBottom: 20 },
  metricsContainer: { display: "flex", gap: 20, marginBottom: 20 },
  card: { padding: 20, background: "#fff", borderRadius: 8, minWidth: 200 },
};

export default Dashboard;
