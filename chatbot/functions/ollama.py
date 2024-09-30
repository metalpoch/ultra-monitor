import ollama

def Chatbot(mensaje):
    response= ollama.chat(model='gemma2:2b',messages=[
        {
            'role': 'user',
            'content': mensaje
        }
    ])
    return response['message']['content']