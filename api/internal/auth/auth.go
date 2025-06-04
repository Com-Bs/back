package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	// 1) Parse y verificar firma + algoritmo
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Asegurarnos de que el método de firmado sea HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firmado inesperado: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("falló verificación de firma o token malformado: %w", err)
	}

	// 2) jwt.Parse ya revisa exp/nbf internamente, pero confirmamos que token.Valid sea true
	if !token.Valid {
		return "", fmt.Errorf("token inválido o expirado")
	}

	// 3) Para asegurarnos explícitamente de 'exp', podemos hacer un chequeo extra:
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("no se pudieron leer los claims")
	}
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return "", fmt.Errorf("token expirado")
		}
	} else {
		return "", fmt.Errorf("claim 'exp' no presente o con formato incorrecto")
	}

	// Extract username from claims
	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in token claims")
	}

	return username, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
