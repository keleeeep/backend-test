/*
 * @Author: Adrian Faisal
 * @Date: 14/10/21 13.44
 */

package user

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/keleeeep/test/internal/pkg/model"
	"github.com/keleeeep/test/internal/pkg/resource/db"
	"math/rand"
	"time"
)

type Usecase interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, user *model.User) (*model.TokenResponse, error)
	CheckToken(tokenString *string) (interface{}, error)
}

type usecase struct {
	dbResource db.Persistent
	secretKey string
}

func NewUsecase(dbResource db.Persistent, secretKey string) Usecase {
	return &usecase{dbResource: dbResource, secretKey: secretKey}
}

func (uc *usecase) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	existUser, err := uc.dbResource.FindUser(ctx, user.Name, "name")
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	if existUser != nil {
		if user.Name == existUser.Name {
			return nil, fmt.Errorf("username already exist")
		}
	}

	user.Password = uc.randString(4)

	ret, err := uc.dbResource.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %v", err)
	}

	return ret, nil
}

func (uc *usecase) Login(ctx context.Context, user *model.User) (*model.TokenResponse, error) {
	existUser, err := uc.dbResource.FindUser(ctx, user.Phone, "phone")
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	if existUser == nil {
		return nil, fmt.Errorf("wrong name")
	}

	if existUser.Password != user.Password {
		return nil, fmt.Errorf("wrong password")
	}

	token, err := uc.createToken(existUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %v", err)
	}

	return &model.TokenResponse{Token: token}, nil
}

func (uc *usecase) CheckToken(tokenString *string) (interface{}, error) {
	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(uc.secretKey), nil
	})
	if err != nil || token == nil {
		return nil, fmt.Errorf("failed check token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return claims, nil
}

func (uc *usecase) randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (uc *usecase) createToken(user *model.User) (string, error) {

	//Creating Access Token
	atClaims := jwt.MapClaims{
		"name":      user.Name,
		"phone":     user.Phone,
		"role":      user.Role,
		"timestamp": user.Timestamp,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(uc.secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}
