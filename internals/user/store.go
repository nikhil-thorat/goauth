package user

import (
	"database/sql"
	"fmt"

	"github.com/nikhil-thorat/goauth/types"
)

const (
	CreateUserQuery                   = "INSERT INTO users (firstName, lastName, email, password, isVerified, verificationCode) VALUES ($1,$2,$3,$4,$5,$6)"
	GetUserByEmailQuery               = "SELECT * FROM users WHERE email = $1"
	GetUserByIDQuery                  = "SELECT * FROM users WHERE id = $1"
	UpdateUserDetailsQuery            = "UPDATE users set firstName = $1, lastName = $2 WHERE id = $3"
	UpdateUserPasswordQuery           = "UPDATE users set password = $1 WHERE id = $2"
	UpdateUserVerificationStatusQuery = "UPDATE users set isVerified = $1 WHERE id = $2"
	DeleteUserQuery                   = "DELETE FROM users WHERE id = $1"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(CreateUserQuery, user.FirstName, user.LastName, user.Email, user.Password, user.IsVerified, user.VerificationCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {

	row := s.db.QueryRow(GetUserByEmailQuery, email)

	var user types.User

	err := scanRowIntoUser(row, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("USER NOT FOUND")
		}

		return nil, fmt.Errorf("ERROR SCANNING USER : %v", err)
	}
	return &user, nil

}

func (s *Store) GetUserByID(userID int) (*types.User, error) {

	row := s.db.QueryRow(GetUserByIDQuery, userID)

	var user types.User

	err := scanRowIntoUser(row, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("USER NOT FOUND")
		}

		return nil, fmt.Errorf("ERROR SCANNING USER : %v", err)
	}
	return &user, nil

}

func (s *Store) UpdateUserDetails(userID int, firstName string, lastName string) error {
	_, err := s.db.Exec(UpdateUserDetailsQuery, firstName, lastName, userID)

	if err != nil {
		return nil
	}

	return nil
}

func (s *Store) UpdateUserPassword(userID int, password string) error {

	_, err := s.db.Exec(UpdateUserPasswordQuery, password, userID)

	if err != nil {
		return nil
	}
	return nil
}

func (s *Store) UpdateUserVerificationStatus(userID int, isVerified bool) error {

	_, err := s.db.Exec(UpdateUserVerificationStatusQuery, isVerified, userID)

	if err != nil {
		return nil
	}
	return nil
}

func (s *Store) DeleteUser(userID int) error {
	res, err := s.db.Exec(DeleteUserQuery, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("USER NOT FOUND")
	}

	return nil
}

func scanRowIntoUser(row *sql.Row, user *types.User) error {
	return row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.VerificationCode,
		&user.CreatedAt,
	)
}
