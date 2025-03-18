/**
 * Handler for error responses to API requests.
 */
export class ErrorHandler {
    public status!: number;
    public message!: string;

    constructor(status: number, message: string) {
        this.status = status;
        this.message = message;
    }
}