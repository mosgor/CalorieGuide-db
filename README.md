# ЗАПРОСЫ И ОТВЕТЫ

*omitempty - поля может быть пропущено*

# Продукты

- Создание продукта `POST /product`

Request example:

```json
{
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "picture": "byte[],omitempty"
}
```

Response example:

```json
{
  "id": "number",
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "author_id": "number",
  "likes": "number,omitempty",
  "picture": "byte[],omitempty"
}
```

- Получение данных по всем продуктам `GET /products`

Если отправить пустой запрос, то режим сортировки по умолчанию будет от новейшего
Значение `likesAsc` означает `ascending` (по возрастанию лайков)
Значение `likesDesc` означает `descending` (по убыванию лайков)

Request example:

```json
{
  "sort": "likesAsc/likesDesc/fromNewest/fromOldest"
}
```

Response example:

```json
{
  "products": [
    {
      "id": "number",
      "food_name": "string",
      "description": "string,omitempty",
      "calories": "number",
      "proteins": "number",
      "carbohydrates": "number",
      "fats": "number",
      "author_id": "number",
      "likes": "number,omitempty",
      "picture": "byte[],omitempty"
    },
    "..."
  ]
}
```

- Полуние данных по одному продукту `GET /product/<product_id>` **Скоро будет добавлено**

Response example:

```json
{
  "id": "number",
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "author_id": "number",
  "likes": "number,omitempty",
  "picture": "byte[],omitempty"
}
```

- Обновление данных продукта `PUT /product/<product_id>` **Скоро будет добавлено**

Request example:

```json
{
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "picture": "byte[],omitempty"
}
```

Response example:

```json
{
  "id": "number",
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "author_id": "number",
  "likes": "number,omitempty",
  "picture": "byte[],omitempty"
}
```

- Удаление продукта `DELETE /product/<product_id>` **Скоро будет добавлено**

Response example:

```json
{
  "id": "number",
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "author_id": "number",
  "likes": "number,omitempty",
  "picture": "byte[],omitempty"
}
```

# Пользователь

- Создание пользователя `POST /user` **Скоро будет добавлено**

Request example:

```json
{
  "name": "string",
  "surname": "string",
  "email": "string",
  "password": "string"
}
```

Response example:

```json
{
  "id": "number",
  "name": "string",
  "surname": "string",
  "email": "string",
  "password": "string"
}
```

- Получение информации о пользователе `GET /user/<user_id>` **Скоро будет добавлено**

Response example:

```json
{
  "id": "number",
  "name": "string",
  "surname": "string",
  "email": "string",
  "password": "string"
}
```

- Обновление информации о пользователе `PUT /user/<user_id>` **Скоро будет добавлено**

Request example:

```json
{
  "name": "string",
  "surname": "string",
  "email": "string",
  "password": "string"
}
```

Response example:

```json
{
  "id": "number",
  "name": "string",
  "surname": "string",
  "email": "string",
  "password": "string"
}
```

- Удаление пользователя `DELETE /user/<user_id>` **Скоро будет добавлено**

Response example:

```json
{
  "id": "number",
  "name": "string",
  "surname": "string",
  "email": "string",
  "password": "string"
}
```

# Диета

- Получение информации о диете `GET /user/<user_id>/diet` **Скоро будет добавлено**

Response example: 

```json
{
  "id": "number",
  "breakfast_id": "number",
  "lunch_id": "number",
  "dinner_id": "number"
}
```

- Обновление информации о диете `PUT /user/<user_id>/diet` **Скоро будет добавлено**

Request example:

```json
{
  "breakfast_id": "number,omitempty",
  "lunch_id": "number,omitempty",
  "dinner_id": "number,omitempty"
}
```

Response example:

```json
{
  "id": "number",
  "breakfast_id": "number",
  "lunch_id": "number",
  "dinner_id": "number"
}
```

# Цели

- Получение информации о цели `GET /user/<user_id>/goal` **Скоро будет добавлено**

Response example:

```json
{
  "id": "number",
  "calories_goal": "number",
  "fats_goal": "number",
  "proteins_goal": "number",
  "carbohydrates_goal": "number"
}
```

- Обновление информации о цели `PUT /user/<user_id>/goal` **Скоро будет добавлено**

Request example:

```json
{
  "calories_goal": "number,omitempty",
  "fats_goal": "number,omitempty",
  "proteins_goal": "number,omitempty",
  "carbohydrates_goal": "number,omitempty"
}
```

Response example:

```json
{
  "id": "number",
  "calories_goal": "number",
  "fats_goal": "number",
  "proteins_goal": "number",
  "carbohydrates_goal": "number"
}
```
