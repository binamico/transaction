package transaction

import (
	"context"

	"gorm.io/gorm"
)

// Session управляет запуском сессии на application уровне.
type Session interface {
	Begin(context.Context) (ctx context.Context, solver SessionSolver, err error)
}

// SessionSolver управляет результатом транзакции в рамках сессии.
type SessionSolver interface {
	Rollback() error
	Commit() error
}

// GORMSessionDB управляет доступом к сессионному соединению с БД.
type GORMSessionDB interface {
	DB(context.Context) *gorm.DB
}

// GORMSessionAdapter реализует интерфейс Session и инкапсулирует взаимодействие
// с зависимостями, зависящими от сессии.
type GORMSessionAdapter struct {
	injector *GORMInjector
}

// NewGORMSessionAdapter создает новый инстанс SessionAdapter.
func NewGORMSessionAdapter(db *gorm.DB) *GORMSessionAdapter {
	return &GORMSessionAdapter{
		injector: NewGORMInjector(db),
	}
}

// Begin запускает сессию для зависимых от сессии объектов.
func (s *GORMSessionAdapter) Begin(
	ctx context.Context,
) (context.Context, SessionSolver, error) {
	return s.injector.Inject(ctx)
}

// DB возвращает соединение с БД из сессии.
func (s *GORMSessionAdapter) DB(ctx context.Context) *gorm.DB {
	db, _ := s.injector.ExtractGormDB(ctx)
	return db
}
