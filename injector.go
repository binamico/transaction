package transaction

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type transactionKey struct{}

// GORMInjector надстройка, которая позволяет открывать транзакцию,
// передавая контроль над результатом стороннему потребителю,
// передавая транзакцию через контекст.
type GORMInjector struct {
	db *gorm.DB
}

// NewGORMInjector создает новый экземпляр GORMInjector.
func NewGORMInjector(db *gorm.DB) *GORMInjector {
	return &GORMInjector{
		db: db,
	}
}

// Inject начинает транзакцию и запечатывает ее в контекст.
func (c *GORMInjector) Inject(ctx context.Context) (context.Context, *Manager, error) {
	if _, ok := c.ExtractGormDB(ctx); ok {
		return ctx, noopSolver(), nil
	}

	tx := c.db.Begin(&sql.TxOptions{})
	if err := tx.Error; err != nil {
		return nil, nil, err
	}

	ctx = context.WithValue(ctx, transactionKey{}, tx)
	return ctx, newSolver(tx), nil
}

// ExtractGormDB возвращает транзакцию из контеста или создает новую если
// в контексте транзакция не найдена.
func (c *GORMInjector) ExtractGormDB(ctx context.Context) (db *gorm.DB, isTx bool) {
	session := &gorm.Session{
		NewDB:   true,
		Context: ctx,
	}
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if !ok {
		// Если транзакции не было в контексте, возвращаем соединение
		// с БД изолированное от остальных транзакций.
		return c.db.Session(session), false
	}
	return tx.Session(session), ok
}
