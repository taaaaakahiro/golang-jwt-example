package server

import (
	"bytes"
	"context"
	"encoding/json"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const collection = "users"

func TestServer_ListUsers(t *testing.T) {
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

	// insert seeds
	seeds := []interface{}{
		entity.User{UserID: "userId1", Name: "user1", Password: "pass1"},
		entity.User{UserID: "userId2", Name: "user2", Password: "pass2"},
		entity.User{UserID: "userId3", Name: "user3", Password: "pass3"},
	}
	db.InsertMany(ctx, seeds)

	t.Run("All", func(t *testing.T) {
		var buf io.Reader
		req, err := http.NewRequest("GET", testServer.URL+`/user/all`, buf)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		res.Body.Close()
		assert.NotEmpty(t, body)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		var got []*entity.User
		err = json.Unmarshal(body, &got)
		assert.NoError(t, err)

		assert.Len(t, got, 3)
		assert.Equal(t, "userId1", got[0].UserID)
		assert.Equal(t, "user1", got[0].Name)
		assert.Equal(t, "pass1", got[0].Password)

		assert.Equal(t, "userId2", got[1].UserID)
		assert.Equal(t, "user2", got[1].Name)
		assert.Equal(t, "pass2", got[1].Password)

		assert.Equal(t, "userId3", got[2].UserID)
		assert.Equal(t, "user3", got[2].Name)
		assert.Equal(t, "pass3", got[2].Password)
	})
}
func TestServer_Login(t *testing.T) {
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
		data := input.LoginInfo{
			LoginID:  "userId1",
			Password: "NewPassword1",
		}
		d, err := json.Marshal(data)
		assert.NoError(t, err)
		buf := bytes.NewBuffer(d)
		req, err := http.NewRequest("POST", testServer.URL+`/user/login`, buf)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

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

	})
}
