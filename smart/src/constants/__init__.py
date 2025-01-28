SQL_INFORMATION_SCHEMA = """
SELECT 
    c.table_name,
    c.column_name,
    c.data_type,
    CONCAT(ccu.table_name, '(', ccu.column_name, ')') AS relacionada
FROM
    information_schema.columns c
LEFT JOIN
    information_schema.key_column_usage kcu 
ON
    c.column_name = kcu.column_name AND c.table_name = kcu.table_name
LEFT JOIN
    information_schema.table_constraints tc 
ON
    kcu.constraint_name = tc.constraint_name AND tc.constraint_type = 'FOREIGN KEY'
LEFT JOIN
    information_schema.constraint_column_usage ccu 
ON
    kcu.constraint_name = ccu.constraint_name
WHERE
    c.table_schema = 'public'
ORDER BY c.table_name
"""

PROMPT_1 = """Eres un experto en bases de datos Postgresql con mas de 5 años de experiencia, el schema de la base de datos es como el que se muestra en el siguiente csv:
"""

PROMPT_2 = """
Recuerda solo debes usar el schema que se te ha indicado, sin olvidar usar la sentencia inner en caso de trabajar con varias tablas.
La respuesta que debes dar es solo la sentencia sql, nada mas que eso.
"""

PROMPT_3 = """
El cliente indiciará una instrucción en lenguaje natural, como experto debes transformarlo en una sentencia sql que sea capaz de dar respuesta a la pregunta del cliente de la manera mas acertada posible, No olvides que la respuesta que debes dar es solo la sentencia SQL y cuando el cliente se refiera a entrante quiere hacer referencia a los datos in, y si es saliente out. No olvides el inner en la sentenia SQL.
"""
