# Планировщик задач (Task Scheduler)

![Go](https://img.shields.io/badge/Go-1.26.5-00ADD8?logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-modernc%2Fsqlite-003B57?logo=sqlite&logoColor=white)
![CI](https://github.com/bek63-dev/final-project-13-14-sprint/actions/workflows/tests.yml/badge.svg)

Веб-сервер планировщика задач — аналог TODO-листа с поддержкой периодических задач, поиска, JWT-аутентификации и готовым фронтендом. Задачи можно назначать на конкретную дату или на правило повторения (каждые N дней, по дням недели, по дням месяца или ежегодно), а выполненная периодическая задача автоматически переносится на следующую дату по своему правилу.

---
## Содержание

- [Возможности](#возможности)
- [Стек технологий](#стек-технологий)
- [Структура проекта](#структура-проекта)
- [Быстрый старт](#быстрый-старт)
- [Переменные окружения](#переменные-окружения)
- [API](#api)
- [Правила повторения задач](#правила-повторения-задач)
- [Аутентификация](#аутентификация)
- [Тестирование](#тестирование)
- [CI/CD](#cicd)
- [Docker](#docker)

---

## Возможности

- CRUD-операции над задачами: добавление, получение списка, получение одной задачи, редактирование, удаление.
- Отметка задачи выполненной: одноразовая задача удаляется, периодическая — переносится на следующую дату по правилу повторения.
- Полная поддержка правил повторения: `d` (интервал в днях), `y` (ежегодно), `w` (дни недели) и `m` (дни месяца, включая последний/предпоследний день и фильтр по месяцам).
- Поиск задач по подстроке в заголовке/комментарии или по конкретной дате (`GET /api/tasks?search=...`).
- Опциональная аутентификация по паролю с выдачей JWT-токена и middleware-проверкой защищённых маршрутов.
- Конфигурация хоста, порта, пути к БД, пароля и секретного ключа через `.env`/переменные окружения (`TODO_HOST`, `TODO_PORT`, `TODO_DBFILE`, `TODO_PASSWORD`, `TODO_SECRET_KEY`).
- Хранение данных в SQLite через чистый Go-драйвер `modernc.org/sqlite` — сборка без cgo.
- Готовый фронтенд (Svelte-сборка) с формой добавления/редактирования задач, календарём, поиском и страницей входа.
- Docker-сборка (multi-stage) со статическим бинарником и volume-хранением базы данных на хосте.
- CI-пайплайн на GitHub Actions, поднимающий сервер и прогоняющий тесты при каждом push.

---

## Стек технологий

| Категория      | Технология                                             |
|-----------------|--------------------------------------------------------|
| Язык             | Go 1.26.5                                              |
| HTTP             | стандартный `net/http` + `http.ServeMux`               |
| База данных      | SQLite (`modernc.org/sqlite`, без cgo)                 |
| Аутентификация   | JWT (`github.com/golang-jwt/jwt/v5`), SHA-256          |
| Конфигурация     | `.env` через `github.com/joho/godotenv`                |
| Фронтенд         | Svelte (собранный в `scripts.min.js`), Axios           |
| Тесты            | `go test`, `github.com/stretchr/testify`, `sqlx`       |
| CI/CD            | GitHub Actions                                         |
| Контейнеризация  | Docker (multi-stage build)                             |

---

## Структура проекта

```
final-project-13-14-sprint/
├── .github/
│   └── workflows/              # CI-пайплайн (tests.yml)
├── cmd/
│   └── main.go                 # точка входа: сборка зависимостей и запуск сервера
├── internal/
│   ├── api/                    # HTTP-обработчики и модель Task уровня API
│   │   ├── addtask.go          # POST /api/task
│   │   ├── auth.go             # POST /api/signin, middleware Auth
│   │   ├── donetask.go         # POST /api/task/done, DELETE /api/task
│   │   ├── edittask.go         # GET/PUT /api/task
│   │   ├── gettasks.go         # GET /api/tasks
│   │   ├── handlers.go         # Handler, маршрутизация /api/task по методу
│   │   ├── nextdate.go         # GET /api/nextdate
│   │   ├── response.go         # writeJSON / writeError
│   │   └── validate.go         # checkDate, isMethodAllowed
│   ├── auth/
│   │   └── jwt.go              # генерация и проверка JWT-токенов
│   ├── config/
│   │   └── config.go           # загрузка конфигурации из .env/окружения
│   ├── database/
│   │   ├── database.go         # подключение к SQLite, применение схемы
│   │   └── task.go             # CRUD-методы над таблицей scheduler
│   ├── domain/                 # бизнес-логика вычисления следующей даты
│   │   ├── date.go
│   │   ├── next_date.go
│   │   ├── repeat_rule.go
│   │   ├── rule_daily.go
│   │   ├── rule_monthly.go
│   │   ├── rule_weekly.go
│   │   └── rule_yearly.go
│   ├── router/
│   │   └── router.go           # регистрация маршрутов и middleware
│   └── server/
│       └── server.go           # обёртка над http.Server с таймаутами
├── tests/                      # интеграционные тесты API и БД
├── web/                        # статические файлы фронтенда (index.html, login.html, css, js)
├── .env                        # переменные окружения (не коммитится в реальном .gitignore)
├── Dockerfile                  # multi-stage сборка: golang:1.26.5-alpine → alpine:3.24
├── go.mod                      # модуль программы и его зависимости
└── README.md
```

Проект организован по слоям: `api` — транспорт (разбор запросов, валидация, JSON-ответы), `domain` — бизнес-правила вычисления дат, `database` — доступ к SQLite, `config`/`server`/`router` — сборка и запуск приложения. Такое разделение позволяет менять формат API или способ хранения, не затрагивая остальные слои.

---

## Быстрый старт

### Требования

- Go 1.26.5 (или совместимая версия)
- Git

### Установка и запуск

```bash
git clone https://github.com/bek63-dev/final-project-13-14-sprint.git
cd final-project-13-14-sprint
go mod tidy
```

Создайте в корне проекта файл `.env` (пример ниже) и запустите сервер:

```bash
go run ./cmd/main.go
```

При старте сервер подключается к файлу базы данных (создаёт его и таблицу `scheduler`, если их ещё нет) и начинает слушать указанный адрес. В консоли появится сообщение вида:

```
SERVER: Сервер запущен на :7540
```

Откройте в браузере:

```
http://localhost:7540/
```

Если задан `TODO_PASSWORD`, сервер сначала перенаправит на `http://localhost:7540/login.html` для ввода пароля.

---

## Переменные окружения

Конфигурация читается из файла `.env` (через `godotenv`) и/или из переменных окружения процесса; переменная окружения имеет приоритет, если она задана и непуста.

| Переменная         | Обязательна | Значение по умолчанию | Назначение                                                                 |
|--------------------|:-----------:|:----------------------:|------------------------------------------------------------------------------|
| `TODO_HOST`         | нет          | `""` (все интерфейсы)   | Хост, на котором слушает сервер                                              |
| `TODO_PORT`         | нет          | `7540`                  | Порт HTTP-сервера                                                            |
| `TODO_DBFILE`       | нет          | `scheduler.db`          | Путь к файлу базы данных SQLite                                              |
| `TODO_PASSWORD`     | нет          | `""` (аутентификация выключена) | Пароль для входа; если пуст — API не защищено, JWT не требуется   |
| `TODO_SECRET_KEY`   | да, если задан пароль | —              | Секретный ключ подписи JWT-токена                                            |

Пример `.env`:

```env
TODO_HOST=localhost
TODO_PORT=7540
TODO_DBFILE=scheduler.db
TODO_PASSWORD=password123
TODO_SECRET_KEY=kj45ls12
```
---

## API

Все ответы возвращаются в формате `application/json; charset=UTF-8`. Маршруты `/api/task`, `/api/tasks` и `/api/task/done` защищены middleware `Auth` — если задан `TODO_PASSWORD`, для них требуется валидный JWT-токен в cookie `token`; в противном случае они доступны без аутентификации.

| Метод  | Путь                       | Защищён | Описание                                                            |
|--------|----------------------------|:-------:|------------------------------------------------------------------------|
| GET    | `/`                         | нет      | Отдаёт статические файлы фронтенда из `./web`                          |
| GET    | `/api/nextdate`             | нет      | Вычисляет следующую дату по `now`, `date` и `repeat`                    |
| POST   | `/api/signin`               | нет      | Вход по паролю, возвращает JWT-токен                                    |
| POST   | `/api/task`                 | да       | Добавляет задачу                                                        |
| GET    | `/api/task?id=`             | да       | Возвращает задачу по идентификатору                                     |
| PUT    | `/api/task`                 | да       | Редактирует задачу                                                      |
| DELETE | `/api/task?id=`             | да       | Удаляет задачу                                                          |
| GET    | `/api/tasks?search=`        | да       | Возвращает список ближайших задач (до 50), с фильтром по дате/тексту    |
| POST   | `/api/task/done?id=`        | да       | Отмечает задачу выполненной (удаляет или переносит на следующую дату)   |

### Примеры запросов

**Добавление задачи** — `POST /api/task`

```json
{
  "date": "20240201",
  "title": "Подвести итог",
  "comment": "Мой комментарий",
  "repeat": "d 5"
}
```

Ответ:

```json
{"id": "186"}
```

**Список задач** — `GET /api/tasks`

```json
{
  "tasks": [
    {"id": "171", "date": "20240131", "title": "Заголовок задачи", "comment": "", "repeat": ""}
  ]
}
```

**Вход** — `POST /api/signin`

```json
{"password": "password123"}
```

Ответ:

```json
{"token": "eyJhbGciOiJIUzI1NiIsIn..."}
```

Любая ошибка возвращается в едином формате:

```json
{"error": "текст ошибки"}
```

---

## Правила повторения задач

Правило хранится в поле `repeat` (не более 128 символов) и разбирается функцией `domain.NextDate(now, date, repeat)`:

| Правило        | Значение                                                                                  |
|-----------------|--------------------------------------------------------------------------------------------|
| *(пусто)*       | Задача одноразовая — при выполнении удаляется                                              |
| `d <1-400>`     | Перенос на указанное число дней (например, `d 7` — раз в неделю)                            |
| `y`             | Перенос на год вперёд                                                                        |
| `w <1-7,...>`   | Перенос на ближайший из указанных дней недели (1 — понедельник, 7 — воскресенье)             |
| `m <1-31,-1,-2>[ <1-12,...>]` | Перенос на указанные дни месяца (`-1` — последний день, `-2` — предпоследний), опционально ограниченный списком месяцев |

Дата в прошлом при добавлении/редактировании задачи автоматически заменяется: на сегодняшнюю дату для одноразовых задач или на ближайшую подходящую дату по правилу — для периодических.

---

## Аутентификация

- Если `TODO_PASSWORD` не задан — все API-маршруты доступны без проверки.
- Если пароль задан:
    1. Клиент отправляет `POST /api/signin` с паролем.
    2. Сервер сверяет пароль с `TODO_PASSWORD`, генерирует JWT (payload содержит SHA-256-хэш пароля, срок жизни токена — 8 часов) и возвращает его в поле `token`.
    3. Фронтенд сохраняет токен в cookie `token` и передаёт его при каждом запросе к защищённым маршрутам.
    4. Middleware `Auth` проверяет подпись, срок действия токена и соответствие хэша пароля текущему значению `TODO_PASSWORD` — при смене пароля старые токены автоматически становятся недействительными.
    5. При отсутствии или невалидности токена сервер возвращает `401 Unauthorized`.

---

## Тестирование

Интеграционные тесты лежат в директории `tests` и обращаются к запущенному серверу по HTTP, а также напрямую к файлу базы данных.

Перед запуском тестов:

1. Запустите сервер (`go run ./cmd/main.go`).
2. При необходимости настройте `tests/settings.go`:

| Переменная      | Значение по умолчанию | Назначение                                                        |
|------------------|:----------------------:|------------------------------------------------------------------|
| `Port`            | `7540`                  | Порт, на который отправляются тестовые запросы                    |
| `DBFile`          | `../scheduler.db`       | Путь к файлу БД (можно переопределить через `TODO_DBFILE`)         |
| `FullNextDate`    | `true`                  | Включает проверку правил `w` и `m` в `TestNextDate`                |
| `Search`          | `true`                  | Включает проверку поиска в `TestTasks`                             |
| `Token`           | `` (пусто)              | JWT-токен для тестирования с включённой аутентификацией            |

Если сервер запущен с непустым `TODO_PASSWORD`, получите токен через `/api/signin` и присвойте его переменной `Token`, иначе защищённые маршруты вернут `401`.

Запуск отдельных наборов тестов:

```bash
go test -run ^TestApp$ ./tests        # раздача статики фронтенда
go test -run ^TestDB$ ./tests         # прямая работа с таблицей scheduler
go test -run ^TestNextDate$ ./tests   # вычисление следующей даты
go test -run ^TestAddTask$ ./tests    # добавление задачи
go test -run ^TestTasks$ ./tests      # список и поиск задач
go test -run ^TestTask$ ./tests       # получение задачи по id
go test -run ^TestEditTask$ ./tests   # редактирование задачи
go test -run ^TestDone$ ./tests       # отметка о выполнении
go test -run ^TestDelTask$ ./tests    # удаление задачи
```

Или все тесты сразу:

```bash
go test ./tests
```

---

## CI/CD

Workflow `.github/workflows/tests.yml` запускается на каждый `push`:

1. Выполняет checkout репозитория и устанавливает Go `1.26.5`.
2. Устанавливает зависимости (`go mod tidy`).
3. Создаёт `.env` с тестовыми значениями (`TODO_HOST`, `TODO_PORT`, `TODO_DBFILE`, `TODO_PASSWORD`, `TODO_SECRET_KEY`).
4. Запускает сервер в фоне (`go run ./cmd/main.go &`), выжидает 100 секунд для полной инициализации.
5. Прогоняет тесты (`go test -v -run ^TestApp$ ./tests`).

---

## Docker

Приложение можно собрать в Docker-образ и запустить в контейнере, который подключается к SQLite базе данных, физически хранящейся на хосте.

### Требования

- Установленный и запущенный [Docker Desktop](https://www.docker.com/products/docker-desktop/) (или Docker Engine на Linux).

### Структура Docker-файлов

В корне проекта находятся:

```
├── Dockerfile
├── .dockerignore
```

`Dockerfile` использует multi-stage сборку:

- **Стадия `builder`** (образ `golang:1.26.5-alpine`) — компилирует статический бинарник без cgo (это возможно благодаря `modernc.org/sqlite` — чистой Go-реализации SQLite-драйвера).
- **Финальная стадия** (образ `alpine:3.24`) — минимальный runtime-образ с исполняемым файлом и директорией `web`; исходный код и инструменты сборки в него не попадают.

### Сборка образа

Из корня проекта (там, где лежит `Dockerfile`):

```bash
docker build -t scheduler:latest .
```

Флаг `-t scheduler:latest` задаёт имя и тег образа.

### Запуск контейнера
База данных SQLite монтируется с хоста через `-v` (bind mount), чтобы данные сохранялись между пересозданиями контейнера.

Ниже приведены команды для трёх окружений: **Linux/macOS**, **Windows (Git Bash)** и **Windows (PowerShell)**. Выберите раздел, соответствующий вашей ОС и терминалу.

#### Linux / macOS

**Создайте директорию на хосте для базы данных:**

```bash
mkdir -p ~/scheduler-data
```

**Запустите контейнер, подставив свои значения вместо плейсхолдеров:**

```bash
docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=<ваш_пароль> -e TODO_SECRET_KEY=<ваш_секретный_ключ> -e TODO_DBFILE=/app/data/scheduler.db -v ~/scheduler-data:/app/data scheduler:latest
```

Например, с конкретными значениями (используются только как демонстрация — при развёртывании у себя замените на собственные):

```bash
docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=password123 -e TODO_SECRET_KEY=kj45ls12 -e TODO_DBFILE=/app/data/scheduler.db -v ~/scheduler-data:/app/data scheduler:latest
```

#### Windows — Git Bash

**Создайте директорию на хосте для базы данных:**

```bash
mkdir -p ~/scheduler-data
```

**Запустите контейнер, подставив свои значения вместо плейсхолдеров:**

```bash
MSYS_NO_PATHCONV=1 docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=<ваш_пароль> -e TODO_SECRET_KEY=<ваш_секретный_ключ> -e TODO_DBFILE=/app/data/scheduler.db -v ~/scheduler-data:/app/data scheduler:latest
```

Например, с конкретными значениями:

```bash
MSYS_NO_PATHCONV=1 docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=password123 -e TODO_SECRET_KEY=kj45ls12 -e TODO_DBFILE=/app/data/scheduler.db -v ~/scheduler-data:/app/data scheduler:latest
```

> ⚠️ Команду важно вводить **одной строкой**, без переносов `\`. Также обязателен префикс `MSYS_NO_PATHCONV=1` — без него Git Bash автоматически преобразует путь `/app/data/scheduler.db` в Windows-путь вида `C:/Program Files/Git/app/data/scheduler.db`, из-за чего сервер не сможет открыть файл базы данных.

#### Windows — PowerShell

**Создайте директорию на хосте для базы данных:**

```powershell
mkdir $HOME/scheduler-data
```

**Запустите контейнер, подставив свои значения вместо плейсхолдеров:**

```powershell
docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=<ваш_пароль> -e TODO_SECRET_KEY=<ваш_секретный_ключ> -e TODO_DBFILE=/app/data/scheduler.db -v ${HOME}/scheduler-data:/app/data scheduler:latest
```

Например, с конкретными значениями:

```powershell
docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=password123 -e TODO_SECRET_KEY=kj45ls12 -e TODO_DBFILE=/app/data/scheduler.db -v ${HOME}/scheduler-data:/app/data scheduler:latest
```

> ⚠️ Команду также рекомендуется вводить одной строкой, без переносов ``` ` ``` — при вставке многострочной команды символ переноса иногда теряется, из-за чего часть флагов может не примениться.

#### Разбор параметров (общий для всех ОС)

| Параметр | Назначение |
|---|---|
| `-d` | запуск в фоновом режиме |
| `--name scheduler` | имя контейнера для удобного управления |
| `-p 7540:7540` | проброс порта: `<порт на хосте>:<порт в контейнере>` |
| `-e TODO_PASSWORD=<ваш_пароль>` | пароль для аутентификации в приложении — задайте своё значение |
| `-e TODO_SECRET_KEY=<ваш_секретный_ключ>` | секретный ключ для подписи JWT-токенов — задайте своё значение |
| `-e TODO_DBFILE=/app/data/scheduler.db` | путь к файлу БД **внутри контейнера** (должен совпадать с точкой монтирования `-v`) |
| `-v ~/scheduler-data:/app/data` (Linux/macOS, Git Bash) или `-v ${HOME}/scheduler-data:/app/data` (PowerShell) | монтирование директории хоста в директорию контейнера `/app/data` |

> 💡 Если не передать `TODO_PASSWORD` и `TODO_SECRET_KEY` через `-e`, будут использованы пустые значения по умолчанию, заданные в Dockerfile (`ENV TODO_PASSWORD=""` и `ENV TODO_SECRET_KEY=""`). Пустой пароль означает, что аутентификация в приложении будет отключена (см. `config.go`: пустая строка в `TODO_PASSWORD` трактуется приложением как "пароль не задан", и middleware `Auth` пропускает запросы без проверки).

### Проверка работы

**Убедитесь, что контейнер запущен и порт проброшен:**

```bash
docker ps
```

В колонке `PORTS` должно быть указано `0.0.0.0:7540->7540/tcp` (а не просто `7540/tcp`) — это подтверждает, что порт проброшен на хост, а не просто открыт внутри контейнера.

**Проверьте логи сервера:**

```bash
docker logs -f scheduler
```

Ожидаемый вывод:

```
Connecting to database "/app/data/scheduler.db"
Initializing database "/app/data/scheduler.db"
SERVER: Сервер запущен на :7540
```

**Откройте приложение в браузере:**

```
http://localhost:7540
```

(если возникают проблемы с `localhost`, попробуйте `http://127.0.0.1:7540`)

Введите пароль, который вы задали в `TODO_PASSWORD` при запуске контейнера.

**Убедитесь, что база данных действительно хранится на хосте:**

Linux / macOS / Git Bash:
```bash
ls -la ~/scheduler-data/
```

PowerShell:
```powershell
ls $HOME/scheduler-data/
```

Там должен появиться файл `scheduler.db`. Он сохранится даже после удаления контейнера.

### Управление контейнером

```bash
docker stop scheduler       # остановить контейнер
docker start scheduler      # запустить существующий контейнер снова
docker restart scheduler    # перезапустить контейнер
docker rm -f scheduler      # остановить и удалить контейнер
docker logs -f scheduler    # смотреть логи в реальном времени
```

При повторном запуске контейнера с тем же volume (`-v ~/scheduler-data:/app/data`) все ранее созданные задачи сохранятся, так как база данных хранится не внутри контейнера, а на хосте.

### Пересборка после изменений в коде

Если код приложения изменился, образ нужно пересобрать:

**Linux / macOS:**

```bash
docker rm -f scheduler
docker build -t scheduler:latest .
docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=<ваш_пароль> -e TODO_SECRET_KEY=<ваш_секретный_ключ> -e TODO_DBFILE=/app/data/scheduler.db -v ~/scheduler-data:/app/data scheduler:latest
```

**Windows — Git Bash:**

```bash
docker rm -f scheduler
docker build -t scheduler:latest .
MSYS_NO_PATHCONV=1 docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=<ваш_пароль> -e TODO_SECRET_KEY=<ваш_секретный_ключ> -e TODO_DBFILE=/app/data/scheduler.db -v ~/scheduler-data:/app/data scheduler:latest
```

**Windows — PowerShell:**

```powershell
docker rm -f scheduler
docker build -t scheduler:latest .
docker run -d --name scheduler -p 7540:7540 -e TODO_PASSWORD=<ваш_пароль> -e TODO_SECRET_KEY=<ваш_секретный_ключ> -e TODO_DBFILE=/app/data/scheduler.db -v ${HOME}/scheduler-data:/app/data scheduler:latest
```

Существующие данные в директории `scheduler-data` при этом не пострадают.