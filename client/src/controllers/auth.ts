import { AuthService } from "../services/auth";
import { ErrorHandler } from "../lib/errors";
import { ResponseHandler } from "../lib/response";
import { FailedMessage } from "../constant/message";
import type { Auth } from "../models/auth";
import type { User } from "../models/user";
import { atom } from "nanostores";

export const user = atom<User | null>(null);

export class AuthController {
    static async login(email: string, password: string): Promise<ResponseHandler> {
        const data: Auth = {
            email: email,
            password: password
        }
        const response = await AuthService.login(data);
        if (response.status === 200) {
            const data: User = {
                p00: response.info.p00,
                email: response.info.email,
                change_password: response.info.change_password,
                fullname: response.info.fullname,
                is_admin: response.info.is_admin
            }
            user.set(data);
            return new ResponseHandler(true, "", response.info.token);
        }
        else if (response.status === 401) return new ResponseHandler(false, FailedMessage.loginFailed, null, response.info as ErrorHandler);
        else {
            console.error('Error Status:', response.info?.status, 'ErrorMessage:', response.info!.message);
            return new ResponseHandler(false, FailedMessage.loginFailed, null, response.info as ErrorHandler);
        }
    }
}