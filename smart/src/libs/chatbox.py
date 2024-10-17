import re

import ollama

from src import constants


class AI:
    """
    Class for interacting with an artificial intelligence model.

    This class allows for creating prompts and querying the specified model.
    """

    def __init__(self, model: str) -> None:
        """
        Initializes the AI class with the given model.

        Args:
            model (str): The name of the artificial intelligence model to use.
        """
        self.model = model

    def __create_prompt(self, schemas: str, body: str) -> str:
        """
        Creates a formatted prompt for querying the model.

        Args:
            schemas (str): Schemas in CSV format to be included in the prompt.
            body (str): Additional content to be included in the prompt.

        Returns:
            str: The complete prompt that will be sent to the model.
        """
        p1 = f"{constants.PROMPT_1}\n\n```csv\n{schemas}\n```\n"
        p2 = f"{constants.PROMPT_2}\n"
        p3 = f"\n{constants.PROMPT_3}"
        return p1 + p2 + body + p3

    def query(self, schemas: str, body: str) -> str:
        """
        Queries the model using the provided schemas and body content.

        Args:
            schemas (str): Schemas in CSV format to be sent in the query.
            body (str): Additional content to be sent in the query.

        Returns:
            str: The response from the model to the executed query.
        """
        res = ollama.chat(
            model=self.model,
            messages=[
                {
                    "role": "user",
                    "content": self.__create_prompt(schemas=schemas, body=body),
                }
            ],
        )
        return res["message"]["content"]

    def sql_extract(self, msg: str) -> str:
        """
        Extracts an SQL query from the given message.

        Searches for a block of SQL code within the message and returns it as a string.

        Args:
            msg (str): The message from which to extract the SQL query.

        Returns:
            str: The extracted SQL query or an empty string if none is found.
        """
        pattern = r"```sql(.*?)```"
        result = re.search(pattern, msg, re.DOTALL)
        return "" if not result else result.group(1).replace("\n", " ").strip()
