import { useState } from "react";
import { api, register } from "../services/api";

export default function Login({ onLogin }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isRegister, setIsRegister] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      if (isRegister) {
        const res = await register(username, password);
        if (!res.success) {
          setError(res.message);
          return;
        }
        setIsRegister(false);
        setError("Usuário criado. Faça login.");
        return;
      }

      const res = await api.post("/login", { username, password });
      localStorage.setItem("token", res.data.token);
      onLogin();
    } catch {
      setError("Usuário ou senha inválidos");
    }
  };

  return (
    <div style={styles.container}>
      <form onSubmit={handleSubmit} style={styles.box}>
        <h2>{isRegister ? "Cadastrar" : "Login"}</h2>

        <input
          placeholder="Usuário"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />

        <input
          type="password"
          placeholder="Senha"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />

        <button type="submit">
          {isRegister ? "Cadastrar" : "Entrar"}
        </button>

        <button
          type="button"
          onClick={() => {
            setIsRegister(!isRegister);
            setError("");
          }}
          style={{ background: "transparent", color: "#38bdf8" }}
        >
          {isRegister ? "Já tenho conta" : "Criar conta"}
        </button>

        {error && <p style={{ color: "red" }}>{error}</p>}
      </form>
    </div>
  );
}

const styles = {
  container: {
    height: "100vh",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    background: "#0f172a",
  },
  box: {
    background: "#020617",
    padding: 30,
    borderRadius: 8,
    display: "flex",
    flexDirection: "column",
    gap: 10,
    color: "white",
    minWidth: 300,
  },
};