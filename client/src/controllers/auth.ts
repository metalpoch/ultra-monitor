import { AuthHandler } from "../handlers/auth";
import { ErrorHandler } from "../lib/errors";
import { ResponseHandler } from "../lib/response";
import { FailedMessage } from "../constant/message";
import type { Auth } from "../models/auth";

export class AuthController {
    static async login(email: string, password: string): Promise<ResponseHandler> {
        const data: Auth = {
            email: email,
            password: password
        }
        const response = await AuthHandler.login(data);
        if (response.status === 200) return new ResponseHandler(true, "", response.info);
        else if (response.status === 401) return new ResponseHandler(false, FailedMessage.loginFailed, null, response.info as ErrorHandler);
        else {
            console.error('Error Status:', response.info?.status, 'ErrorMessage:', response.info!.message);
            return new ResponseHandler(false, FailedMessage.loginFailed, null, response.info as ErrorHandler);
        }
    }
}