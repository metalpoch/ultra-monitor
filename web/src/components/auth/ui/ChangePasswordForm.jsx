import { useState } from "react";

const BASE_URL = `${import.meta.env.PUBLIC_URL || ""}/api/auth`

export default function ChangePasswordForm({token}) {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [resetError, setResetError] = useState(null);
  const [resetLoading, setResetLoading] = useState(false);

  const handleReset = async (e) => {
    e.preventDefault();
    setResetLoading(true);
    setResetError(null);

    const response = await fetch(`${BASE_URL}/reset_password`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        password: currentPassword,
        new_password: newPassword,
        password_confirm: confirmPassword,
      }),
    })
    .then((res) => {
        if (res.status === 200) {
            sessionStorage.setItem("access_token", `Bearer ${token}`);
            window.location.href = "/";
          return
        }
        return res.json()
      })
      .catch((error) => error)

    if (response.error) {
      setResetError(response.error || "Error al cambiar la contraseña");
      setResetLoading(false);
    }
  };

if (resetLoading) {
    return <span className="mx-auto py-20 loader"></span>;
  }

  return (
    <form
      method="PATCH"
      onSubmit={handleReset}
      className="flex flex-col gap-2 mt-8"
    >
      <label className="text-gray-200 mt-4">Contraseña actual</label>
      <section className="flex items-center gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500">
        <input
          className="w-full bg-[#f1729] focus:outline-none"
          type="password"
          value={currentPassword}
          onChange={(e) => setCurrentPassword(e.target.value)}
          placeholder="Contraseña actual"
          required
          autoFocus
        />
      </section>

      <label className="text-gray-200 mt-4">Nueva contraseña</label>
      <section className="flex items-center gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500">
        <input
          className="w-full bg-[#f1729] focus:outline-none"
          type="password"
          value={newPassword}
          onChange={(e) => setNewPassword(e.target.value)}
          placeholder="Nueva contraseña"
          required
        />
      </section>

      <label className="text-gray-200 mt-4">Confirmar nueva contraseña</label>
      <section className="flex items-center gap-2 border border-[hsl(217,33%,20%)] bg-[#0f1729] p-2 rounded-xs focus-within:ring-2 focus-within:ring-blue-500">
        <input
          className="w-full bg-[#f1729] focus:outline-none"
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          placeholder="Confirmar nueva contraseña"
          required
        />
      </section>

      {resetError && <p className="text-center text-red-300">{resetError}</p>}
      <button
        type="submit"
        className="mt-4 w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition"
        disabled={resetLoading}
      >
        Cambiar contraseña
      </button>
    </form>
  );
}
