package client

import (
	"CalorieGuide-db/internal/client"
	"CalorieGuide-db/internal/config"
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"CalorieGuide-db/internal/meal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type FindMailRequest struct {
	Mail     string `json:"email"`
	Password string `json:"password"`
}

type FindMailResponse struct {
	client.Client
	client.Diet
	client.Goal
	BearerToken string
}

type clientFull struct {
	client.Client
	client.Diet
	client.Goal
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
			Picture:  req.Picture,
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
		goal, err := repository.FindGoalById(r.Context(), cli.Id)
		if err != nil {
			log.Error("Failed to find goal", slg.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		diet, err := repository.FindDietById(r.Context(), cli.Id)
		if err != nil {
			log.Error("Failed to find diet", slg.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		claims := map[string]interface{}{"id": cli.Id, "email": cli.Email, "password": cli.Password}
		_, tokenString, err := config.GetToken(log).Encode(claims)
		if err != nil {
			log.Error("Failed to get token", slg.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := FindMailResponse{cli, diet, goal, tokenString}
		w.Header().Set("Content-Type", "application/json")
		render.JSON(w, r, resp)
	}
}

func NewUpdate(log *slog.Logger, repository client.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "client.handlers.NewUpdate"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req clientFull
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body for client", slg.Err(err))
			return
		}
		clientId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get client Id", slg.Err(err))
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if clientId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req.Client.Id = clientId
		err = repository.UpdateClient(r.Context(), req.Client)
		if err != nil {
			log.Error("Failed to update client", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = repository.UpdateDiet(r.Context(), req.Diet, req.Client.Id)
		if err != nil {
			log.Error("Failed to update diet", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = repository.UpdateGoal(r.Context(), req.Goal, req.Client.Id)
		if err != nil {
			log.Error("Failed to update goal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		req.Client, _ = repository.FindById(r.Context(), req.Client.Id)
		req.Diet, _ = repository.FindDietById(r.Context(), req.Client.Id)
		req.Goal, _ = repository.FindGoalById(r.Context(), req.Goal.Id)
		render.JSON(w, r, req)
	}
}

func NewDelete(log *slog.Logger, repository client.Repository, fdRepo food.Repository, mlRepo meal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "client.handlers.NewDelete"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		clientId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get client Id", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if clientId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		cl, err := repository.FindById(r.Context(), clientId)
		if err != nil {
			log.Error("Failed to find client", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		goal, err := repository.FindGoalById(r.Context(), clientId)
		if err != nil {
			log.Error("Failed to find goal", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		diet, err := repository.FindDietById(r.Context(), clientId)
		if err != nil {
			log.Error("Failed to find diet", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = repository.Delete(r.Context(), clientId, fdRepo, mlRepo)
		if err != nil {
			log.Error("Failed to delete client", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, clientFull{cl, diet, goal})
	}
}

func NewFindMealLikes(log *slog.Logger, repository client.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "client.handlers.NewFindMealLikes"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		clientId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get client Id", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if clientId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		cl, err := repository.FindMealLikes(r.Context(), clientId)
		if err != nil {
			log.Error("Failed to find meal likes", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, struct {
			Meals []int `json:"meals,omitempty"`
		}{Meals: cl})
	}
}

func NewFindFoodLikes(log *slog.Logger, repository client.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "client.handlers.NewFindFoodLikes"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		clientId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get client Id", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if clientId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		cl, err := repository.FindFoodLikes(r.Context(), clientId)
		if err != nil {
			log.Error("Failed to find food likes", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, struct {
			Products []int `json:"products,omitempty"`
		}{Products: cl})
	}
}
