# StudyBuddyBot

StudyBuddyBot is a Telegram bot written in Golang using the pure Telegram Bot API. It is designed to facilitate homework tracking, schedule retrieval, mini-games, and other useful features within Telegram group chats.

## Commands

### Mini-Game "Cultivating Your Farm":

- `/dick` - Cultivate your virtual farm.
- `/top_dick` - Check who has the most impressive "farm" in the chat.
- `/duel [@username]` - Challenge or accept a duel with a chat member.
- `/hp` - Check and replenish your life points.

### Mini-Game "Gay of the Day" (Admin-exclusive functionality):

- `/gay` - Discover whose day is particularly lucky today.
- `/top_gay` - View the top participants based on their "luckiness."

### Auction Mini-Game:

- `/start_auction {time}` - Initiate an auction in the chat. Time is an optional parameter for the auction duration.
- `/deposit [value]` - Deposit units from your "farm" into the auction.
- `/auction` - See the current participants in the auction.

### Educational Commands:

- `/add_calendar [google_calendar_id]` - Link a calendar to the chat.
- `/schedule` - View the schedule based on the linked calendar.
- `/add` - Add homework to the list.
- `/get {number} {subject}` - Retrieve the latest homework entries.
- `/delete [id]` - Remove a homework entry from the list.

### Miscellaneous Bot Commands:

- `/my_stats` - Check your personal statistics in the chat.
- `/chat_stats` - View overall chat statistics.

## Installation and Execution

1. Clone the repository: `git clone https://github.com/arevbond/StudyBuddyBot`
2. Set up the following environment variables:
    - `TELEGRAM_TOKEN`: Your Telegram bot token.
    - `CONFIG_PATH`: Path to your bot configuration file.
    - `ADMINS_ID`: IDs of administrators who can access admin commands.
3. Run the bot: `go run main.go`

### Important Notes:

- To use the bot, create a Telegram bot and obtain its token via BotFather.
- For Google Calendar integration, create a project on Google Cloud Platform and obtain an API key.
