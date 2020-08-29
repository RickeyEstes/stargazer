package crawler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type databaseClient struct {
	db *mongo.Database
}

func (c databaseClient) init() error {
	coStargazers := c.db.Collection("stargazers")
	coUsers := c.db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := coStargazers.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"data.user.login": -1},
	}); err != nil {
		return errors.WithStack(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := coUsers.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"login": -1},
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c databaseClient) getRepository(path string) (*repository, error) {
	co := c.db.Collection("repositories")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var r repository
	if err := co.FindOne(ctx, bson.M{"path": path}).Decode(&r); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return &r, nil
}

func (c databaseClient) insertRepository(r *repository) error {
	co := c.db.Collection("repositories")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.ID = primitive.NewObjectID()
	_, err := co.InsertOne(ctx, r)
	return errors.WithStack(err)
}

func (c databaseClient) updateRepository(r *repository) error {
	co := c.db.Collection("repositories")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := co.UpdateOne(ctx, bson.M{"_id": r.ID}, bson.M{"$set": r})
	return errors.WithStack(err)
}

func (c databaseClient) countStargazers(repositoryID primitive.ObjectID) (int, error) {
	co := c.db.Collection("stargazers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := co.CountDocuments(ctx, bson.M{"_repository_id": repositoryID})
	return int(count), errors.WithStack(err)
}

func (c databaseClient) getStargazers() ([]stargazer, error) {
	co := c.db.Collection("stargazers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := co.Find(ctx, bson.M{}, &options.FindOptions{
		Sort: bson.M{"data.starred_at": -1},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var ss []stargazer
	for cur.Next(context.Background()) {
		var s stargazer
		if err := cur.Decode(&s); err != nil {
			return nil, errors.WithStack(err)
		}
		ss = append(ss, s)
	}

	return ss, nil
}

func (c databaseClient) deleteStargazers(repositoryID primitive.ObjectID, page int) error {
	co := c.db.Collection("stargazers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := co.DeleteMany(ctx, bson.M{"_repository_id": repositoryID, "page": page})
	return errors.WithStack(err)
}

func (c databaseClient) insertStargazers(ss []stargazer) error {
	co := c.db.Collection("stargazers")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := range ss {
		ss[i].ID = primitive.NewObjectID()
		if _, err := co.InsertOne(ctx, ss[i]); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c databaseClient) getUser(login string) (*user, error) {
	co := c.db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u user
	if err := co.FindOne(ctx, bson.M{"login": login}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return &u, nil
}

func (c databaseClient) insertUser(u *user) error {
	co := c.db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u.ID = primitive.NewObjectID()
	_, err := co.InsertOne(ctx, u)
	return errors.WithStack(err)
}

func (c databaseClient) updateUser(u *user) error {
	co := c.db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := co.UpdateOne(ctx, bson.M{"_id": u.ID}, bson.M{"$set": u})
	return errors.WithStack(err)
}
