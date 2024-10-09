import ollama


def Chatbot(msg: str):
    """Run the promt in the ollama AI."""
    response = ollama.chat(
        model="gemma2:2b", messages=[{"role": "user", "content": msg}]
    )
    return response["message"]["content"]
