package docs_repo

import (
	"context"
	"rest_with_mongo/db/docs"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DocsRepo struct {
	collection *mongo.Collection
}

func NewDocsRepo(db *mongo.Database) *DocsRepo {
	return &DocsRepo{collection: db.Collection("docs")}
}

func (repo *DocsRepo) Create(ctx context.Context, doc *docs.Doc) (*docs.Doc, error) {
	
	_, err := repo.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (repo *DocsRepo) GetByAuthor(ctx context.Context, id primitive.ObjectID) ([]docs.Doc, error) {
	var docs []docs.Doc

	cursor, err := repo.collection.Find(ctx, bson.M{"created_by": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (repo *DocsRepo) GetAll(ctx context.Context) ([]docs.Doc, error) {
	var docs []docs.Doc

	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}

func (repo *DocsRepo) GetById(ctx context.Context, id primitive.ObjectID) (*docs.Doc, error) {
	var doc docs.Doc

	err := repo.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (repo *DocsRepo) DeleteById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

func (repo *DocsRepo) DeleteAll(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	result, err := repo.collection.DeleteMany(ctx, bson.M{"created_by": userID})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (repo *DocsRepo) UpdateDoc(ctx context.Context, docID primitive.ObjectID, fieldName string, fieldValue string) (bool, error) {
	update := bson.M{
		"$set": bson.M{
			fieldName:  fieldValue,
			"updatedAt": time.Now(),
		},
	}

	result, err := repo.collection.UpdateByID(ctx, docID, update)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}
