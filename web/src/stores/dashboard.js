import { atom } from "nanostores";

export const selectedLevel = atom("");
export const selectedRegion = atom("");
export const selectedState = atom("");

// Store para datos diarios de tráfico
export const dailyTrafficData = atom({
  data: [], // Array de objetos con volume_in y volume_out diarios
  loading: true,
  error: null,
  lastState: null, // Último estado cargado
  lastRegion: null // Última región cargada
});

// Store para datos de usuarios
export const userData = atom({
  activeUsers: 0,
  loading: true,
  error: null,
  lastState: null, // Último estado cargado
  lastRegion: null // Última región cargada
});
