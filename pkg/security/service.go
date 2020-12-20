package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

//Service ...
type Service struct {
	pool *pgxpool.Pool
}

// ErrNoSuchUser ...
var ErrNoSuchUser = errors.New("no such user")

// ErrInvalidPassword ...
var ErrInvalidPassword = errors.New("invalid password")

// ErrInternal ...
var ErrInternal = errors.New("internal error")

// ErrExpireToken ...
var ErrExpireToken = errors.New("ExpireToken error")

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

//TokenForCustomer ...
func (s *Service) TokenForCustomer(
	ctx context.Context,
	phone string,
	password string,
) (token string, err error) {
	var hash string
	var id int64
	sql := `Select id, password FROM customers where phone = $1`
	err = s.pool.QueryRow(ctx, sql, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}
	token = hex.EncodeToString(buffer)
	sqlInsert := `INSERT INTO customers_tokens (token, customer_id) VALUES ($1, $2)`
	a, err := s.pool.Exec(ctx, sqlInsert, token, id)
	log.Print(a)
	return token, nil
}

//AuthenticateCustomer ...
func (s *Service) AuthenticateCustomer(ctx context.Context, tkn string) (int64, error) {
	var id int64
	var expire time.Time
	err := s.pool.QueryRow(ctx, `select customer_id, expire from customers_tokens where token=$1`, tkn).Scan(&id, &expire)
	if err == pgx.ErrNoRows {
		log.Println("1")
		return 0, ErrNoSuchUser
	}
	if err != nil {
		log.Println("2")
		return 0, ErrInternal
	}

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	timeEnd := expire.Format("2006-01-02 15:04:05")

	if timeNow > timeEnd {
		return 0, ErrExpireToken
	}

	return id, nil
}

//TokenForManagers ...
func (s *Service) TokenForManager(
	ctx context.Context,
	phone string,
	password string,
) (token string, err error) {
	var hash string
	var id int64
	sql := `Select id, password FROM managers where phone = $1`
	err = s.pool.QueryRow(ctx, sql, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}
	token = hex.EncodeToString(buffer)
	sqlInsert := `INSERT INTO managers_tokens (token, manager_id) VALUES ($1, $2)`
	a, err := s.pool.Exec(ctx, sqlInsert, token, id)
	log.Print(a)
	return token, nil
}

//AuthenticateCustomer ...
func (s *Service) AuthenticateManagers(ctx context.Context, tkn string) (int64, error) {
	var id int64
	var expire time.Time
	err := s.pool.QueryRow(ctx, `select manager_id, expire from managers_tokens where token=$1`, tkn).Scan(&id, &expire)
	if err == pgx.ErrNoRows {
		log.Println("1")
		return 0, ErrNoSuchUser
	}
	if err != nil {
		log.Println("2")
		return 0, ErrInternal
	}

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	timeEnd := expire.Format("2006-01-02 15:04:05")

	if timeNow > timeEnd {
		return 0, ErrExpireToken
	}

	return id, nil
}

func (s *Service) TokenWithOut(ctx context.Context, id int64) (string, error) {
	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", err
	}
	token := hex.EncodeToString(buffer)
	sqlInsert := `INSERT INTO managers_tokens (token, manager_id) VALUES ($1, $2)`
	a, err := s.pool.Exec(ctx, sqlInsert, token, id)
	if err != nil {
		log.Println(a, err)
		return "", err
	}
	return token, nil
}
