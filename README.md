# project


## Тесты

Для запуска тестов использовать
```cmd
go test ./...
```

Для запуска интеграционных тестов
```cmd
go test -tags=integration ./...
```

Для создания покрытия unit тестов
```cmd
go test -race --coverprofile=coverage.out -covermode=atomic ./...
```

Для отображения покрытия unit тестов
```cmd
go tool cover --html=coverage.out
```

Для создания покрытия интеграционных тестов
```cmd
go test -tags=integration -race --coverprofile=coverage-integration.out -covermode=atomic ./internal/storage/postgres/repository/...
```

Для отображения покрытия интеграционных тестов
```cmd
go tool cover --html=coverage-integration.out
```

При необходимости можно добавлять или убирать -race и соответственно -covermode=atomic в тестах