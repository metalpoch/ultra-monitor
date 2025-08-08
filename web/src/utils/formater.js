export function removeAccentsAndToUpper(text) {
  const textWithoutAccents = text
    .normalize("NFD")
    .replace(/[\u0300-\u036f]/g, "");
  return textWithoutAccents.toUpperCase();
}
