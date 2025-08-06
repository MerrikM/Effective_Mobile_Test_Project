package service

import (
	"Effective_Mobile_Test_Project/internal/model"
	"Effective_Mobile_Test_Project/internal/repository"
	"Effective_Mobile_Test_Project/internal/util"
	"context"
	"log"
	"time"
)

type SubscriptionService struct {
	*repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscription *model.SubscriptionDetails) error {
	err := s.SubscriptionRepository.SaveSubscription(ctx, s.Database, subscription)
	if err != nil {
		return util.LogError("не удалось создать подписку", err)
	}

	log.Printf("подписка сохранена, детали подписки: %v", subscription)
	return nil
}

func (s *SubscriptionService) GetSubscriptionsByUserUUID(ctx context.Context, uuid_id string) ([]model.SubscriptionDetails, error) {
	subscriptions, err := s.SubscriptionRepository.GetSubscriptionsByUserUUID(ctx, s.Database, uuid_id)
	if err != nil {
		return nil, util.LogError("не удалось найти подписку", err)
	}

	log.Printf("список подписок пользователя с uuid=%s: %v", uuid_id, subscriptions)
	return subscriptions, nil
}

func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id int) (*model.SubscriptionDetails, error) {
	subscription, err := s.SubscriptionRepository.GetSubscriptionByID(ctx, s.Database, id)
	if err != nil {
		return nil, util.LogError("не удалось найти подписку", err)
	}

	log.Printf("подписка с id=%d: %v", id, subscription)
	return subscription, nil
}

func (s *SubscriptionService) GetSubscriptionsCostByUserDetails(
	ctx context.Context,
	userID *string,
	serviceName *string,
	startPeriod time.Time,
	endPeriod time.Time,
) (int, error) {
	totalCost, err := s.SubscriptionRepository.GetTotalSubscriptionCost(ctx, s.Database, userID, serviceName, startPeriod, endPeriod)
	if err != nil {
		return 0, util.LogError("не удалось получить общую стоимость подписок", err)
	}

	log.Printf("общая стоимость подписок пользователя с uuid=%d : %d", userID, totalCost)
	return totalCost, nil
}

func (s *SubscriptionService) UpdateSubscriptionByID(ctx context.Context, subscription *model.SubscriptionDetails, id int) error {
	err := s.SubscriptionRepository.UpdateSubscriptionByID(ctx, s.Database, subscription, id)
	if err != nil {
		return util.LogError("не удалось обновить подписку", err)
	}

	log.Printf("подписка с id=%d успешно обновлена", subscription.ID)
	return nil
}

func (s *SubscriptionService) DeleteSubscriptionByID(ctx context.Context, id int) error {
	err := s.SubscriptionRepository.DeleteSubscriptionByID(ctx, s.Database, id)
	if err != nil {
		return util.LogError("не удалось удалить подписку", err)
	}

	log.Printf("подписка с id=%d успешно удалена", id)
	return nil
}
