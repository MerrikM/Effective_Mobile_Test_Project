package repository

import (
	"Effective_Mobile_Test_Project/internal/config"
	"Effective_Mobile_Test_Project/internal/model"
	"Effective_Mobile_Test_Project/internal/util"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

type SubscriptionRepository struct {
	*config.Database
}

func NewSubscriptionRepository(database *config.Database) *SubscriptionRepository {
	return &SubscriptionRepository{database}
}

func (repo *SubscriptionRepository) SaveSubscription(ctx context.Context, exec sqlx.ExtContext, subscription *model.SubscriptionDetails) error {
	query := `INSERT INTO subscriptions 
        (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := exec.ExecContext(
		ctx,
		query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate.ToTime(),
		subscription.EndDate.ToTime(),
	)

	if err != nil {
		return util.LogError("ошибка при вставке подписки", err)
	}
	return nil
}

func (repo *SubscriptionRepository) GetSubscriptionByID(ctx context.Context, exec sqlx.ExtContext, id int) (*model.SubscriptionDetails, error) {
	query := `SELECT * FROM subscriptions WHERE id=$1`

	var returnedOrder model.SubscriptionDetails
	err := sqlx.GetContext(ctx, exec, &returnedOrder, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, util.LogError("не удалось найти подписку по ее id", err)
		}
		return nil, util.LogError("ошибка получения таблицы подписок", err)
	}

	return &returnedOrder, nil
}

func (repo *SubscriptionRepository) GetSubscriptionsByUserUUID(ctx context.Context, exec sqlx.ExtContext, uuid string) ([]model.SubscriptionDetails, error) {
	query := `SELECT * FROM subscriptions WHERE user_id=$1`

	var subscriptions []model.SubscriptionDetails
	err := sqlx.SelectContext(ctx, exec, &subscriptions, query, uuid)
	if err != nil {
		return nil, util.LogError("ошибка получения подписок по uuid пользователя", err)
	}

	return subscriptions, nil
}

func (repo *SubscriptionRepository) GetTotalSubscriptionCost(
	ctx context.Context,
	exec sqlx.ExtContext,
	userID string,
	serviceName *string,
	startPeriod time.Time,
	endPeriod time.Time,
) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0) AS total
		FROM subscriptions
		WHERE 
			($1::uuid IS NULL OR user_id = $1::uuid) AND
			($2::text IS NULL OR service_name = $2::text) AND
			start_date <= $4 AND (end_date IS NULL OR end_date >= $3)
	`

	var total int
	err := sqlx.GetContext(ctx, exec, &total, query, userID, serviceName, startPeriod, endPeriod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, util.LogError("не удалось найти подписку по ее id", err)
		}
		return 0, util.LogError("ошибка подсчета общей стоимости подписок", err)
	}

	return total, nil
}

func (repo *SubscriptionRepository) UpdateSubscriptionByID(ctx context.Context, exec sqlx.ExtContext, subscription *model.SubscriptionDetails, id int) error {
	query := `UPDATE subscriptions
			SET service_name = $1,
			    price = $2,
			    user_id = $3,
			    start_date = $4,
			    end_date = $5
			WHERE id = $6
			`
	res, err := exec.ExecContext(ctx, query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate.ToTime(),
		subscription.EndDate.ToTime(),
		id,
	)
	if err != nil {
		return util.LogError("ошибка при обновлении подписки", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return util.LogError("не удалось получить количество затронутых строк", err)
	}

	if rowsAffected == 0 {
		return util.LogError("подписка с таким ID не найдена", sql.ErrNoRows)
	}

	return nil
}

func (repo *SubscriptionRepository) DeleteSubscriptionByID(ctx context.Context, exec sqlx.ExtContext, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := exec.ExecContext(ctx, query, id)
	if err != nil {
		return util.LogError("ошибка при удалении подписки по ID", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return util.LogError("не удалось получить количество удалённых строк", err)
	}

	if rowsAffected == 0 {
		return util.LogError("подписка с таким ID не найдена", sql.ErrNoRows)
	}

	return nil
}
