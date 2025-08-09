package security

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func NewJWTManager(secret, accessTTL, refreshTTL string) (*JWTManager, error) {
	accessDuration, err := time.ParseDuration(accessTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid access token TTL: %w", err)
	}

	refreshDuration, err := time.ParseDuration(refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token TTL: %w", err)
	}

	return &JWTManager{
		secretKey:       secret,
		accessTokenTTL:  accessDuration,
		refreshTokenTTL: refreshDuration,
	}, nil
}

func (j *JWTManager) GenerateTokenPair(userID uint, email string) (string, string, error) {
	accessToken, err := j.generateToken(userID, email, "access", j.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.generateToken(userID, email, "refresh", j.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *JWTManager) generateToken(userID uint, email, tokenType string, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatUint(uint64(userID), 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
