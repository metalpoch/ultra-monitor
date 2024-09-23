export class ErrorHandler {
    public status!: number;
    public message!: string;
    public solution?: string;

    constructor(status: number, message: string) {
        this.status = status;
        this.message = message;
        
        // TODO: Add solution for each error
    }
}