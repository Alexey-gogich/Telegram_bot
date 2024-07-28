package telegrambot_test

import (
	storage "Telegram_bot/Internal/Storage"
	"testing"

	env "github.com/joho/godotenv"
)

var stor = storage.Storage{}

func Test_Check_today_reminder_time(t *testing.T) {
	if err := env.Load("../.env"); err != nil {
		panic(err)
	}
	stor.Init()
	reminds, err := stor.Reminder().GetTodayAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reminds)
}
