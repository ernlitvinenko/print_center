package models

type PaginationRequest struct {
	Limit  int32 `json:"limit" query:"limit"`
	Offset int32 `json:"offset" query:"offset"`
}

type OrderResponse struct {
	ID               int32  `json:"id"`
	DateFrom         string `json:"date_from"`
	DateTill         string `json:"date_till"`
	ManagerID        int32  `json:"manager_id"`
	CounterpartiesID int32  `json:"counterparties_id"`
	StatusID         int16  `json:"status_id"`
	Priority         int16  `json:"priority"`
}

type OrderListResponse struct {
	Success bool            `json:"success"`
	Orders  []OrderResponse `json:"orders"`
	Total   int             `json:"total"`
	Error   string          `json:"error,omitempty"`
}

// CreateOrderRequest запрос на создание заказа
type CreateOrderRequest struct {
	DateTill         string `json:"date_till" binding:"required"`
	CounterpartiesID int32  `json:"counterparties_id" binding:"required"`
	StatusID         *int16 `json:"status_id"`
	Priority         *int16 `json:"priority"`
}

// OrderItemRequest элемент заказа
type OrderItemRequest struct {
	NomenclatureID int32 `json:"nomenclature_id" binding:"required"`
	SizeID         int32 `json:"size_id" binding:"required"`
	MaterialID     int32 `json:"material_id" binding:"required"`
	PlanningCount  int32 `json:"planning_count" binding:"required"`
	TotalCount     int32 `json:"total_count" binding:"required"`
}

// CreateOrderWithItemsRequest заказ с элементами
type CreateOrderWithItemsRequest struct {
	DateTill         string               `json:"date_till" binding:"required"`
	CounterpartiesID int32                `json:"counterparties_id" binding:"required"`
	StatusID         *int16               `json:"status_id"`
	Priority         *int16               `json:"priority"`
	Items            []OrderItemRequest   `json:"items"`
}

// CreateOrderResponse ответ после создания заказа
type CreateOrderResponse struct {
	Success bool          `json:"success"`
	Order   OrderResponse `json:"order"`
	Items   []OrderItem   `json:"items,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// OrderItem элемент заказа в ответе
type OrderItem struct {
	ID             int32 `json:"id"`
	NomenclatureID int32 `json:"nomenclature_id"`
	OrderID        int32 `json:"order_id"`
	SizeID         int32 `json:"size_id"`
	MaterialID     int32 `json:"material_id"`
	PlanningCount  int32 `json:"planning_count"`
	TotalCount     int32 `json:"total_count"`
}
