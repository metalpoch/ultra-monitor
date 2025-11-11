import { atom } from 'nanostores';

// Store for PDF header configuration
export const pdfHeaderConfig = atom({
  headers: [],
  columnCount: 0
});

// Helper function to extract text from React elements
export function extractHeaderText(headerElement) {
  if (!headerElement) return [];

  const headers = [];

  // Handle the header structure based on different filter types
  if (headerElement.props && headerElement.props.children) {
    const children = Array.isArray(headerElement.props.children)
      ? headerElement.props.children
      : [headerElement.props.children];

    children.forEach(child => {
      if (child.type === 'tr') {
        const thElements = Array.isArray(child.props.children)
          ? child.props.children
          : [child.props.children];

        thElements.forEach(th => {
          if (th.type === 'th') {
            const text = th.props.children;
            if (text) {
              headers.push(String(text));
            }
          }
        });
      }
    });
  }

  return headers;
}

// Function to calculate column count based on filter type
export function getColumnCount(filterType) {
  switch (filterType) {
    case 'gpon':
      return 7; // Puerto, Prom. Entrante, Max. Entrante, Prom. Saliente, Max. Saliente, Capacidad, Uso %
    case 'ip':
      return 10; // Puerto, Prom. Entrante, Max. Entrante, Prom. Saliente, Max. Saliente, Capacidad, Uso, Activo, Cortado, En progreso
    case 'state':
      return 11; // OLT, Agregador, Prom. Entrante, Max. Entrante, Prom. Saliente, Max. Saliente, Capacidad, Uso, Activo, Cortado, En progreso
    case 'region':
      return 10; // Estado, Prom. Entrante, Max. Entrante, Prom. Saliente, Max. Saliente, Capacidad, Uso, Activo, Cortado, En progreso
    default:
      return 0;
  }
}

