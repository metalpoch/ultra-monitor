import { atom } from "nanostores";
import { IAService } from "../services/ia";

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

        const response = await IAService.postQuestion(question);

        if (response.status === 200) {
            let answer = response.info.output;
            answerIA.set(answer);
            loadingIA.set(false);
        } else {
            answerIA.set("No se pudo responder a tu pregunta. Inténtalo más tarde.");
            loadingIA.set(false);
        }
    }
}