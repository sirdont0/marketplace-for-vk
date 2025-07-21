# Marketplace REST API

---

## 🚀 Быстрый старт для проверяющего

1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/yourusername/marketplace-for-vk.git
   cd marketplace-for-vk
   ```
2. Создайте файл .env на основе примера:
   ```sh
   cp .env.example .env
   ```
3. Запустите проект через Docker Compose:
   ```sh
   docker-compose up --build
   ```
4. Проверьте работу API через curl или Postman:
   - Зарегистрируйте пользователя
   - Авторизуйтесь и получите токен
   - Создайте объявление
   - Получите ленту объявлений

(Примеры curl-запросов и структура ответов — ниже по тексту README)

---

## 🛠️ Основные технологии

- Go (Gin, bcrypt, JWT)
- PostgreSQL
- Docker, Docker Compose

---

## 📚 Описание API

### Авторизация и регистрация

#### POST /api/auth/register
Регистрация пользователя

- Request:
   {
    "login": "user1",
    "password": "123456"
  }
  - Response:
   {
    "id": 1,
    "login": "user1"
  }
  
#### POST /api/auth/login
Авторизация пользователя

- Request:
   {
    "login": "user1",
    "password": "123456"
  }
  - Response:
   {
    "token": "JWT_TOKEN"
  }
  
- Токен передаётся в заголовке:
   Authorization: Bearer JWT_TOKEN
  
---

### Работа с объявлениями

#### POST /api/ads
Создать объявление (только для авторизованных)

- Request:
   {
    "title": "iPhone 15 Pro",
    "text": "Новый в коробке, гарантия.",
    "image_url": "https://example.com/iphone.jpg",
    "price": 999.99
  }
  - Response:
   {
    "id": 1,
    "title": "iPhone 15 Pro",
    "text": "Новый в коробке, гарантия.",
    "image_url": "https://example.com/iphone.jpg",
    "price": 999.99,
    "author_id": 1,
    "created_at": "2024-07-21T18:00:00Z"
  }
  
- Ограничения:
  - Заголовок: 5–100 символов
  - Текст: 10–1000 символов
  - Цена: > 0
  - Адрес изображения: валидный URL

---

#### GET /api/ads
Получить ленту объявлений

- Параметры запроса:
  - page (int, default: 1) — номер страницы
  - limit (int, default: 10) — количество на странице
  - sort (string, default: created_at_desc) — сортировка: created_at_desc, price_asc, price_desc
  - min_price (float) — фильтр по минимальной цене
  - max_price (float) — фильтр по максимальной цене

- Response:
   {
    "data": [
      {
        "id": 1,
        "title": "iPhone 15 Pro",
        "text": "Новый в коробке, гарантия.",
        "image_url": "https://example.com/iphone.jpg",
        "price": 999.99,
        "author": "user1",
        "is_mine": true,
        "created_at": "2024-07-21T18:00:00Z"
      }
      // ...
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 1
    }
  }
  
- Особенности:
  - Для авторизованных пользователей поле is_mine показывает, принадлежит ли объявление текущему пользователю.

---

## 📝 Валидация и ограничения

- Логин: 3–50 символов, уникальный
- Пароль: 3–100 символов
- Заголовок объявления: 5–100 символов
- Текст объявления: 10–1000 символов
- Цена: > 0
- Адрес изображения: валидный URL

---

## 🗄️ Миграции базы данных

Миграции SQL находятся в папке migrations/.  
При первом запуске контейнеров таблицы создаются автоматически.

---

## 🛡️ Безопасность

- Пароли хранятся только в виде bcrypt-хеша.
- Авторизация через JWT.
- Защищённые эндпоинты требуют передачи токена.
- Поддержка HTTPS (по желанию).

---

## 🧑‍💻 Примеры запросов

Регистрация:
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login":"user1", "password":"123456"}'

Авторизация:
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login":"user1", "password":"123456"}'

Создать объявление:
curl -X POST http://localhost:8080/api/ads \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"iPhone 15 Pro","text":"Новый в коробке","image_url":"https://example.com/iphone.jpg","price":999.99}'

Получить ленту:
curl http://localhost:8080/api/ads

---

## 🤝 Контакты

Автор: Дмитрий Томчук (tg: @tomchukd)

---
