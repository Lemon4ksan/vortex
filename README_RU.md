<div align="center">

# vortex

### AST-тулчейн без аллокаций и движок суверенных контрактов

_«Лучше немного скопировать, чем добавить немного зависимостей — владейте сетевыми контрактами на уровне AST»_

[![Go Version](https://img.shields.io/badge/go-1.27%2B-007d9c?logo=go&logoColor=white&style=flat-square)](https://go.dev/)
[![Go Reference](https://img.shields.io/badge/godoc-reference-007d9c?style=flat-square)](https://pkg.go.dev/github.com/lemon4ksan/vortex)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue?style=flat-square)](LICENSE)
[![Zero-Alloc](https://img.shields.io/badge/memory-0%20B%2Fop%20%7C%200%20allocs-brightgreen?style=flat-square)](docs/VORTEX.md)
[![OpenAPI 3.1](https://img.shields.io/badge/spec-OpenAPI%203.1%20%26%20AsyncAPI-blueviolet?style=flat-square)](pkg/openapi)
[![AST Engine](https://img.shields.io/badge/compiler-Go%20AST%203--Way%20Merge-orange?style=flat-square)](pkg/parser)
[![Oracle Hub](https://img.shields.io/badge/attestation-Universal%20Oracle%20v2-yellow?style=flat-square)](pkg/oracle/gen)

**vortex** — это унифицированный декларативный тулчейн контрактов, генератор высокопроизводительного кода и система реверс-инжиниринга сетевого трафика для [`aoni`](https://github.com/lemon4ksan/aoni). Он работает напрямую с абстрактными синтаксическими деревьями (AST) языка Go, используя идиоматичные интерфейсы с Godoc-директивами как единственный источник истины для протоколов REST, WebSocket, SSE, OpenAPI 3.1, AsyncAPI и Protocol Buffers.

#### [English](README.md) • Русский • [Архитектурная спецификация](docs/VORTEX.md) • [Лицензия](LICENSE)

</div>

---

## Манифест Vortex: Суверенные клиенты

На протяжении более двух десятилетий разработка сетевых API-клиентов страдала от фрагментации:
* Стандартный `net/http` оборачивался в неподдерживаемые сторонние форки TLS ради обхода базовых бот-фильтров.
* Хрупкие headless-скрипты браузеров собирались на скорую руку для решения капч, потребляя гигабайты RAM.
* Сторонние тяжеловесные SDK засоряли дерево зависимостей, создавая утечки памяти и ломая совместимость.
* При малейшем изменении внутреннего эндпоинта инженеры неделями ждали обновлений от внешних мейнтейнеров библиотек.

**Суверенная парадигма**: Ни одна инженерная команда не должна зависеть от сторонних API-библиотек. Каждый продакшен-проект обязан владеть собственным суверенным клиентом с нулевыми аллокациями, сгенерированным прямо в своей кодовой базе (`pkg/api/`) из AST-контрактов Go, OpenAPI-схем или живых дампов сетевого трафика (`.har`).

---

## Ключевые возможности

| Возможность | Описание | Команда |
| :--- | :--- | :--- |
| **AST-компиляция** | Генерация Go-клиентов без рефлексии и аллокаций из декларативных интерфейсов. | `vortex gen` |
| **Реверс-инжиниринг трафика** | Импорт дампов `.har` или OpenAPI 3.1 с безопасным 3-сторонним слиянием AST. | `vortex spec import` |
| **Universal Oracle v2** | Компиляция браузерных сайдкаров аттестации для невидимого обхода WAF и Cloudflare. | `vortex oracle` |
| **Контроль качества контрактов** | Статический линтер корректности директив, маршрутов и типов. | `vortex check --fix` |
| **In-Memory мокирование** | Генерация виртуальных in-memory тестовых серверов без внешних зависимостей. | `vortex mock` |
| **Smoke-зондирование** | Мгновенная проверка живых эндпоинтов с выводом таблицы сетевых задержек. | `vortex smoke` |
| **Инспектор трафика (Web UI)** | Локальный веб-интерфейс для анализа WebSocket-фреймов, пейлоадов и TLS JA4. | `vortex traffic` |
| **AST Borrow Checker** | Верификация отсутствия аллокаций в куче и проверка жизненного цикла памяти. | `vortex borrow` |
| **Аппаратные бенчмарки** | Замер пропускной способности, профилей памяти и pprof-диагностика. | `vortex bench` |

---

## Установка

`vortex` требует версию Go `1.27` или выше.

```bash
go install github.com/lemon4ksan/vortex/cmd/vortex@latest
```

Проверка установки:
```bash
vortex --version
```

---

## Примеры использования

### 1. Объявление интерфейса контракта

Опишите контракт API с помощью чистого Go-интерфейса и Godoc-директив `@aoni`:

```go
// Package api объявляет суверенный контракт GitHub API.
package api

import (
	"context"
)

type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Bio   string `json:"bio"`
}

// GitHubService определяет декларативный контракт API.
//
// @aoni:service GitHubService
// @aoni:baseURL https://api.github.com
type GitHubService interface {
	// GetUser возвращает профиль пользователя.
	//
	// @aoni:method GET /users/{username}
	// @aoni:header Accept: application/vnd.github.v3+json
	// @aoni:cache TTL=5m
	GetUser(ctx context.Context, username string) (*User, error)
}
```

### 2. Компиляция клиента без аллокаций (`vortex gen`)

Скомпилируйте интерфейс в производительный клиент на базе `aoni` без рантайм-рефлексии:

```bash
vortex gen ./...
```

Vortex сгенерирует файл `github_service.vortex.go`, реализующий интерфейс со строгой упаковкой данных, пулингом соединений и нулевыми аллокациями на горячем пути.

### 3. Импорт трафика и реверс-инжиниринг (`vortex spec import`)

Превратите реальный трафик браузера или мобильного приложения в типизированные Go-контракты:

```bash
# Запись живого сетевого трафика в архив HAR
vortex traffic record -out=session.har

# Импорт HAR в Go-контракты с автовыводом типов и 3-сторонним AST-слиянием
vortex spec import -spec=session.har -pkg=api -add
```

Vortex автоматически выделит эндпоинты, параметры путей, сформирует Go-структуры моделей запросов/ответов и аккуратно встроит их в ваши файлы, сохранив сделанные вручную правки.

### 4. Браузерный сайдкар аттестации (`vortex oracle`)

Разверните автоматизированный сайдкар для обхода Cloudflare Turnstile или JS-челленджей:

```bash
vortex oracle -name=portal https://portal.example.com
```

Сайдкар Oracle управляет пулом Chromium-вкладок, решает челленджи, извлекает куки `cf_clearance` и через IPC поставляет свежие токены авторизации прямо в ваш Go-клиент.

---

## Справочник CLI-команд

### Основные команды

```bash
vortex                                 # Автопилот: аудит, синхронизация и компиляция всех контрактов
vortex gen [путь]                      # Компиляция Go-клиентов без аллокаций из AST
vortex check [путь]                    # Статический линтинг и валидация контрактов (--fix для автопочинки)
vortex mock [путь]                     # Генерация виртуального in-memory сервера для тестов
vortex smoke [путь]                    # Зондирование живых эндпоинтов с таблицей задержек
vortex env [путь]                      # Поиск переменных ${VAR} и создание .env шаблонов
```

### Доменные хабы

```bash
vortex spec import -spec=api.json      # Импорт OpenAPI 3.1 / HAR в контракты Go
vortex spec export -out=openapi.json   # Экспорт Go-контрактов в спецификацию OpenAPI 3.1
vortex spec diff -spec=v2.json         # Анализ несовместимых изменений (breaking changes)
vortex oracle -name=bot https://...    # Компиляция браузерного сайдкара аттестации
vortex traffic                         # Запуск веб-интерфейса инспектора трафика
vortex ast split                       # Декомпозиция монолитных интерфейсов
vortex perf bench                      # Запуск бенчмарков и анализ аллокаций памяти
vortex borrow ./...                    # Запуск AST borrow checker для проверки нулевых аллокаций
```

### Управление воркспейсом

```bash
vortex init [имя]                      # Инициализация .vortex.yml или шаблона контракта
vortex status                          # Круговой обзор синхронизации контрактов и кода
vortex doctor                          # Диагностика здоровья тулчейна и git-синхронизации
vortex explain @aoni:cache             # Справка и примеры по Godoc-директивам
vortex clean                           # Очистка временных файлов, профилей и артефактов
```

---

## Лицензия

Проект распространяется под лицензией **BSD 3-Clause License**. Подробности в файле [LICENSE](LICENSE).

Copyright (c) 2026 Lemon4ksan. All rights reserved.
