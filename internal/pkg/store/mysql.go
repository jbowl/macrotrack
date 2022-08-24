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
	DB *sql.DB
}

// macro_uuid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
// date  TIMESTAMP WITH TIME ZONE default CURRENT_TIMESTAMP
const tableCreationQueryMysql = `CREATE TABLE IF NOT EXISTS macros
(
	macro_uuid BINARY(16) PRIMARY KEY,	
    carbs int NOT NULL,
	protein int NOT NULL,
	fat int NOT NULL,
	alcohol int NOT NULL,
	date datetime default CURRENT_TIMESTAMP
	
)`

func (s *MysqlDB) createTable() error {

	//	query := `CREATE TABLE IF NOT EXISTS product(product_id int primary key auto_increment, product_name text,
	//        product_price int, created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP)`

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

	dsn := "root:mysql@tcp(127.0.0.1:3306)/mysql"
	s.DB, err = sql.Open("mysql", dsn)

	//	s.DB, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/test")

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

	retUUID := uuid.NullUUID{}

	sqlStatement := `
			INSERT INTO macros (carbs, protein, fat, alcohol)
			VALUES ($1, $2, $3, $4 ) RETURNING macro_uuid`

	fmt.Println(sqlStatement)

	err := s.DB.QueryRow(sqlStatement, m.Carbs, m.Protein, m.Fat, m.Alcohol).Scan(&retUUID)

	return retUUID.UUID.String(), err

	return "", nil
}

func (s *MysqlDB) Read(uuid.UUID) (*types.Macro, error) {
	return nil, nil
}

func (s *MysqlDB) Update() error {
	return nil
}

func (s *MysqlDB) Delete(uuid.UUID) error {
	return nil
}
