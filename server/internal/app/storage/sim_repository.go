package storage

import (
	"server/internal/app/helper"
	"server/internal/app/models"
)

type SimRepository struct {
	storage *Storage
}

func (r *SimRepository) Create(sim *models.SimInfo, p *models.Phone) (*models.SimInfo, error) {
	if helper.IsEmptySimSlot(*sim) {
		return nil, nil
	}

	err := r.storage.db.QueryRow(`INSERT INTO sim_cards (phone_id, phone_number, operator) 
										VALUES ($1, $2, $3) 
										ON CONFLICT (phone_number) DO UPDATE
		                    			SET phone_id = $1
		                    			RETURNING sim_card_id`,
		p.Id, sim.PhoneNumber, sim.Operator).Scan(&sim.Id)
	if err != nil {
		return nil, err
	}

	return sim, nil
}

func (r *SimRepository) RemovePhoneId(phoneId int) {
	r.storage.db.QueryRow(`UPDATE sim_cards
									SET phone_id = null
									WHERE phone_id = $1`, phoneId)
}
