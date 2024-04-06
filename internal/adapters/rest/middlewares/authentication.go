package middlewares

import (
	"brick/internal/pkg/constvar"
	"brick/internal/pkg/serror"
	"brick/internal/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	Log *logrus.Logger
}

func NewAuthMiddleware(log *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		Log: log,
	}
}

func (m *AuthMiddleware) VerifyAuth(authType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := m.extractToken(ctx)
		if tokenString == "" {
			m.Log.Errorf("[Auth][Middleware] while ExtractToken: %s", "token nout found")
			serror.AbortWithSerror(ctx, http.StatusUnauthorized, 1, constvar.UNAUTHORIZED_ERROR, "authorization token not provided")
			return
		}

		jwtSecret, err := m.getJwtSecret(authType)
		if err != nil {
			m.Log.Errorf("[Auth][Middleware] while getJwtSecret: %s", "authType not defined")
			serror.AbortWithSerror(ctx, http.StatusUnauthorized, 1, constvar.UNAUTHORIZED_ERROR, err.Error())
			return
		}

		verifiedToken, err := m.VerifyToken(tokenString, jwtSecret)
		if err != nil {
			m.Log.Errorf("[Auth][Middleware] while VerifyToken: %s", "token not valid")
			serror.AbortWithSerror(ctx, http.StatusUnauthorized, 1, constvar.UNAUTHORIZED_ERROR, err.Error())
			return
		}

		if claims, ok := verifiedToken.Claims.(jwt.MapClaims); ok && verifiedToken.Valid {
			if m.isTokenExpired(claims) {
				m.Log.Errorf("[Auth][Middleware] while check Token: %s", "invalid or expired token")
				serror.AbortWithSerror(ctx, http.StatusUnauthorized, 1, constvar.UNAUTHORIZED_ERROR, "invalid or expired token")
				return
			}
			m.setContextClaims(ctx, claims, authType)
			ctx.Next()
			return
		}
		m.Log.Errorf("[Auth][Middleware] while VerifyAuth: %s", "invalid or expired token")
		serror.AbortWithSerror(ctx, http.StatusUnauthorized, 1, constvar.UNAUTHORIZED_ERROR, "invalid or expired token")
	}
}

func (m *AuthMiddleware) getJwtSecret(authType string) (string, error) {
	if authType == "common" {
		return utils.ReadStringEnvKey("COMMON_SECRET", true)
	}
	return utils.ReadStringEnvKey("API_KEY_SECRET", true)
}

func (m *AuthMiddleware) setContextClaims(ctx *gin.Context, claims jwt.MapClaims, authType string) {
	userId := claims["user_id"].(float64)
	strUserId := strconv.FormatFloat(userId, 'f', -1, 64)
	ctx.Set("user_id", strUserId)
	if authType != "common" {
		clientID := claims["iss"].(string)
		ctx.Set("client_id", clientID)
	}
}

func (m *AuthMiddleware) VerifyToken(tokenString string, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer ")
	}
	return ""
}

func (m *AuthMiddleware) isTokenExpired(claims jwt.MapClaims) bool {
	if exp, ok := claims["exp"].(float64); ok {
		return time.Unix(int64(exp), 0).Before(time.Now())
	}
	return true
}
