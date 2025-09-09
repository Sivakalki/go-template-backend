package docs_service

import (
	"context"
	"errors"
	"fmt"
	"rest_with_mongo/db/docs"
	context_keys "rest_with_mongo/utils/context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocsRepo interface {
	Create(ctx context.Context, doc *docs.Doc) (*docs.Doc, error)
	GetByAuthor(ctx context.Context, id primitive.ObjectID) ([]docs.Doc, error)
	GetAll(ctx context.Context) ([]docs.Doc, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*docs.Doc, error)
	DeleteById(ctx context.Context, id primitive.ObjectID) (bool, error)
	DeleteAll(ctx context.Context, userID primitive.ObjectID) (int64, error)
	UpdateDoc(ctx context.Context, docID primitive.ObjectID, fieldName string, fieldValue string) (bool, error)
}

type DocsService struct {
	docsRepo DocsRepo
}

func NewDocsService(docsRepo DocsRepo) *DocsService {
	return &DocsService{docsRepo: docsRepo}
}

func GetUserIDFromCtx(ctx context.Context) (primitive.ObjectID, error) {
	val := ctx.Value(context_keys.UserIDKey)
	if val == nil {
		return primitive.NilObjectID, errors.New("user not authenticated")
	}

	uidStr, ok := val.(string)
	if !ok || uidStr == "" {
		return primitive.NilObjectID, errors.New("invalid user id in context")
	}

	userID, err := primitive.ObjectIDFromHex(uidStr)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid user id format")
	}

	return userID, nil
}


type InputDoc struct{
	Title string `json:"title"`
	Content string `json:"content"`
}


func (s *DocsService) CreateDoc(ctx context.Context, doc InputDoc) (*docs.Doc, error) {

	if(doc.Title == ""){
		return nil, errors.New("Title of the document should not be empty")
	}
	
	authorId, err := GetUserIDFromCtx(ctx)
	if err!=nil{
		return nil, err
	}

	
	mainDoc := &docs.Doc{
		ID:        primitive.NewObjectID(),
		Title:     doc.Title,
		Content:   doc.Content,
		CreatedBy: authorId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}


	return s.docsRepo.Create(ctx, mainDoc)
}


func (s *DocsService) GetDocsByAuthor(ctx context.Context) ([]docs.Doc, error) {
	authorID, err := GetUserIDFromCtx(ctx)
	if err != nil{
		return nil, err
	}

	return s.docsRepo.GetByAuthor(ctx, authorID)
}

func (s *DocsService) GetAllDocs(ctx context.Context) ([]docs.Doc, error) {
	return s.docsRepo.GetAll(ctx)
}

func (s *DocsService) GetDocByID(ctx context.Context, id primitive.ObjectID) (*docs.Doc, error) {
	return s.docsRepo.GetById(ctx, id)
}

// Delete a document by ID
func (s *DocsService) DeleteDocByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	fmt.Println("entered hre")
	authorId, err := GetUserIDFromCtx(ctx)
	fmt.Println("fetched author id", authorId)
	if err!=nil{
		return false, err
	}
	doc, err := s.GetDocByID(ctx, id)
	fmt.Println("fetching the doc author id", doc.CreatedBy)
	if err!=nil{
		return false, err
	}
	fmt.Println(doc.CreatedBy, authorId, " are the ids of the docs created by")
	if doc.CreatedBy != authorId{
		return false, errors.New("You are not eligible to delete the document")
	}
	return s.docsRepo.DeleteById(ctx, id)
}


func (s *DocsService) DeleteAllDocsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	
	return s.docsRepo.DeleteAll(ctx, userID)
}


func (s *DocsService) UpdateDocField(ctx context.Context, docID primitive.ObjectID, fieldName string, fieldValue string) (bool, error) {
	authorId, err := GetUserIDFromCtx(ctx)
	if err!=nil{
		return false, err
	}
	doc, err := s.GetDocByID(ctx, docID)
	if err!=nil{
		return false, err
	}
	if doc.CreatedBy != authorId{
		return false, errors.New("you are not eligible to update this document")

	}
	return s.docsRepo.UpdateDoc(ctx, docID, fieldName, fieldValue)
}