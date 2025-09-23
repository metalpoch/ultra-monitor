import { useState } from "react";
import ChangeIcon from "../icons/ChangeIcon";

export default function ChangePasswordButton({ id, handleSubmit }) {
  const [openModal, setOpenModal] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [password, setPassword] = useState("")
  const [password_confirm, setPasswordConfirm] = useState("")

  const open = () => {
    setShowModal(true);
    setTimeout(() => setOpenModal(true), 10);

    setPassword("")
    setPasswordConfirm("")
  };

  const close = () => {
    setOpenModal(false);
    setTimeout(() => setShowModal(false), 300);
  };
 
  return (
    <>
      <button
        className="rounded-full p-2 duration-150 hover:bg-slate-950"
        onClick={open}
        aria-label="Abrir modal cambiar contraseña"
      >
        <ChangeIcon width={12} height={12} />
      </button>

      {showModal && <div
        className={`fixed inset-0 flex items-center justify-center bg-[#00000099] bg-opacity-80 backdrop-blur-sm z-1 transition-opacity duration-300 ${openModal ? "opacity-100" : "opacity-0"
          }`}
        onClick={close}
      >
        <div
          className="bg-[#121b31] p-6 rounded shadow-lg relative w-96"
          onClick={(e) => e.stopPropagation()}
        >
          <button
            className="absolute top-3 right-3 text-2xl font-bold"
            onClick={close}
            aria-label="Cerrar modal"
          >
            &times;
          </button>

          <h2 className="text-xl font-semibold mb-4">Contraseña temporal</h2>

          <form onSubmit={(e) => {
            handleSubmit(e, id, password, password_confirm)
            close()
          }} className="flex flex-col gap-4">
            <input
              type="password"
              placeholder="Contraseña"
              value={password}
              onChange={({ target }) => setPassword(target.value)}  
              className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
            <input
              type="password"
              placeholder="Confirmar contraseña"
              value={password_confirm}
              onChange={({ target }) => setPasswordConfirm(target.value)}
              className="border border-gray-300 rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
            <button
              type="submit"
              className="bg-blue-600 text-white rounded px-4 py-2 hover:bg-blue-700 transition"
            >
              Enviar
            </button>
          </form>
        </div>
      </div>}
    </>
  );
}
