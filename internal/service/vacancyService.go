package service

import (
	"cmd/app/bot.go/internal/domain"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type VacancyService struct {
	Vacancies []domain.Vacancy
}

func (vs *VacancyService) ParsPage() []domain.Vacancy {
	var newVacancies []domain.Vacancy

	res, err := http.Get("https://www.aviasales.ru/about/vacancies")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch page: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(`Failed to parse HTML document: %v`, err)
	}

	doc.Find("a.vacancies_vacancy").Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Find("p.vacancies_vacancy__name").Text())
		description := strings.TrimSpace(s.Find("div.team").Text())

		link, exists := s.Attr("href")
		if !exists {
			log.Printf("Vacancy %d: no link found", i)
			return
		}

		vacancy := domain.Vacancy{
			Name:        name,
			Description: description,
			Link:        "https://www.aviasales.ru" + link,
		}

		newVacancies = append(newVacancies, vacancy)
	})

	return newVacancies
}

//For local debug
// func (vs *VacancyService) ParsPage() []domain.Vacancy {
// 	var newVacancies []domain.Vacancy

// 	// Проверяем и выводим текущую директорию
// 	dir, err := os.Getwd()
// 	if err != nil {
// 		log.Fatalf("Failed to get current working directory: %v", err)
// 	}
// 	log.Println("Current working directory:", dir)

// 	// Открываем локальный файл HTML
// 	file, err := os.Open("../../vacancies.html")
// 	if err != nil {
// 		log.Fatalf("Failed to open file %s: %v", err)
// 	}
// 	defer file.Close()

// 	// Парсим документ из файла
// 	doc, err := goquery.NewDocumentFromReader(file)
// 	if err != nil {
// 		log.Fatalf("Failed to parse HTML document: %v", err)
// 	}

// 	// Извлечение вакансий
// 	doc.Find("a.vacancies_vacancy").Each(func(i int, s *goquery.Selection) {
// 		name := strings.TrimSpace(s.Find("p.vacancies_vacancy__name").Text())
// 		description := strings.TrimSpace(s.Find("div.team").Text())

// 		link, exists := s.Attr("href")
// 		if !exists {
// 			log.Printf("Vacancy %d: no link found", i)
// 			return
// 		}

// 		vacancy := domain.Vacancy{
// 			Name:        name,
// 			Description: description,
// 			Link:        "https://www.aviasales.ru" + link,
// 		}

// 		newVacancies = append(newVacancies, vacancy)
// 	})

// 	// Логируем извлеченные вакансии
// 	log.Printf("Parsed %d vacancies from %s", len(newVacancies))

// 	return newVacancies
// }

func (vs *VacancyService) GetNewVacancies(newVacancies []domain.Vacancy) []domain.Vacancy {
	var diff []domain.Vacancy

	// Создаем карту существующих вакансий по их ссылкам
	existing := make(map[string]struct{})
	for _, v := range vs.Vacancies {
		existing[v.Link] = struct{}{}
	}

	// Ищем вакансии, которых еще нет в текущем списке
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
