package orders

// ---------------------------------------------------------------------------
// Request/Response типы
// ---------------------------------------------------------------------------

// OrderResponse — заказ в ответе API.
type OrderResponse struct {
	ID               int32  `json:"id"`
	DateFrom         string `json:"date_from"`
	DateTill         string `json:"date_till"`
	ManagerID        int32  `json:"manager_id"`
	CounterpartiesID int32  `json:"counterparties_id"`
	StatusID         int16  `json:"status_id"`
	Priority         int16  `json:"priority"`
}

// OrderItemResponse — элемент заказа.
type OrderItemResponse struct {
	ID             int32 `json:"id"`
	NomenclatureID int32 `json:"nomenclature_id"`
	OrderID        int32 `json:"order_id"`
	SizeID         int32 `json:"size_id"`
	MaterialID     int32 `json:"material_id"`
	PlanningCount  int32 `json:"planning_count"`
	TotalCount     int32 `json:"total_count"`
}

// OrderListResponse — список заказов.
type OrderListResponse struct {
	Success bool            `json:"success"`
	Orders  []OrderResponse `json:"orders"`
	Total   int             `json:"total"`
	Error   string          `json:"error,omitempty"`
}

// OrderItemRequest — элемент при создании заказа.
type OrderItemRequest struct {
	NomenclatureID int32 `json:"nomenclature_id"`
	SizeID         int32 `json:"size_id"`
	MaterialID     int32 `json:"material_id"`
	PlanningCount  int32 `json:"planning_count"`
	TotalCount     int32 `json:"total_count"`
}

// CreateOrderRequest — запрос на создание заказа.
type CreateOrderRequest struct {
	DateTill         string             `json:"date_till"`
	CounterpartiesID int32              `json:"counterparties_id"`
	StatusID         *int16             `json:"status_id"`
	Priority         *int16             `json:"priority"`
	Items            []OrderItemRequest `json:"items"`
}

// CreateOrderResponse — ответ после создания заказа.
type CreateOrderResponse struct {
	Success bool              `json:"success"`
	Order   OrderResponse     `json:"order"`
	Items   []OrderItemResponse `json:"items,omitempty"`
	Error   string            `json:"error,omitempty"`
}
