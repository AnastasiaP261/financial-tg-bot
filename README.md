# Телеграм бот "Финансовый помощник"

## Команды

- **/category <название категории>** - добавление новой категории для пользователя  

- **/add <сумма>** - добавляет новую трату без категории, в качестве даты берет текущую  

- **/add <сумма> <категория>** - добавляет новую трату в категорию, в качестве даты берет текущую

- **/add <сумма> <категория> <dd.mm.yyyy>** - добавляет новую трату в категорию и выставляет соответствующую дату

- **/report <week|month|year>** - собирает отчет за указанный промежуток. Отчет отправляет в виде текста и круговой диаграммы. Понимает разное количество дней в месяцах и високосные годы.

## Архитектура

```
.
├── bin
├── cmd
├── go.mod
├── go.sum
├── internal
│        ├── clients
│        │        └── tg                    - клиент телеграма
│        ├── config
│        └── model
│            ├── chart_drawing              - модель рисовальщика, здесь лежит логика по рисованию диаграмм для отчетов
│            ├── messages                   - выполняет функции контроллера и отлавливает команды
│            ├── normalize                  - требуется для нормализации входящих от пользователя данных
│            ├── purchases                  - основная бизнес-логика, здесь описана логика добавления трат, категорий и составления отчетов
│            └── store                      - хранилище (будет заменено на БД)
├── Makefile
├── README.md
└── secrets                                 - секреты приложения
```