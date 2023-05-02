package storage

import "server/internal/app/models"

type PhoneRepository struct {
	storage *Storage
}

func (r *PhoneRepository) Create(p *models.Phone) (*models.Phone, error) {
	//r.storage.db.QueryRow("INSERT INTO phones ()")
	return nil, nil
}
