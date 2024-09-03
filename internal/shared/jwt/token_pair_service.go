package jwt

import (
	"errors"
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	jwt "gitlab.com/golib4/jwt-resolver/jwt"
	"time"
)

type TokenPairService struct {
	params      Params
	jwtResolver *jwt.JwtResolver
}

func NewTokenPairService(
	params Params,
	jwtResolver *jwt.JwtResolver,
) *TokenPairService {
	return &TokenPairService{
		params:      params,
		jwtResolver: jwtResolver,
	}
}

type Params struct {
	AccessTokenLifetimeInMinutes  int
	RefreshTokenLifetimeInMinutes int
}

type AccessTokenClaims struct {
	UserId int64 `json:"user_id"`
	jwt2.StandardClaims
}

type RefreshTokenClaims struct {
	UserId int64 `json:"user_id"`
	jwt2.StandardClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func (s TokenPairService) GenerateTokenPair(userId int64) (*TokenPair, error) {
	accessTokenData := AccessTokenClaims{
		UserId: userId,
		StandardClaims: jwt2.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(s.params.AccessTokenLifetimeInMinutes) * time.Minute).Unix(),
			Issuer:    "authorization",
		},
	}
	accessToken, err := s.jwtResolver.Create(accessTokenData)
	if err != nil {
		return nil, fmt.Errorf("can not create access token: %s", err)
	}

	refreshTokenData := RefreshTokenClaims{
		UserId: userId,
		StandardClaims: jwt2.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(s.params.RefreshTokenLifetimeInMinutes) * time.Minute).Unix(),
			Issuer:    "authorization",
		},
	}
	refreshToken, err := s.jwtResolver.Create(refreshTokenData)
	if err != nil {
		return nil, fmt.Errorf("can not create refresh token: %s", err)
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, err
}

func (s TokenPairService) CheckAccessToken(accessToken string) (*AccessTokenClaims, error) {
	claims := AccessTokenClaims{}
	err := s.jwtResolver.Parse(accessToken, &claims)
	if err != nil {
		return nil, fmt.Errorf("can not parse token: %s", err)
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("token is expired")
	}

	return &claims, nil
}

func (s TokenPairService) CheckRefreshToken(refreshToken string) (*RefreshTokenClaims, error) {
	claims := RefreshTokenClaims{}
	err := s.jwtResolver.Parse(refreshToken, &claims)
	if err != nil {
		return nil, fmt.Errorf("can not parse token: %s", err)
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, errors.New("token is expired")
	}

	return &claims, nil
}
