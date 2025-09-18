package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

var (
	stateStore              = map[string]stateMeta{} // stores state and verifier, lets keep it simple for now
	ErrStateMismatch        = errors.New("state mismatches or expired")
	ErrAuthenticationFailed = errors.New("authentication failed")
)

type stateMeta struct {
	verifier  string
	createdAt time.Time
}

func (s *Service) Start(ctx context.Context, isSignUp bool) string {
	state := newState()
	verifier := oauth2.GenerateVerifier()
	challenge := oauth2.S256ChallengeFromVerifier(verifier)
	stateStore[state] = stateMeta{verifier: verifier, createdAt: time.Now()}
	startURL := s.Oauth.Google.AuthCodeURL(state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"))
	return startURL
}

func (s *Service) Callback(ctx context.Context, code, state string) (*session.Cookie, error) {
	defer delete(stateStore, state)
	v, ok := stateStore[state]
	if !ok || state == "" || v.createdAt.Before(time.Now().Add(-time.Minute)) {
		return nil, ErrStateMismatch
	}
	exchangeToken, err := s.Oauth.Google.Exchange(ctx, code, oauth2.VerifierOption(v.verifier))
	if err != nil {
		return nil, errors.Join(err, ErrAuthenticationFailed)
	}
	tokenString := exchangeToken.Extra("id_token")
	token, err := validateToken(tokenString.(string))
	if err != nil {
		return nil, err
	}
	user, err := s.Repositories.Users.GetUserByEmail(ctx, token.Email)
	if err == nil {
		sess := session.Cache.NewSession(ctx, uuid.MustParse(user.UserUUID))
		session.Cache.Set(sess)
		return session.NewCookie(sess.ID), nil
	}
	if errors.Is(err, model.ErrNoSuchEmail) {
		errIn := s.RegisterNewUser(ctx, &storage.User{
			Firstname: token.GivenName,
			Lastname:  token.Surname,
			Email:     token.Email,
			Role:      model.RoleUserStd,
			Status:    model.StatusNew,
			Picture:   pgtype.Text{String: token.Picture, Valid: token.Picture != ""},
			AuthType:  model.AuthTypeOauth2,
			Password:  newState(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if errIn != nil {
			return nil, errIn
		}
		user2, errIn := s.Repositories.Users.GetUserByEmail(ctx, token.Email)
		if errIn != nil {
			return nil, errIn
		}
		sess := session.Cache.NewSession(ctx, uuid.MustParse(user2.UserUUID))
		session.Cache.Set(sess)
		return session.NewCookie(sess.ID), nil
	}
	return nil, err
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func b64URLEncode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func newState() string {
	return b64URLEncode(randBytes(32))
}

func validateToken(tokenString string) (*gwt.Token, error) {
	jwksURL := "https://www.googleapis.com/oauth2/v3/certs"
	jwksGoogle, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, err
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
	token, err := parser.ParseWithClaims(tokenString, &gwt.Token{}, jwksGoogle.Keyfunc)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*gwt.Token)
	if !ok {
		return nil, errors.New("token doesn't match")
	}
	return claims, nil
}
