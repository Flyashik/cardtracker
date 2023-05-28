package storage

import (
	"server/internal/app/models"
)

type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(u *models.User) (*models.User, error) {
	err := r.storage.db.QueryRow(`INSERT INTO users (email,name,code,password,role)
VALUES ($1, $2, $3, $4, $5)`,
		u.Email, u.Name, u.Code, u.Password, u.Role).Scan(&u.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) SelectByEmail(email string) (*models.User, error) {
	u := &models.User{}

	err := r.storage.db.QueryRow("SELECT * FROM users WHERE email = $1 LIMIT 1",
		email).Scan(
		&u.Id,
		&u.Name,
		&u.Code,
		&u.Email,
		&u.Password,
		&u.Role)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) SelectByCode(code int) (*models.User, error) {
	u := &models.User{}

	err := r.storage.db.QueryRow("SELECT * FROM users WHERE code = $1 LIMIT 1",
		code).Scan(
		&u.Id,
		&u.Name,
		&u.Code,
		&u.Email,
		&u.Password,
		&u.Role)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) CheckCodeExists(code int) bool {
	res := false

	err := r.storage.db.QueryRow("select exists(select 1 from users where code = $1)",
		code).Scan(
		&res)

	if err != nil {
		return false
	}

	return res
}

func (r *UserRepository) SelectAll() ([]models.User, error) {
	rows, err := r.storage.db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		err := rows.Scan(
			&u.Id,
			&u.Name,
			&u.Code,
			&u.Email,
			&u.Password,
			&u.Role,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
