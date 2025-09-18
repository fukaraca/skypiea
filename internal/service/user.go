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

func (s *Service) GetAllUsers(ctx context.Context) ([]*storage.User, error) {
	return s.Repositories.Users.GetAllUsers(ctx)
}

func (s *Service) ChangePassword(ctx context.Context, email, newPass string) error {
	return s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		u, err := reg.Users.GetUserByEmail(ctx, email)
		if err != nil {
			return err
		}
		if u.AuthType == model.AuthTypeOauth2 && newPass == "iForgotMyPassword" {
			return model.ErrOauth2UserDoesNotForget
		}
		if u.Role == model.RoleAdmin && newPass == "iForgotMyPassword" {
			return model.ErrRealAdminDoesNotForget
		}
		hPass, err := encryption.HashPassword(newPass)
		if err != nil {
			return err
		}
		return reg.Users.ChangePassword(ctx, uuid.MustParse(u.UserUUID), hPass)
	})
}

func (s *Service) UpdateRole(ctx context.Context, userUUID, role string) error {
	return s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		return reg.Users.ChangeRole(ctx, uuid.MustParse(userUUID), role)
	})
}

func (s *Service) SupportedModels(ctx context.Context, userID uuid.UUID) []string {
	// TODO user tier based supported model...
	return s.GeminiClient.GetAllSupportedModels()
}

func (s *Service) UpdateUserProfile(ctx context.Context, userNew *storage.User) error {
	userOld, err := s.Repositories.Users.GetUserByUUID(ctx, uuid.MustParse(userNew.UserUUID))
	if err != nil {
		return err
	}
	userOld.Firstname = userNew.Firstname
	userOld.Lastname = userNew.Lastname
	userOld.PhoneNumber = userNew.PhoneNumber
	userOld.AboutMe = userNew.AboutMe

	return s.Repositories.DoInTx(ctx, middlewares.GetLoggerFromContext(ctx), func(reg *storage.Registry) error {
		return reg.Users.UpdateUser(ctx, userOld)
	})
}

func (s *Service) GetAdoptionStatistics(ctx context.Context) ([]*storage.AdoptionStat, error) {
	return s.Repositories.Users.GetAdoptionStatistics(ctx)
}
