package handlers

import (
	"cmd/app/bot.go/internal/domain"
	"cmd/app/bot.go/internal/service"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func VacancyHandler(bot *tgbotapi.BotAPI, vs *service.VacancyService, request domain.User) {
	// Вызываем метод парсинга вакансий
	vs.ParsPage()

	// Проверяем, есть ли вакансии
	if len(vs.Vacancies) == 0 {
		responseMessage := "К сожалению, вакансии не найдены."
		msg := tgbotapi.NewMessage(request.ChatId, responseMessage)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	// Формируем ответное сообщение
	responseMessage := "Актуальные вакансии:\n\n"
	for _, vacancy := range vs.Vacancies {
		responseMessage += fmt.Sprintf(
			"Название: %s\nОписание: %s\nСсылка: %s\n\n",
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
