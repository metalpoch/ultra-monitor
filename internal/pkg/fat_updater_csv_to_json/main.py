import sys
import re

import pandas as pd


def extract_plan_speed(plan_string: str) -> str:
    """
    Extrae la velocidad de bajada/subida de un string de plan.

    Args:
        plan_string (str): String que contiene información del plan

    Returns:
        str: String con formato 'bajada/subida' o cadena vacía si no se encuentra
    """
    if pd.isna(plan_string) or not isinstance(plan_string, str):
        return ""

    # Buscar patrones como "300/300", "100/50", "15/6", etc.
    pattern = r'(\d+)/(\d+)'
    match = re.search(pattern, plan_string)

    if match:
        return match.group(0)
    return ""


def proccess_data(df: pd.DataFrame) -> pd.DataFrame:
    """
    Procesa un DataFrame agrupando por columnas específicas y calculando
    el número de registros por cada valor específico de STATUS en cada grupo.

    Args:
        df (pd.DataFrame): DataFrame de entrada que debe contener al menos
            las columnas: REGION, ESTADO, MUNICIPIO, PARROQUIA, SLOT, PUERTO,
            IP_OLT, ODN, FAT, STATUS y PLAN.

    Returns:
        pd.DataFrame: Nuevo DataFrame agrupado por las columnas especificadas
        con columnas adicionales para el conteo de cada valor de STATUS:
        - ACTIVOS: conteo de STATUS=='ACTIVO' por grupo.
        - EN_PROCESO: conteo de STATUS=='EN PROCESO' por grupo.
        - CORTADO: conteo de STATUS=='CORTADO' por grupo.
        - APROVISIONADO_OFFLINE: conteo de STATUS=='APROVISIONADO OFFLINE' por grupo.
        - plans: string con formato 'countxplan;countxplan;...' que agrupa los planes
    """
    # Primero extraemos las velocidades de los planes
    df['PLAN_SPEED'] = df['PLAN'].apply(extract_plan_speed)

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

    def aggregate_plans(plans_series):
        """Agrega los planes en formato 'countxplan;countxplan;...'"""
        plan_counts = plans_series.value_counts()
        plan_strings = []
        for plan, count in plan_counts.items():
            if plan:  # Solo incluir planes no vacíos
                plan_strings.append(f"{count}x{plan}")
        return ";".join(plan_strings)

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
            plans=pd.NamedAgg(
                column="PLAN_SPEED", aggfunc=aggregate_plans
            )
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
        # Solo cargar las columnas necesarias para optimizar recursos
        required_columns = [
            "REGION", "ESTADO", "MUNICIPIO", "PARROQUIA", "SLOT", "PUERTO",
            "IP_OLT", "HOSTNAME", "ODN", "FAT", "STATUS", "PLAN"
        ]
        df = pd.read_csv(filename, sep=";", usecols=required_columns)
        data = proccess_data(df)
        print(data.to_json(orient="records"))
    except Exception as e:
        print(f"error manipulando el csv de fats {type(e).__name__}: {e}")


if __name__ == "__main__":
    main()
