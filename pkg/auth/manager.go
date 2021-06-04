package auth

import (
	"net/http"
	"strings"

	"github.com/JesusG2000/hexsatisfaction_purchase/pkg/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const authorizationHeader = "Authorization"

// TokenManager provides logic for a JWT token generation and parsing.
type TokenManager interface {
	NewJWT(userID string) (string, error)
	Parse(accessToken string) (string, error)
	UserIdentity(next http.Handler) http.Handler
}

// Manager manages a JWT token.
type Manager struct {
	signingKey string
}

// NewManager is a Manager constructor.
func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty secret key")
	}

	return &Manager{signingKey: signingKey}, nil
}

// NewJWT creates a new JWT token.
func (m *Manager) NewJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject: userID,
	})

	return token.SignedString([]byte(m.signingKey))
}

// Parse parses the JWT token.
func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(err, "unexpected signing method")
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "couldn't parse token")
	}

	subClaims, ok := token.Claims.(jwt.MapClaims)["sub"]
	if !ok {
		return "", errors.New("empty claims")
	}
	return subClaims.(string), nil
}

// UserIdentity checks validation of the token.
func (m *Manager) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorizationHeader)
		if header == "" {
			middleware.JSONError(w, errors.New("empty auth header"), http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			middleware.JSONError(w, errors.New("invalid auth header"), http.StatusUnauthorized)
			return
		}
		_, err := m.Parse(headerParts[1])
		if err != nil {
			middleware.JSONError(w, err, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
