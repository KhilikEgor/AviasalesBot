package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KhilikEgor/AviasalesBot/internal/domain"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
)

type VacancyService struct {
	DB        *gorm.DB
	Vacancies []domain.Vacancy
}

func (vs *VacancyService) ParsPage() ([]domain.Vacancy, error) {
	var newVacancies []domain.Vacancy

	// HTTP запрос
	res, err := http.Get("https://www.aviasales.ru/about/vacancies")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vacancies page: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Парсинг HTML документа
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML document: %v", err)
	}

	// Сканируем вакансии на странице
	doc.Find("a.vacancies_vacancy").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find("p.vacancies_vacancy__name").Text())
		description := strings.TrimSpace(s.Find("div.team").Text())
		link, exists := s.Attr("href")
		if !exists {
			log.Printf("Vacancy %d: no link found", i)
			return
		}

		// Создаём структуру вакансии
		vacancy := domain.Vacancy{
			Name:        name,
			Description: description,
			Link:        "https://www.aviasales.ru" + link,
			Active:      true,
			PublishDate: time.Now(),
		}

		// Проверяем наличие вакансии в базе
		var existingVacancy domain.Vacancy
		result := vs.DB.Where("link = ?", vacancy.Link).First(&existingVacancy)

		if result.Error == nil {
			// Обновляем существующую вакансию
			existingVacancy.Name = vacancy.Name
			existingVacancy.Description = vacancy.Description
			existingVacancy.Active = true
			existingVacancy.PublishDate = vacancy.PublishDate
			vs.DB.Save(&existingVacancy)
		} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Добавляем новую вакансию
			vs.DB.Create(&vacancy)
			newVacancies = append(newVacancies, vacancy)
		} else {
			log.Printf("Error querying database: %v", result.Error)
		}
	})

	return newVacancies, nil
}



func (vs *VacancyService) GetNewVacancies(newVacancies []domain.Vacancy) []domain.Vacancy {
	var diff []domain.Vacancy

	existing := make(map[string]struct{})
	for _, v := range vs.Vacancies {
		existing[v.Link] = struct{}{}
	}

	for _, nv := range newVacancies {
		if _, found := existing[nv.Link]; !found {
			diff = append(diff, nv)
		}
	}

	log.Printf("Found %d new vacancies", len(diff))
	return diff
}

func (vs *VacancyService) UpdateVacancies(newVacancies []domain.Vacancy) {
	log.Printf("Updating vacancies. Old count: %d, New count: %d", len(vs.Vacancies), len(newVacancies))
	vs.Vacancies = newVacancies
} 

func (vs *VacancyService) SaveUser(user domain.User) error {
	var existingUser domain.User

	result := vs.DB.First(&existingUser, user.ChatId)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if existingUser.ChatId == 0 {
		user.Notification = true
		if err := vs.DB.Create(&user).Error; err != nil {
			return err
		}
	} else {
		existingUser.UserName = user.UserName
		existingUser.Notification = true

		if err := vs.DB.Model(&existingUser).Updates(domain.User{
			UserName:    user.UserName,
			Notification: true,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (vs *VacancyService) OffNotifications(user domain.User) error {
	var existingUser domain.User

	result := vs.DB.First(&existingUser, user.ChatId)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return result.Error
	}

	existingUser.Notification = false
	if err := vs.DB.Save(&existingUser).Error; err != nil {
		return err
	}
	return nil
}

func (vs *VacancyService) GetAllUsers() ([]domain.User, error) {
    var users []domain.User
    if err := vs.DB.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (vs *VacancyService) GetAllVacancies() ([]domain.Vacancy, error){
	var vacancies []domain.Vacancy
	if err := vs.DB.Find(&vacancies).Error; err != nil {
		return nil, err
	}
	return vacancies, nil
}