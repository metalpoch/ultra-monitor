import ollama


def chatbot(msg: str):
    response = ollama.chat(
        model="gemma2:2b", messages=[{"role": "user", "content": msg}]
    )
    return response["message"]["content"]
