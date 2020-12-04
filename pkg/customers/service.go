package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

//ErrNotFound ...
var ErrNotFound = errors.New("item not found")

//ErrInternal ...
var ErrInternal = errors.New("internal error")

//Service ..
type Service struct {
	pool *pgxpool.Pool
}

//NewService ..
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

//Customer ...
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

//All ....
func (s *Service) All(ctx context.Context) (cs []*Customer, err error) {

	sqlStatement := `select * from customers`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

//AllActive ....
func (s *Service) AllActive(ctx context.Context) (cs []*Customer, err error) {

	sqlStatement := `select * from customers where active=true`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
		}
		cs = append(cs, item)
	}

	return cs, nil
}

//ByID ...
func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	sqlStatement := `select * from customers where id=$1`
	err := s.pool.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//ChangeActive ...
func (s *Service) ChangeActive(ctx context.Context, id int64, active bool) (*Customer, error) {
	item := &Customer{}

	sqlStatement := `update customers set active=$2 where id=$1 returning *`
	err := s.pool.QueryRow(ctx, sqlStatement, id, active).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//Delete ...
func (s *Service) Delete(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	sqlStatement := `delete from customers  where id=$1 returning *`
	err := s.pool.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

//Save ...
func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error) {

	item := &Customer{}

	if customer.ID == 0 {
		sqlStatement := `insert into customers(name, phone) values($1, $2) returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
	} else {
		sqlStatement := `update customers set name=$1, phone=$2 where id=$3 returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}
