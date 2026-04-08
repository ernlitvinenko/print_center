package services

import (
	"backend/core/repositories"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type OrderService struct {
	repo *repositories.Queries
	db   repositories.DBTX
}

func NewOrderService(repo *repositories.Queries, db repositories.DBTX) *OrderService {
	return &OrderService{
		repo: repo,
		db:   db,
	}
}

func (s *OrderService) GetOrders(ctx context.Context, profileID, limit, offset int32) ([]repositories.Order, error) {
	//	 Получаем роль пользователя
	roles, err := s.repo.GetProfileRoles(ctx, profileID)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, errors.New("No roles for user")
	}

	for _, role := range roles {
		if role.ID == 1 {
			return s.GetOrdersForAdmin(ctx, profileID, limit, offset)
		}
	}
	return s.GetOrdersForManager(ctx, profileID, limit, offset)
}

// GetOrdersForAdmin возвращает заказы для администратора (все + без менеджера)
func (s *OrderService) GetOrdersForAdmin(ctx context.Context, profileID int32, limit, offset int32) ([]repositories.Order, error) {
	orders, err := s.repo.ListOrdersForAdmin(ctx, repositories.ListOrdersForAdminParams{
		ManagerID: profileID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, nil
}

// GetOrdersForManager возвращает заказы менеджера
func (s *OrderService) GetOrdersForManager(ctx context.Context, managerID, limit, offset int32) ([]repositories.Order, error) {
	orders, err := s.repo.ListAllOrdersToManager(ctx, repositories.ListAllOrdersToManagerParams{
		ManagerID: managerID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, nil
}

// CreateOrderWithItems создаёт заказ с элементами в транзакции
func (s *OrderService) CreateOrderWithItems(
	ctx context.Context,
	managerID int32,
	dateTill time.Time,
	counterpartiesID int32,
	statusID int16,
	priority int16,
	items []repositories.AddItemToOrderParams,
) (repositories.Order, []repositories.OrderItem, error) {
	// Приводим db к pgx.Conn для начала транзакции
	conn, ok := s.db.(interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	})
	if !ok {
		return repositories.Order{}, nil, fmt.Errorf("database connection does not support transactions")
	}

	// Начинаем транзакцию
	tx, err := conn.Begin(ctx)
	if err != nil {
		return repositories.Order{}, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.WithTx(tx)

	// Создаём заказ
	order, err := qtx.CreateOrder(ctx, repositories.CreateOrderParams{
		DateTill:         pgtype.Timestamptz{Time: dateTill, Valid: true},
		ManagerID:        managerID,
		CounterpartiesID: counterpartiesID,
		StatusID:         statusID,
		Priority:         priority,
	})
	if err != nil {
		return repositories.Order{}, nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Добавляем элементы заказа
	createdItems := make([]repositories.OrderItem, 0, len(items))
	for _, item := range items {
		createdItem, err := qtx.AddItemToOrder(ctx, item)
		if err != nil {
			return repositories.Order{}, nil, fmt.Errorf("failed to add order item: %w", err)
		}
		createdItems = append(createdItems, createdItem)
	}

	// Коммитим транзакцию
	if err := tx.Commit(ctx); err != nil {
		return repositories.Order{}, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, createdItems, nil
}
