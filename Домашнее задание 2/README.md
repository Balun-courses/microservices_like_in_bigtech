# Домашнее задание 2


## Задание
1. Создать репозитории (или моно репозиторий) для сервисов из [Домашнего задания 1](../%D0%94%D0%BE%D0%BC%D0%B0%D1%88%D0%BD%D0%B5%D0%B5%20%D0%B7%D0%B0%D0%B4%D0%B0%D0%BD%D0%B8%D0%B5%201/README.md).
2. Для каждого сервиса создать `main` с http-сервером c _liveness/readiness probe_ (можно использовать как стандартную библиотеку golang `net/http` так и любой понравившийся вам фреймворк)
3. Написать `Dockerfile` для каждого сервиса (естественно это все должно собираться).
4. Написать инструкцию или скрипт для того, чтобы можно было поднять все сервисы в контейнерах локально. (_Подсказка_: для удобства локальной разработки лучше всего воспользоваться [docker-compose](https://docs.docker.com/compose/) и _Makefile_)
5. ⭐ Реализовать стратегии деплоя `blue-green` и `canary` с помощью стандартных средств kubernetes.

## Полезные ссылки
* [Стандартный макет Go проекта](https://github.com/golang-standards/project-layout/blob/master/README_ru.md)
* [go HTTP сервер](https://pkg.go.dev/net/http#hdr-Servers)
* [Руководство по Docker Compose для начинающих](https://habr.com/ru/companies/ruvds/articles/450312/)