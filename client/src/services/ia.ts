import type { QuestionIASchema } from "../schemas/question.ia"
import type { AnswerIASchema } from "../schemas/answer.ia"

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
    static async postQuestion(question: QuestionIASchema): Promise<{ status: (number | null), info: AnswerIASchema }> {
        const response = await fetch(`${this.url}/chatbox`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(question)
        });
        return { status: response.status, info: await response.json() }
    }
}