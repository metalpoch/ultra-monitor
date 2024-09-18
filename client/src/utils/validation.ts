export function isValidEmail(email: string) {
    email = email.trim();
    const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return regex.test(email);
}

export function isValidPassword(password:string) {
    password = password.trim();
    return password.length >= 6;
}
