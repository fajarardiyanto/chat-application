package auth

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"net/http"
	"strings"
	"time"
)

func CreateToken(user model.TokenModel) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["payload"] = model.TokenModel{
		UserId:    user.UserId,
		AccountId: user.AccountId,
		Role:      user.Role,
		SessionId: user.SessionId,
	}
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.GetConfig().ApiSecret))

}

func TokenValid(r *http.Request) error {
	tokenString := ExtractToken(r)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().ApiSecret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(r *http.Request) (*model.TokenModel, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().ApiSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		var user model.TokenModel

		userMarshal, err := json.Marshal(claims["payload"].(map[string]interface{}))
		if err != nil {
			config.GetLogger().Error(err)
			return nil, err
		}

		if err = json.Unmarshal(userMarshal, &user); err != nil {
			config.GetLogger().Error(err)
			return nil, err
		}

		return &user, nil
	}
	return nil, nil
}
