package client

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mosgor/CalorieGuide-db/app/internal/client"
	"github.com/mosgor/CalorieGuide-db/app/internal/config"
	"github.com/mosgor/CalorieGuide-db/app/internal/food"
	food2 "github.com/mosgor/CalorieGuide-db/app/internal/food/handlers"
	"github.com/mosgor/CalorieGuide-db/app/internal/lib/logger/slg"
	"github.com/mosgor/CalorieGuide-db/app/internal/meal"
	meal2 "github.com/mosgor/CalorieGuide-db/app/internal/meal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

// @Description Запрос для регистрации пользователя
type RegistrationRequest struct {
	client.Client `json:"client"`
	client.Goal   `json:"goal"`
}

// @Description Запрос на вход в аккаунт
type FindMailRequest struct {
	Mail     string `json:"email"`
	Password string `json:"password"`
}

// @Description Ответ при входе в систему
type FindMailResponse struct {
	client.Client `json:"client"`
	client.Diet   `json:"diet"`
	client.Goal   `json:"goal"`
	BearerToken   string `json:"bearer_token"`
}

type clientFull struct {
	client.Client
	client.Diet
	client.Goal
}

// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя
// @Tags clients
// @Accept json
// @Produce json
// @Param request body RegistrationRequest true "Данные пользователя"
// @Success 201 {object} client.Client
// @Failure 400 {object} string "Ошибка при создании пользователя"
// @Router /user [post]
func NewAdd(log *slog.Logger, repository client.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "client.handlers.NewAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("Failed to read request body", slog.Any("error", err))
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// 2. Логируем содержимое (внимание на безопасность, см. ниже)
		log.Debug("Request body content", slog.String("body", string(bodyBytes)))

		// 3. Восстанавливаем тело запроса, чтобы DecodeJSON мог его прочитать
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req RegistrationRequest
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to parse request body", slg.Err(err))
			return
		}
		err = repository.Create(r.Context(), &req.Client)
		if err != nil {
			log.Error("Failed to create client", slg.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = repository.UpdateGoal(r.Context(), req.Goal, req.Client.Id)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, req)
	}
}

// @Summary Вход пользователя
// @Description Проверяет email и пароль, возвращает токен
// @Tags clients
// @Accept json
// @Produce json
// @Param request body FindMailRequest true "Email и пароль"
// @Success 200 {object} FindMailResponse
// @Failure 400 {object} string "Неверный email или пароль"
// @Router /login [post]
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

// @Summary Обновление данных пользователя
// @Description Обновляет данные пользователя по ID
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body clientFull true "Данные пользователя"
// @Success 200 {object} clientFull
// @Failure 400 {object} string "Ошибка при обновлении"
// @Failure 401 {object} string "Нет доступа"
// @Router /user/{id} [put]
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

// @Summary Удаление пользователя
// @Description Удаляет пользователя по ID
// @Tags clients
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} clientFull
// @Failure 400 {object} string "Ошибка при удалении"
// @Failure 401 {object} string "Нет доступа"
// @Router /user/{id} [delete]
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

// @Summary Получить понравившиеся приёмы пищи
// @Description Возвращает список понравившихся приёмов пищи пользователя
// @Tags clients
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} meal2.FindAllResponse
// @Failure 400 {object} string "Ошибка при получении"
// @Failure 401 {object} string "Нет доступа"
// @Router /user/{id}/meals [get]
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
		render.JSON(w, r, meal2.FindAllResponse{Meals: cl})
	}
}

// @Summary Получить понравившиеся продукты
// @Description Возвращает список понравившихся продуктов пользователя
// @Tags clients
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} food2.FindAllResponse
// @Failure 400 {object} string "Ошибка при получении"
// @Failure 401 {object} string "Нет доступа"
// @Router /user/{id}/products [get]
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
		render.JSON(w, r, food2.FindAllResponse{Products: cl})
	}
}
