# Telegram бот ICS
## Статус

В разработке

## Команды бота

| Команда                           | Описание                                           |
|-----------------------------------|----------------------------------------------------|
| `/xkcd`                           | случайная картина из [xkcd.com](https://xkcd.com/) |
| `/today`, `/tomorrow`, `/lessons` | узнать расписание на `сегодня`, `завтра`, `неделю` |
| `/dick`, `/top_dick`              | игра: по выращиванию своего хозяйства              |
| `/gay`, `/top_gay`                | игра: узнать кто gay дня                           |

## Tasks:
- Расписание через файл
  - Реализовать новую таблицу в БД для хранения ID группы - расписания
  - Добавить hhtp-сервер с post ручкой с аунтентификацией и добавлением json`а 
