package storage

import "server/internal/app/models"

type UserPhoneRepository struct {
	storage *Storage
}

func (r *UserPhoneRepository) CreateRelation(userId int, phoneId int) error {
	row := r.storage.db.QueryRow(`INSERT INTO user_phone (user_id, phone_id) 
								 SELECT u.user_id, p.phone_id
								 FROM users u, phones p
								 WHERE u.user_id = $1 AND p.phone_id = $2
								 ON CONFLICT (phone_id) DO UPDATE
								 SET user_id = $1;`, userId, phoneId)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (r *UserPhoneRepository) SelectUsersWithPhones() ([]models.UserPhone, error) {
	rows, err := r.storage.db.Query(`SELECT u.user_id, u.name, u.email, u.code, p.phone_id
										   FROM users u
										   JOIN user_phone up ON u.user_id = up.user_id
										   JOIN phones p ON up.phone_id = p.phone_id `)
	if err != nil {
		return nil, err
	}

	var usersPhones []models.UserPhone

	for rows.Next() {
		var u models.User
		var up models.UserPhone
		var p int

		err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.Code,
			&p,
		)
		if err != nil {
			return nil, err
		}

		up.User = u
		up.Phones = append(up.Phones, p)
		usersPhones = append(usersPhones, up)
	}

	return usersPhones, nil
}
