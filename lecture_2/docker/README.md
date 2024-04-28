# Сборка приложения с помощью Docker

## Сборка
1. `docker image build --tag main --file Dockerfile.1 .` - собираем образ используя Dockerfile.1
1. `docker images main` - просмотр информации о нашем образе

## Запуск 
1. `docker container run --rm --name main-container -p 8080:8080 main:latest` - запуск образа в контейнере
1. `docker ps` - просмотр запущенных контейнеров


## Литература
1. https://docs.docker.com/language/golang/build-images/
1. [framework echo](https://echo.labstack.com/)
1. [CMD vs ENTRYPOINT](https://habr.com/ru/companies/slurm/articles/329138/)
1. [multi-stage builds](https://docs.docker.com/build/building/multi-stage/)