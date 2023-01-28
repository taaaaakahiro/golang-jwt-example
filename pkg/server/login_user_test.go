package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func TestServer_GetLoginUser(t *testing.T) {
	ctx := context.Background()
	// init db
	database := getDatabase()
	db := database.Collection(collection)
	opt := &options.DeleteOptions{}
	_, _ = db.DeleteMany(ctx, opt)
	// cleanup db
	t.Cleanup(func() {
		opt := &options.DeleteOptions{}
		_, _ = db.DeleteMany(ctx, opt)
	})

	pwHash, err := bcrypt.GenerateFromPassword([]byte("NewPassword1"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// insert seeds
	seeds := []interface{}{
		entity.User{UserID: "userId1", Name: "user1", Password: string(pwHash)},
	}
	db.InsertMany(ctx, seeds)

	t.Run("ok", func(t *testing.T) {
		accessToken := getAccessToken()
		data := input.LoginInfo{
			LoginID:  "userId1",
			Password: "NewPassword1",
		}
		d, err := json.Marshal(data)
		assert.NoError(t, err)
		buf := bytes.NewBuffer(d)
		req, err := http.NewRequest("GET", testServer.URL+`/user/auth`, buf)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		client := &http.Client{}
		res, err := client.Do(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		res.Body.Close()
		assert.NotEmpty(t, body)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		assert.NotEmpty(t, string(body))

		var user entity.User
		err = json.Unmarshal(body, &user)
		assert.NoError(t, err)
		assert.Equal(t, "userId1", user.UserID)
		assert.NotEmpty(t, user.Password)
		assert.Equal(t, "user1", user.Name)

	})
}

func getAccessToken() string {
	c := context.Background()
	pwHash, _ := bcrypt.GenerateFromPassword([]byte("NewPassword1"), bcrypt.DefaultCost)

	// insert seeds
	seeds := []interface{}{
		entity.User{UserID: "userId1", Name: "user1", Password: string(pwHash)},
	}
	database := getDatabase()
	db := database.Collection(collection)

	db.InsertMany(c, seeds)
	data := input.LoginInfo{
		LoginID:  "userId1",
		Password: "NewPassword1",
	}
	d, _ := json.Marshal(data)
	buf := bytes.NewBuffer(d)
	req, _ := http.NewRequest("POST", testServer.URL+`/user/login`, buf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, _ := client.Do(req)

	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	var aToken string
	_ = json.Unmarshal(body, &aToken)

	return aToken
}
