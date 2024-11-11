import { ErrorHandler } from "../lib/errors";
import type { QuestionIA } from "../models/question.ia"
import type { AnswerIA } from "../models/answer.ia"

export class IAService {
    private static url: string = import.meta.env.PUBLIC_API_SMART;

    static async postQuestion(question: QuestionIA): Promise<{ status: (number | null), info: AnswerIA }> {
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