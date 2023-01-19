package persistence

import (
	"context"
	"fmt"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepo_ListUsers(t *testing.T) {
	// start test
	tests := []struct {
		name    string
		want    []*entity.User
		wantErr error
	}{
		{
			name: "ok",
			want: []*entity.User{
				{UserID: "userId1", Name: "user1", Password: "pass1"},
				{UserID: "userId2", Name: "user2", Password: "pass2"},
				{UserID: "userId3", Name: "user3", Password: "pass3"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		// init db
		opt := &options.DeleteOptions{}
		_, _ = userRepo.database.DeleteMany(ctx, opt)
		// cleanup db
		t.Cleanup(func() {
			opt := &options.DeleteOptions{}
			_, _ = userRepo.database.DeleteMany(ctx, opt)
		})

		// insert seeds
		seeds := []interface{}{
			entity.User{UserID: "userId1", Name: "user1", Password: "pass1"},
			entity.User{UserID: "userId2", Name: "user2", Password: "pass2"},
			entity.User{UserID: "userId3", Name: "user3", Password: "pass3"},
		}
		userRepo.database.InsertMany(ctx, seeds)

		t.Run(tt.name, func(t *testing.T) {
			got, err := userRepo.ListUsers(ctx)
			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

func TestUserRepo_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   input.User
		want    entity.User
		wantErr error
	}{
		{
			name: "ok",
			input: input.User{
				Name:     "NewUser1",
				Password: "NewPassword1",
			},
			want: entity.User{
				Name: "NewUser1",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		// init db
		opt := &options.DeleteOptions{}
		_, _ = userRepo.database.DeleteMany(ctx, opt)
		// cleanup db
		t.Cleanup(func() {
			opt := &options.DeleteOptions{}
			_, _ = userRepo.database.DeleteMany(ctx, opt)
		})

		t.Run(tt.name, func(t *testing.T) {
			insertID, err := userRepo.CreateUser(ctx, tt.input)
			assert.NoError(t, err)

			var got entity.User
			err = userRepo.database.FindOne(ctx, bson.M{"_id": insertID}).Decode(&got)
			fmt.Println(got.Password)

			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			opt := cmpopts.IgnoreFields(entity.User{}, "Password")
			if diff := cmp.Diff(tt.want, got, opt); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

			// hash password check
			err = bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(tt.input.Password))
			assert.NoError(t, err)

		})
	}

}
