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
			return false, jwt.ErrTokenExpired
		}
		return false, err
	}

	if !parsedToken.Valid {
		return false, jwt.ErrTokenUnverifiable
	}
	return true, nil
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

	idStr, ok := refreshClaims["id"].(string)
	if !ok {
		return "", errors.New("invalid id in refresh token")
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return "", errors.New("malformed UUID in refresh token")
	}

	accessToken, _ := jwt.Parse(atoken, func(t *jwt.Token) (interface{}, error) {
		return p.key, nil
	}, jwt.WithLeeway(5*time.Minute))

	if accessToken != nil && accessToken.Valid {
		accessClaims, ok := accessToken.Claims.(jwt.MapClaims)
		if ok && accessClaims["id"] != nil {
			accessIDStr, ok := accessClaims["id"].(string)
			if !ok {
				return "", errors.New("invalid id in access token")
			}
			accessID, err := uuid.Parse(accessIDStr)
			if err != nil {
				return "", errors.New("malformed UUID in access token")
			}
			if accessID != id {
				return "", errors.New("token user mismatch")
			}
		}
	}

	newAccessToken, _, err := p.GenToken(ctx, id)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (p *Provider) ExtractUserID(tokenStr string) (uuid.UUID, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if idStr, ok := claims["id"].(string); ok {
			return uuid.Parse(idStr)
		}
	}

	return uuid.Nil, jwt.ErrTokenMalformed
}
