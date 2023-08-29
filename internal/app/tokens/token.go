package tokens

import "github.com/golang-jwt/jwt"

// MySignKey token signature
var MySignKey = []byte("super-secret-auth-key")

// GenerateToken creates an authorization token for the user
func GenerateToken(email, password string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(MySignKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
