package storage

import (
	"backend/core/repositories"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// OrderStorage — хранилище для работы с заказами.
type OrderStorage interface {
	Storage
	GetOrders(ctx context.Context, profileID, limit, offset int32) ([]repositories.Order, error)
	CreateOrderWithItems(ctx context.Context, managerID int32, dateTill time.Time, counterpartiesID int32, statusID int16, priority int16, items []repositories.AddItemToOrderParams) (repositories.Order, []repositories.OrderItem, error)
}

// GetOrders возвращает заказы в зависимости от роли пользователя.
func (s *storage) GetOrders(ctx context.Context, profileID, limit, offset int32) ([]repositories.Order, error) {
	roles, err := s.queries.GetProfileRoles(ctx, profileID)

	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("no roles for user")
	}

	// role_id = 1 — администратор
	for _, role := range roles {
		if role.ID == 1 {
			return s.getOrdersForAdmin(ctx, profileID, limit, offset)
		}
	}

	return s.getOrdersForManager(ctx, profileID, limit, offset)
}

// getOrdersForAdmin — все заказы + без менеджера.
func (s *storage) getOrdersForAdmin(ctx context.Context, profileID, limit, offset int32) ([]repositories.Order, error) {
	return s.queries.ListAllOrders(ctx, repositories.ListAllOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
}

// getOrdersForManager — только заказы менеджера.
func (s *storage) getOrdersForManager(ctx context.Context, managerID, limit, offset int32) ([]repositories.Order, error) {
	return s.queries.ListAllOrdersToManager(ctx, repositories.ListAllOrdersToManagerParams{
		ManagerID: pgtype.Int4{Int32: managerID, Valid: true},
		Limit:     limit,
		Offset:    offset,
	})
}

// CreateOrderWithItems создаёт заказ с элементами в транзакции.
func (s *storage) CreateOrderWithItems(
	ctx context.Context,
	managerID int32,
	dateTill time.Time,
	counterpartiesID int32,
	statusID int16,
	priority int16,
	items []repositories.AddItemToOrderParams,
) (repositories.Order, []repositories.OrderItem, error) {
	tx, err := s.beginTx(ctx)
	if err != nil {
		return repositories.Order{}, nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	// Создаём заказ
	order, err := qtx.CreateOrder(ctx, repositories.CreateOrderParams{
		DateTill:         pgtype.Timestamptz{Time: dateTill, Valid: true},
		ManagerID:        pgtype.Int4{Int32: managerID, Valid: managerID != 0},
		CounterpartiesID: counterpartiesID,
		StatusID:         statusID,
		Priority:         priority,
	})
	if err != nil {
		return repositories.Order{}, nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Добавляем элементы
	createdItems := make([]repositories.OrderItem, 0, len(items))
	for _, item := range items {
		createdItem, err := qtx.AddItemToOrder(ctx, item)
		if err != nil {
			return repositories.Order{}, nil, fmt.Errorf("failed to add order item: %w", err)
		}
		createdItems = append(createdItems, createdItem)
	}

	if err := tx.Commit(ctx); err != nil {
		return repositories.Order{}, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, createdItems, nil
}

// beginTx начинает транзакцию.
func (s *storage) beginTx(ctx context.Context) (pgx.Tx, error) {
	return s.conn.Begin(ctx)
}
