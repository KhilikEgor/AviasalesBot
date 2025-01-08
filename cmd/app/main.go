package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

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
	BotToken      = flag.String("tg.token", "", "token for telegram")
	WebHookServer = flag.String("tg.webhook", "", "web hook server")
)


func startAviasalesBot() error {
	flag.Parse()

	db.Connect()
	err := db.DB.AutoMigrate(&domain.User{})
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(*BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	webhookConfig, err := tgbotapi.NewWebhook(*WebHookServer + "/" + *BotToken)
	if err != nil {
		log.Panicf("Error creating webhook: %s", err)
	}

	_, err = bot.Request(webhookConfig)
	if err != nil {
		log.Panicf("Error setting webhook: %s", err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	http.HandleFunc("/"+*BotToken, func(w http.ResponseWriter, r *http.Request) {
		update := tgbotapi.Update{}
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Failed to parse update", http.StatusBadRequest)
			log.Println("Failed to parse update:", err)
			return
		}
		handleUpdate(update, bot, w)
	})

	// Запуск веб-сервера
	go func() {
		log.Println("Starting webhook server on port 8081...")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()

	log.Println("Bot is now running...")
	select {}

	return nil
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, w http.ResponseWriter) {
    log.Printf("Received update: %+v", update)
    if update.Message == nil || update.Message.Chat == nil {
        log.Println("Invalid update received, skipping")
        http.Error(w, "Invalid update received", http.StatusBadRequest)
        return
    }
    
    txt := update.Message.Text
    log.Printf("Message text: %s", txt)

	user := domain.User{
		ChatId:   update.Message.Chat.ID,
		UserName: update.Message.Chat.UserName,
	}

	vacancyService := &service.VacancyService{
		DB: db.DB,
	}

	handlers.StartVacancyChecker(bot, vacancyService)

	switch txt {
	case startCommand:
		handlers.WelcomeMessageHandler(bot, vacancyService, user)
	case allVacanciesCommand:
		handlers.GetAllVacancyHandler(bot, vacancyService, user)
	case offNotificationsCommand:
		handlers.OffUserNotifications(bot, vacancyService, user)
	default:
		handlers.DefaultMessagesHandler(bot, user)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	err := startAviasalesBot()
	if err != nil {
		log.Fatalf("Error starting bot: %s", err)
	}
}
