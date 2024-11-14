package cookie

import (
	"errors"
	"net/http"

	"github.com/MomsEngineer/urlshortener/internal/adapters/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var log = logger.Create(logger.DebugLevel)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const SecretKey = "token"

func PublicCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieName := "token"
		cookies, err := c.Request.Cookie(cookieName)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			log.Error("Failed to get a request cookie", err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}

		var userID string
		if errors.Is(err, http.ErrNoCookie) {
			log.Error("No cookie", err)
			userID = uuid.NewString()
		} else if cookies != nil {
			userID, err = checkCookie(cookies.Value)
			if err != nil {
				log.Debug("Invalid cookie", err)
				userID = uuid.NewString()
			}
		}

		if userID == "" {
			log.Error("User id does not exist", nil)
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		setCookie(c, userID, cookieName)

		c.Next()
	}
}

func AuthCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieName := "token"
		cookies, err := c.Request.Cookie(cookieName)
		if err != nil {
			log.Error("Failed to get a request cookie", err)
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		userID, err := checkCookie(cookies.Value)
		if err != nil || userID == "" {
			log.Error("User id does not exist", err)
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.SetCookie(cookieName, cookies.Value, 3600, "/", "", false, true)

		c.Next()
	}
}

func setCookie(c *gin.Context, userID, cookieName string) {
	cookie, err := buildJWTString(userID)
	if err != nil {
		log.Error("Failed to create a cookie", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.SetCookie(cookieName, cookie, 3600, "/", "", false, true)
}

func checkCookie(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		log.Error("Invalid or expired token:", err)
		return "", err
	}

	return claims.UserID, nil
}

func buildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
