package persistence

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUserRepo_GetUser(t *testing.T) {
	// start test
	tests := []struct {
		name    string
		userID  string
		want    *entity.User
		wantErr error
	}{
		{
			name:    "ok",
			userID:  "id1",
			want:    &entity.User{UserID: "id1", Name: "user1", Password: "pass1"},
			wantErr: nil,
		},
		{
			name:    "ok",
			userID:  "id2",
			want:    &entity.User{UserID: "id2", Name: "user2", Password: "pass2"},
			wantErr: nil,
		},
		{
			name:    "ng: not found",
			userID:  "notFoundId",
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "ng: empty",
			userID:  "",
			want:    nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		ctx := context.Background()
		// init db
		op := &options.DeleteOptions{}
		_, _ = userRepo.col.DeleteMany(ctx, op)

		// cleanup db
		t.Cleanup(func() {
			op = &options.DeleteOptions{}
			_, _ = userRepo.col.DeleteMany(ctx, op)
		})

		// insert seeds
		seeds := []interface{}{
			entity.User{UserID: "id1", Name: "user1", Password: "pass1"},
			entity.User{UserID: "id2", Name: "user2", Password: "pass2"},
		}
		userRepo.col.InsertMany(ctx, seeds)

		t.Run(tt.name, func(t *testing.T) {
			got, err := userRepo.GetUser(ctx, tt.userID)
			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

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
		op := &options.DeleteOptions{}
		_, _ = userRepo.col.DeleteMany(ctx, op)
		// cleanup db
		t.Cleanup(func() {
			op = &options.DeleteOptions{}
			_, _ = userRepo.col.DeleteMany(ctx, op)
		})

		// insert seeds
		seeds := []interface{}{
			entity.User{UserID: "userId1", Name: "user1", Password: "pass1"},
			entity.User{UserID: "userId2", Name: "user2", Password: "pass2"},
			entity.User{UserID: "userId3", Name: "user3", Password: "pass3"},
		}
		userRepo.col.InsertMany(ctx, seeds)

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
	type args struct {
		c  context.Context
		in input.User
	}

	tests := []struct {
		name   string
		args   args
		expect func(t *testing.T, a args)
	}{
		{
			name: "ok",
			args: args{
				c: context.Background(),
				in: input.User{
					UserID:   "id1",
					Name:     "NewUser1",
					Password: "NewPassword1",
				},
			},
			expect: func(t *testing.T, a args) {
				insID, err := userRepo.CreateUser(a.c, a.in)
				assert.NoError(t, err)

				var got entity.User
				err = userRepo.col.FindOne(a.c, bson.M{"_id": insID}).Decode(&got)
				// inserted data
				assert.Equal(t, "id1", got.UserID)
				assert.NotEmpty(t, got.Password)
				assert.Equal(t, "NewUser1", got.Name)
				// hash password check
				err = bcrypt.CompareHashAndPassword([]byte(got.Password), []byte(a.in.Password))
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		// init db
		op := &options.DeleteOptions{}
		_, err := userRepo.col.DeleteMany(tt.args.c, op)
		assert.NoError(t, err)
		// cleanup db
		t.Cleanup(func() {
			op = &options.DeleteOptions{}
			_, err = userRepo.col.DeleteMany(tt.args.c, op)
			assert.NoError(t, err)
		})

		t.Run(tt.name, func(t *testing.T) {
			tt.expect(t, tt.args)
		})
	}

}
