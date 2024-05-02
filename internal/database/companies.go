package database

import (
	"fmt"

	"github.com/OlafStar/jobboard-api/internal/types"
	"github.com/google/uuid"
)

func (s *Store) GetJWTCompany(email string) (*types.JWTCompany, error) {
	var jwtCompany types.JWTCompany
	err := s.db.QueryRow("SELECT email, name, password_hash, country, city, address FROM companies WHERE email = ?", email).Scan(
		&jwtCompany.Email, 
		&jwtCompany.Name, 
		&jwtCompany.PasswordHash,
		&jwtCompany.Country, 
		&jwtCompany.City, 
		&jwtCompany.Address,
	)
	if err != nil {
		return nil, err
	}
	return &jwtCompany, nil
}

func (s *Store) DoesCompanyExist(email string) (bool, error) {
	var exists int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM companies WHERE email = ?", email).Scan(&exists); err != nil || exists > 0 {
			if exists > 0 {
					return true, fmt.Errorf("user already exists")
			}
			return true, fmt.Errorf("internal server error")
	}

	return false, nil
}

func (s *Store) InsertCompany(company types.JWTCompany) error {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	stmt, err := s.db.Prepare(`
		INSERT INTO companies (
			id, name, email, password_hash, country, city, address
		) VALUES (?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		newUUID.String(), company.Name, company.Email, company.PasswordHash,
		company.Country, company.City, company.Address,
	)
	if err != nil {
		return err
	}

	return nil
}