package transaction

import "gorm.io/gorm"

// Manager управляет транзакцией.
type Manager struct {
	rollback func() error
	commit   func() error
}

// Rollback откатывает транзакцию.
func (s *Manager) Rollback() error {
	return s.rollback()
}

// Commit завершает транзакцию.
func (s *Manager) Commit() error {
	return s.commit()
}

// newGORMSolver создает экземпляр Manager,
// связанный с переданным объектом транзакции GORM.
func newGORMSolver(tx *gorm.DB) *Manager {
	return &Manager{
		rollback: func() error {
			return tx.Rollback().Error
		},
		commit: func() error {
			return tx.Commit().Error
		},
	}
}

// noopSolver создает экземпляр Manager,
// который не выполняет никаких операций при Rollback и Commit.
func noopSolver() *Manager {
	return &Manager{
		rollback: func() error {
			return nil
		},
		commit: func() error {
			return nil
		},
	}
}
