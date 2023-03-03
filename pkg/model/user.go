package model

import (
	"context"
	"fmt"

	"example.com/pkg/domain"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	userStruct = sqlbuilder.NewStruct(new(user))
	userTable  = "users"
)

type user struct {
	ID          pgtype.Int4 `db:"id"`
	Username    pgtype.Text `db:"username" fieldtag:"credentials"`
	Password    pgtype.Text `db:"passwrd" fieldtag:"credentials"`
	FirstName   pgtype.Text `db:"firstname" fieldtag:"details"`
	LastName    pgtype.Text `db:"lastname" fieldtag:"details"`
	Email       pgtype.Text `db:"email" fieldtag:"details"`
	Status      pgtype.Text `db:"status" fieldtag:"details"`
	Description pgtype.Text `db:"description" fieldtag:"details"`
}

func userToUserRepo(u *domain.User) user {
	res := user{}
	res.ID.Scan(int64(u.ID))
	res.Username.Scan(u.Username)
	res.Password.Scan(u.Password)
	res.FirstName.Scan(u.FirstName)
	res.LastName.Scan(u.LastName)
	res.Email.Scan(u.Email)
	res.Status.Scan(u.Status)
	res.Description.Scan(u.Description)
	return res
}

func userRepoToUser(u *user) domain.User {
	return domain.User{
		ID:          u.ID.Int32,
		Username:    u.Username.String,
		Password:    u.Password.String,
		FirstName:   u.FirstName.String,
		LastName:    u.LastName.String,
		Email:       u.Email.String,
		Status:      u.Status.String,
		Description: u.Description.String,
	}
}

func (db *UserDB) CreateUser(ctx context.Context, user *domain.User) (int32, error) {
	u := userToUserRepo(user)
	sb := userStruct.WithTag("credentials", "details").InsertInto(userTable, u)
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	sql += " RETURNING id"
	row := db.pool.QueryRow(ctx, sql, args...)
	var id int32
	err := row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error scanning returning id: %w", err)
	}
	return id, nil
}

func (db *UserDB) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	sb := userStruct.SelectFrom(userTable)
	sb.Where(sb.Equal("username", username))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	var user user
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(userStruct.Addr(&user)...)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	res := userRepoToUser(&user)
	return &res, nil
}

func (db *UserDB) GetUserById(ctx context.Context, id int32) (*domain.User, error) {
	sb := userStruct.SelectFrom(userTable)
	sb.Where(sb.Equal("id", id))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	var user user
	row := db.pool.QueryRow(ctx, sql, args...)
	err := row.Scan(userStruct.Addr(&user)...)
	if err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}
	res := userRepoToUser(&user)
	return &res, nil
}

func (db *UserDB) DeleteUserByUsername(ctx context.Context, username string) error {
	sb := userStruct.DeleteFrom(userTable)
	sb.Where(sb.Equal("username", username))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing delete query: %w", err)
	}
	return nil
}

func (db *UserDB) UpdateUserDetails(ctx context.Context, user *domain.User) error {
	u := userToUserRepo(user)
	sb := userStruct.WithTag("details").Update(userTable, u)
	sb.Where(sb.Equal("id", user.ID))
	sql, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("error executing update query: %w", err)
	}
	return err
}
