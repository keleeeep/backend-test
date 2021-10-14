/*
 * @Author: Fredy Gunawan
 * @Date: 14/10/21 13.48
 */

package fetch

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/keleeeep/test/internal/pkg/model"
	"github.com/keleeeep/test/internal/pkg/resource/db"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Usecase interface {
	Fetch(tokeString *string) ([]*model.Data, error)
	Aggregate(tokenString *string) ([]*model.Data, error)
	CheckToken(tokenString *string) (interface{}, error)
}

type usecase struct {
	dbResource db.Persistent
	secretKey  string
}

func NewUsecase(dbResource db.Persistent, secretKey string) Usecase {
	return &usecase{dbResource: dbResource, secretKey: secretKey}
}

func (uc *usecase) Fetch(tokenString *string) ([]*model.Data, error) {
	c := cache.New(5*time.Minute, 10*time.Minute)

	checkToken, err := uc.CheckToken(tokenString)
	if err != nil || checkToken == nil {
		return nil, fmt.Errorf("forbidden access to this resource on the server is denied")
	}

	res, err := http.Get("https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list") //use package "net/http"
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var m []*model.Data
	if err := json.Unmarshal(body, &m); err != nil {
		log.Printf("failed to unmarshal: %v", err)
		return nil, fmt.Errorf("failed to unmarshal data")
	}

	for _, v := range m {
		key := fmt.Sprintf("uuid:%s", v.Uuid)

		data, found := c.Get(key)
		if found {
			var n model.Data
			bodyBytes, _ := json.Marshal(data)
			if err = json.Unmarshal(bodyBytes, &n); err != nil {
				return nil, fmt.Errorf("failed to unmarshal data")
			}
			v = &n
			continue
		}

		if v.Price != "" && v.Price != "dfhdfh" && v.Price != "A" { //i have no idea why the data type of price is string ini here
			i, err := strconv.ParseFloat(v.Price, 64)
			if err != nil {
				log.Printf("invalid etd time: %v", err)
				return nil, fmt.Errorf("failed to convert from string to integer")
			}

			//TODO: do the conversion here, i have registered to get api key but still error "invalid API Key"
			v.Usd = i / 14100.45
		}

		c.Set(key, v, 5*time.Minute)
	}

	return m, nil
}

func (uc *usecase) Aggregate(tokenString *string) ([]*model.Data, error) {
	//not enough time
	return nil, nil
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
