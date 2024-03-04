package user

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"template/internal/model/entities"
	"template/internal/repository"
	"template/pkg/customerr"
)

type User struct {
	db *sqlx.DB
}

func InitUserRepo(db *sqlx.DB) repository.User {
	return User{
		db: db,
	}
}

func (u User) Create(ctx context.Context, userCreate entities.UserCreate) (int, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var userID int

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userCreate.Password), 14)

	row := tx.QueryRowContext(ctx, `INSERT INTO users (name, hashed_password) VALUES ($1, $2) RETURNING id;`,
		userCreate.Name, hashedPassword)

	err = row.Scan(&userID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.UserRollbackError
		}
		return 0, customerr.UserScanError
	}

	err = tx.Commit()
	if err != nil {
		return 0, customerr.UserCommitError
	}

	return userID, nil
}

func (u User) Get(ctx context.Context, userID int) (entities.User, error) {
	var user entities.User

	err := u.db.QueryRowContext(ctx, `SELECT id, name FROM users WHERE users.id = $1`,
		userID).Scan(&user.ID, &user.Name)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (u User) GetHashedPassword(ctx context.Context, name string) (int, string, error) {
	var hashedPassword string
	var userID int

	err := u.db.QueryRowContext(ctx, `SELECT id, hashed_password FROM users WHERE users.name = $1`,
		name).Scan(&userID, &hashedPassword)

	if err != nil {
		return 0, "", err
	}

	return userID, hashedPassword, nil
}

func (u User) Delete(ctx context.Context, userID int) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM users WHERE users.id = $1;`, userID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return err
		}
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return err
		}
		return err
	}
	if count != 1 {
		err = errors.New("count error")
		if rbErr := tx.Rollback(); rbErr != nil {
			return err
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
