// https://www.cdata.com/kb/tech/snowflake-odbc-go-linux.rst
// xtore client  to  xstore server which is container based snowflake

package store

import (
	"context"
	"fmt"
	"macrotrack/internal/pkg/types"

	"github.com/google/uuid"

	"database/sql"

	"github.com/snowflakedb/gosnowflake"
)

//const create_snwflk = `INSERT INTO MACROS(id, carbs, protein, fat, alcohol) VALUES (?, ?, ?, ? , ?);`
//const readall_snwflk = `SELECT id, carbs, protein, fat, alcohol, date FROM MACROS;`

type snwflk struct {
	DSN string
	Db  *sql.DB
}

const tableCreationQuery_snwflk = `CREATE TABLE IF NOT EXISTS macros
(
   carbs int NOT NULL,
	protein int NOT NULL,
	fat int NOT NULL,
	alcohol int NOT NULL,
	date  TIMESTAMP WITH TIME ZONE default CURRENT_TIMESTAMP
);`

func (s *snwflk) Init() error {

	config := gosnowflake.Config{Account: "VXB59306",
		User:     "JOELH",
		Password: "R4nd0mCr4p",
		Database: "MACROS",
		Schema:   "MACROS_SCHEMA",
		Role:     "ACCOUNTADMIN",
		Region:   "us-east-1",
	}

	dsn, err := gosnowflake.DSN(&config)

	if err != nil {
		return err
		//log.Fatal(err)
	}

	//var err error
	//s.Db, err = sql.Open("snowflake", "user:password@my_organization-my_account/mydb")
	//s.Db, err = sql.Open("snowflake", "JOELH:R4nd0mCr4p@VXB59306/macros")
	s.Db, err = sql.Open("snowflake", dsn)
	if err != nil {
		return err
		//log.Fatal(err)
	}
	//defer db.Close()

	return nil
}

func (s *snwflk) Open() error {
	return nil
}

func (s *snwflk) Create(m types.Macro) (string, error) {
	//	stmt, err := s.Db.Prepare(create_snwflk)
	//	if err != nil {
	//		fmt.Printf("--> Prrepare query Error occured ")

	//		return "", err
	//	}
	//
	//	defer stmt.Close()

	//res, err := stmt.Exec(m.Carbs, m.Protein, m.Fat, m.Alcohol)

	var dbContext = context.Background()

	u := uuid.New()

	id := u.String()

	newRecord, err := s.Db.ExecContext(
		dbContext,
		createMS, // TODO: using same string as MS
		id,
		m.Carbs,
		m.Protein,
		m.Fat,
		m.Alcohol,
	)

	if err != nil {
		return "", err
	}

	fmt.Println(newRecord.RowsAffected())

	return id, err
}

func (s *snwflk) ReadAll() ([]types.Macro, error) {
	var dbContext = context.Background()
	err := s.Db.PingContext(dbContext)
	if err != nil {
		return nil, err
	}

	m := make([]types.Macro, 0)

	data, queryErr := s.Db.QueryContext(dbContext, readAllMS) // TODO: same string as MS
	if queryErr != nil {
		return nil, queryErr
	}

	for data.Next() {
		t := types.Macro{}

		var id string

		nErr := data.Scan(&id, &t.Carbs, &t.Protein, &t.Fat, &t.Alcohol, &t.Date)
		if nErr != nil {
			return nil, nErr
		}

		t.ID, nErr = uuid.Parse(id)
		// ignore err here
		if nErr != nil {
			fmt.Println(nErr)
		}

		m = append(m, t)
	}

	return m, err
}

func (s *snwflk) Read(u uuid.UUID) (*types.Macro, error) {

	//	binaryUUID, err := u.MarshalBinary()

	strUUID := u.String()

	retMacro := types.Macro{}
	// Execute the query

	var readUUID string

	err := s.Db.QueryRow(readByIDMS, strUUID).Scan(&readUUID, &retMacro.Carbs, &retMacro.Protein, &retMacro.Fat, &retMacro.Alcohol, &retMacro.Date)
	if err != nil {
		return nil, err
	}

	retMacro.ID, err = uuid.Parse(readUUID)

	if err != nil {
		return nil, err
	}

	return &retMacro, err

}

func (s *snwflk) Update(_uuid uuid.UUID, m types.Macro) error {
	return nil
}

func (s *snwflk) Delete(u uuid.UUID) (int64, error) {

	var dbContext = context.Background()

	//	u := uuid.New()

	id := u.String()

	newRecord, err := s.Db.ExecContext(
		dbContext,
		deleteByIDMS,
		id,
	)

	if err != nil {
		return -1, err
	}

	return newRecord.RowsAffected()

}
