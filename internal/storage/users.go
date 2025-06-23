package storage

import (
	"context"
	"errors"
	"time"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/encryption"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	addUserPG = `INSERT INTO users(user_uuid,firstname,lastname,email,role,status,password, phone_number)
					VALUES ($1,$2,$3,$4,$5,$6,$7,$8);`
	updateUserPG = `UPDATE users SET firstname = $1, lastname = $2, status = $3, 
                 	phone_number = $4, about_me = $5, summary = $6, updated_at = NOW()
					WHERE id = $7;`
	getUserByUUIDPG  = `SELECT * FROM users WHERE user_uuid = $1;`
	getUsersPG       = `SELECT * FROM users ORDER BY updated_at DESC;`
	getUserByEmailPG = `SELECT * FROM users WHERE email = $1;`
	getPassPG        = `SELECT password FROM users WHERE email = $1;`
	updatePasswordPG = `UPDATE users SET password = $1, updated_at = NOW() where user_uuid = $2;` //nolint: gosec
	updateUserRolePG = `UPDATE users SET role = $1, updated_at = NOW() where user_uuid = $2;`     //nolint: gosec
	deleteUserByUUID = `DELETE FROM users WHERE user_uuid = $1;`
)

// TODO: this is conventional and not needed right now. It is better created and implemented on caller side
type UsersRepo interface {
	AddUser(context.Context, *User) (uuid.UUID, error)
	GetUserByUUID(context.Context, uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetHPassword(ctx context.Context, username string) (string, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, hPassword string) error
	DeleteUsersByUUID(context.Context, uuid.UUID) error
	UpdateUser(ctx context.Context, user *User) error
	ChangeRole(ctx context.Context, userID uuid.UUID, role string) error
	GetAdoptionStatistics(ctx context.Context) ([]*AdoptionStat, error)
}

type User struct {
	ID          int
	UserUUID    string
	Firstname   string
	Lastname    string
	Email       string
	PhoneNumber pgtype.Text
	AboutMe     pgtype.Text
	Summary     pgtype.Text
	Role        string
	Status      string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time // TODO: maybe not needed for now but nullable types matter for pgx>> pgtype..
}

func (u *User) Convert() *model.User {
	return &model.User{
		ID:          u.ID,
		UUID:        uuid.MustParse(u.UserUUID),
		Firstname:   u.Firstname,
		Lastname:    u.Lastname,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber.String,
		AboutMe:     u.AboutMe.String,
		Summary:     u.Summary.String,
		Role:        u.Role,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
	}
}

type usersRepoPgx struct {
	dbConn
}

func NewUsersRepo(dia Dialect, conn dbConn) *usersRepoPgx {
	switch dia {
	case DialectPostgres, DialectPgx:
		return &usersRepoPgx{conn}
	}
	return nil
}

func (u *usersRepoPgx) AddUser(ctx context.Context, user *User) (uuid.UUID, error) {
	hashed, err := encryption.HashPassword(user.Password)
	if err != nil {
		return uuid.UUID{}, err
	}
	user.Password = hashed
	uid := uuid.New()

	_, err = u.Exec(ctx, addUserPG, uid.String(),
		user.Firstname, user.Lastname, user.Email, user.Role, user.Status, user.Password, user.PhoneNumber)
	if err != nil {
		return uuid.UUID{}, err
	}
	// TODO check conflict
	return uid, nil
}

func (u *usersRepoPgx) GetUserByUUID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var out User
	row := u.QueryRow(ctx, getUserByUUIDPG, userID.String())
	if err := row.Scan(&out.ID, &out.UserUUID, &out.Firstname, &out.Lastname, &out.Email,
		&out.Password, &out.Role, &out.Status, &out.CreatedAt, &out.UpdatedAt, &out.PhoneNumber, &out.AboutMe, &out.Summary); err != nil {
		return nil, err
	}
	return &out, nil
}

func (u *usersRepoPgx) GetAllUsers(ctx context.Context) ([]*User, error) {
	out := make([]*User, 0)
	rows, err := u.Query(ctx, getUsersPG)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	for rows.Next() {
		var t User
		err = rows.Scan(&t.ID, &t.UserUUID, &t.Firstname, &t.Lastname, &t.Email, &t.Password, &t.Role, &t.Status, &t.CreatedAt, &t.UpdatedAt, &t.PhoneNumber, &t.AboutMe, &t.Summary)
		if err != nil {
			return nil, err
		}
		out = append(out, &t)
	}

	return out, nil
}

func (u *usersRepoPgx) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var out User
	row := u.QueryRow(ctx, getUserByEmailPG, email)
	if err := row.Scan(&out.ID, &out.UserUUID, &out.Firstname, &out.Lastname, &out.Email,
		&out.Password, &out.Role, &out.Status, &out.CreatedAt, &out.UpdatedAt, &out.PhoneNumber, &out.AboutMe, &out.Summary); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNoSuchEmail
		}
		return nil, err
	}
	return &out, nil
}

func (u *usersRepoPgx) GetHPassword(ctx context.Context, email string) (string, error) {
	var out string
	row := u.QueryRow(ctx, getPassPG, email)
	return out, row.Scan(&out)
}

func (u *usersRepoPgx) ChangePassword(ctx context.Context, userID uuid.UUID, hPassword string) error {
	_, err := u.Exec(ctx, updatePasswordPG, hPassword, userID.String())
	return err
}

func (u *usersRepoPgx) ChangeRole(ctx context.Context, userID uuid.UUID, role string) error {
	_, err := u.Exec(ctx, updateUserRolePG, role, userID.String())
	return err
}

func (u *usersRepoPgx) UpdateUser(ctx context.Context, user *User) error {
	_, err := u.Exec(ctx, updateUserPG, user.Firstname, user.Lastname, user.Status, user.PhoneNumber, user.AboutMe, user.Summary, user.ID)
	return err
}

func (u *usersRepoPgx) DeleteUsersByUUID(ctx context.Context, userID uuid.UUID) error {
	_, err := u.Exec(ctx, deleteUserByUUID, userID.String())
	return err
}
