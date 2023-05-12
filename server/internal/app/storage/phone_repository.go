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
										ON CONFLICT (model_tag) DO UPDATE
										SET manufacturer = EXCLUDED.manufacturer 
										RETURNING phone_id`,
		p.Manufacturer, p.ModelTag, p.ModelNumber, p.OsVersion, p.ApiVersion, p.Cpu, p.Firmware, p.Bootloader, pq.StringArray(p.SupportedArchs), p.SimSlots, p.SdSlots).Scan(&p.Id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *PhoneRepository) SelectByModelTag(modelTag string) (*models.Phone, error) {
	p := &models.Phone{}

	err := r.storage.db.QueryRow("SELECT * FROM phones WHERE model_tag = $1 LIMIT 1",
		modelTag).Scan(
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
