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
	Products []food.WithLike `json:"products,omitempty"`
}

type AddResponse struct {
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

// @Summary Получить список продуктов
// @Description Возвращает список продуктов с возможностью сортировки
// @Tags products
// @Accept json
// @Produce json
// @Param request body FindAllRequest false "Параметры сортировки"
// @Success 200 {object} FindAllResponse
// @Failure 400 {object} string "Ошибка при получении"
// @Router /products [post]
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
			Products: foods,
		})
	}
}

// @Summary Добавить новый продукт
// @Description Создает новый продукт
// @Tags products
// @Accept json
// @Produce json
// @Param request body food.Food true "Данные продукта"
// @Success 201 {object} AddResponse
// @Failure 400 {object} string "Ошибка при создании"
// @Failure 401 {object} string "Нет доступа"
// @Router /product [post]
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

// @Summary Получить продукт по ID
// @Description Возвращает продукт по его ID
// @Tags products
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} food.Food
// @Failure 400 {object} string "Продукт не найден"
// @Router /products/{id} [get]
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

// @Summary Обновить продукт
// @Description Обновляет продукт по ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "ID продукта"
// @Param request body food.Food true "Обновленные данные продукта"
// @Success 200 {object} food.Food
// @Failure 400 {object} string "Ошибка при обновлении"
// @Failure 401 {object} string "Нет доступа"
// @Router /products/{id} [put]
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

// @Summary Удалить продукт
// @Description Удаляет продукт по ID
// @Tags products
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} food.Food
// @Failure 400 {object} string "Ошибка при удалении"
// @Failure 401 {object} string "Нет доступа"
// @Router /products/{id} [delete]
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

// @Summary Лайк продукта
// @Description Ставит или убирает лайк продукту
// @Tags products
// @Accept json
// @Produce json
// @Param request body LikeResponse true "Данные лайка"
// @Success 200 {object} LikeResponse
// @Failure 400 {object} string "Ошибка при лайке"
// @Failure 401 {object} string "Нет доступа"
// @Router /products/like [post]
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

// @Summary Поиск продуктов
// @Description Ищет продукты по ключевому слову
// @Tags products
// @Accept json
// @Produce json
// @Param request body SearchRequest true "Поисковый запрос"
// @Success 200 {object} FindAllResponse
// @Failure 400 {object} string "Ошибка при поиске"
// @Router /products/search [post]
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		foods, err := repository.Search(r.Context(), req.Word, req.UserId)
		if err != nil {
			log.Error("Failed to get foods")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, FindAllResponse{Products: foods})
	}
}
