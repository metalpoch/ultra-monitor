from functions.ollama import *
from functions.text import *
from database.postgre import *

try:
    schema = Get_database_schema()
except Exception as e:
    print("Error al conectar a la base de datos:", e)

print("introduce el mensaje que desee")
mensaje=input()

promp= UnionString(mensaje,schema)

ia=Chatbot(promp)

resustado=extraer_sql(ia)
print("SOY LA RESPUEATA DE LA IA: ",resustado)
print()

Ejecute_SQL(resustado)