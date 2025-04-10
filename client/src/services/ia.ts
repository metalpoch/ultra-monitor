import type { QuestionIASchema } from "../schemas/question.ia";
import type { AnswerIASchema } from "../schemas/answer.ia";
import { ErrorHandler } from "../lib/errors";

/**
 * @class Handler of the IA requests for the API.
 */
export class IAService {
    private static url: string = import.meta.env.PUBLIC_API_SMART;

    /**
     * Request API to post a question to the IA.
     * 
     * @param {QuestionIASchema} question Question to be posted.
     */
    static async postQuestion(question: QuestionIASchema): Promise<{ status: (number | null), info: (any | ErrorHandler) }> {
        try {
            const response = await fetch(`${this.url}/chatbox`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(question)
            });
            if (response.ok) return { status: response.status, info: await response.json() as AnswerIASchema };
            else return { status: response.status, info: new ErrorHandler(response.status, response.statusText) };
        } catch(err) {
            return { status: null, info: new ErrorHandler(500, (err as Error).name + ": " + (err as Error).message) };
        }
    }
}