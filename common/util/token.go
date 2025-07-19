package util

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TODO: convert into env
const (
	TOKEN_SIGNATURE                = "signature-token"
	REFRESH_TOKEN_SIGNATURE        = "signature-refresh-token"
	RESET_PASSWORD_TOKEN_SIGNATURE = "signature-reset-password-token"
)

type AccessTokenClaims struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	jwt.Claims
}

type RefreshTokenClaims struct {
	Id uuid.UUID `json:"id"`
	jwt.Claims
}

type ResetPasswordTokenClaims struct {
	Email string `json:"email"`
	jwt.Claims
}

func RemoveBearer(token string) (string, bool) {
	return strings.CutPrefix(token, "Bearer ")
}

func GenerateAccessToken(c AccessTokenClaims) (string, error) {
	claims := jwt.MapClaims{
		"id":    c.Id.String(),
		"name":  c.Name,
		"email": c.Email,
		"exp":   time.Now().Add(time.Second * time.Duration(200)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(TOKEN_SIGNATURE))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(c RefreshTokenClaims) (string, error) {
	claims := jwt.MapClaims{
		"id":  c.Id.String(),
		"exp": time.Now().Add(time.Hour * time.Duration(2)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(REFRESH_TOKEN_SIGNATURE))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateResetPasswordToken(c ResetPasswordTokenClaims) (string, error) {
	claims := jwt.MapClaims{
		"email": c.Email,
		"exp":   time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(RESET_PASSWORD_TOKEN_SIGNATURE))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (AccessTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(TOKEN_SIGNATURE), nil
	})
	if err != nil {
		return AccessTokenClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		new_claims := AccessTokenClaims{}

		if id, ok := claims["id"].(string); ok {
			if id, err := uuid.Parse(id); err == nil {
				new_claims.Id = id
			} else {
				return AccessTokenClaims{}, errors.New("invalid uuid format in id claims")
			}
		} else {
			return AccessTokenClaims{}, errors.New("missing id claims")
		}

		if name, ok := claims["name"].(string); ok {
			new_claims.Name = name
		} else {
			return AccessTokenClaims{}, errors.New("missing name claims")
		}

		if email, ok := claims["email"].(string); ok {
			new_claims.Email = email
		} else {
			return AccessTokenClaims{}, errors.New("missing email claims")
		}

		return new_claims, nil
	}

	return AccessTokenClaims{}, errors.New("invalid token")
}

func ValidateRefreshToken(tokenString string) (RefreshTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(REFRESH_TOKEN_SIGNATURE), nil
	})
	if err != nil {
		return RefreshTokenClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		new_claims := RefreshTokenClaims{}

		if id, ok := claims["id"].(string); ok {
			if id, err := uuid.Parse(id); err == nil {
				new_claims.Id = id
			} else {
				return RefreshTokenClaims{}, errors.New("invalid uuid format in id claims")
			}
		} else {
			return RefreshTokenClaims{}, errors.New("missing id claims")
		}

		return new_claims, nil
	}

	return RefreshTokenClaims{}, errors.New("invalid refresh token")
}

func ValidateResetPasswordToken(tokenString string) (ResetPasswordTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(RESET_PASSWORD_TOKEN_SIGNATURE), nil
	})
	if err != nil {
		return ResetPasswordTokenClaims{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return ResetPasswordTokenClaims{
			Email: claims["email"].(string),
		}, nil
	}

	return ResetPasswordTokenClaims{}, errors.New("invalid reset password token")
}
