/**
 * Get the string of a month.
 * 
 * @param {number} month Month to consult.
 */
export function getMonthString(month: number): string {
    const months = ['Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio', 'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'];
    return months[month];
}

/**
 * Get the past year.
 */
export function getPastYear(): number {
    const now = new Date();
    const year = now.getFullYear() - 1;
    return year;
}

/**
 * Get the past month.
 */
export function getPastMonth(): number {
    const now = new Date();
    let month = now.getMonth();
    if (month == 0) month = 12;
    if (month == 1) month = 1;
    return month;
}