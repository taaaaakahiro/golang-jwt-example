package persistence

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
	"testing"

	"github.com/google/go-cmp/cmp"
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
