// maria db client

//$ mariadb --host 127.0.0.1 --port 3306 --user root --password mysql
// mysql running in container

// https://www.mysqltutorial.org/mysql-uuid/   UUID as pkey
// https://golangbot.com/mysql-create-table-insert-row/

package store

import (
	"context"
	"fmt"
	"log"
	"macrotrack/internal/pkg/types"
	"time"

	"github.com/google/uuid"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlDB struct {
	DB  *sql.DB
	DSN string
}

const createMS = `INSERT INTO macros (id, carbs, protein, fat, alcohol) VALUES (?, ?, ?, ?,?);`
const readAllMS = `SELECT id, carbs, protein, fat, alcohol, date FROM macros`
const readByIDMS = `SELECT id, carbs, protein, fat, alcohol, date FROM macros WHERE id = ?`
const updateByIDMS = `UPDATE macros SET carbs=?, protein=?, fat=?, alcohol=? WHERE id = ?`
const deleteByIDMS = `DELETE from macros WHERE id = ?;`

// SELECT BIN_TO_UUID(id) id , carbs, protein, fat, alcohol, `date`
//
//	FROM macros.macros;
const tableCreationQueryMysql = `CREATE TABLE IF NOT EXISTS macros
(
	id BINARY(16) PRIMARY KEY,	
    carbs int NOT NULL,
	protein int NOT NULL,
	fat int NOT NULL,
	alcohol int NOT NULL,
	date datetime default CURRENT_TIMESTAMP	
)`

func (s *MysqlDB) createTable() error {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := s.DB.ExecContext(ctx, tableCreationQueryMysql)
	if err != nil {
		log.Printf("Error %s when creating product table", err)
		return err
	}

	fmt.Println(res)

	return err

}

//	func dsn(dbName string) string {
//	  return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
//	}
func (s *MysqlDB) Init() error {

	var err error

	s.DSN = "gouser:gopwd@tcp(127.0.0.1:3306)/macros"

	//dsn := "root:mysql@tcp(127.0.0.1:3306)/mysql"
	s.DB, err = sql.Open("mysql", s.DSN)

	if err != nil {
		return err

	}

	err = s.createTable()

	// if there is an error opening the connection, handle it
	return err
}

func (s *MysqlDB) Open() error {
	return nil
}

func (s *MysqlDB) Create(m types.Macro) (string, error) {

	retUUID := uuid.New()
	b, err := retUUID.MarshalBinary()

	if err != nil {
		return retUUID.String(), err
	}

	// fails every time, or more likely driver bug reports failure when works successfully
	row := s.DB.QueryRow(createMS, b, m.Carbs, m.Protein, m.Fat, m.Alcohol)
	row.Err()

	// secondary scan to check success and get key of new row
	row = s.DB.QueryRow("select id from macros where id = ?", b)
	var uuidBytes []byte

	err = row.Scan(&uuidBytes)
	if err != nil {
		return retUUID.String(), err
	}

	// assert uuid == retUUID

	uuid, err := uuid.FromBytes(uuidBytes)
	return uuid.String(), err

}

// ReadAll - returns array of types.Macro or error
func (s *MysqlDB) ReadAll() ([]types.Macro, error) {

	m := make([]types.Macro, 0)

	data, err := s.DB.Query(readAllMS)

	if err != nil {
		return m, err
	}

	for data.Next() {

		t := types.Macro{}

		if err := data.Scan(&t.ID, &t.Carbs, &t.Protein, &t.Fat, &t.Alcohol, &t.Date); err != nil {
			return m, err
		}

		m = append(m, t)
	}

	return m, nil
}

func (s *MysqlDB) Read(u uuid.UUID) (*types.Macro, error) {

	binaryUUID, err := u.MarshalBinary()

	if err != nil {
		return nil, err
	}

	retMacro := types.Macro{}
	// Execute the query
	err = s.DB.QueryRow(readByIDMS, binaryUUID).Scan(&retMacro.ID, &retMacro.Carbs, &retMacro.Protein, &retMacro.Fat, &retMacro.Alcohol, &retMacro.Date)

	if err != nil {
		return nil, err
	}

	return &retMacro, err
}

func (s *MysqlDB) Update(u uuid.UUID, m types.Macro) error {

	binaryUUID, err := u.MarshalBinary()

	if err != nil {
		return err
	}

	res, err := s.DB.Exec(updateByIDMS, m.Carbs, m.Protein, m.Fat, m.Alcohol, binaryUUID)
	if err != nil {
		return err
	}

	fmt.Println(res.RowsAffected())

	return err

}

func (s *MysqlDB) Delete(u uuid.UUID) (int64, error) {

	binaryUUID, err := u.MarshalBinary()

	if err != nil {
		return -1, err
	}
	result, err := s.DB.Exec(deleteByIDMS, binaryUUID)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
