# Телеграм бот "Финансовый помощник"

## Запуск

`docker-compose up -d pg`, а затем `make run && make dev-db-data` - запустит бота локально. Вся остальная инфраструктура (например БД) будет запущена в контейнерах вокруг.   
`make docker-run` - запустит приложение полностью в контейнерах.

## Команды

- **/category <название категории>** - добавление новой категории для пользователя  

- **/add <сумма>** - добавляет новую трату без категории, в качестве даты берет текущую  

- **/add <сумма> <категория>** - добавляет новую трату в категорию, в качестве даты берет текущую

- **/add <сумма> <категория> <dd.mm.yyyy>** - добавляет новую трату в категорию и выставляет соответствующую дату

- **/report <week|month|year>** - собирает отчет за указанный промежуток. Отчет отправляет в виде текста и круговой диаграммы. Понимает разное количество дней в месяцах и високосные годы.

- **/currency <RUB|USD|EUR|CNY>** - сменить основную валюту пользователя. После этой команды отчеты и добавления трат будут в этой валюте. По умолчанию у каждого пользователя установлена RUB.  
- 
- **/limit <сумма>** - установить лимит на траты в календарный месяц. Для снятия лимита отправить -1.

## Архитектура

```
.
├── bin
├── cmd
├── go.mod
├── go.sum
├── internal
│        ├── clients
│        │        ├── tg                    - клиент телеграма
│        │        └── fixer                 - клиент fixer - апи для получения курсов валют
│        ├── config
│        ├── model
│        │        ├── chart_drawing         - модель рисовальщика, здесь лежит логика по рисованию диаграмм для отчетов
│        │        ├── db                    - база данных
│        │        ├── exchange-rates        - модель курсов валют, оборачивает склиент fixer в необходимую нам бизнес-логику
│        │        ├── messages              - выполняет функции контроллера и отлавливает команды
│        │        ├── normalize             - требуется для нормализации входящих от пользователя данных
│        │        └── purchases             - основная бизнес-логика, здесь описана логика добавления трат, категорий и составления отчетов
│        ├── test_data                      - тестовые данные, в том числе фикстуры для интеграционных тестов
│        └── migrations                     - миграции для базы данных
├── Makefile
├── README.md
└── config                                 - секреты приложения
```

![Описание архитектуры](./images/arch.jpg)  
[Ссылка на доску miro](https://miro.com/app/board/uXjVPJQpCgA=/)  

## Масштабирование  

![Масштабирование приложения на 1000 пользователей](./images/1000_users.jpg)  

![Масштабирование приложения на 1000 пользователей](./images/100_000_users.jpg)  

![Масштабирование приложения на 1000 пользователей](./images/1_000_000_users.jpg)
[Ссылка на доску miro](https://miro.com/app/board/uXjVPJQpCgA=/)
