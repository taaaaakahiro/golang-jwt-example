package persistence

import (
	"context"
	"fmt"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
				{ID: 1, Name: "user1", Password: "pass1"},
				{ID: 2, Name: "user2", Password: "pass2"},
				{ID: 3, Name: "user3", Password: "pass3"},
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
			entity.User{ID: 1, Name: "user1", Password: "pass1"},
			entity.User{ID: 2, Name: "user2", Password: "pass2"},
			entity.User{ID: 3, Name: "user3", Password: "pass3"},
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
				Name:     "user1",
				Password: "pass1",
			},
			want:    entity.User{Name: "user1", Password: "pass1"},
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
			if err != nil {
				log.Fatal(err)
			}

			var got entity.User
			err = userRepo.database.FindOne(ctx, bson.M{"_id": insertID}).Decode(&got)
			fmt.Println(got)

			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

		})
	}

}
