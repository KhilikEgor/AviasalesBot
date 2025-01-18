package main

import (
	"flag"
	"log"

	"github.com/KhilikEgor/AviasalesBot/internal/db"
	"github.com/KhilikEgor/AviasalesBot/internal/domain"
	"github.com/KhilikEgor/AviasalesBot/internal/handlers"
	"github.com/KhilikEgor/AviasalesBot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	startCommand            = "/start"
	allVacanciesCommand     = "Все вакансии"
	offNotificationsCommand = "Отключить уведомления"
)

var (
	BotToken = flag.String("tg.token", "", "token for telegram")
)

func startAviasalesBot() error {
	flag.Parse()

	db.Connect()
	err := db.DB.AutoMigrate(&domain.User{}, domain.Vacancy{})
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(*BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	vacancyService := &service.VacancyService{
		DB: db.DB,
	}

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Chat == nil {
			continue
		}

		txt := update.Message.Text

		user := domain.User{
			ChatId:   update.Message.Chat.ID,
			UserName: update.Message.Chat.UserName,
		}

		handlers.StartVacancyChecker(bot, vacancyService)

		switch txt {
		case startCommand:
			handlers.WelcomeMessageHandler(bot, vacancyService, user)
			continue
		case allVacanciesCommand:
			handlers.GetAllVacancyHandler(bot, vacancyService, user)
			continue
		case offNotificationsCommand:
			handlers.OffUserNotifications(bot, vacancyService, user)
			continue
		default:
			handlers.DefaultMessagesHandler(bot, user)
			continue
		}
	}
	return nil
}

func main() {
	err := startAviasalesBot()
	if err != nil {
		log.Fatalf("Error starting echo bot: %s", err)
	}
}
