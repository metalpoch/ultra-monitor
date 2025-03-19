import re
from typing import Sequence

import ollama

from src import constants


class AI:
    """
    Class for interacting with an artificial intelligence model.

    This class allows for creating prompts and querying the specified model.
    """

    def __init__(self, model: str, schemas: str) -> None:
        """
        Initializes the AI class with the given model.

        Args:
            model (str): The name of the artificial intelligence model to use.
        """
        self.model = model
        self.messages = [
            {
                "role": "system",
                "content": f"{constants.PROMPT_1}\n```csv\n{schemas}\n```",
            },
            {
                "role": "system",
                "content": constants.PROMPT_2,
            },
            {
                "role": "system",
                "content": constants.PROMPT_3,
            },
        ]

    def query(self, body: str) -> str:
        """
        Queries the model using the provided schemas and body content.

        Args:
            schemas (str): Schemas in CSV format to be sent in the query.
            body (str): Additional content to be sent in the query.

        Returns:
            str: The response from the model to the executed query.
        """
        self.messages.append({"role": "user", "content": body})

        sequence = self.messages  # type: Sequence
        res = ollama.chat(model=self.model, messages=sequence)
        self.messages.append(
            {"role": "assistant", "content": res["message"]["content"]}
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
