// https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql

package store

import (
	"database/sql"
	"errors"
	"fmt"
	"macrotrack/internal/pkg/types"

	"github.com/google/uuid"

	_ "github.com/lib/pq" // <------------ here
)

type postgres struct {
	DSN string
	Db  *sql.DB
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS macros
(
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    carbs int NOT NULL,
	protein int NOT NULL,
	fat int NOT NULL,
	alcohol int NOT NULL,
	date  TIMESTAMP WITH TIME ZONE default CURRENT_TIMESTAMP
)`

func ensureTableExists(db *sql.DB) error {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		return err
	}

	return nil
}

func (s *postgres) Init() error {

	var err error

	//dsn := "host=localhost port=5432 user=postgres password=postgres dbname=macros sslmode=disable"

	s.Db, err = sql.Open("postgres", s.DSN)

	if err != nil {
		fmt.Print("error=", err)
		return err
	}

	if s.Db == nil {
		return errors.New("null db")
	}

	err = ensureTableExists(s.Db)
	if err != nil {
		fmt.Print("error=", err)
		return err
	}

	return nil
}

func (s *postgres) Open() error {
	return nil
}
func (s *postgres) Create(m types.Macro) (string, error) {

	retUUID := uuid.NullUUID{}

	sqlStatement := `
			INSERT INTO macros (carbs, protein, fat, alcohol)
			VALUES ($1, $2, $3, $4 ) RETURNING id`

	fmt.Println(sqlStatement)

	err := s.Db.QueryRow(sqlStatement, m.Carbs, m.Protein, m.Fat, m.Alcohol).Scan(&retUUID)

	return retUUID.UUID.String(), err
}

func (s *postgres) ReadAll() ([]types.Macro, error) {

	macros := make([]types.Macro, 0)

	rows, err := s.Db.Query("SELECT * FROM macros")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var m types.Macro
		if err := rows.Scan(&m.ID, &m.Carbs, &m.Protein, &m.Fat, &m.Alcohol, &m.Date); err != nil {
			return nil, err
		}
		macros = append(macros, m)
	}

	return macros, nil
}

func (s *postgres) Read(_uuid uuid.UUID) (*types.Macro, error) {

	retMacro := types.Macro{}

	sqlStatement := `SELECT carbs, protein, fat, alcohol, date FROM macros WHERE id=$1;`

	err := s.Db.QueryRow(sqlStatement, _uuid).Scan(&retMacro.Carbs, &retMacro.Protein, &retMacro.Fat, &retMacro.Alcohol, &retMacro.Date)

	return &retMacro, err
}

func (s *postgres) Update(_uuid uuid.UUID, m types.Macro) error {

	//u_id, err := strconv.Atoi(vars["macro_uuid"])

	//_, err :=
	//		db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
	//			p.Name, p.Price, p.ID)

	sqlStatement := `
UPDATE macros SET carbs=$1, protein=$2, fat=$3, alcohol=$4 WHERE id = $5 RETURNING id`

	//retMacro := &types.Macro{m.Carbs, m.Protein, m.Fat, m.Alcohol, m.Date}
	fmt.Println(sqlStatement)

	retUUID := uuid.NullUUID{}

	err := s.Db.QueryRow(sqlStatement, m.Carbs, m.Protein, m.Fat, m.Alcohol, _uuid).Scan(&retUUID)

	return err
}
func (s *postgres) Delete(u uuid.UUID) (int64, error) {

	result, err := s.Db.Exec("DELETE FROM macros WHERE id=$1", u)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

//read env variables``
