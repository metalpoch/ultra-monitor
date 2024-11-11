import { atom } from "nanostores";
import { IAService } from "../services/ia";
import type { QuestionIA } from "../models/question.ia";

export const questionIA = atom<string | null>(null);
export const answerIA = atom<string | null>(null);

export class IAController {
    static async postQuestion(question: string): Promise<void> {
        questionIA.set(question);
        let current_question: QuestionIA = { message: question }
        let response = await IAService.postQuestion(current_question)
        answerIA.set(response.info.answer);
    }
}