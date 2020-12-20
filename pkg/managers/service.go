package managers

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/SsSJKK/crud/cmd/app/middleware"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

//ErrInternal ...
var ErrInternal = errors.New("internal error")

// ErrExpireToken ...
var ErrExpireToken = errors.New("ExpireToken error")

//Service ...
type Service struct {
	pool *pgxpool.Pool
}

//NewService ..
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

//Managers ...
type Managers struct {
	ID       int64     `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
	Roles    []string  `json:"roles"`
	Active   bool      `json:"active"`
	Created  time.Time `json:"created"`
}

//Product ...
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
}

////Sales ...
//type Sales struct {
//	ID         int64     `json:"id"`
//	ManagerID  int64     `json:"manager_id"`
//	CustomerID int64     `json:"customer_id"`
//	Created    time.Time `json:"created"`
//}
//
////SalePositions ...
//type SalePositions struct {
//	ID        int64     `json:"id"`
//	SaleID    int64     `json:"sale_id"`
//	ProductID int64     `json:"product_id"`
//	Price     int64     `json:"price"`
//	Qty       int64     `json:"qty"`
//	Created   time.Time `json:"created"`
//}

//SalePositions ...
type SalePositions struct {
	ID         int64 `json:"id"`
	CustomerID int64 `json:"customer_id"`
	Positions  []struct {
		ID        int64 `json:"id"`
		ProductID int64 `json:"product_id"`
		Qty       int64 `json:"qty"`
		Price     int64 `json:"price"`
	} `json:"positions"`
}

//Registration ...
func (s *Service) Registration(ctx context.Context, item *Managers) (*Managers, error) {
	sql := `INSERT INTO managers (	name,	phone,	password,	roles  ) VALUES ($1, $2,$3,	$4) returning *`

	manager := &Managers{}
	var id int64
	idM, err := middleware.Authentication(ctx)
	if err != nil {
		return nil, err
	}
	err = s.pool.QueryRow(ctx, `select id from managers WHERE roles[2] = 'ADMIN' and id = $1`, idM).Scan(&id)
	if err == pgx.ErrNoRows {
		log.Println("No priv", err)
		return nil, err
	}

	if err != nil {
		log.Println("No priv", err)
		return nil, err
	}

	err = s.pool.QueryRow(ctx, sql, item.Name, item.Phone, item.Password, item.Roles).Scan(
		&manager.ID,
		&manager.Name,
		&manager.Phone,
		&manager.Password,
		&manager.Roles,
		&manager.Active,
		&manager.Created,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return manager, nil

}

//IDByToken ...
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {
	var id int64
	var expire time.Time
	err := s.pool.QueryRow(ctx,
		`select manager_id, expire from managers_tokens where token = $1`, token).Scan(&id, &expire)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	timeNow := time.Now().Format("2006-01-02 15:04:05")
	timeEnd := expire.Format("2006-01-02 15:04:05")

	if timeNow > timeEnd {
		return 0, ErrExpireToken
	}

	if err != nil {
		return 0, ErrInternal
	}

	return id, nil
}

//ChangeProduct ...
func (s *Service) ChangeProduct(ctx context.Context, item *Product) (*Product, error) {
	product := &Product{}
	if item.ID == 0 {
		sql := `INSERT INTO products (name, price, qty)
		VALUES (
			$1,
			$2,
			$3)
		 returning id, name, price, qty`
		err := s.pool.QueryRow(ctx, sql, item.Name, item.Price, item.Qty).Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Qty,
		)
		if err != nil {
			return nil, err
		}
		return product, nil
	}
	sql := `UPDATE products 
	set
	name = $1,
	price = $2,
	qty = $3
	where id = $4 returning id, name, price, qty`
	err := s.pool.QueryRow(ctx, sql, item.Name, item.Price, item.Qty, item.ID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Qty,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

//MakeSele ...
func (s *Service) MakeSele(ctx context.Context, saleP *SalePositions, idManager int64) error {
	var idSale int64
	sqlSales := `INSERT INTO sales (manager_id, customer_id)
	VALUES (
		$1,
		$2
	  ) RETURNING id;`
	err := s.pool.QueryRow(ctx, sqlSales, idManager, saleP.CustomerID).Scan(&idSale)

	if err != nil {
		return err
	}
	sqlSalePositions := `INSERT INTO sale_positions (sale_id, product_id, price, qty)
	VALUES (
		$1,
		$2,
		$3,
		$4
	  );`
	sqlUpdate := `UPDATE products
	SET
	qty = $1
	WHERE id=$2`
	for _, v := range saleP.Positions {
		var qty int64
		err = s.pool.QueryRow(ctx, `select qty from products where id = $1`, v.ProductID).Scan(&qty)
		if err != nil {
			return err
		}
		_, err = s.pool.Exec(ctx, sqlUpdate, qty-v.Qty, v.ProductID)
		if err != nil {
			return err
		}
		_, err = s.pool.Exec(ctx, sqlSalePositions, idSale, v.ProductID, v.Price, v.Qty)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	return nil
}

//GetSales ...
func (s *Service) GetSales(ctx context.Context, id int64) (int64, error) {
	var getSales struct {
		id  int64
		sum int64
	}
	sql := `SELECT s.manager_id, sum(sp.qty * sp.price) FROM sales s
	JOIN sale_positions as sp on sp.sale_id = s.id
	where s.manager_id = $1
	GROUP by s.manager_id;`

	err := s.pool.QueryRow(ctx, sql, id).Scan(&getSales.id, &getSales.sum)
	if err != nil {
		return 0, err
	}
	return getSales.sum, nil

}
