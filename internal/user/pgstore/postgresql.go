package pgstore

import (
	"bulletin-board/internal/ad"
	"bulletin-board/internal/user"
	"bulletin-board/pkg/postgresql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
)

type repository struct {
	client postgresql.Client
}

func (r repository) GetAll(ctx context.Context) ([]user.User, error) {
	q := `
		select id, name, birthday
		from users`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]user.User, 0)

	for rows.Next() {
		var user user.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Birthday); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r repository) GetByID(ctx context.Context, id int) (user.User, error) {
	q := `
		select id, name, birthday
		from users
		where id = $1`
	var usr user.User
	err := r.client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Name, &usr.Birthday)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

func (r repository) GetUsersAds(ctx context.Context, userId int) ([]ad.Ad, error) {
	q := `
		select id, title, description, price, user_id
		from ads
		where user_id = $1`

	ads := make([]ad.Ad, 0)

	rows, err := r.client.Query(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ad ad.Ad
		err = rows.Scan(&ad.ID, &ad.Title, &ad.Description, &ad.Price, &ad.UserID)

		if err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ads, nil
}

func (r repository) Create(ctx context.Context, newUser user.User) (user.User, error) {
	q := `
		insert into users (name, birthday, contact) 
		values ($1, $2, $3)
		returning id, name, birthday, contact`

	err := r.client.QueryRow(ctx, q, newUser.Name, newUser.Birthday, newUser.Contact).
		Scan(&newUser.ID, &newUser.Name, &newUser.Birthday, &newUser.Contact)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf("SQL error: %s, Detail: %s, Where: %s, Code: %s, SQL State: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			fmt.Println(newErr)
			return user.User{}, newErr
		}
		return user.User{}, err
	}
	return newUser, nil
}

func (r repository) Update(ctx context.Context, newUser user.User, id int) (user.User, error) {
	q := `
		update users
		set 
		    name = $1,
		    birthday = $2,
		    contact = $3
		where id = $4
		returning id, name, birthday, contact`

	err := r.client.QueryRow(ctx, q, newUser.Name, newUser.Birthday, newUser.Contact, id).
		Scan(&newUser.ID, &newUser.Name, &newUser.Birthday, &newUser.Contact)

	if err != nil {
		return user.User{}, nil
	}

	return newUser, nil
}

func (r repository) Delete(ctx context.Context, id int) error {
	q := `
		delete from users
		where id = $1`
	tag, err := r.client.Exec(ctx, q, id)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}
	return nil
}

func NewRepository(client postgresql.Client) user.Repository {
	return repository{client: client}
}
