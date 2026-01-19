import { useState } from "react";
import { api } from "../services/api";

export default function Register() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState("");

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      await api.post("/register", { username, password });
      setMsg("Usuário criado com sucesso!");
    } catch {
      setMsg("Erro ao cadastrar");
    }
  };

  return (
    <form onSubmit={handleRegister}>
      <input placeholder="Usuário" onChange={e => setUsername(e.target.value)} />
      <input type="password" placeholder="Senha" onChange={e => setPassword(e.target.value)} />
      <button type="submit">Cadastrar</button>
      {msg && <p>{msg}</p>}
    </form>
  );
}
