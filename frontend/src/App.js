import { useEffect, useState } from "react";

function App() {
  const [metrics, setMetrics] = useState([]);

  useEffect(() => {
    fetch("http://localhost:5002/metrics")
      .then(res => res.json())
      .then(data => setMetrics(data));
  }, []);

  return (
    <div>
      <h1>Mini Analytics Dashboard</h1>

      {metrics.map((m, i) => (
        <div key={i}>
          {m.date} | {m.status} | {m.payment_method} | {m.total_orders} | R$ {m.total_revenue}
        </div>
      ))}
    </div>
  );
}

export default App;
