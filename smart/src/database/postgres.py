import psycopg2
from psycopg2.extensions import connection

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
        self.uri = uri

    def connect(self) -> connection:
        """
        Establishes a connection to the PostgreSQL database.

        Returns:
            connection: A connection object to the PostgreSQL database.

        Raises:
            psycopg2.Error: If there is an error while connecting to the database.
        """
        try:
            conn_dict = psycopg2.connect(self.uri)
        except psycopg2.Error as e:
            raise e
        else:
            return conn_dict

    def csv_schemas(self, conn: connection) -> str:
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
        cursor = conn.cursor()
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
