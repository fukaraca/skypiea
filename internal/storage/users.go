package storage

import (
	"context"
	"time"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/encryption"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	addUserPG = `INSERT INTO users(user_uuid,firstname,lastname,email,role,status,password)
					VALUES ($1,$2,$3,$4,$5,$6,$7);`
	getUserByUUIDPG  = `SELECT * FROM users WHERE user_uuid = $1;`
	getUserByEmailPG = `SELECT * FROM users WHERE email = $1;`
	getPassPG        = `SELECT password FROM users WHERE email = $1;`
	updatePasswordPG = `UPDATE users SET password = $1 where user_uuid = $2;` //nolint: gosec
)

type UsersRepo interface {
	AddUser(context.Context, *User) (*uuid.UUID, error)
	GetUserByUUID(context.Context, uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetHPassword(context.Context, string) (string, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, hPassword string) error
}

type User struct {
	ID        int
	UserUUID  string
	Firstname string
	Lastname  string
	Email     string
	Role      string
	Status    string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time // TODO: maybe not needed for now but nullable types matter for pgx>> pgtype..
}

func (u *User) Convert() *model.User {
	return &model.User{
		ID:        u.ID,
		UUID:      uuid.MustParse(u.UserUUID),
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}

type usersRepoPgx struct {
	*pgxpool.Pool
}

func NewUsersRepo(db *DB) UsersRepo {
	switch db.Dialect {
	case DialectPostgres, DialectPgx:
		return &usersRepoPgx{db.Pool}
	}
	return nil
}

func (u *usersRepoPgx) AddUser(ctx context.Context, user *User) (*uuid.UUID, error) {
	hashed, err := encryption.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed
	uid := uuid.New()

	_, err = u.Exec(ctx, addUserPG, uid.String(),
		user.Firstname, user.Lastname, user.Email, user.Role, user.Status, user.Password)
	if err != nil {
		return nil, err
	}
	// TODO check conflict
	return &uid, nil
}

func (u *usersRepoPgx) GetUserByUUID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var out User
	row := u.QueryRow(ctx, getUserByUUIDPG, userID.String())
	if err := row.Scan(&out.ID, &out.UserUUID, &out.Firstname, &out.Lastname, &out.Email,
		&out.Password, &out.Role, &out.Status, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

func (u *usersRepoPgx) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var out User
	row := u.QueryRow(ctx, getUserByEmailPG, email)
	if err := row.Scan(&out.ID, &out.UserUUID, &out.Firstname, &out.Lastname, &out.Email,
		&out.Password, &out.Role, &out.Status, &out.CreatedAt, &out.UpdatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

func (u *usersRepoPgx) GetHPassword(ctx context.Context, username string) (string, error) {
	var out string
	row := u.QueryRow(ctx, getPassPG, username)
	return out, row.Scan(&out)
}

func (u *usersRepoPgx) ChangePassword(ctx context.Context, userID uuid.UUID, hPassword string) error {
	_, err := u.Exec(ctx, updatePasswordPG, hPassword, userID.String())
	return err
}
