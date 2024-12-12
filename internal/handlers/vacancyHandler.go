package handlers

import (
	"cmd/app/bot.go/internal/domain"
	"cmd/app/bot.go/internal/service"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetAllVacancyHandler(bot *tgbotapi.BotAPI, vs *service.VacancyService, request domain.User) {
	// Вызываем метод парсинга вакансий
	allVacancy := vs.ParsPage()

	// Проверяем, есть ли вакансии
	if len(allVacancy) == 0 {
		responseMessage := "К сожалению, вакансии не найдены."
		msg := tgbotapi.NewMessage(request.ChatId, responseMessage)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	// Формируем ответное сообщение
	responseMessage := "Горячие вакансии 🔥🔥🔥\n\n"
	for _, vacancy := range allVacancy {
		responseMessage += fmt.Sprintf(
			"%s\n%s\n%s\n\n",
			vacancy.Name, vacancy.Description, vacancy.Link,
		)
	}

	// Убедимся, что сообщение не превышает лимита Telegram (4096 символов)
	if len(responseMessage) > 4096 {
		responseMessage = responseMessage[:4093] + "..." // Урезаем сообщение
	}

	// Отправляем сообщение пользователю
	msg := tgbotapi.NewMessage(request.ChatId, responseMessage)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func WelcomeMessageHandler(bot *tgbotapi.BotAPI, vs *service.VacancyService, request domain.User) {
	vs.ParsPage()

	if len(vs.Vacancies) == 0 {
		responseMessage := "К сожалению, вакансии не найдены."
		msg := tgbotapi.NewMessage(request.ChatId, responseMessage)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	responseMessage := "Отлично! Теперь как появится новая 🔥ГОРЯЧАЯ вакансия, ты узнаешь один из первых\n\nА пока можешь отдыхать я сделаю все сам!"

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Все вакансии"),
			tgbotapi.NewKeyboardButton("Создать подписку"),
		),
	)

	msg := tgbotapi.NewMessage(request.ChatId, responseMessage)
	msg.ReplyMarkup = replyKeyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func StartVacancyChecker(bot *tgbotapi.BotAPI, vs *service.VacancyService, request domain.User) {
	initialVacancies := vs.ParsPage()
	vs.UpdateVacancies(initialVacancies)

	go func() {
		for {
			newVacancies := vs.ParsPage()

			diff := vs.GetNewVacancies(newVacancies)

			if len(diff) > 0 {
				// Уведомляем пользователя о новых вакансиях
				for _, vacancy := range diff {
					message := fmt.Sprintf(
						"🔥 Новая вакансия!\n\n%s\n%s\n%s\n",
						vacancy.Name, vacancy.Description, vacancy.Link,
					)
					msg := tgbotapi.NewMessage(request.ChatId, message)
					_, err := bot.Send(msg)
					if err != nil {
						log.Printf("Failed to send message: %v", err)
					}
				}

				vs.UpdateVacancies(newVacancies)
			}

			time.Sleep(600 * time.Second)
		}
	}()
}

func DefaultMessagesHandler(bot *tgbotapi.BotAPI, request domain.User) {
	sticker := tgbotapi.NewSticker(request.ChatId, tgbotapi.FileID("CAACAgIAAxkBAAENQDBnS2_zZGpxdw7SwmUrGzDLcmNofwACw0IAAtAnyEqlQ3xNhpVNmTYE"))
	_, err := bot.Send(sticker)
	if err != nil {
		log.Printf("Failed to send sticker: %s", err)
	}

	message := "Сейчас бот это не умеет делать. Можешь отправить свои предложения мне в личку @khilikegor"
	msg := tgbotapi.NewMessage(request.ChatId, message)
	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}