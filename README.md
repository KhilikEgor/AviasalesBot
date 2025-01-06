# Aviasales Vacancies Bot

https://t.me/AviasalesVacanciesBot

**AviasalesVacanciesBot** - данный Telegram-бот я создал для упрощения поиска и отслеживания вакансий в Aviasales (прям как для авиабилетов)

## Технологический стек
* Язык программирования: Go (Golang)
* База данных: PostgreSQL
* ORM: GORM
* Контейнеризация: Docker
* Платформа: Telegram Bot API

## E2E
Пользователь активирует бота командой /start, после чего начинает получать  уведомления о новых вакансиях. Пользователь может в любой момент запросить список всех текущих вакансий или отключить уведомления, если они больше не нужны.

## Развертывание Бота
Бот развернут и функционирует на моем личном сервере. Процесс развертывания автоматизирован с использованием Docker. Образ бота лежит в общем доступе в [Dockerhub](https://hub.docker.com/r/khilik/server-bot-aviasales/tags)

## TODO
- [x] Провести рефакторинг кода
- [x] Сделать автоматизацию CI/CD
- [ ] Добавить БД для сущности вакансии
- [ ] Начать писать тесты, сделать покрытие хотя бы 50-60%
- [ ] Подумать 🤔. В случе возникновения ошибки сделать автоподъем сервера, типа healthcheck. Возможно можно это сделать в докере.