package service

import (
	"cmd/app/bot.go/internal/domain"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type VacancyService struct {
	Vacancies []domain.Vacancy
}

func (vs *VacancyService) ParsPage() {
	res, err := http.Get("https://www.aviasales.ru/about/vacancies")
	if err != nil{
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Failed to fetch page: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil{
		log.Fatal(err)
	}

	doc.Find("a.vacancies_vacancy").Each(func(i int, s *goquery.Selection) {
		name := s.Find("p.vacancies_vacancy__name").Text()

		description := s.Find("div.team").Text()

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

		vs.Vacancies = append(vs.Vacancies, vacancy)
	})
}
