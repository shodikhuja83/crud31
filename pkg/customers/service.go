package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

)

var ErrNotFound = errors.New("item not found")

var ErrInternal = errors.New("internal error")

type Service struct {
	pool  *pgxpool.Pool
}
func NewService(pool  *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Customer struct {
	ID      	int64     `json:"id"`
	Name    	string    `json:"name"`
	Phone   	string    `json:"phone"`
	Password 	string	  `json:"password"`
	Active		bool 	  `json:"active"`
	Created 	time.Time `json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}

	err := s.pool.QueryRow(ctx, `
		SELECT id, name, phone, active, created FROM customers WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
func (s *Service) All(ctx context.Context) (items []*Customer, err error) {

	rows, err:= s.pool.Query(ctx, `
		SELECT * FROM customers
	`)

	for rows.Next(){
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
		if err != nil {
			log.Print(err)
		}

		items = append(items, item)
	}
	return items, nil
}
func (s *Service) AllActive(ctx context.Context) (items []*Customer, err  error) {

	rows, err:= s.pool.Query(ctx, `
		SELECT * FROM customers WHERE active
	`)

	for rows.Next(){
		item := &Customer{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
		if err != nil {
			log.Print(err)
		}

		items = append(items, item)
	}
	return items, nil
}

// //Save method
func (s *Service) Save(ctx context.Context, customer *Customer) (c *Customer, err error) {

	item := &Customer{}

	if customer.ID == 0 {
		sqlStatement := `insert into customers(name, phone, password) values($1, $2, $3) returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.Password).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Password, 
			&item.Active, 
			&item.Created)
	} else {
		sqlStatement := `update customers set name=$1, phone=$2, password=$3 where id=$4 returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, customer.Name, customer.Phone, customer.Password, customer.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Password, 
			&item.Active, 
			&item.Created)
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

func (s *Service) RemoveById(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
	DELETE FROM customers WHERE id=$1 RETURNING id,name,phone,active,created 
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}

func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = false WHERE id = $1 RETURNING id, name, phone, active, created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
func (s *Service) UnBlockByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
		UPDATE customers SET active = true WHERE id = $1 RETURNING id, name, phone, active, created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
