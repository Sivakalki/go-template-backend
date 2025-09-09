package server

import (
	"context"
	"fmt"
	"net/http"
	"rest_with_mongo/http/handlers/docshandlers"
	user_handlers "rest_with_mongo/http/handlers/userhandlers"
	middleware_logger "rest_with_mongo/http/middleware"
	context_keys "rest_with_mongo/utils/context"
	"rest_with_mongo/utils/jwt"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	router *chi.Mux
	userHandler *user_handlers.UserHandler
	docsHandler *docs_handlers.DocsHandler
	jwtGen *jwt.ApxJwt

}


func NewServer(logger *zap.Logger, handler *user_handlers.UserHandler,docsHandler *docs_handlers.DocsHandler ,jwtGen *jwt.ApxJwt) *Server{
	r := chi.NewRouter()
	r.Use(middleware_logger.ZapLogger(logger))
	
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, 
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s := &Server{
		logger: logger,
		userHandler: handler,
		docsHandler: docsHandler,
		router: r,
		jwtGen: jwtGen,
	}

	s.routes()

	return s;
}


func (s *Server) routes(){
	s.router.Use(middleware.StripSlashes)
	s.router.Post("/auth/register", s.userHandler.Register)
	s.router.Post("/auth/login", s.userHandler.Login)
	

	s.router.Route("/docs",func(r chi.Router){
		r.Use(s.AuthMiddleware)
		r.Post("/", s.docsHandler.CreateDoc)
		r.Get("/", s.docsHandler.GetAllDocs)
		r.Get("/mine", s.docsHandler.GetDocsByAuthor)
		r.Get("/{id}", s.docsHandler.GetDocByID)
		r.Delete("/{id}", s.docsHandler.DeleteDocByID)
		r.Delete("/all", s.docsHandler.DeleteAllDocsByUser)
		r.Patch("/{id}", s.docsHandler.UpdateDocField)
	})
}

func (s *Server) Start(addr string) error{
	s.logger.Info("server starting ",zap.String("addr", addr))
	return http.ListenAndServe(addr, s.router)
}



func(s *Server) AuthMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie,err := r.Cookie("token")
		fmt.Println(cookie.Value, err)
		if err!=nil{
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, err :=s.jwtGen.Decode(cookie.Value)
		if err!=nil{
			http.Error(w, "unable to decode", http.StatusUnauthorized)
			return
		}	

		userID, ok := claims["user_id"].(string)
        if !ok {
            http.Error(w, "invalid token payload", http.StatusUnauthorized)
            return
        }

		ctx := context.WithValue(r.Context(), context_keys.UserIDKey,userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

