package service

import (
	"context"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/server/middlewares"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/encryption"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/google/uuid"
)

func (s *Service) SignIn(ctx context.Context, email, password string) (*session.Cookie, error) {
	var user *storage.User
	var err error
	err = s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		user, err = reg.Users.GetUserByEmail(ctx, email)
		return err
	})
	if err != nil {
		return nil, err
	}

	if !encryption.CheckPasswordHash(password, user.Password) {
		return nil, model.ErrIncorrectCred
	}
	sess := session.Cache.NewSession(ctx, uuid.MustParse(user.UserUUID))
	session.Cache.Set(sess)
	return session.NewCookie(sess.ID), nil
}

func (s *Service) RegisterNewUser(ctx context.Context, user *storage.User) error {
	return s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		_, errInner := reg.Users.AddUser(ctx, user)
		return errInner
	})
}

func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*storage.User, error) {
	var user *storage.User
	var err error
	err = s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		user, err = reg.Users.GetUserByUUID(ctx, userID)
		return err
	})
	return user, err
}

func (s *Service) ChangePassword(ctx context.Context, email, newPass string) error {
	return s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		u, err := reg.Users.GetUserByEmail(ctx, email)
		if err != nil {
			return err
		}
		hPass, err := encryption.HashPassword(newPass)
		if err != nil {
			return err
		}
		return reg.Users.ChangePassword(ctx, uuid.MustParse(u.UserUUID), hPass)
	})
}
