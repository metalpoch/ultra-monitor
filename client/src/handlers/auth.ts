import type { Auth } from "../models/auth";
import { ErrorHandler } from "../lib/errors";

export class AuthHandler {
    private static url: string = import.meta.env.PUBLIC_API;

    static async login(data: Auth): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/auth/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            });
            const token = await response.json();
            if (response.status === 200) return { status: response.status, info: token }
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) }
        } catch (err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) }
        }
    }
}