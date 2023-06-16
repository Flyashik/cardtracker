package storage

import (
	"server/internal/app/models"
)

type NotificationRepository struct {
	storage *Storage
}

func (r *NotificationRepository) Create(n *models.Notification) (*models.Notification, error) {
	err := r.storage.db.QueryRow(`INSERT INTO notifications (model_number,notification_source,sender,body,timestamp)
										VALUES ($1, $2, $3, $4, $5) RETURNING notification_id`,
		n.ModelNumber, n.Source, n.Sender, n.Body, n.Timestamp).Scan(&n.Id)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (r *NotificationRepository) SelectByModelTag(tag string) ([]models.Notification, error) {
	rows, err := r.storage.db.Query(`SELECT * FROM notifications WHERE model_number = $1`, tag)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var n models.Notification

		err := rows.Scan(
			&n.Id,
			&n.ModelNumber,
			&n.Source,
			&n.Sender,
			&n.Body,
			&n.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, n)
	}

	return notifications, nil
}
