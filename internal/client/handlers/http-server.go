package client

import (
	"CalorieGuide-db/internal/client"
	"CalorieGuide-db/internal/config"
	"CalorieGuide-db/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type FindMailRequest struct {
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type FindMailResponse struct {
	client.Client
	BearerToken string
}

func NewAdd(log *slog.Logger, repository client.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "client.handlers.NewAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req client.Client
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			return
		}
		err = repository.Create(r.Context(), &req)
		if err != nil {
			log.Error("Failed to create client", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, client.Client{
			Id:       req.Id,
			Name:     req.Name,
			Surname:  req.Surname,
			Email:    req.Email,
			Password: req.Password,
		})
	}
}

func FindEmail(log *slog.Logger, repository client.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "client.handlers.FindEmail"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req FindMailRequest
		var cli client.Client
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			return
		}
		cli, err = repository.FindByEmail(r.Context(), req.Mail)
		if err != nil {
			log.Error("Failed to find by email", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Password != cli.Password {
			log.Error("Passwords does not match")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		claims := map[string]interface{}{"id": cli.Id, "email": cli.Email, "password": cli.Password}
		_, tokenString, err := config.GetToken(log).Encode(claims)
		if err != nil {
			log.Error("Failed to get token", slg.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := FindMailResponse{cli, tokenString}
		w.Header().Set("Content-Type", "application/json")
		render.JSON(w, r, resp)
	}
}
