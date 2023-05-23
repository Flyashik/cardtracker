package storage

import (
	"server/internal/app/helper"
	"server/internal/app/models"
)

type SdRepository struct {
	storage *Storage
}

func (r *SdRepository) Create(sd *models.SdInfo, p *models.Phone) (*models.SdInfo, error) {
	if helper.IsEmptySdSlot(*sd) {
		return nil, nil
	}

	err := r.storage.db.QueryRow(`INSERT INTO sd_cards (phone_id, sd_manufacturer_id, serial_no, total_space, used_space, free_space) 
										VALUES ($1, $2, $3, $4, $5, $6) 
										ON CONFLICT (serial_no) DO UPDATE
										SET phone_id = $1
										RETURNING sd_card_id`,
		p.Id, sd.SdManufacturerId, sd.SerialNo, sd.TotalSpace, sd.UsedSpace, sd.FreeSpace).Scan(&sd.Id)
	if err != nil {
		return nil, err
	}

	return sd, nil
}

func (r *SdRepository) RemovePhoneId(phoneId int) {
	r.storage.db.QueryRow(`UPDATE sd_cards
										SET phone_id = null
										WHERE phone_id = $1`, phoneId)
}

func (r *SdRepository) SelectAll() ([]models.SdInfo, error) {
	rows, err := r.storage.db.Query(`SELECT * FROM sd_cards`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sdCards []models.SdInfo

	for rows.Next() {
		var sd models.SdInfo

		err := rows.Scan(
			&sd.Id,
			&sd.PhoneId,
			&sd.SdManufacturerId,
			&sd.SerialNo,
			&sd.TotalSpace,
			&sd.UsedSpace,
			&sd.FreeSpace,
		)
		if err != nil {
			return nil, err
		}

		sdCards = append(sdCards, sd)
	}

	return sdCards, nil
}