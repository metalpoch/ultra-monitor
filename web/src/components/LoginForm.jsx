import { useState } from 'react';
import PasswordIcon from './icons/PasswordIcon';
import UserIcon from './icons/UserIcon';

const BASE_URL = import.meta.env.PUBLIC_AUTH_URL

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
  const [username, setUsername] = useState(null);
  const [password, setPassword] = useState(null);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    setError(null);
    e.preventDefault();
    const response = await userLogin(username, password);
    if (response.token) {
      const token = response.token;
      sessionStorage.setItem("access_token", JSON.stringify(token));
      window.location.href = "/"
    } else {
      setError(response.error);
    }
  }

  return (
    <>
      <form method="POST" onSubmit={handleSubmit} className="flex flex-col gap-1">
        <label htmlFor="username" className="text-gray-200">Usuario</label>
        <section className="flex items-center gap-1 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500">
          <UserIcon className="w-5" />
          <input
            className="w-full bg-[#f1729] focus:outline-none"
            id="username"
            type="text"
            onChange={({ target }) => setUsername(target.value)}
            placeholder="Ingresa tu usuario"
            required
          />
        </section>

        <label htmlFor="password" className="text-gray-200 mt-4">Contraseña</label>
        <section className="flex items-center gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500">
          <PasswordIcon className="w-5" />
          <input
            className="w-full bg-[#f1729] focus:outline-none"
            id="password"
            type="password"
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

