package food

import (
	"CalorieGuide-db/internal/food"
	"CalorieGuide-db/internal/lib/logger/slg"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type FindAllRequest struct {
	Sort      string `json:"sort,omitempty"`
	TwoDecade int    `json:"two-decade,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
}

type FindAllResponse struct {
	//response.Response
	Products []food.WithLike `json:"products,omitempty"`
}

type AddResponse struct {
	//response.Response
	Product food.Food `json:"product"`
}

type LikeResponse struct {
	UserId int    `json:"user_id"`
	FoodId int    `json:"product_id"`
	Action string `json:"action,omitempty"`
}

type SearchRequest struct {
	Word   string `json:"word"`
	UserId int    `json:"user,omitempty"`
}

func NewFindAll(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewFindAll"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req FindAllRequest
		var sortType = "fromNewest"
		var twoDecade = 1
		var userId = 0
		if r.Body != http.NoBody {
			err := render.DecodeJSON(r.Body, &req)
			if err != nil {
				log.Error("Failed to parse json")
				return
			}
			if req.Sort != "" {
				sortType = req.Sort
			}
			if req.TwoDecade != 0 {
				twoDecade = req.TwoDecade
			}
			if req.UserID != 0 {
				userId = req.UserID
			}
		}
		foods, err := repository.FindAll(r.Context(), sortType, twoDecade, userId)
		if err != nil {
			log.Error("Failed to get all foods")
			return
		}
		render.JSON(w, r, FindAllResponse{
			//Response: response.Ok(),
			Products: foods,
		})
	}
}

func NewAdd(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req food.Food
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorFromClaims := int(claims["id"].(float64))
		if req.AuthorId != authorFromClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = repository.Create(r.Context(), &req)
		if err != nil {
			log.Error("Failed to create food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, req)
	}
}

func NewFindOne(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewFindOne"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		prodId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get food Id", slg.Err(err))
			return
		}
		product, err := repository.FindOne(r.Context(), prodId)
		if err != nil {
			log.Error("Failed to get food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, product)
	}
}

func NewUpdate(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewUpdate"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req food.Food
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			return
		}
		prodId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get food Id", slg.Err(err))
			return
		}
		fd, err := repository.FindOne(r.Context(), prodId)
		if err != nil {
			log.Error("There is no such product", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if fd.AuthorId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req.Id = prodId
		err = repository.Update(r.Context(), req)
		if err != nil {
			log.Error("Failed to update food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product, err := repository.FindOne(r.Context(), prodId)
		if err != nil {
			log.Error("Failed to update food", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, product)
	}
}

func NewDelete(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewDelete"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		prodId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("Failed to get food Id", slg.Err(err))
			return
		}
		fd, err := repository.FindOne(r.Context(), prodId)
		if err != nil {
			log.Error("There is no such product", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		authorIdClaims := int(claims["id"].(float64))
		if fd.AuthorId != authorIdClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err = repository.Delete(r.Context(), prodId)
		if err != nil {
			log.Error("There is no such product", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, fd)
	}
}

func NewLike(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewLike"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req LikeResponse
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, claims, _ := jwtauth.FromContext(r.Context())
		idClaims := int(claims["id"].(float64))
		if req.UserId != idClaims {
			log.Error("Error with authentication")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		liked, err := repository.Like(r.Context(), req.FoodId, req.UserId)
		if err != nil {
			log.Error("Some problem with like", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if liked {
			req.Action = "liked"
		} else {
			req.Action = "disliked"
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, req)
	}
}

func NewSearch(log *slog.Logger, repository food.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "food.handlers.NewSearch"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req SearchRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
		}
		foods, err := repository.Search(r.Context(), req.Word, req.UserId)
		if err != nil {
			log.Error("Failed to get foods")
			return
		}
		render.JSON(w, r, FindAllResponse{Products: foods})
	}
}
