import { useState, useEffect } from "react";
import useFetch from "../../hooks/useFetch";
import ChangeStatusDisableButton from "./ChangeStatusDisableUserButton";
import ChangePasswordButton from "./ChangePasswordButton"

const URL = `${import.meta.env.PUBLIC_URL || ""}/api/auth/`;
const TOKEN = sessionStorage.getItem("access_token").replace("Bearer ", "");

export default function UsersTable() {
  const { data: fetchData, loading, error, status } = useFetch(URL, { headers: { Authorization: `Bearer ${TOKEN}` } });
  const [users, setUsers] = useState([]);

  useEffect(() => {
    if (fetchData) {
      setUsers(fetchData);
    }
  }, [fetchData ]);

  if (status === 401 || status === 403) {
    sessionStorage.removeItem("access_token");
    window.location.href = "/";
  }

  const statusHandler = async (id, isDisabled) => {
    const method = isDisabled ? "POST" : "DELETE";
    try {
      const response = await fetch(URL + id, { method, headers: { Authorization: `Bearer ${TOKEN}` } });
      if (response.status === 204) {
        setUsers((prevUsers) =>
          prevUsers.map((user) =>
            user.id === id ? { ...user, is_disabled: true } : user
          )
        );
      } else if (response.status === 200) {
        setUsers((prevUsers) =>
          prevUsers.map((user) =>
            user.id === id ? { ...user, is_disabled: false } : user
          )
        );
      }
    } catch (err) {
      alert(err.message);
    }
  };

  const handleSubmit = async (e, id, password, password_confirm) => {
    e.preventDefault();
    try {
      const response = await fetch(URL + `temporal_passw/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${TOKEN}` },
        body: JSON.stringify({ password, password_confirm })
      });

      if (response.status === 200) {
        setUsers((prevUsers) =>
          prevUsers.map((user) =>
            user.id === id ? { ...user, change_password: true } : user
          )
        );
      } else {
        response.json().then(res => alert(res.error))
      }
    } catch (err) {
      alert(err.message);
    }
  };

  return (
    <table className="min-w-full table-auto border-collapse">
      <thead className="sticky top-0 bg-[#121b31] pb-1 text-left">
        <tr>
          <th></th>
          <th>ID</th>
          <th>Nombre Completo</th>
          <th>Usuario</th>
          <th>Admin</th>
          <th>Cambiar Contraseña</th>
          <th>Fecha Creación</th>
        </tr>
      </thead>
      <tbody>
        {users && users.map((user) => (
          <tr key={user.id} className={`duration-150 hover:bg-slate-800 ${user.is_disabled ? "text-slate-500" : ""}`}>
            <td className="py-1 flex items-center gap-2">
              <ChangeStatusDisableButton handler={statusHandler} id={user.id} isDisable={user.is_disabled} />
            </td>
            <td>{user.id}</td>
            <td>{user.fullname}</td>
            <td>{user.username}</td>
            <td>{user.is_admin ? "Sí" : "No"}</td>
            <td className="flex gap-5">
              <ChangePasswordButton id={user.id} handleSubmit={handleSubmit} />
              {user.change_password ? "Sí" : "No"}
            </td>
            <td>{new Date(user.created_at).toLocaleDateString()}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
