package handler

import (
	"Effective_Mobile_Test_Project/internal/model"
	"Effective_Mobile_Test_Project/internal/service"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

type SubscriptionHandler struct {
	*service.SubscriptionService
}

func NewSubscriptionHandler(s *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{s}
}

// CreateUpdateSubscriptionRequest
// Структура запроса на создание подписки для документации
// (для документации)
type CreateUpdateSubscriptionRequest struct {
	ServiceName string              `json:"service_name" example:"Yandex Plus" description:"Название сервиса подписки"`
	Price       int                 `json:"price" example:"400" description:"Цена подписки"`
	UserID      string              `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" description:"UUID пользователя"`
	StartDate   model.DayMonthYear  `json:"start_date" example:"07-01-2025" description:"Дата начала подписки"`
	EndDate     *model.DayMonthYear `json:"end_date,omitempty" example:"10-12-2027" description:"Дата окончания подписки"`
}

// SubscriptionCreateResponse
// Структура для ответа при создании подписки
type SubscriptionCreateResponse struct {
	Message      string                    `json:"message" example:"подписка успешно создана"`
	Subscription model.SubscriptionDetails `json:"subscription"`
}

// TotalCostResponse
// Структура для вывода инф-ии по общей стоимости всех подписок
type TotalCostResponse struct {
	UserID    string `json:"user_id"`
	TotalCost int    `json:"общая_стоимость"`
}

// SubscriptionUpdateResponse
// Структура ответа для обновления инф-ии по подписке
type SubscriptionUpdateResponse struct {
	Message string `json:"message" example:"подписка успешно обновлена"`
}

// Create godoc
// @Summary      Создание подписки
// @Description  Создаёт новую подписку с указанием пользователя, сервиса, стоимости и периода
// @Tags         Подписки
// @Accept       json
// @Produce      json
// @Param        subscription  body      CreateUpdateSubscriptionRequest         true  "Детали подписки"
// @Success      201           {object}  SubscriptionCreateResponse
// @Failure      400           {string}  string  "неверный формат запроса"
// @Failure      500           {string}  string  "ошибка создания подписки"
// @Router       /subscriptions/create [post]
func (handler *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.SubscriptionDetails

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	err := handler.CreateSubscription(r.Context(), &input)
	if err != nil {
		log.Println(err)
		http.Error(w, "ошибка создания подписки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(SubscriptionCreateResponse{
		Message:      "подписка успешно создана",
		Subscription: input,
	})
}

// GetByUserUUID godoc
// @Summary      Получение подписок пользователя
// @Description  Возвращает список всех подписок по UUID пользователя
// @Tags         Подписки
// @Produce      json
// @Param        uuid  path      string  true  "UUID пользователя"
// @Success      200   {array}   model.SubscriptionDetails
// @Failure      404   {string}  string  "не удалось получить подписки"
// @Router       /subscriptions/user/{uuid} [get]
func (handler *SubscriptionHandler) GetByUserUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	subscriptions, err := handler.GetSubscriptionsByUserUUID(r.Context(), uuid)
	if err != nil {
		log.Println(err)
		http.Error(w, "не удалось получить подписки", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(subscriptions)
}

// GetByID godoc
// @Summary      Получение подписки по ID
// @Description  Возвращает подписку по её уникальному идентификатору
// @Tags         Подписки
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  model.SubscriptionDetails
// @Failure      400  {string}  string  "неверный ID"
// @Failure      404  {string}  string  "не удалось поулчить подписку"
// @Router       /subscriptions/get/{id} [get]
func (handler *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	subscription, err := handler.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, "не удалось поулчить подписку", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(subscription)
}

// GetTotalCost godoc
// @Summary      Получение общей стоимости подписок пользователя
// @Description  Возвращает сумму всех подписок пользователя с возможностью фильтрации по сервису и диапазону дат
// @Tags         Подписки
// @Produce      json
// @Param        user_id      query     string  true   "UUID пользователя"
// @Param        service_name query     string  false  "Название сервиса (опционально)"
// @Param        start_date   query     string  false  "Дата начала (формат DD-MM-YYYY, по умолчанию 01-01-2000)"
// @Param        end_date     query     string  false  "Дата окончания (формат DD-MM-YYYY)"
// @Success      200  {object}  TotalCostResponse
// @Failure      400  {string}  string  "ошибка параметров запроса"
// @Failure      500  {string}  string  "ошибка сервера"
// @Router       /subscriptions/total-cost [get]
func (handler *SubscriptionHandler) GetTotalCost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	if userID == "" {
		http.Error(w, "обязательный параметр user_id отсутствует", http.StatusBadRequest)
		return
	}

	serviceName := query.Get("service_name")
	startStr := query.Get("start_date")
	endStr := query.Get("end_date")

	var startDate model.DayMonthYear
	var endDate model.DayMonthYear
	var err error

	if startStr != "" {
		err = startDate.UnmarshalJSON([]byte(`"` + startStr + `"`))
		if err != nil {
			http.Error(w, "неверный формат start_date, ожидается DD-MM-YYYY", http.StatusBadRequest)
			return
		}
	} else {
		startDate = model.DayMonthYear(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	}

	if endStr != "" {
		err = endDate.UnmarshalJSON([]byte(`"` + endStr + `"`))
		if err != nil {
			http.Error(w, "неверный формат end_date, ожидается DD-MM-YYYY", http.StatusBadRequest)
			return
		}
	} else {
		endDate = model.DayMonthYear(time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC))
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	total, err := handler.GetSubscriptionsCostByUserDetails(
		r.Context(),
		userID,
		serviceNamePtr,
		startDate.ToTime(),
		endDate.ToTime(),
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "не удалось получить подписки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(TotalCostResponse{
		UserID:    userID,
		TotalCost: total,
	})
}

// UpdateByID godoc
// @Summary      Обновить подписку по ID
// @Description  Обновляет информацию о подписке по заданному ID
// @Tags         Подписки
// @Accept       json
// @Produce      json
// @Param        id            path      int                       true  "ID подписки"
// @Param        subscription  body      CreateUpdateSubscriptionRequest true  "Обновлённые данные подписки"
// @Success      200           {object}  SubscriptionUpdateResponse
// @Failure      400           {string}  string  "неверный ID или формат запроса"
// @Failure      500           {string}  string  "не удалось обновить информацию по подписке"
// @Router       /subscriptions/update/{id} [put]
func (handler *SubscriptionHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	var input model.SubscriptionDetails
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "проверьте правильность переданных данных", http.StatusBadRequest)
		return
	}

	err = handler.UpdateSubscriptionByID(r.Context(), &input, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "не удалось обновить информацию по подписке", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SubscriptionUpdateResponse{
		Message: "подписка успешно обновлена",
	})
}

// DeleteByID godoc
// @Summary      Удалить подписку по ID
// @Description  Удаляет подписку по её идентификатору
// @Tags         Подписки
// @Param        id   path      int  true  "ID подписки"
// @Success      204  "подписка успешно удалена"
// @Failure      400  {string}  string  "неверный ID"
// @Failure      500  {string}  string  "не удалось удалить подписку"
// @Router       /subscriptions/delete/{id} [delete]
func (handler *SubscriptionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	err = handler.DeleteSubscriptionByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, "не удалось удалить подписку", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
