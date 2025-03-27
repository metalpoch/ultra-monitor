/**
 * Check if the email is valid.
 * 
 * @param {string} email Email to check.
 * @returns {boolean} True if the email is valid, false otherwise.
 */
export function isValidEmail(email: string): boolean {
    email = email.trim();
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
}


/**
 * Check if the password is valid.
 * 
 * @param {string} password Password to check.
 * @returns {boolean} True if the password is valid, false otherwise.
 */
export function isValidPassword(password: string): boolean {
    // const regex = /[ `!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]/;
    password = password.trim();
    // return password.length >= 0 && regex.test(password);
    return password.length >= 0;
}
