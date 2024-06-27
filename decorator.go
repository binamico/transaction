package transaction

import (
	"context"
	"fmt"

	"github.com/binamico/logging"
	"go.opentelemetry.io/otel/trace"
)

// Decorator обертка для выполнения операций в рамках транзакции.
type Decorator struct {
	session Session
	logger  *logging.Logger
	tracer  trace.Tracer
}

// NewDecorator возвращает новый экземпляр Decorator.
func NewDecorator(
	session Session,
	l *logging.Logger,
) *Decorator {
	return &Decorator{
		session: session,
		logger:  l,
	}
}

// RunInTx выполняет функцию f в транзакции.
func (d *Decorator) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	ctx, span := d.tracer.Start(ctx, "decorator.RunInTx")
	defer span.End()

	ctx, tx, err := d.session.Begin(ctx)
	if err != nil {
		d.logger.WithError(err).Error("cannot begin transaction")
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	if err = f(ctx); err != nil {
		if err = tx.Rollback(); err != nil {
			d.logger.WithError(err).Error("cannot rollback transaction")
		}
		return fmt.Errorf("execute function in transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		d.logger.WithError(err).Error("cannot commit transaction")
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	return nil
}
