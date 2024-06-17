import os
from groq import Groq
import argparse
from dotenv import load_dotenv
from pathlib import Path

env_path = Path('../..') / '.env'
load_dotenv(dotenv_path=env_path)

client = Groq(
    api_key=os.getenv("GROQ_API_KEY"),
)

def main():

    parser = argparse.ArgumentParser(description='Пример скрипта для обработки аргументов командной строки.')

    parser.add_argument('input', type=str, help='Входная строка для обработки')

    args = parser.parse_args()

    chat_completion = client.chat.completions.create(
        messages=[
            {
                "role": "user",
                "content": "You can only answer in Russian! And imagine that you are a bot named Arkady, who works in the Telegram messenger" + args.input,
            }
        ],
        model="llama3-8b-8192",
    )

    print(chat_completion.choices[0].message.content)

if __name__ == '__main__':
    main()
