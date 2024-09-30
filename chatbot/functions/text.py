import os
import re
from dotenv import load_dotenv
load_dotenv(override=True)

def UnionString(mensaje_user:str, schema:str)->str:
    m1= os.getenv('MENSAJE1')
    m2= os.getenv('MENSAJE2')
    m3= os.getenv('MENSAJE3')
    m4= os.getenv('MENSAJE4')
    if(mensaje_user==""):
        return "Error en el mensaje del cliente"
    
    return m1 +schema+". "+ m2  + mensaje_user+ ". " + m3 + m4


def extraer_sql(cadena):
    
    patron = r'```sql(.*?)```'
    resultado = re.search(patron, cadena, re.DOTALL)
    
    if resultado:
        return resultado.group(1).strip()  
    else:
        return None  