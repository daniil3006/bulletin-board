package pgstore

import (
	"bulletin-board/internal/ad"
	"bulletin-board/pkg/postgresql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
}

func (r repository) GetAll(ctx context.Context) ([]ad.Ad, error) {
	q := `
		select id, title, description, price, user_id 
		from ads`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ads := make([]ad.Ad, 0)

	for rows.Next() {
		var ad ad.Ad

		err = rows.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.Price, &ad.UserID)
		if err != nil {
			return nil, err
		}

		ads = append(ads, ad)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ads, nil
}

func (r repository) GetByID(ctx context.Context, ID int) (ad.Ad, error) {
	q := `
		select id, title, description, price, user_id 
		from ads 
		where id = $1`
	var returnedAd ad.Ad
	err := r.client.QueryRow(ctx, q, ID).Scan(&returnedAd.ID, &returnedAd.Title, &returnedAd.Description, &returnedAd.Price, &returnedAd.UserID)
	if err != nil {
		return ad.Ad{}, err
	}
	return returnedAd, nil
}

func (r repository) Create(ctx context.Context, newAd ad.Ad) (ad.Ad, error) {
	q := `
		insert into ads (title, description, price, user_id) 
		values ($1, $2, $3, $4)
		returning id, title, description, price, user_id`
	err := r.client.QueryRow(ctx, q, newAd.Title, newAd.Description, newAd.Price, newAd.UserID).
		Scan(&newAd.ID, &newAd.Title, &newAd.Description, &newAd.Price, &newAd.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL error: %s, Detail: %s, Where: %s, Code: %s, SQL State: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			fmt.Println(newErr)
			return ad.Ad{}, newErr
		}
		return ad.Ad{}, err
	}
	return newAd, nil
}

func (r repository) Update(ctx context.Context, newAd ad.Ad, id int) (ad.Ad, error) {
	q := `
		update ads
		set
			title = $1,
			description = $2,
			price = $3,
			user_id = $4
		where id = $5
		returning id, title, description, price, user_id`
	err := r.client.QueryRow(ctx, q, newAd.Title, newAd.Description, newAd.Price, newAd.UserID, id).
		Scan(&newAd.ID, &newAd.Title, &newAd.Description, &newAd.Price, &newAd.UserID)
	if err != nil {
		return ad.Ad{}, err
	}
	return newAd, nil
}

func (r repository) Delete(ctx context.Context, id int) error {
	q := `
		delete from ads
		where id = $1`
	tag, err := r.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ad.ErrNotFound
	}

	return nil
}

func NewRepository(client postgresql.Client) ad.Repository {
	return &repository{client: client}
}
