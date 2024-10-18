import requests


class Telegram:
    def __init__(self, bot_id: str, chat_id: str) -> None:
        self.chat_id = chat_id
        self.bot_id = bot_id

    def __send_message_payload(
        self, module: str, category: str, event: str, msg: str
    ) -> dict[str, str | bool]:
        text = f"""<b>Tracker Error</b>

<b>🧩 Module:</b> {module}

<b>🗃 Category:</b> {category}

<b>⚠ Event:</b> {event}

<b>💬 Message:</b> {msg}"""

        return {
            "text": text,
            "parse_mode": "HTML",
            "chat_id": self.chat_id,
            "disable_web_page_preview": True,
        }

    def send_message(self, module: str, category: str, event: str, msg: str) -> dict:
        payload = self.__send_message_payload(module, category, event, msg)
        response = requests.post(
            f"https://api.telegram.org/bot{self.bot_id}/sendMessage", json=payload
        )
        return response.json()
