from datetime import datetime

import requests


class Telegram:
    def __init__(
        self, bot_id: str, chat_id: str, disable_notification: bool = True
    ) -> None:
        self.chat_id = chat_id
        self.bot_id = bot_id
        self.disable_notification = disable_notification

    def __send_message_payload(
        self, module: str, category: str, event: str, msg: str
    ) -> dict[str, str | bool]:
        text = f"""<pre><code class="language-Tracker Error">
📅 Date: {datetime.now().strftime("%d/%m/%Y %H:%M:%S")}

🧩 Module: {module}

🗃 Category: {category}

💬 Message: {msg}

⚠ Event: {event}</code></pre>"""

        return {
            "text": text,
            "parse_mode": "HTML",
            "chat_id": self.chat_id,
            "disable_web_page_preview": True,
            "disable_notification": self.disable_notification,
        }

    def send_message(self, module: str, category: str, event: str, msg: str) -> dict:
        payload = self.__send_message_payload(module, category, event, msg)
        response = requests.post(
            f"https://api.telegram.org/bot{self.bot_id}/sendMessage", json=payload
        )
        return response.json()
