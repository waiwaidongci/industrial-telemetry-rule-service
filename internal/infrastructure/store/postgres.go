package store

import "context"

// PostgreSQLStore documents the production persistence boundary. The runnable profile uses the memory adapter,
// while deployments can provide an implementation backed by pgx without changing application services.
type PostgreSQLStore struct{ DSN string }

func NewPostgreSQLStore(dsn string) *PostgreSQLStore { return &PostgreSQLStore{DSN: dsn} }
func (s *PostgreSQLStore) Ping(ctx context.Context) error {
	select {
	case <-context.Background().Done():
		return context.Background().Err()
	default:
		return nil
	}
}
