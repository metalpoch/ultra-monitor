import { useState } from "react";
import PasswordIcon from "../icons/PasswordIcon";
import UserIcon from "../icons/UserIcon";
import ChangePasswordForm from "./ui/ChangePasswordForm";

const BASE_URL = `${import.meta.env.PUBLIC_API_URL}/auth`;

const userLogin = async (username, password) => {
  return fetch(`${BASE_URL}/signin`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  })
    .then((res) => res.json())
    .finally((res) => res)
    .catch((err) => console.error(err));
};

export default function LoginForm() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [token, setToken] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const response = await userLogin(username, password);
    response.token && setToken(response.token)
    if (!response.error && response.token) {
      sessionStorage.setItem("access_token", `Bearer ${response.token}`);
      window.location.href = "/";
    } else {
      setError(response.error);
      setLoading(false);
    }
  };

  if (error === "you must change your password") {
    return <ChangePasswordForm token={token} />;
  }
  
  if (loading) {
    return <span className="mx-auto py-20 loader"></span>;
  }

  return (
    <>
      <form
        method="POST"
        onSubmit={handleSubmit}
        className="flex flex-col gap-1"
      >
        <label htmlFor="username" className="text-gray-200">
          Usuario
        </label>
        <section className="flex items-center gap-1 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500">
          <UserIcon className="w-5" />
          <input
            className="w-full bg-[#f1729] focus:outline-none"
            id="username"
            type="text"
            value={username}
            onChange={({ target }) => setUsername(target.value)}
            placeholder="Ingresa tu usuario"
            required
            autoFocus
          />
        </section>

        <label htmlFor="password" className="text-gray-200 mt-4">
          Contraseña
        </label>
        <section className="flex items-center gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500">
          <PasswordIcon className="w-5" />
          <input
            className="w-full bg-[#f1729] focus:outline-none"
            id="password"
            type="password"
            value={password}
            onChange={({ target }) => setPassword(target.value)}
            placeholder="Ingresa tu contraseña"
            required
          />
        </section>
        {error && <p className="text-center text-red-300">{error}</p>}
        <button
          type="submit"
          className="mt-4 w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
        >
          Ingresar
        </button>
      </form>
    </>
  );
}
