package protection

import (
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/serror"
	"brick/internal/pkg/utils"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateSalt() (string, *serror.Error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", serror.NewError(500, 1, constvar.SERVER_INFO_ERROR, err.Error())
	}
	return hex.EncodeToString(salt), nil
}

func HashPassword(password, salt string) string {
	sha256Hasher := sha256.New()
	sha256Hasher.Write([]byte(salt))
	sha256Hasher.Write([]byte(password))
	return hex.EncodeToString(sha256Hasher.Sum(nil))
}

func ComparePassword(storedHash, storedSalt, providedPassword string) bool {
	hash := HashPassword(providedPassword, storedSalt)

	return hash == storedHash
}

func GenerateToken(userID int) (string, *serror.Error) {
	jwtExpirationTime, err := utils.ReadIntEnvKey("COMMON_EXPIRY_TIME", true)
	if err != nil {
		return "", serror.NewError(500, 1, constvar.SERVER_INFO_ERROR, err.Error())
	}

	jwtDuration := time.Duration(jwtExpirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * jwtDuration).Unix(),
	})

	jwtSecret, err := utils.ReadStringEnvKey("COMMON_SECRET", true)
	if err != nil {
		return "", serror.NewError(500, 1, constvar.SERVER_INFO_ERROR, err.Error())
	}
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", serror.NewError(500, 1, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return tokenString, nil
}

func GenerateApiToken(clientID string) (string, *serror.Error) {
	apiKeyExpirationValue, err := utils.ReadIntEnvKey("API_KEY_EXPIRY_TIME", true)
	if err != nil {
		return "", serror.NewError(500, 1, constvar.SERVER_INFO_ERROR, err.Error())
	}
	apiKeyExpirationDuration := time.Duration(apiKeyExpirationValue)
	expirationTime := time.Now().Add(time.Minute * apiKeyExpirationDuration)

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	apiKeySecret, err := utils.ReadStringEnvKey("API_KEY_SECRET", true)

	if err != nil {
		return "", serror.NewError(500, 1, constvar.SERVER_INFO_ERROR, err.Error())
	}
	tokenString, err := token.SignedString([]byte(apiKeySecret))
	if err != nil {
		return "", serror.NewError(500, 1, constvar.SERVER_INFO_ERROR, err.Error())
	}

	return tokenString, nil
}
