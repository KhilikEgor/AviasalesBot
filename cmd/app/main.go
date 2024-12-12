package main

import (
	"cmd/app/bot.go/internal/domain"
	"cmd/app/bot.go/internal/handlers"
	"cmd/app/bot.go/internal/service"
	"flag"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	BotToken = flag.String("tg.token", "", "token for telegram")
)

func startAviasalesBot() error {
	flag.Parse()

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(*BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	vacancyService := &service.VacancyService{}

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Chat == nil {
			continue
		}

		txt := update.Message.Text
		var responseText string

		user := domain.User{
			ChatId: update.Message.Chat.ID,
		}

		handlers.StartVacancyChecker(bot, vacancyService, user)

		switch txt {
		case "/start":
			handlers.WelcomeMessageHandler(bot, vacancyService, user)
			continue
		case "Все вакансии":
			handlers.GetAllVacancyHandler(bot, vacancyService, user)
			continue
		case "Создать подписку":
			fmt.Println("Hello World")
		default:
			handlers.DefaultMessagesHandler(bot, user)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Failed to send message: %s", err)
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	}

	return nil
}

func main() {
	err := startAviasalesBot()
	if err != nil {
		log.Fatalf("Error starting echo bot: %s", err)
	}
}
