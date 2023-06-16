package storage

import (
	"github.com/lib/pq"
	"server/internal/app/models"
)

type PhoneRepository struct {
	storage *Storage
}

func (r *PhoneRepository) Create(p *models.Phone) (*models.Phone, error) {
	err := r.storage.db.QueryRow(`INSERT INTO phones (manufacturer, model_tag, model_number, os_version, api_version, cpu, firmware, bootloader, supported_archs, sim_slots, sd_slots) 
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
										ON CONFLICT (model_number) DO UPDATE
										SET manufacturer = EXCLUDED.manufacturer 
										RETURNING phone_id`,
		p.Manufacturer, p.ModelTag, p.ModelNumber, p.OsVersion, p.ApiVersion, p.Cpu, p.Firmware, p.Bootloader, pq.StringArray(p.SupportedArchs), p.SimSlots, p.SdSlots).Scan(&p.Id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PhoneRepository) SelectByModelNumber(modelNumber string) (*models.Phone, error) {
	p := &models.Phone{}

	err := r.storage.db.QueryRow("SELECT * FROM phones WHERE model_tag = $1 LIMIT 1",
		modelNumber).Scan(
		&p.Id,
		&p.Manufacturer,
		&p.ModelTag,
		&p.ModelNumber,
		&p.OsVersion,
		&p.ApiVersion,
		&p.Cpu,
		&p.Firmware,
		&p.Bootloader,
		pq.Array(&p.SupportedArchs),
		&p.SimSlots,
		&p.SdSlots,
	)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PhoneRepository) SelectAll() ([]models.Phone, error) {
	rows, err := r.storage.db.Query(`SELECT * FROM phones`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phones []models.Phone

	for rows.Next() {
		var p models.Phone

		err := rows.Scan(
			&p.Id,
			&p.Manufacturer,
			&p.ModelTag,
			&p.ModelNumber,
			&p.OsVersion,
			&p.ApiVersion,
			&p.Cpu,
			&p.Firmware,
			&p.Bootloader,
			pq.Array(&p.SupportedArchs),
			&p.SimSlots,
			&p.SdSlots,
		)
		if err != nil {
			return nil, err
		}

		phones = append(phones, p)
	}

	return phones, nil
}

func (r *PhoneRepository) Delete(id int) error {
	err := r.storage.db.QueryRow(`DELETE FROM phones WHERE phone_id = $1`, id).Err()
	if err != nil {
		return err
	}

	return nil
}