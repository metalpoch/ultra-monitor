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

PROMP_1 = "Eres un experto en bases de datos Postgresql con mas de 5 años de experiencia, el schema de la base de datos es como el que se muestra en el siguiente csv:"

PROMP_2 = "El cliente indiciará una instrucción en lenguaje natural, como experto debes transformarlo en un script sql que sea capaz de dar respuesta a la pregunta del cliente de la manera mas acertada posible, la pregunta del cliente es: "

PROMP_3 = """
Recuerda solo debes usar el schema que se te ha indicado.
La respuesta que debes dar es solo el script sql nada mas que eso
"""

PROMP_4 = "Eres un experto en bases de datos Postgresql con mas de 5 años de experiencia, tu trabajo es verificar si el siguiente codigo sql contiene alguna de las siguientes instrucciones: INSERT, UPDATE, DELETE, REVOKE, DROP TABLE, ALTER TABLE, CREATE TABLE. El codigo SQL que debes revisar es: "

PROMP_5 = "La respuesta que debes dar es un True si contiene alguna de las instrucciones y False si no contiene ninguna."

PROMP_6 = "Solo debes responder True o False segun sea el caso nasa mas."
