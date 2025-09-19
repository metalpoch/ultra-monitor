import { useState } from "react"
import InputField from "../ui/InputField"

const URL = `${import.meta.env.PUBLIC_URL || ""}/api/auth/signup`;
const TOKEN = sessionStorage.getItem("access_token").replace("Bearer ", "");

export default function CreateUserForm(){
  const [data, setData] = useState(undefined)
  const [error, setError] = useState(undefined)
  const [loading, setLoading] = useState(false)
  const [status, setStatus] = useState(undefined)
  const [p00, setP00] = useState("")
  const [fullname, setFullname] = useState("")
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [password_confirm, setPasswordConfirm] = useState("")
 
 if (status === 401) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(undefined);
    setData(undefined);

    try {
      const response = await fetch(URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${TOKEN}`,
        },
        body: JSON.stringify({
          p00: +p00,
          fullname,
          username,
          password,
          password_confirm,
        }),
      });

      setStatus(response.status);

      const responseData = await response.json();

      if (!response.ok) {
        setError(responseData);
      } else {
        setData(responseData);
      }
    } catch (err) {
      setError(err.message || "Error inesperado");
    } finally {
      setLoading(false);
    }
  };

  return <form onSubmit={handleSubmit}>
    <InputField
      id="p00"
      type="number"
      label="P00"
      value={p00}
      onChange={({target}) => setP00(target.value)}
    />
    <InputField
      id="fullname"
      label="Nombre Completo"
      value={fullname}
      onChange={({target}) => setFullname(target.value)}
    />
   <InputField
      id="username"
      label="Usuario corporativo"
      value={username}
      onChange={({target}) => setUsername(target.value)}
    />
    <InputField
      id="password"
      type="password"
      label="Contrasena"
      value={password}
      onChange={({target}) => setPassword(target.value)}
    />
    <InputField
      id="password_confirm"
      type="password"
      label="Confirmar contrasena"
      value={password_confirm}
      onChange={({target}) => setPasswordConfirm(target.value)}
    />

    <button
      type="submit"
      className="mt-4 w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
      >
       Crear
    </button>

    {error && <p className="text-red-400">{error.error}</p>}
    {data && <p>{data}</p>}

  </form>
}   
