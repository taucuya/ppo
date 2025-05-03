package auth_prov

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Provider struct {
	key       []byte
	aduration time.Duration
	rduration time.Duration
}

func New(key []byte, adur time.Duration, rdur time.Duration) *Provider {
	return &Provider{key: key, aduration: adur, rduration: rdur}
}

func (p *Provider) GenToken(ctx context.Context, id uuid.UUID) (atoken string, rtoken string, err error) {
	accessClaims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(p.aduration).Unix(),
		"iat": time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	atoken, err = accessToken.SignedString(p.key)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(p.rduration).Unix(),
		"iat": time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	rtoken, err = refreshToken.SignedString(p.key)
	if err != nil {
		return "", "", err
	}

	return atoken, rtoken, nil
}

func (p *Provider) VerifyToken(ctx context.Context, token string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return p.key, nil
	}, jwt.WithLeeway(5*time.Minute))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return true, nil
		}
		return false, err
	}

	if !parsedToken.Valid {
		return false, jwt.ErrTokenUnverifiable
	}

	return false, nil
}

func (p *Provider) RefreshToken(ctx context.Context, atoken string, rtoken string) (string, error) {
	refreshToken, err := jwt.Parse(rtoken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return p.key, nil
	}, jwt.WithLeeway(5*time.Minute))
	if err != nil || !refreshToken.Valid {
		return "", errors.New("invalid refresh token")
	}

	refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || refreshClaims["id"] == nil {
		return "", errors.New("invalid refresh token claims")
	}
	id := uuid.UUID(refreshClaims["id"].(uuid.UUID))

	accessToken, _ := jwt.Parse(atoken, func(t *jwt.Token) (interface{}, error) {
		return p.key, nil
	}, jwt.WithLeeway(5*time.Minute))

	if accessToken != nil && accessToken.Valid {
		accessClaims, ok := accessToken.Claims.(jwt.MapClaims)
		if ok && accessClaims["id"] != nil && accessClaims["id"].(uuid.UUID) != uuid.UUID(id) {
			return "", errors.New("token user mismatch")
		}
	}

	newAccessToken, _, err := p.GenToken(ctx, id)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
