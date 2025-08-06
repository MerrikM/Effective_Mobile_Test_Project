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

func (handler *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.SubscriptionDetails

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	err := handler.CreateSubscription(r.Context(), &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "подписка успешно создана",
		"subscription": input,
	})
}

func (handler *SubscriptionHandler) GetByUserUUID(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	subscriptions, err := handler.GetSubscriptionsByUserUUID(r.Context(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(subscriptions)
}

func (handler *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	subscription, err := handler.GetSubscriptionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(subscription)
}

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
		endDate = model.DayMonthYear(time.Now())
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	userIDPtr := &userID

	total, err := handler.GetSubscriptionsCostByUserDetails(
		r.Context(),
		userIDPtr,
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
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":         userID,
		"общая_стоимость": total,
	})
}

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "подписка успешно обновлена",
	})
}

func (handler *SubscriptionHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	err = handler.DeleteSubscriptionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
