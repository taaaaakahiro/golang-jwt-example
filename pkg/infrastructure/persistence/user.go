package persistence

import (
	"context"
	"golang-jwt-example/pkg/domain/entity"
	"golang-jwt-example/pkg/domain/input"
	"golang-jwt-example/pkg/domain/repository"
	"log"

	errs "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const collection = "users"

type UserRepository struct {
	col *mongo.Collection
}

var _ repository.IUserRepository = (*UserRepository)(nil)

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		col: db.Collection(collection),
	}
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*entity.User, error) {
	var user entity.User
	flt := bson.D{
		primitive.E{Key: "user_id", Value: userID},
	}
	opt := options.FindOne()
	if err := r.col.FindOne(ctx, flt, opt).Decode(&user); err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, nil
		default:
			return nil, errs.WithStack(err)
		}
	}

	return &user, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]*entity.User, error) {
	users := make([]*entity.User, 0)
	srt := bson.D{
		primitive.E{Key: "_id", Value: 1},
	}
	opt := options.Find().SetSort(srt)
	flt := bson.D{
		primitive.E{},
	}
	cur, err := r.col.Find(ctx, flt, opt)
	if err != nil {
		return nil, errs.WithStack(err)
	}

	if err = cur.All(ctx, &users); err != nil {
		log.Printf("error %+v", err)
		return nil, errs.WithStack(err)
	}

	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (r *UserRepository) CreateUser(ctx context.Context, in input.User) (interface{}, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.WithStack(err)
	}

	data := input.User{
		UserID:   in.UserID,
		Name:     in.Name,
		Password: string(password),
	}
	//opts := options.InsertOneOptions{}
	id, err := r.col.InsertOne(ctx, data, nil)
	if err != nil {
		return nil, errs.WithStack(err)
	}

	return id.InsertedID, nil
}
