import type { Register } from "../schemas/register";
import { ErrorHandler } from "../lib/errors";

export class RegisterService {
    private static url: string = import.meta.env.PUBLIC_API_AUTH;

    static async newUser(data: Register): Promise<{ status: (number | null), info: (null | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/auth/signup`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
            if (response.ok) return { status: response.status, info: null };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch (err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }
}