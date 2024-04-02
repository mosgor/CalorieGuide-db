# Запросы и ответы

*omitempty - поля может быть пропущено*

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