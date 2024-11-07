import { atom } from "nanostores";

export const questionIA = atom<string | null>(null);
export const answerIA = atom<string | null>(null);

export class IAController {
    static async postQuestion(question: string): Promise<void> {
        questionIA.set(question);
        answerIA.set(question);
    }
}