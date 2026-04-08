package order

import (
	"backend/core/models"
	"backend/core/repositories"
	"backend/core/services"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

var orderSvc *services.OrderService

// InitOrderService инициализирует сервис заказов
func InitOrderService(svc *services.OrderService) {
	orderSvc = svc
}

// AddHandlers регистрирует order routes
func AddHandlers(e *echo.Echo, jwtMiddleware echo.MiddlewareFunc) error {
	g := e.Group("/api/v1/orders", jwtMiddleware)
	g.GET("", GetOrders)
	g.POST("", CreateOrder)
	return nil
}

// CreateOrder создаёт новый заказ с элементами
func CreateOrder(c *echo.Context) error {
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.CreateOrderResponse{
			Success: false,
			Error:   "user not authenticated",
		})
	}

	req := &models.CreateOrderWithItemsRequest{}
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.CreateOrderResponse{
			Success: false,
			Error:   "invalid request body",
		})
	}

	// Валидация
	if req.DateTill == "" {
		return c.JSON(http.StatusBadRequest, models.CreateOrderResponse{
			Success: false,
			Error:   "date_till is required",
		})
	}

	dateTill, err := time.Parse(time.RFC3339, req.DateTill)
	if err != nil {
		dateTill, err = time.Parse("2006-01-02", req.DateTill)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.CreateOrderResponse{
				Success: false,
				Error:   "invalid date_till format, use RFC3339 or YYYY-MM-DD",
			})
		}
	}

	statusID := int16(1)
	if req.StatusID != nil {
		statusID = *req.StatusID
	}

	priority := int16(0)
	if req.Priority != nil {
		priority = *req.Priority
	}

	// Конвертируем элементы в параметры для sqlc
	items := make([]repositories.AddItemToOrderParams, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, repositories.AddItemToOrderParams{
			NomenclatureID: item.NomenclatureID,
			OrderID:        0, // Будет установлен после создания заказа
			SizeID:         item.SizeID,
			MaterialID:     item.MaterialID,
			PlanningCount:  item.PlanningCount,
			TotalCount:     item.TotalCount,
		})
	}

	// Создаём заказ с элементами
	order, createdItems, err := orderSvc.CreateOrderWithItems(
		context.Background(),
		userID,
		dateTill,
		req.CounterpartiesID,
		statusID,
		priority,
		items,
	)
	if err != nil {
		log.Logger.Error().Err(err).Int32("manager_id", userID).Msg("Failed to create order")
		return c.JSON(http.StatusInternalServerError, models.CreateOrderResponse{
			Success: false,
			Error:   "failed to create order",
		})
	}

	// Формируем ответ
	orderResponse := models.OrderResponse{
		ID:               order.ID,
		DateFrom:         order.DateFrom.Time.Format(time.RFC3339),
		DateTill:         order.DateTill.Time.Format(time.RFC3339),
		ManagerID:        order.ManagerID,
		CounterpartiesID: order.CounterpartiesID,
		StatusID:         order.StatusID,
		Priority:         order.Priority,
	}

	itemResponses := make([]models.OrderItem, 0, len(createdItems))
	for _, item := range createdItems {
		itemResponses = append(itemResponses, models.OrderItem{
			ID:             item.ID,
			NomenclatureID: item.NomenclatureID,
			OrderID:        item.OrderID,
			SizeID:         item.SizeID,
			MaterialID:     item.MaterialID,
			PlanningCount:  item.PlanningCount,
			TotalCount:     item.TotalCount,
		})
	}

	log.Logger.Info().Int32("order_id", order.ID).Int32("manager_id", userID).Msg("Order created")

	return c.JSON(http.StatusCreated, models.CreateOrderResponse{
		Success: true,
		Order:   orderResponse,
		Items:   itemResponses,
	})
}

// GetOrders возвращает список заказов (для администратора — все + без менеджера)
func GetOrders(c *echo.Context) error {
	// Получаем ID пользователя из контекста (устанавливается JWT middleware)
	userID, ok := c.Get("user_id").(int32)
	if !ok {
		return c.JSON(http.StatusUnauthorized, models.OrderListResponse{
			Success: false,
			Error:   "user not authenticated",
		})
	}

	// Парсим параметры пагинации
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Получаем заказы
	orders, err := orderSvc.GetOrders(context.Background(), userID, int32(limit), int32(offset))
	if err != nil {
		log.Logger.Error().Err(err).Int32("user_id", userID).Msg("Failed to get orders for admin")
		return c.JSON(http.StatusForbidden, models.OrderListResponse{
			Success: false,
			Error:   "access denied or failed to fetch orders",
		})
	}

	// Конвертируем в response формат
	orderResponses := make([]models.OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, models.OrderResponse{
			ID:               order.ID,
			DateFrom:         order.DateFrom.Time.Format("2006-01-02T15:04:05Z"),
			DateTill:         order.DateTill.Time.Format("2006-01-02T15:04:05Z"),
			ManagerID:        order.ManagerID,
			CounterpartiesID: order.CounterpartiesID,
			StatusID:         order.StatusID,
			Priority:         order.Priority,
		})
	}

	return c.JSON(http.StatusOK, models.OrderListResponse{
		Success: true,
		Orders:  orderResponses,
		Total:   len(orderResponses),
	})
}
