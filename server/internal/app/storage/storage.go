package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	config          *DbConfig
	db              *sql.DB
	phoneRepository *PhoneRepository
	simRepository   *SimRepository
	sdRepository    *SdRepository
	userRepository  *UserRepository
}

func New(config *DbConfig) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.config.DbURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) Phone() *PhoneRepository {
	if s.phoneRepository != nil {
		return s.phoneRepository
	}

	s.phoneRepository = &PhoneRepository{
		storage: s,
	}

	return s.phoneRepository
}

func (s *Storage) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		storage: s,
	}

	return s.userRepository
}

func (s *Storage) Sim() *SimRepository {
	if s.simRepository != nil {
		return s.simRepository
	}

	s.simRepository = &SimRepository{
		storage: s,
	}

	return s.simRepository
}

func (s *Storage) SdCard() *SdRepository {
	if s.sdRepository != nil {
		return s.sdRepository
	}

	s.sdRepository = &SdRepository{
		storage: s,
	}

	return s.sdRepository
}
