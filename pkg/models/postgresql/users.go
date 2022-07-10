package postgresql

import (
	"database/sql"
	"strings"

	"github.com/SmoothWay/snippetBox/pkg/models"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users(name, email, hashed_password, created)
		VALUES($1,$2,$3, now())`

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if postgreErr, ok := err.(*pq.Error); ok {
			if postgreErr.Code == "23505" && strings.Contains(postgreErr.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err

}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	// stmt := fmt.Sprintf("SELECT id, hashed_password FROM users WHERE email = '%s'", email)
	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
