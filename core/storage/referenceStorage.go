package storage

import (
	"backend/core/repositories"
	"context"
)

// ReferenceStorage — хранилище для справочников.
type ReferenceStorage interface {
	Storage
	ListCounterparties(ctx context.Context) ([]repositories.Counterparty, error)
	ListNomenclatures(ctx context.Context) ([]repositories.Nomenclature, error)
	ListMaterials(ctx context.Context) ([]repositories.Material, error)
	ListSizes(ctx context.Context) ([]repositories.Size, error)
	ListStatuses(ctx context.Context) ([]repositories.Status, error)
	ListRoles(ctx context.Context) ([]repositories.Role, error)
}

// ListCounterparties возвращает всех контрагентов.
func (s *storage) ListCounterparties(ctx context.Context) ([]repositories.Counterparty, error) {
	return s.queries.ListAllCounterparties(ctx)
}

// ListNomenclatures возвращает всю номенклатуру.
func (s *storage) ListNomenclatures(ctx context.Context) ([]repositories.Nomenclature, error) {
	return s.queries.ListAllNomenclatures(ctx)
}

// ListMaterials возвращает все материалы.
func (s *storage) ListMaterials(ctx context.Context) ([]repositories.Material, error) {
	return s.queries.ListAllMaterials(ctx)
}

// ListSizes возвращает все размеры.
func (s *storage) ListSizes(ctx context.Context) ([]repositories.Size, error) {
	return s.queries.ListAllSizes(ctx)
}

// ListStatuses возвращает все статусы.
func (s *storage) ListStatuses(ctx context.Context) ([]repositories.Status, error) {
	return s.queries.ListAllStatuses(ctx)
}

// ListRoles возвращает все роли.
func (s *storage) ListRoles(ctx context.Context) ([]repositories.Role, error) {
	return s.queries.ListAllRoles(ctx)
}
