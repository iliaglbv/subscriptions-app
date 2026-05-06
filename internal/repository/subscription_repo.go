// internal/repository/subscription_repo.go
package repository

import (
	"database/sql"
	"time"

	"subscriptions-app/internal/models"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *models.CreateSubscriptionRequest, userID int64) (int64, error) {
	nextPaymentDate, _ := time.Parse("2006-01-02", sub.NextPaymentDate)

	var id int64
	err := r.db.QueryRow(
		`INSERT INTO subscriptions 
		 (user_id, name, cost, currency, billing_cycle, next_payment_date, category, is_active, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		userID, sub.Name, sub.Cost, sub.Currency, sub.BillingCycle,
		nextPaymentDate, sub.Category, true, time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SubscriptionRepository) GetAllByUserID(userID int64) ([]*models.Subscription, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, cost, currency, billing_cycle, 
		 next_payment_date, category, is_active, created_at
		 FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.Name, &sub.Cost, &sub.Currency,
			&sub.BillingCycle, &sub.NextPaymentDate, &sub.Category,
			&sub.IsActive, &sub.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) GetByID(id, userID int64) (*models.Subscription, error) {
	var sub models.Subscription
	err := r.db.QueryRow(
		`SELECT id, user_id, name, cost, currency, billing_cycle,
		 next_payment_date, category, is_active, created_at
		 FROM subscriptions WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(
		&sub.ID, &sub.UserID, &sub.Name, &sub.Cost, &sub.Currency,
		&sub.BillingCycle, &sub.NextPaymentDate, &sub.Category,
		&sub.IsActive, &sub.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *SubscriptionRepository) Update(id, userID int64, sub *models.CreateSubscriptionRequest) error {
	nextPaymentDate, _ := time.Parse("2006-01-02", sub.NextPaymentDate)

	_, err := r.db.Exec(
		`UPDATE subscriptions SET 
		 name = $1, cost = $2, currency = $3, billing_cycle = $4,
		 next_payment_date = $5, category = $6
		 WHERE id = $7 AND user_id = $8`,
		sub.Name, sub.Cost, sub.Currency, sub.BillingCycle,
		nextPaymentDate, sub.Category, id, userID,
	)

	return err
}

func (r *SubscriptionRepository) Delete(id, userID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM subscriptions WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	return err
}

func (r *SubscriptionRepository) GetStatsByCategory(userID int64) (map[string]float64, error) {
	rows, err := r.db.Query(
		`SELECT category, SUM(cost) as total 
		 FROM subscriptions 
		 WHERE user_id = $1 AND is_active = true 
		 GROUP BY category`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]float64)
	for rows.Next() {
		var category string
		var total float64
		if err := rows.Scan(&category, &total); err != nil {
			continue
		}
		stats[category] = total
	}

	return stats, nil
}

func (r *SubscriptionRepository) GetTotalMonthlyCost(userID int64) (float64, error) {
	var total float64
	err := r.db.QueryRow(
		`SELECT COALESCE(SUM(
			CASE 
				WHEN billing_cycle = 'monthly' THEN cost
				WHEN billing_cycle = 'yearly' THEN cost / 12
				ELSE cost
			END
		), 0)
		FROM subscriptions 
		WHERE user_id = $1 AND is_active = true`,
		userID,
	).Scan(&total)

	return total, err
}
