package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	models "Telegram_bot/Internal/Models"
	stor "Telegram_bot/Internal/Storage"

	telegram_bot "github.com/go-telegram-bot-api/telegram-bot-api"
	env "github.com/joho/godotenv"
)

var Bot *telegram_bot.BotAPI
var Storage *stor.Storage
var Reminder_read_chanel = make(chan models.Reminder, 8)
var Users_messages = make(map[int64]chan string)
var Map_Mutex = sync.RWMutex{}

func main() {
	err := env.Load()
	if err != nil {
		panic(err)
	}

	Bot, _ = telegram_bot.NewBotAPI(os.Getenv("BOT_KEY"))
	if Bot == nil {
		panic("Вот не инициализирван.")
	}

	Storage, err = stor.NewStorage()
	if err != nil {
		log.Println("Storage initialize error:")
		log.Fatal(err)
	}

	fmt.Printf("Authorized on account %s\n", Bot.Self.UserName)

	// go Check_today_reminder_time()
	// go Check_reminder_time()

	u := telegram_bot.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			if Users_messages[update.Message.Chat.ID] != nil {
				Users_messages[update.Message.Chat.ID] <- update.Message.Text
			} else {
				go Main_bot_thread(update.Message.Chat.ID)
				Users_messages[update.Message.Chat.ID] = make(chan string, 1)
				Users_messages[update.Message.Chat.ID] <- update.Message.Text
			}
		}
	}
}
