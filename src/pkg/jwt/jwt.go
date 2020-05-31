package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	// ErrMissingHeader means the `Authorization` header was empty.
	ErrMissingHeader = errors.New("The length of the `Authorization` header is zero.")
)

// SignContext is the context of the JSON web token.
type SignContext struct {
	ID int64 //id
}

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

// Parse validates the token with the specified secret,
// and returns the context if the token was valid.
func Parse(tokenString string, secret string) (*SignContext, error) {
	ctx := &SignContext{}

	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))

	// Parse error.
	if err != nil {
		return ctx, err

		// Read the token if it's valid.
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.ID = int64(claims["id"].(float64))
		return ctx, nil
		// Other errors.
	} else {
		return ctx, err
	}
}

// ParseRequest gets the token from the header and
// pass it to the Parse function to parses the token.
func ParseRequest(c *gin.Context, secret string) (*SignContext, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return &SignContext{}, ErrMissingHeader
	}
	//截取加密串
	var t string
	// Parse the header to get the token part.
	if n, err := fmt.Sscanf(header, "Bearer %s", &t); err != nil || n == 0 {
		return &SignContext{}, ErrMissingHeader
	}
	return Parse(t, secret)
}

// Sign signs the context with the specified secret.
// c 需要签名的内容  secret 加密秘钥 exp有效期
func Sign(c SignContext, secret string, exp int64) (tokenString string, err error) {
	// The token content.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  c.ID,
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": exp,
	})
	// Sign the token with the specified secret.
	tokenString, err = token.SignedString([]byte(secret))
	return
}
