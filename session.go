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

// SessionDB управляет доступом к сессионному соединению с БД.
type SessionDB interface {
	DB(context.Context) *gorm.DB
}

// SessionAdapter реализует интерфейс Session и инкапсулирует взаимодействие
// с зависимостями, зависящими от сессии.
type SessionAdapter struct {
	injector *GORMInjector
}

// NewSessionAdapter создает новый инстанс SessionAdapter.
func NewSessionAdapter(db *gorm.DB) *SessionAdapter {
	return &SessionAdapter{
		injector: NewGORMInjector(db),
	}
}

// Begin запускает сессию для зависимых от сессии объектов.
func (s *SessionAdapter) Begin(
	ctx context.Context,
) (context.Context, SessionSolver, error) {
	return s.injector.Inject(ctx)
}

// DB возвращает соединение с БД из сессии.
func (s *SessionAdapter) DB(ctx context.Context) *gorm.DB {
	db, _ := s.injector.ExtractGormDB(ctx)
	return db
}
