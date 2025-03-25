import { atom } from "nanostores";
import { IAService } from "../services/ia";
import type { QuestionIASchema } from "../schemas/question.ia";

export const questionIA = atom<string | null>(null);
export const answerIA = atom<string | null>(null);
export const loadingIA = atom<boolean>(false);

/**
 * @class Controller for all requests to the IA.
 */
export class IAController {

    /**
     * Send a question to the IA.
     * 
     * @param {string} question Text of the question.
     */
    static async postQuestion(question: string): Promise<void> {
        loadingIA.set(true);
        questionIA.set(question);
        let currentQuestion: QuestionIASchema = { message: question }
        let response = await IAService.postQuestion(currentQuestion)
        if (response.status === 200) {
            answerIA.set(response.info.answer);
            loadingIA.set(false);
        } else {
            answerIA.set("No se pudo responder a tu pregunta. Inténtalo más tarde.");
            loadingIA.set(false);
        }
    }
}