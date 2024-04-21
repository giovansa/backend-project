package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) RegisterUser(ctx context.Context, input RegisterUser) (string, error) {
	res, err := r.Db.ExecContext(ctx, qInsertUser, input.ID, input.Phone, input.Name, input.Password)
	if err != nil {
		return "", err
	}

	if count, _ := res.RowsAffected(); count < 1 {
		return "", errors.New("fail register")
	}

	return input.ID, nil
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	stmt, err := r.Db.PrepareContext(ctx, qGetUserByPhone)
	if err != nil {
		return User{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(phone).
		Scan(&user.UserID, &user.Phone, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("user not exists")
		}
		return User{}, err
	}
	return user, err
}

func (r *Repository) IncrSuccessLogin(ctx context.Context, phone string) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(ctx, qIncrementLoginCount, phone)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("login tracking failed")
	}

	return tx.Commit()
}

func (r *Repository) UpdateUser(ctx context.Context, input UpdateUser, identifier string) error {
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	column := 1
	query := `UPDATE users SET `

	updateCol := make([]string, 0)
	valueUpdate := make([]interface{}, 0)
	if len(strings.TrimSpace(input.Phone)) > 0 {
		valueUpdate = append(valueUpdate, input.Phone)
		updateCol = append(updateCol, fmt.Sprintf("phone = $%d", column))
		column++
	}
	if len(strings.TrimSpace(input.Name)) > 0 {
		valueUpdate = append(valueUpdate, input.Name)
		updateCol = append(updateCol, fmt.Sprintf("name = $%d", column))
		column++
	}

	valueUpdate = append(valueUpdate, time.Now())
	updateCol = append(updateCol, fmt.Sprintf("updated_at = $%d", column))
	column++

	query += strings.Join(updateCol, ",")
	query += fmt.Sprintf(" WHERE phone = $%d", column)

	valueUpdate = append(valueUpdate, identifier)
	res, err := tx.ExecContext(ctx, query, valueUpdate...)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected < 1 {
		return errors.New("update user error")
	}

	return tx.Commit()
}
