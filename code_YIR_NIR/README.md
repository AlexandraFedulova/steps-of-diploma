# Социометрия

Веб-приложение для анализа **социометрических графов**: загрузка социоматрицы в формате CSV, построение ориентированного взвешенного графа, визуализация, **degree centrality** (по входящим связям) и выделение подгрупп как **клик взаимного положительного выбора**.

Типичный сценарий: исследователь загружает матрицу выборов/отношений между респондентами → получает интерактивную социограмму с раскраской по сообществам и метриками по вершинам.

## Стек

| Часть | Технологии |
|-------|------------|
| **Frontend** | React 19, TypeScript, Vite, React Router, Cytoscape.js (fcose) |
| **Backend** | Go 1.23, стандартная библиотека `net/http` |
| **Shared** | TypeScript-пакет `@sociometry/shared` (типы API и DTO графа) |
| **Инфраструктура** | pnpm workspaces, Docker, nginx (production frontend) |
| **Тесты** | Go `testing`, Vitest, React Testing Library |

## Структура репозитория

```
.
├── frontend/                 # SPA: загрузка CSV, граф, интерактив
│   ├── Dockerfile
│   └── nginx.conf            # прокси /api → backend
├── backend/
│   ├── cmd/server/           # точка входа HTTP API
│   ├── internal/             # handlers, router, middleware
│   └── pkg/
│       ├── graph/            # граф, centrality, BFS, клики, Louvain (вспомог.)
│       └── sociomatrix/      # парсер CSV
│           └── testdata/     # примеры матриц
├── shared/                   # общие TypeScript-типы
├── docker-compose.yml
└── README.md
```

---

## Требования

| Инструмент | Версия | Проверка |
|------------|--------|----------|
| Node.js | ≥ 20 | `node -v` |
| pnpm | ≥ 9 | `pnpm -v` |
| Go | ≥ 1.23 | `go version` |
| Docker (опционально) | — | `docker compose version` |

Установка pnpm через Corepack:

```bash
corepack enable
corepack prepare pnpm@9.15.0 --activate
```

---

## Запуск

### Вариант A — Docker (рекомендуется для первого знакомства)

Из **корня** репозитория (нужен запущенный Docker Desktop):

```bash
docker compose up --build
```

| Сервис | URL |
|--------|-----|
| Приложение (UI) | http://localhost:3000 |
| Backend API | http://localhost:8080 |

Проверка API:

```bash
curl http://localhost:8080/health
```

Остановка: `Ctrl+C` или `docker compose down`.

---

### Вариант B — локальная разработка

Нужны **два терминала** (backend и frontend). Все команды ниже — из **корня** репозитория, если не указано иное.

#### 1. Зависимости

```bash
pnpm install
pnpm --filter @sociometry/shared build
```

> Frontend импортирует `@sociometry/shared`. После изменений в `shared/src/` пересоберите пакет той же командой.

#### 2. Backend (терминал 1)

```bash
cd backend
go run ./cmd/server
```

Ожидаемый вывод: `backend listening on :8080`

```bash
curl http://localhost:8080/health
# {"status":"ok","service":"sociometry-backend"}
```

Порт задаётся переменной `ADDR` (по умолчанию `:8080`):

```bash
ADDR=:9090 go run ./cmd/server
```

#### 3. Frontend (терминал 2)

```bash
pnpm --filter frontend dev
```

Откройте **http://localhost:5173**

- **/** — главная (проверка связи с backend)
- **/upload** — загрузка CSV → граф (клики + centrality)
- **/graph** — визуализация графа после успешной загрузки

Vite проксирует `/api/*` → `http://localhost:8080/*` (префикс `/api` снимается).

Альтернатива (только frontend + watch shared, **без** backend):

```bash
pnpm dev
```

---

## Работа в UI (пример сценария)

1. Запустите backend и frontend (или `docker compose up`).
2. Откройте http://localhost:5173/upload (или :3000 в Docker).
3. Перетащите или выберите CSV, например `backend/pkg/sociomatrix/testdata/health.csv`.
4. После анализа откроется страница **Граф**: номера на узлах, имя респондента — по клику, zoom/pan, tooltip при наведении.

---

## Формат CSV (социоматрица)

Файл — **квадратная** матрица в кодировке UTF-8.

### Структура

| Элемент | Описание |
|---------|----------|
| **Строка 1** | Заголовок: первая ячейка пустая, далее **метки столбцов** (имена респондентов) |
| **Строки 2…n+1** | `метка_строки, значение₁, значение₂, …` — столько значений, сколько столбцов |
| **Метки** | Метки строк и столбцов должны **совпадать** и идти в **одном порядке** |
| **Размер** | Число строк данных = числу столбцов (квадрат n×n) |

Связь читается как **«от строки → к столбцу»** (ориентированное ребро).

### Значения ячеек

| Запись в CSV | Смысл | Вес в графе |
|--------------|-------|-------------|
| пусто, `0`, `.` | Нет связи | ребро не создаётся |
| `+` | Положительная | +1 |
| `++` | Сильная положительная | +2 |
| `-` | Отрицательная | −1 |
| `--` | Сильная отрицательная | −2 |
| положительное число (`1`, `2.5`) | Положительная | значение числа |
| отрицательное число (`-2`, `-1`) | Отрицательная | модуль числа со знаком минус |

### Минимальный пример

Файл `backend/pkg/sociomatrix/testdata/sample.csv`:

```csv
,A,B,C
A,0,+,-
B,1,0,-2
C,.,++,0
```

Расшифровка:

- `A → B`: `+` (вес +1)
- `A → C`: `-` (вес −1)
- `B → A`: `1` (вес +1)
- `B → C`: `-2` (вес −2)
- `C → B`: `++` (вес +2)
- ячейки `0`, `.` — нет исходящей связи

### Пример с буквенными метками (больше респондентов)

Начало `health.csv`:

```csv
,A,B,C,D,E,...
A,0,+,.,+,1,...
B,++,0,-1,-,.,...
C,++,-,0,.,-1,...
```

Полный файл: `backend/pkg/sociomatrix/testdata/health.csv` (15×15).

### Типичные ошибки при разборе

- Несовпадение меток строк и столбцов
- Разное число строк и столбцов (не квадрат)
- Дублирующиеся метки
- Некорректное значение ячейки (не из таблицы выше)
- Пустой файл или только строка заголовка

---

## Примеры

### Тестовые CSV в репозитории

| Файл | Описание |
|------|----------|
| `sample.csv` | Маленькая матрица 3×3 для отладки |
| `communities.csv` | Три явных клики по 3 респондента (для проверки раскраски) |
| `health.csv` | Средняя матрица с `+`, `++`, `-`, `--`, числами |
| `frgm.csv` | Пример разрозненного коллектива |
| `test1.csv` | Реальные данные социометрии |
| `conflict.csv` | Пример с конфликтными связями |

Путь: `backend/pkg/sociomatrix/testdata/`.

### HTTP API (curl)

Команды с путём `@pkg/...` выполняйте из каталога **`backend/`**.

**Health-check:**

```bash
curl http://localhost:8080/health
```

**Загрузка CSV → полный граф (узлы, рёбра, centrality, клики, матрица весов):**

```bash
cd backend
curl -s -X POST \
  -F "file=@pkg/sociomatrix/testdata/sample.csv" \
  http://localhost:8080/upload | python3 -m json.tool
```

Фрагмент ответа:

```json
{
  "ok": true,
  "filename": "sample.csv",
  "graph": {
    "directed": true,
    "nodeCount": 3,
    "edgeCount": 5,
    "nodes": [{ "id": "A", "label": "A" }],
    "edges": [{ "source": "A", "target": "B", "weight": 1 }],
    "centrality": [{ "nodeId": "B", "inWeightSum": 3, "value": 1.5 }],
    "community": [{ "nodeId": "A", "communityId": 0 }]
  }
}
```

**Поиск клик-сообществ (отдельный эндпоинт):**

```bash
curl -s -X POST \
  -F "file=@pkg/sociomatrix/testdata/frgm.csv" \
  http://localhost:8080/communities | python3 -m json.tool
```

### Эндпоинты API

| Метод | Путь | Тело | Ответ |
|-------|------|------|--------|
| `GET` | `/health` | — | `{ "status": "ok", "service": "..." }` |
| `POST` | `/upload` | `multipart/form-data`, поле `file` | `{ "ok", "graph": GraphPayload }` |
| `POST` | `/communities` | `multipart/form-data`, поле `file` | `{ "ok", "cliques": ... }` |

В Docker и dev frontend обращается к API через префикс **`/api`** (например `/api/upload`).

---

## Сборка и тесты

### Production-сборка

```bash
# Backend
cd backend && go build -o bin/server ./cmd/server

# Frontend
pnpm --filter @sociometry/shared build
pnpm --filter frontend build
```

Артефакт UI: `frontend/dist/`.

### Тесты

```bash
# Backend (парсер CSV, граф, centrality, BFS, клики)
cd backend && go test ./...

# Frontend
cd frontend && pnpm test
```

---

## Частые проблемы

| Симптом | Решение |
|---------|---------|
| `pnpm: command not found` | Установите pnpm через Corepack (см. выше) |
| Ошибки TypeScript про `@sociometry/shared` | `pnpm --filter @sociometry/shared build` |
| Порт 8080 занят | `lsof -i :8080` → `kill <PID>` или `ADDR=:9090` для backend |
| UI: «Сервер недоступен» | Запущен ли backend? `curl http://localhost:8080/health` |
| UI без API при открытии `index.html` | Используйте Vite (`pnpm dev`) или Docker/nginx, не `file://` |
| `curl` не находит CSV | Выполняйте команду из `backend/` или укажите полный путь к файлу |
| Docker: daemon not running | Запустите Docker Desktop |

---

## Шпаргалка команд

```bash
# Локально (из корня)
pnpm install
pnpm --filter @sociometry/shared build
pnpm --filter frontend dev          # терминал 2
cd backend && go run ./cmd/server   # терминал 1

# Docker
docker compose up --build

# Тесты
cd backend && go test ./...
cd frontend && pnpm test
```

---

## Типы и контракт API

- TypeScript: `shared/src/index.ts` → пакет `@sociometry/shared`
- Go DTO графа: `backend/pkg/graph/dto.go` (`GraphPayload`)

После изменения shared пересоберите: `pnpm --filter @sociometry/shared build`.
