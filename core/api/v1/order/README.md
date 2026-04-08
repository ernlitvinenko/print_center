# Orders API

## GET /api/v1/orders

Возвращает список заказов. Требует JWT токен.

### Для администратора (role_id = 1)

Выводит:
- Все заказы, у которых **нет менеджера** (`manager_id IS NULL`)
- Все заказы, где администратор указан как менеджер

### Для менеджера

Выводит только заказы, назначенные на этого менеджера.

### Query Parameters

| Параметр | Тип    | По умолчанию | Описание         |
|----------|--------|--------------|------------------|
| limit    | int32  | 20           | Количество записей |
| offset   | int32  | 0            | Смещение         |

### Headers

```
Authorization: Bearer <jwt_token>
```

### Request Example

```bash
curl http://localhost:8000/api/v1/orders?limit=10&offset=0 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Response (Success - 200 OK)

```json
{
  "success": true,
  "orders": [
    {
      "id": 1,
      "date_from": "2025-04-01T10:00:00Z",
      "date_till": "2025-04-10T18:00:00Z",
      "manager_id": 3,
      "counterparties_id": 5,
      "status_id": 1,
      "priority": 0
    }
  ],
  "total": 1
}
```

---

## POST /api/v1/orders

Создаёт новый заказ с элементами. Требует JWT токен.

### Request Body

```json
{
  "date_till": "2025-05-01T18:00:00Z",
  "counterparties_id": 5,
  "status_id": 1,
  "priority": 2,
  "items": [
    {
      "nomenclature_id": 1,
      "size_id": 2,
      "material_id": 3,
      "planning_count": 100,
      "total_count": 100
    },
    {
      "nomenclature_id": 2,
      "size_id": 1,
      "material_id": 4,
      "planning_count": 50,
      "total_count": 50
    }
  ]
}
```

### Поля запроса

| Поле                | Тип      | Обязательное | Описание                          |
|---------------------|----------|--------------|-----------------------------------|
| date_till           | string   | ✅           | Дата окончания (RFC3339 или YYYY-MM-DD) |
| counterparties_id   | int32    | ✅           | ID контрагента                    |
| status_id           | int16    | ❌ (1)       | ID статуса                        |
| priority            | int16    | ❌ (0)       | Приоритет заказа                  |
| items               | array    | ❌           | Список элементов заказа           |
| items[].nomenclature_id | int32 | ✅        | ID номенклатуры                   |
| items[].size_id     | int32    | ✅           | ID размера                        |
| items[].material_id | int32    | ✅           | ID материала                      |
| items[].planning_count | int32 | ✅           | Плановое количество               |
| items[].total_count | int32    | ✅           | Фактическое количество            |

### Response (Success - 201 Created)

```json
{
  "success": true,
  "order": {
    "id": 10,
    "date_from": "2025-04-08T12:00:00Z",
    "date_till": "2025-05-01T18:00:00Z",
    "manager_id": 3,
    "counterparties_id": 5,
    "status_id": 1,
    "priority": 2
  },
  "items": [
    {
      "id": 20,
      "nomenclature_id": 1,
      "order_id": 10,
      "size_id": 2,
      "material_id": 3,
      "planning_count": 100,
      "total_count": 100
    },
    {
      "id": 21,
      "nomenclature_id": 2,
      "order_id": 10,
      "size_id": 1,
      "material_id": 4,
      "planning_count": 50,
      "total_count": 50
    }
  ]
}
```

### Response (Error - 400 Bad Request)

```json
{
  "success": false,
  "error": "date_till is required"
}
```

### Response (Error - 401 Unauthorized)

```json
{
  "success": false,
  "error": "user not authenticated"
}
```

### Response (Error - 500 Internal Server Error)

```json
{
  "success": false,
  "error": "failed to create order"
}
```
