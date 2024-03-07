package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"template/internal/model/entities"
	"template/internal/repository"
	"template/pkg/utils"
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
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.TransactionErr, Err: err})
	}

	var userID int

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userCreate.Password), 14)

	row := tx.QueryRowContext(ctx, `INSERT INTO users (name, hashed_password) VALUES ($1, $2) RETURNING id;`,
		userCreate.Name, hashedPassword)

	err = row.Scan(&userID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ScanErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return userID, nil
}

func (u User) Get(ctx context.Context, userID int) (entities.User, error) {
	var user entities.User

	err := u.db.QueryRowContext(ctx, `SELECT id, name FROM users WHERE users.id = $1`,
		userID).Scan(&user.ID, &user.Name)

	if err != nil {
		return entities.User{}, utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
	}

	return user, nil
}

func (u User) GetHashedPassword(ctx context.Context, name string) (int, string, error) {
	var hashedPassword string
	var userID int

	err := u.db.QueryRowContext(ctx, `SELECT id, hashed_password FROM users WHERE users.name = $1`,
		name).Scan(&userID, &hashedPassword)

	if err != nil {
		return 0, "", utils.ErrNormalizer(utils.ErrorPair{Message: utils.ScanErr, Err: err})
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
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.ExecErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: err},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: err})
	}
	if count != 1 {
		err = errors.New("count error")
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.ErrNormalizer(
				utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)},
				utils.ErrorPair{Message: utils.RollbackErr, Err: rbErr},
			)
		}
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.RowsErr, Err: fmt.Errorf(utils.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return utils.ErrNormalizer(utils.ErrorPair{Message: utils.CommitErr, Err: err})
	}

	return nil
}
