package docs_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"rest_with_mongo/db/docs"
	docs_service "rest_with_mongo/services/docs"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocsService interface {
	CreateDoc(ctx context.Context, doc docs_service.InputDoc) (*docs.Doc, error)
	GetDocsByAuthor(ctx context.Context) ([]docs.Doc, error)
	GetAllDocs(ctx context.Context) ([]docs.Doc, error)
	GetDocByID(ctx context.Context, id primitive.ObjectID) (*docs.Doc, error)
	DeleteDocByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	DeleteAllDocsByUser(ctx context.Context, userID primitive.ObjectID) (int64, error)
	UpdateDocField(ctx context.Context, docID primitive.ObjectID, fieldName string, fieldValue string) (bool, error)
}

type DocsHandler struct {
	docsService DocsService
}

func NewDocsHandler(docsService DocsService) *DocsHandler {
	return &DocsHandler{docsService: docsService}
}


func (h *DocsHandler) CreateDoc(w http.ResponseWriter, r *http.Request) {
	var input docs_service.InputDoc
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	newDoc, err := h.docsService.CreateDoc(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(newDoc)
}

func (h *DocsHandler) GetAllDocs(w http.ResponseWriter, r *http.Request) {
	allDocs, err := h.docsService.GetAllDocs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(allDocs)
}

func (h *DocsHandler) GetDocsByAuthor(w http.ResponseWriter, r *http.Request) {
	userDocs, err := h.docsService.GetDocsByAuthor(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(userDocs)
}

func (h *DocsHandler) GetDocByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("running this line")
	idHex := chi.URLParam(r, "id")
	if idHex == "" {
		http.Error(w, "missing document id", http.StatusBadRequest)
		return
	}

	docID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	doc, err := h.docsService.GetDocByID(r.Context(), docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(doc)
}

func (h *DocsHandler) DeleteDocByID(w http.ResponseWriter, r *http.Request) {
	idHex := chi.URLParam(r, "id")
	if idHex == "" {
		http.Error(w, "missing document id", http.StatusBadRequest)
		return
	}

	fmt.Println("getting the id of the document from query paarms", idHex)

	docID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	fmt.Println("calling the function")

	deleted, err := h.docsService.DeleteDocByID(r.Context(), docID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	fmt.Println("came out of the function")

	json.NewEncoder(w).Encode(map[string]any{"deleted": deleted})
}

func (h *DocsHandler) DeleteAllDocsByUser(w http.ResponseWriter, r *http.Request) {
	// userID is taken from context by service
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	count, err := h.docsService.DeleteAllDocsByUser(r.Context(), objID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"deleted_count": count})
}

type UpdateRequest struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (h *DocsHandler) UpdateDocField(w http.ResponseWriter, r *http.Request) {
	idHex := r.URL.Query().Get("id")
	if idHex == "" {
		http.Error(w, "missing document id", http.StatusBadRequest)
		return
	}

	docID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	updated, err := h.docsService.UpdateDocField(r.Context(), docID, req.Field, req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"updated": updated})
}
