# graphql


1. `go get github.com/99designs/gqlgen`
1. `go run github.com/99designs/gqlgen init`
1. `mkdir -p ./cmd/server && mv server.go ./cmd/server`
1. Изменить/создать схему `graph/*.graphql`
1. Отредактироварать `gqlgen.yml` в разделе `schema`
1. `go run github.com/99designs/gqlgen generate`
1. `go build ./cmd/server`
1. `./server`


## Полезные материалы
* https://www.howtographql.com/graphql-go/1-getting-started/
* https://gqlgen.com/getting-started/
* https://habr.com/ru/companies/ruvds/articles/444346/