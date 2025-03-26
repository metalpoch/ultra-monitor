from os import error

import psycopg2

from src import constants


class Postgres:
    """
    Class for managing PostgreSQL database connections and operations.

    This class allows for connecting to a PostgreSQL database and retrieving
    schema information in CSV format.
    """

    def __init__(self, uri: str) -> None:
        """
        Initializes the Postgres class with the given database URI.

        Args:
            uri (str): The connection URI for the PostgreSQL database.
        """
        try:
            self.uri = uri
            self.conn = psycopg2.connect(uri)
        except psycopg2.Error as e:
            raise e

    def csv_schemas(self) -> str:
        """
        Retrieves the database schema information and formats it as a CSV string.

        This method queries the information schema of the connected database,
        excluding the "users" table, and returns the schema details in CSV format.

        Args:
            conn (connection): The active connection to the PostgreSQL database.

        Returns:
            str: A CSV string containing table, column, type, and relation information.
        """
        csv = "table,column,type,relation"
        cursor = self.conn.cursor()
        cursor.execute(constants.SQL_INFORMATION_SCHEMA)

        for res in cursor.fetchall():
            table, column_name, column_type, relation = res
            if table == "users":
                continue

            column_type = column_type.replace("timestamp with time zone", "datetime")
            relation = (
                str(relation).replace("()", "").replace("(", ".").replace(")", "")
            )

            csv += f"\n{table},{column_name},{column_type},{relation}"

        return csv

    def execute_sql(
        self, query: str
    ) -> tuple[list[dict], BaseException] | tuple[list[dict], None]:
        results = []
        cursor = self.conn.cursor()
        try:
            cursor.execute(query)
        except BaseException as e:
            self.conn.rollback()
            return results, e

        column_name = [desc[0] for desc in cursor.description]  # type: ignore

        rows = cursor.fetchall()
        for row in rows:
            results.append(dict(zip(column_name, row)))

        return results, None
