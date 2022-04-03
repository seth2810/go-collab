# Go collab

Разработать веб страницу с многострочным текстовом полем с функцией совместного редактирования.
Аналог code.yandex-team.ru или google docs

- Стек: Golang + websocket
- Ввод текста не должен мешать вводу текста другим пользователям
- На странице должен быть индикатор количества подключенных пользователей

## Системные требования

* Yarn
* Node.js
* Golang 1.18
* [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli#download-and-install)

## Установка

```sh
make install
make start # run server http://localhost:8080
```

## Docker

```sh
docker-compose up -d # run server http://localhost:8080
```
