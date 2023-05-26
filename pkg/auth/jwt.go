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

func CreateTokenAgent(user model.AgentTokenModel) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["payload"] = model.AgentTokenModel{
		UserId:    user.UserId,
		AccountId: user.AccountId,
		Role:      user.Role,
		SessionId: user.SessionId,
	}
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.GetConfig().ApiSecret))
}

func CreateTokenContact(user model.ContactTokenModel) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["payload"] = model.ContactTokenModel{
		SourceId: user.SourceId,
		InboxId:  user.InboxId,
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

func ExtractTokenAgent(r *http.Request) (*model.AgentTokenModel, error) {
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
		var user model.AgentTokenModel

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

func ExtractTokenContact(r *http.Request) (*model.ContactTokenModel, error) {
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
		var user model.ContactTokenModel

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
