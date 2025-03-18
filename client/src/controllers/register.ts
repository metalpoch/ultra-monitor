import type { Register } from "../schemas/register";
import { RegisterHandler } from "../handlers/register";
import { SuccessMessage, FailedMessage } from "../constant/message";
import { ResponseHandler } from "../lib/response";
import type { ErrorHandler } from "../lib/errors";

export class RegisterController {
    static async newUser(p00: number, email: string, password: string, passwordConfirm: string, names: string, lastnames: string): Promise<ResponseHandler> {
        const register: Register = {
            p00: p00,
            email: email,
            password: password,
            password_confirm: passwordConfirm,
            fullname: `${names} ${lastnames}`
        }
        const response = await RegisterHandler.newUser(register);
        if (response.status === 200 || response.status === 201) return new ResponseHandler(true, SuccessMessage.newUserSuccess);
        else {
            console.error('Error Status:', response.info?.status, 'ErrorMessage:', response.info!.message);
            return new ResponseHandler(false, FailedMessage.newUserFailed, null, response.info as ErrorHandler);
        }
    }
}