import jsPDF from 'jspdf';

/**
 * Generates a PDF from table data in landscape orientation
 * @param {Array} tableData - The data from the traffic table
 * @param {Array} headers - The table headers configuration
 * @param {Object} filters - The current filter values (date, region, state, etc.)
 * @returns {jsPDF} The generated PDF document
 */
export function generateTablePDF(tableData, headers, filters) {
  // Create PDF in landscape orientation (A4 landscape)
  const doc = new jsPDF({
    orientation: 'landscape',
    unit: 'mm',
    format: 'a4'
  });

  const pageWidth = doc.internal.pageSize.getWidth();
  const pageHeight = doc.internal.pageSize.getHeight();
  const margin = 15;
  const tableStartY = 60;
  const rowHeight = 8;
  const headerHeight = 12;

  // Add title and date
  doc.setFontSize(16);
  doc.setFont('helvetica', 'bold');
  doc.text('Reporte de Tráfico', pageWidth / 2, 20, { align: 'center' });

  // Add filter information
  doc.setFontSize(10);
  doc.setFont('helvetica', 'normal');

  let filterText = '';
  if (filters.region) filterText += `Región: ${filters.region} | `;
  if (filters.state) filterText += `Estado: ${filters.state} | `;
  if (filters.ip) filterText += `OLT: ${filters.ip} | `;
  if (filters.gpon) filterText += `GPON: ${filters.gpon} | `;

  if (filters.initDate && filters.endDate) {
    const initDate = new Date(filters.initDate).toLocaleDateString();
    const endDate = new Date(filters.endDate).toLocaleDateString();
    filterText += `Período: ${initDate} - ${endDate}`;
  }

  if (filterText) {
    doc.text(filterText, pageWidth / 2, 30, { align: 'center' });
  }

  // Add generation date
  const generationDate = new Date().toLocaleString();
  doc.text(`Generado: ${generationDate}`, pageWidth / 2, 40, { align: 'center' });

  // Calculate column widths based on content
  const columnCount = headers.length;
  const availableWidth = pageWidth - (margin * 2);
  const columnWidth = availableWidth / columnCount;

  // Draw table headers
  doc.setFillColor(41, 41, 41); // Dark gray background
  doc.rect(margin, tableStartY, availableWidth, headerHeight, 'F');

  doc.setFontSize(9);
  doc.setFont('helvetica', 'bold');
  doc.setTextColor(255, 255, 255); // White text

  // Draw header text
  headers.forEach((header, index) => {
    const x = margin + (index * columnWidth) + (columnWidth / 2);
    doc.text(header, x, tableStartY + (headerHeight / 2), { align: 'center' });
  });

  // Draw table rows
  let currentY = tableStartY + headerHeight;

  tableData.forEach((row, rowIndex) => {
    // Check if we need a new page
    if (currentY + rowHeight > pageHeight - margin) {
      doc.addPage();
      currentY = margin;

      // Redraw headers on new page
      doc.setFillColor(41, 41, 41);
      doc.rect(margin, currentY, availableWidth, headerHeight, 'F');
      doc.setTextColor(255, 255, 255);

      headers.forEach((header, index) => {
        const x = margin + (index * columnWidth) + (columnWidth / 2);
        doc.text(header, x, currentY + (headerHeight / 2), { align: 'center' });
      });

      currentY += headerHeight;
    }

    // Alternate row colors for better readability
    if (rowIndex % 2 === 0) {
      doc.setFillColor(240, 240, 240); // Light gray
    } else {
      doc.setFillColor(255, 255, 255); // White
    }

    doc.rect(margin, currentY, availableWidth, rowHeight, 'F');

    // Draw row data
    doc.setFontSize(8);
    doc.setFont('helvetica', 'normal');
    doc.setTextColor(0, 0, 0); // Black text

    Object.values(row).forEach((value, colIndex) => {
      const x = margin + (colIndex * columnWidth) + (columnWidth / 2);
      const text = String(value || '');

      // Truncate long text
      const maxChars = Math.floor(columnWidth / 2); // Approximate character limit
      const displayText = text.length > maxChars ? text.substring(0, maxChars - 3) + '...' : text;

      doc.text(displayText, x, currentY + (rowHeight / 2), { align: 'center' });
    });

    currentY += rowHeight;
  });

  return doc;
}

/**
 * Extracts table data and headers from the DOM table element
 * @param {HTMLElement} tableElement - The table DOM element
 * @returns {Object} Object containing headers and data
 */
export function extractTableData(tableElement) {
  if (!tableElement) return { headers: [], data: [] };

  const headers = [];
  const data = [];

  // Extract headers from thead
  const headerRows = tableElement.querySelectorAll('thead tr');
  headerRows.forEach(row => {
    const thElements = row.querySelectorAll('th');
    thElements.forEach(th => {
      const text = th.textContent?.trim() || '';
      if (text && !headers.includes(text)) {
        headers.push(text);
      }
    });
  });

  // Extract data from tbody
  const dataRows = tableElement.querySelectorAll('tbody tr');
  dataRows.forEach(row => {
    const rowData = {};
    const tdElements = row.querySelectorAll('td');

    tdElements.forEach((td, index) => {
      const header = headers[index] || `Columna ${index + 1}`;
      rowData[header] = td.textContent?.trim() || '';
    });

    if (Object.keys(rowData).length > 0) {
      data.push(rowData);
    }
  });

  return { headers, data };
}

