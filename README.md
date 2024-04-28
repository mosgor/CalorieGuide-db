# ЗАПРОСЫ И ОТВЕТЫ

*omitempty - поле может быть пропущено*

## Продукты

- Создание продукта `POST /product`  
**Требуется авторизация по токену**

Request example:

```json
{
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "author_id": "number",
  "picture": "int8[],omitempty"
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
  "picture": "int8[],omitempty"
}
```

- Получение данных по всем продуктам `GET /products`

Если отправить пустой запрос, то режим сортировки по умолчанию будет от новейшего, а для вывода будет выбрана первая двадцатка
Значение `likesAsc` означает `ascending` (по возрастанию лайков)  
Значение `likesDesc` означает `descending` (по убыванию лайков)

Request example:

```json
{
  "sort": "likesAsc/likesDesc/fromNewest/fromOldest,omitempty",
  "two-decade": "number,omitempty"
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
      "picture": "int8[],omitempty"
    },
    "..."
  ]
}
```

- Полуние данных по одному продукту `GET /products/<product_id>`

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
  "picture": "int8[],omitempty"
}
```

- Обновление данных продукта `PUT /products/<product_id>`  
**Требуется авторизация по токену**

Request example:

```json
{
  "food_name": "string",
  "description": "string,omitempty",
  "calories": "number",
  "proteins": "number",
  "carbohydrates": "number",
  "fats": "number",
  "picture": "int8[],omitempty"
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
  "picture": "int8[],omitempty"
}
```

- Удаление продукта `DELETE /products/<product_id>`  
**Требуется авторизация по токену**

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
  "picture": "int8[],omitempty"
}
```

- Лайк `POST /products/like`  
**Требуется авторизация по токену**

Отправка лайка пользователем уже поставившим его на продукт, удалит его

Request example:

```json
{
  "user_id": "number",
  "product_id": "number"
}
```

Response example:

```json
{
  "user_id": "number",
  "product_id": "number",
  "action": "liked/disliked"
}
```

## Пользователь

- Создание пользователя `POST /user`

Request example:

```json
{
  "user_name": "string",
  "surname": "string",
  "email": "string",
  "password": "string",
  "picture": "int8[],omitempty"
}
```

Response example:

```json
{
  "id": "number",
  "user_name": "string",
  "surname": "string",
  "email": "string",
  "password": "string",
  "picture": "int8[],omitempty"
}
```

- Авторизация пользователя `POST /login`

Request example:

```json
{
  "email": "string",
  "password": "string"
}
```

Response example:

```json
{
  "id": "number",
  "user_name": "string",
  "surname": "string",
  "email": "string",
  "password": "string",
  "picture": "int8[],omitempty",
  "breakfast_id": "number,omitempty",
  "lunch_id": "number,omitempty",
  "dinner_id": "number,omitempty",
  "calories_goal": "number",
  "fats_goal": "number",
  "proteins_goal": "number",
  "carbohydrates_goal": "number",
  "BearerToken": "string"
}
```

- Обновление информации о пользователе `PUT /user/<user_id>`  
**Требуется авторизация по токену**

Request example:

```json
{
  "user_name": "string",
  "surname": "string",
  "email": "string",
  "password": "string",
  "picture": "int8[],omitempty",
  "breakfast_id": "number,omitempty",
  "lunch_id": "number,omitempty",
  "dinner_id": "number,omitempty",
  "calories_goal": "number",
  "fats_goal": "number",
  "proteins_goal": "number",
  "carbohydrates_goal": "number"
}
```

Response example:

```json
{
  "id": "number",
  "user_name": "string",
  "surname": "string",
  "email": "string",
  "password": "string",
  "picture": "int8[],omitempty",
  "breakfast_id": "number,omitempty",
  "lunch_id": "number,omitempty",
  "dinner_id": "number,omitempty",
  "calories_goal": "number",
  "fats_goal": "number",
  "proteins_goal": "number",
  "carbohydrates_goal": "number"
}
```

- Удаление пользователя `DELETE /user/<user_id>`  
**Требуется авторизация по токену**

Response example:

```json
{
  "id": "number",
  "user_name": "string",
  "surname": "string",
  "email": "string",
  "password": "string",
  "picture": "int8[],omitempty",
  "breakfast_id": "number,omitempty",
  "lunch_id": "number,omitempty",
  "dinner_id": "number,omitempty",
  "calories_goal": "number",
  "fats_goal": "number",
  "proteins_goal": "number",
  "carbohydrates_goal": "number"
}
```

## Приёмы пищи

- Получение всех приёмов пищи `GET /meals`

Request example:

```json
{
  "sort": "likesAsc/likesDesc/fromNewest/fromOldest,omitempty",
  "two_decade": "number,omitempty"
}
```

Response example:

```json
{
  "meals": [
    {
      "id": "number",
      "meal_name": "string",
      "total_calories": "number",
      "total_proteins": "number",
      "total_fats": "number",
      "total_carbohydrates": "number",
      "products_id": [{"product_id": "number", "quantity": "number"}, "..."],
      "author_id": "number",
      "description": "string,omitempty",
      "likes": "number,omitempty"
    },
    "..."
  ]
}
```

- Получение приёма пищи `GET /meals/<meal_id>`

Response example:

```json
{
  "id": "number",
  "meal_name": "string",
  "total_calories": "number",
  "total_proteins": "number",
  "total_fats": "number",
  "total_carbohydrates": "number",
  "products_id": [{"product_id": "number", "quantity": "number"}, "..."],
  "author_id": "number",
  "description": "string,omitempty",
  "likes": "number,omitempty"
}
```

- Добавление приёма пищи `POST /meal`  
**Требуется авторизация по токену**

Request example:

```json
{
  "meal_name": "string",
  "products_id": [{"product_id": "number", "quantity": "number"}, "..."],
  "author_id": "number",
  "description": "string,omitempty"
}
```

Response example:

```json
{
  "id": "number",
  "meal_name": "string",
  "total_calories": "number",
  "total_proteins": "number",
  "total_fats": "number",
  "total_carbohydrates": "number",
  "products_id": [{"product_id": "number", "quantity": "number"}, "..."],
  "author_id": "number",
  "description": "string,omitempty",
  "likes": "number,omitempty"
}
```

- Обновление информации о приёме пищи `PUT /meals/<meal_id>`  
**Требуется авторизация по токену**

Request example:

```json
{
  "meal_name": "string",
  "products_id": [{"product_id": "number", "quantity": "number"}, "..."],
  "description": "string,omitempty"
}
```

Response example:

```json
{
  "id": "number",
  "meal_name": "string",
  "total_calories": "number",
  "total_proteins": "number",
  "total_fats": "number",
  "total_carbohydrates": "number",
  "products_id": [{"product_id": "number", "quantity": "number"}, "..."],
  "author_id": "number",
  "description": "string,omitempty",
  "likes": "number,omitempty"
}
```

- Лайк `POST /meals/like`  
**Требуется авторизация по токену**

Request example:

```json
{
  "meal_id": "number",
  "user_id": "number"
}
```

Response example:

```json
{
  "meal_id": "number",
  "user_id": "number",
  "action": "liked/disliked"
}
```

- Удаление приёма пищи `DELETE /meals/<meal_id>`  
**Требуется авторизация по токену**

Response example:

```json
{
  "id": "number",
  "meal_name": "string",
  "total_calories": "number",
  "total_proteins": "number",
  "total_fats": "number",
  "total_carbohydrates": "number",
  "products_id": [{"product_id": "number", "quantity": "number"}, "..."],
  "author_id": "number",
  "description": "string,omitempty",
  "likes": "number,omitempty"
}
```