import { useState } from "react";
import InputField from "../ui/InputField"

const URL = `${import.meta.env.PUBLIC_URL || ""}/api/fat/`;

export default function UpdateFatsForm(){
  const [date, setDate] = useState("");
  const [file, setFile] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [submitSuccess, setSubmitSuccess] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!date || !file) {
      setSubmitError("Por favor, completa todos los campos");
      return;
    }

    const allowedTypes = ["text/csv", "text/plain"];
    const fileExtension = file.name.split('.').pop().toLowerCase();

    if (!allowedTypes.includes(file.type) && !["csv", "txt"].includes(fileExtension)) {
      setSubmitError("Por favor, sube un archivo CSV o TXT válido");
      return;
    }

    setIsSubmitting(true);
    setIsProcessing(true);
    setSubmitError("");
    setSubmitSuccess(false);

    try {
      const formData = new FormData();
      formData.append("date", date);
      formData.append("file", file);

      const token = sessionStorage.getItem("access_token")?.replace("Bearer ", "") || "";

      const response = await fetch(URL, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || errorData.error);
      }

      setSubmitSuccess(true);
      setDate("");
      setFile(null);

      const fileInput = document.getElementById("file");
      if (fileInput) fileInput.value = "";

    } catch (error) {
      setSubmitError(error.message || "Error al subir el archivo");
    } finally {
      setIsSubmitting(false);
      setIsProcessing(false);
    }
  };

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
    setSubmitError("");
  };

  const handleDateChange = (e) => {
    setDate(e.target.value);
    setSubmitError("");
  };

  return (
    <form onSubmit={handleSubmit}>
      <InputField
        id="date"
        type="date"
        label="Fecha perteneciente del reporte"
        value={date}
        onChange={handleDateChange}
        required
      />
      <InputField
        id="file"
        type="file"
        label="Reporte (CSV o TXT)"
        accept=".csv,.txt"
        onChange={handleFileChange}
        required
      />

      {isProcessing && (
        <div className="mt-4 flex justify-center">
          <span className="mx-auto py-20 loader"></span>
        </div>
      )}

      {submitError && (
        <div className="mt-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {submitError}
        </div>
      )}

      {submitSuccess && (
        <div className="mt-4 p-3 bg-green-100 border border-green-400 text-green-700 rounded">
          Archivo subido y actualizado correctamente
        </div>
      )}

      <button
        type="submit"
        disabled={isSubmitting}
        className="mt-4 w-full bg-blue-600 text-white py-2 rounded-md hover:bg-blue-700 transition disabled:bg-blue-400 disabled:cursor-not-allowed"
      >
        {isSubmitting ? "Subiendo..." : "Subir y actualizar"}
      </button>
    </form>
  )
}   
