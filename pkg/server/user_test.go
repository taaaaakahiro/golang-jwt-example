package server

import (
	"context"
	"encoding/json"
	"golang-jwt-example/pkg/domain/entity"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collection = "users"

func Test_Server_ListUsers(t *testing.T) {
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
		entity.User{ID: 1, Name: "user1", Password: "pass1"},
		entity.User{ID: 2, Name: "user2", Password: "pass2"},
		entity.User{ID: 3, Name: "user3", Password: "pass3"},
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
		assert.Equal(t, 1, got[0].ID)
		assert.Equal(t, "user1", got[0].Name)
		assert.Equal(t, "pass1", got[0].Password)

	})

}
