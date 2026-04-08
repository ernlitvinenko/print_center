package orders

import (
	"backend/core/repositories"
	"backend/core/storage"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

var orderStorage storage.OrderStorage
var jwtStorage storage.JWTStorage

// InitStorage связывает handler с реализациями storage.
func InitStorage(ord storage.OrderStorage, jwt storage.JWTStorage) {
	orderStorage = ord
	jwtStorage = jwt
}

// getOrdersHandler возвращает список заказов с пагинацией.
func getOrdersHandler(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, OrderListResponse{
			Success: false,
			Error:   "user not authenticated",
		})
	}

	limit := parseQueryInt(c, "limit", 20)
	offset := parseQueryInt(c, "offset", 0)

	orders, err := orderStorage.GetOrders(c.Request().Context(), userID, int32(limit), int32(offset))
	if err != nil {
		log.Logger.Error().Err(err).Int32("user_id", userID).Msg("Failed to get orders")
		return c.JSON(http.StatusForbidden, OrderListResponse{
			Success: false,
			Error:   "access denied or failed to fetch orders",
		})
	}

	return c.JSON(http.StatusOK, OrderListResponse{
		Success: true,
		Orders:  toOrderResponses(orders),
		Total:   len(orders),
	})
}

// createOrderHandler создаёт заказ с элементами в транзакции.
func createOrderHandler(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, CreateOrderResponse{
			Success: false,
			Error:   "user not authenticated",
		})
	}

	req := &CreateOrderRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, CreateOrderResponse{
			Success: false,
			Error:   "invalid request body",
		})
	}

	dateTill, err := parseDateTill(req.DateTill)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateOrderResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	statusID := int16(1)
	if req.StatusID != nil {
		statusID = *req.StatusID
	}

	priority := int16(0)
	if req.Priority != nil {
		priority = *req.Priority
	}

	items := buildOrderItemParams(req.Items)

	order, createdItems, err := orderStorage.CreateOrderWithItems(
		c.Request().Context(),
		userID,
		dateTill,
		req.CounterpartiesID,
		statusID,
		priority,
		items,
	)
	if err != nil {
		log.Logger.Error().Err(err).Int32("manager_id", userID).Msg("Failed to create order")
		return c.JSON(http.StatusInternalServerError, CreateOrderResponse{
			Success: false,
			Error:   "failed to create order",
		})
	}

	log.Logger.Info().Int32("order_id", order.ID).Int32("manager_id", userID).Msg("Order created")

	return c.JSON(http.StatusCreated, CreateOrderResponse{
		Success: true,
		Order:   toOrderResponse(order),
		Items:   toItemResponses(createdItems),
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func toOrderResponse(o repositories.Order) OrderResponse {
	return OrderResponse{
		ID:               o.ID,
		DateFrom:         o.DateFrom.Time.Format("2006-01-02T15:04:05Z"),
		DateTill:         o.DateTill.Time.Format("2006-01-02T15:04:05Z"),
		ManagerID:        o.ManagerID,
		CounterpartiesID: o.CounterpartiesID,
		StatusID:         o.StatusID,
		Priority:         o.Priority,
	}
}

func toOrderResponses(orders []repositories.Order) []OrderResponse {
	out := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, toOrderResponse(o))
	}
	return out
}

func toItemResponses(items []repositories.OrderItem) []OrderItemResponse {
	out := make([]OrderItemResponse, 0, len(items))
	for _, i := range items {
		out = append(out, OrderItemResponse{
			ID:             i.ID,
			NomenclatureID: i.NomenclatureID,
			OrderID:        i.OrderID,
			SizeID:         i.SizeID,
			MaterialID:     i.MaterialID,
			PlanningCount:  i.PlanningCount,
			TotalCount:     i.TotalCount,
		})
	}
	return out
}

func parseQueryInt(c *echo.Context, key string, defaultValue int) int {
	val, _ := strconv.Atoi(c.QueryParam(key))
	if val <= 0 {
		return defaultValue
	}
	return val
}

func parseDateTill(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, &parseError{"date_till is required"}
	}
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}
	return time.Time{}, &parseError{"invalid date_till format, use RFC3339 or YYYY-MM-DD"}
}

func buildOrderItemParams(items []OrderItemRequest) []repositories.AddItemToOrderParams {
	params := make([]repositories.AddItemToOrderParams, 0, len(items))
	for _, i := range items {
		params = append(params, repositories.AddItemToOrderParams{
			NomenclatureID: i.NomenclatureID,
			OrderID:        0,
			SizeID:         i.SizeID,
			MaterialID:     i.MaterialID,
			PlanningCount:  i.PlanningCount,
			TotalCount:     i.TotalCount,
		})
	}
	return params
}

type parseError struct{ msg string }

func (e *parseError) Error() string { return e.msg }
