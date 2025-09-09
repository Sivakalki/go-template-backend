package user_handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"rest_with_mongo/db/users"
	user_services "rest_with_mongo/services/users"
	"rest_with_mongo/utils/jwt"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user *user_services.InputUser) (*users.User, error)
	Login(ctx context.Context, email string, password string)(string, error)
}


type UserHandler struct{
	userService UserService
	jwtGen *jwt.ApxJwt
}

func NewUserHandler(svc UserService, jwtGen *jwt.ApxJwt)(*UserHandler){
	return &UserHandler{userService: svc, jwtGen: jwtGen}
}

func(handler *UserHandler) Register(w http.ResponseWriter, req *http.Request){
	var user user_services.InputUser
	if err:= json.NewDecoder(req.Body).Decode(&user); err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	created,e := handler.userService.CreateUser(req.Context(),&user)
	if e!=nil{
		
http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(created)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func(handler *UserHandler) Login(w http.ResponseWriter, req *http.Request){
	var reqData LoginRequest
	if err:= json.NewDecoder(req.Body).Decode(&reqData); err!=nil{
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token,err := handler.userService.Login(req.Context(), reqData.Email, reqData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // set true in production with https
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successfull"))
}

