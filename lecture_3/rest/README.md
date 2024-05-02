# REST

## Генерация моделей на примере go-swagger
1. `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`
1. `mkdir -p ./server`
1. `swagger generate model -f ./api/swagger.json -t ./server/`
1. `go mod tidy`

## Полезные ссылки:
* [swagger editor](https://editor.swagger.io/)
* [swagger-api/swagger-codegen](https://github.com/swagger-api/swagger-codegen)
* [go-swagger](https://github.com/go-swagger/go-swagger)
* [swaggo/swag](https://github.com/swaggo/swag)
* [Стандартный макет Go проекта](https://github.com/golang-standards/project-layout/blob/master/README_ru.md)
* [go HTTP сервер](https://pkg.go.dev/net/http#hdr-Servers)
* [echo framework](https://echo-labstack-com.translate.goog/docs/quick-start?_x_tr_sl=en&_x_tr_tl=ru&_x_tr_hl=ru&_x_tr_pto=sc&_x_tr_hist=true)
