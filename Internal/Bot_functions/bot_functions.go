package Bot_functions

import (
	models "Telegram_bot/Internal/Models"
	"fmt"
	reg "regexp"
	"strings"
	"time"

	telegram_bot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Main_bot_thread(chat_id int64) { //Проверить ожидание
	timer := time.Now().Minute() - 45
	for timer >= time.Now().Minute()-60 {
		message := ""
		Map_Mutex.Lock()
		if len(Users_messages[chat_id]) > 0 {
			message = <-Users_messages[chat_id]
		}
		Map_Mutex.Unlock()
		if message != "" {
			message_text := ""
			switch message {
			case "/start":
				message_text = "Приветствую, меня зовут Виктор, чем могу помочь? Для просмотра функций используйте: /help"
				msg := telegram_bot.NewMessage(chat_id, message_text)
				Bot.Send(msg)
			case "/help":
				message_text = "/create_reminder - создание текстового напоминания."
				msg := telegram_bot.NewMessage(chat_id, message_text)
				Bot.Send(msg)
			case "/create_reminder":
				message_text = "Укажите время и текст напоминания в формате:\n00:00, 01.01.06 - текст напоминалки\nИли /exit для выхода"
				msg := telegram_bot.NewMessage(chat_id, message_text)
				Bot.Send(msg)
				Create_Reminder(chat_id, timer)
			default:
				message_text = "Я не знаю такой команды, пожалуйста ознакомтесь со списком возможных команд."
				msg := telegram_bot.NewMessage(chat_id, message_text)
				Bot.Send(msg)
			}
		}
	}
	Users_messages[chat_id] = nil
}

func Create_Reminder(chat_id int64, timer int) {
	message_check := reg.MustCompile(`[0-9]{1,2}:[0-9]{2},\s?[0-9]{2}.[0-9]{2}.[0-9]{2,4}\s?-\s?.+`)
	var mg telegram_bot.MessageConfig //???
	for timer >= time.Now().Minute()-60 {
		message := ""
		Map_Mutex.Lock()
		if len(Users_messages[chat_id]) > 0 {
			message = <-Users_messages[chat_id]
		}
		Map_Mutex.Unlock()
		if message != "" {
			if message == "/exit" {
				mg = telegram_bot.NewMessage(chat_id, "Вы вышли из меню создания уведомления")
				Bot.Send(mg)
				return
			}
			if message_check.MatchString(message) {
				reminder := models.Reminder{}
				pars_data := strings.Split(message, "-")
				reminder.Chat_id = chat_id
				reminder.Text = reg.MustCompile(`^\s`).ReplaceAllString(pars_data[1], "")

				if clear_date := strings.ReplaceAll(pars_data[0], " ", ""); len(clear_date) <= 14 {
					reminder.Time, _ = time.Parse("15:04,02.01.06", clear_date)
				} else if len(clear_date) >= 15 {
					reminder.Time, _ = time.Parse("15:04,02.01.2006", clear_date)
				}

				if reminder.Time.Year() == reminder.Time.Year() && reminder.Time.Month() == reminder.Time.Month() && reminder.Time.Day() == time.Now().Day() {
					Reminder_read_chanel <- reminder
				} else {
					err := Storage.Reminder().Create(&reminder)
					if err != nil {
						mg = telegram_bot.NewMessage(chat_id, "Ошибка сервера")
						Bot.Send(mg)
						fmt.Println(err)
						return
					}
				}
				mg = telegram_bot.NewMessage(chat_id, "Напоминание создано")
				Bot.Send(mg)
				return
			} else {
				mg = telegram_bot.NewMessage(chat_id, "Неверный формат даты/времени")
				Bot.Send(mg)
			}
		}
	}
}

func Check_today_reminder_time() { //Ежедневная докачка сегодняшних уведомлений в память
	for {
		reminds, err := Storage.Reminder().GetTodayAll()
		if err != nil {
			fmt.Println(err)
		}
		for _, remind := range reminds {
			Reminder_read_chanel <- remind
		}
		time.Sleep(24 * time.Hour)
	}
}

func Check_reminder_time() { //Проверка уведомлений раз в минуту
	daily_reminds := []models.Reminder{}
	time.Sleep(3 * time.Second)
	for {
		for len(Reminder_read_chanel) != 0 { // Загрузка в буффер сегодняшних напоминалок.
			daily_reminds = append(daily_reminds, <-Reminder_read_chanel)
		}

		for index := 0; index < len(daily_reminds); index++ {
			if daily_reminds[index].Time.Hour() == time.Now().Hour() && daily_reminds[index].Time.Minute() == time.Now().Minute() { //Несовпадение времени.
				msg := telegram_bot.NewMessage(daily_reminds[index].Chat_id, "Вы просили вам напомнить, что сегодня в "+daily_reminds[index].Time.Format("03:04")+" вы хотели - "+daily_reminds[index].Text)
				Bot.Send(msg)
				daily_reminds = append(daily_reminds[:index], daily_reminds[index+1:]...)
				index--
			}
		}
		time.Sleep(15 * time.Second)
	}
}
