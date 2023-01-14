package persistence

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUserReposirory_ListUsers(t *testing.T) {
	tests := []struct {
		name    string
		want    []*entity.User
		wantErr error
	}{
		{
			name: "ok",
			want: []*entity.User{
				{ID: 1, Name: "hoge", Password: "pass1"},
				{ID: 2, Name: "fuga", Password: "pass2"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userRepo.ListUsers(context.Background())
			if diff := cmp.Diff(tt.wantErr, err); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want, got); len(diff) != 0 {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}

		})

	}

}
