package security

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

//Service ...
type Service struct {
	pool *pgxpool.Pool
}

//NewService ...
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

//Auth ...
func (s *Service) Auth(login, password string) bool {
	sql := `select login, password from managers where login=$1 and password=$2`
	err := s.pool.QueryRow(context.Background(), sql, login, password).Scan(&login, &password)
	log.Print(login, password)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}
