package storage

import (
	models "Telegram_bot/Internal/Models"
	"time"
)

type Reminders struct {
	storage *Storage
}

func (Rem *Reminders) GetTodayAll() ([]models.Reminder, error) {
	reminds_mass := []models.Reminder{}
	rows, err := Rem.storage.db.Query("select * from reminders where time >= $1 and time <=$2", time.Now().Format("2006-01-02"), time.Now().AddDate(0, 0, 1).Format("2006-01-02"))
	if err != nil {
		return reminds_mass, err
	}
	for rows.Next() {
		prod := models.Reminder{}
		rows.Scan(&prod.Id, &prod.Chat_id, &prod.Text, &prod.Time)
		reminds_mass = append(reminds_mass, prod)
	}
	return reminds_mass, nil
}

func (Rem *Reminders) GetUserAll(chat_id int64) ([]models.Reminder, error) {
	reminds_mass := []models.Reminder{}
	rows, err := Rem.storage.db.Query("select * from reminders where chat_id = $1", chat_id)
	if err != nil {
		return reminds_mass, err
	}
	for rows.Next() {
		prod := models.Reminder{}
		rows.Scan(&prod.Id, &prod.Chat_id, &prod.Text, &prod.Time)
		reminds_mass = append(reminds_mass, prod)
	}
	return reminds_mass, nil
}

func (Rem *Reminders) Create(reminder *models.Reminder) error {
	_, err := Rem.storage.db.Exec("insert into reminders (chat_id, text, time) values ($1, $2, $3)", reminder.Chat_id, reminder.Text, reminder.Time)
	if err != nil {
		return err
	}
	return nil
}

func (Rem *Reminders) Delete(id int) error {
	_, err := Rem.storage.db.Exec("delete from reminders where id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
