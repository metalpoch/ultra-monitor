import sys

import pandas as pd


def proccess_data(df: pd.DataFrame) -> pd.DataFrame:
    """
    Procesa un DataFrame agrupando por columnas específicas y calculando
    el número de registros por cada valor específico de STATUS en cada grupo.

    Args:
        df (pd.DataFrame): DataFrame de entrada que debe contener al menos
            las columnas: REGION, ESTADO, MUNICIPIO, PARROQUIA, SLOT, PUERTO,
            IP_OLT, ODN, FAT y STATUS.

    Returns:
        pd.DataFrame: Nuevo DataFrame agrupado por las columnas especificadas
        con columnas adicionales para el conteo de cada valor de STATUS:
        - ACTIVOS: conteo de STATUS=='ACTIVO' por grupo.
        - EN_PROCESO: conteo de STATUS=='EN PROCESO' por grupo.
        - CORTADO: conteo de STATUS=='CORTADO' por grupo.
        - APROVISIONADO_OFFLINE: conteo de STATUS=='APROVISIONADO OFFLINE' por grupo.
        - TOTAL: total de registros por grupo.
    """
    cols = [
        "REGION",
        "ESTADO",
        "MUNICIPIO",
        "PARROQUIA",
        "SLOT",
        "PUERTO",
        "IP_OLT",
        "HOSTNAME",
        "ODN",
        "FAT",
    ]

    grouped = (
        df.groupby(cols)
        .agg(
            ACTIVOS=pd.NamedAgg(
                column="STATUS", aggfunc=lambda x: (x == "ACTIVO").sum()
            ),
            EN_PROCESO=pd.NamedAgg(
                column="STATUS", aggfunc=lambda x: (x == "EN PROCESO").sum()
            ),
            CORTADO=pd.NamedAgg(
                column="STATUS", aggfunc=lambda x: (x == "CORTADO").sum()
            ),
            APROVISIONADO_OFFLINE=pd.NamedAgg(
                column="STATUS", aggfunc=lambda x: (x == "APROVISIONADO OFFLINE").sum()
            ),
        )
        .reset_index()
    )

    return grouped


def main() -> None:
    """
    Función principal que maneja la ejecución del script desde línea de comandos.

    El script espera un único argumento: el nombre de un archivo CSV
    que contiene el reporte de los fats para procesar.

    El archivo CSV debe usar ';' como separador y contener las columnas
    necesarias para el procesamiento.

    La función carga el archivo, llama a 'proccess_data' y imprime el resultado.

    Si no se proporciona exactamente un argumento, muestra un mensaje de error y termina la ejecución.
    """
    if len(sys.argv) != 2:
        sys.exit("se requiere como unico parametro el reporte de los fats")

    try:
        filename = sys.argv[1]
        df = pd.read_csv(filename, sep=";")
        data = proccess_data(df)
        print(data.to_json(orient="records"))
    except Exception as e:
        print(f"error manipulando el csv de fats {type(e).__name__}: {e}")


if __name__ == "__main__":
    main()
