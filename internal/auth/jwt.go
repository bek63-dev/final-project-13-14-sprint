// Package auth отвечает за генерацию и проверку JWT-токенов, которыми
// планировщик задач аутентифицирует пользователя после входа по паролю.
package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenTTL = 48 * time.Hour // - время жизни выданного токена.
)

// Claims — полезная нагрузка JWT-токена. Хранит хэш пароля, а не сам
// пароль, чтобы токен оставался валидным только до смены пароля.
type Claims struct {
	PasswordHash string `json:"password-hash"`
	jwt.RegisteredClaims
}

// hashPassword возвращает SHA-256 хэш пароля.
func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

// GenerateToken создаёт подписанный JWT-токен с хэшем переданного пароля
// в полезной нагрузке. Вызывается обработчиком входа при успешной аутентификации.
func GenerateToken(password string, secretKey string) (string, error) {
	claims := Claims{
		PasswordHash: hashPassword(password),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("GenerateToken: jwtToken.SignedString: %w", err)
	}

	return signedToken, nil
}

// ValidateToken проверяет подпись и срок действия токена, а также сверяет
// хэш пароля из токена с хэшем текущего пароля. Если пароль в TODO_PASSWORD
// изменился — токены, выданные для старого пароля, перестают быть валидными.
func ValidateToken(tokenString, password, secretKey string) bool {
	if tokenString == "" {
		return false
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("ValidateToken: неожиданный метод подписи токена: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return claims.PasswordHash == hashPassword(password)
}
