package store

// https://sqlchoice.azurewebsites.net/en-us/sql-server/developer-get-started/go/rhel/step/2.html
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"macrotrack/internal/pkg/types"

	"github.com/google/uuid"
	_ "github.com/microsoft/go-mssqldb"
)

const create = `INSERT INTO macros(carbs, protein, fat, alcohol) OUTPUT Inserted.id VALUES (@carbs, @protein, @fat , @alcohol);`
const readall = `SELECT * FROM macros;`
const read = `SELECT id, carbs, protein, fat, alcohol, date FROM macros WHERE id=@id;`
const update = `UPDATE macros SET carbs=@carbs, protein=@protein, fat=@fat, alcohol=@alcohol WHERE id = @id`
const delete = `DELETE from macros WHERE id = @id;`

var (
	Password = "___Aa123"
	User     = "SA"
	Port     = "1433"
	Database = "master"
	DSN      string
)

type sqlserver struct {
	DSN string
	Db  *sql.DB
}

// macro_uuid INT IDENTITY(1,1) NOT NULL PRIMARY KEY,
const tableCreationQuery_ss = `CREATE TABLE macros
(
	id uniqueIdentifier DEFAULT NEWID() PRIMARY KEY,	
    carbs int NOT NULL,
	protein int NOT NULL,
	fat int NOT NULL,
	alcohol int NOT NULL,
	date datetime default CURRENT_TIMESTAMP
);`

func ensureTableExistsSqlServer(db *sql.DB) error {
	if _, err := db.Exec(tableCreationQuery_ss); err != nil {
		return err
	}

	return nil
}

func (s *sqlserver) Init() error {

	//DSN = fmt.Sprintf("user id=%s;password=%s;port=%s;database=%s", User, Password, Port, Database)

	var err error

	s.Db, err = sql.Open("sqlserver", s.DSN)
	if err != nil {
		fmt.Println(fmt.Errorf("error opening database: %v", err))
		return err
	}

	if s.Db == nil {
		return errors.New("null db")
	}

	var dbContext = context.Background()

	err = s.Db.PingContext(dbContext)
	if err != nil {
		return err
	}

	err = ensureTableExistsSqlServer(s.Db)
	if err != nil {
		fmt.Print("error=", err)
		//return err
	}

	return nil
}

func (s *sqlserver) Open() error {
	return nil
}
func (s *sqlserver) Create(m types.Macro) (string, error) {

	stmt, err := s.Db.Prepare(create)
	if err != nil {
		fmt.Printf("--> Prrepare query Error occured ")

		return "", err
	}

	defer stmt.Close()

	//res, err := stmt.Exec(m.Carbs, m.Protein, m.Fat, m.Alcohol)

	var dbContext = context.Background()

	newRecord := stmt.QueryRowContext(dbContext,
		sql.Named("carbs", m.Carbs),
		sql.Named("protein", m.Protein),
		sql.Named("fat", m.Fat),
		sql.Named("alcohol", m.Alcohol),
	)
	//var retID int
	//err = newRecord.Scan(&retUUID)
	//err = newRecord.Scan(&retID)

	var b []byte

	err = newRecord.Scan(&b)

	if err != nil {
		fmt.Printf("--> Prrepare query Error occured ")

		return "", err
	}
	// https://github.com/denisenkom/go-mssqldb/issues/56
	// big/little endian fix
	b[0], b[1], b[2], b[3] = b[3], b[2], b[1], b[0]
	b[4], b[5] = b[5], b[4]
	b[6], b[7] = b[7], b[6]

	retUUID, err := uuid.FromBytes(b)

	//idStr, _ := uuid.Parse(theID)

	//return idStr.String(), err

	return retUUID.String(), err

	//return retUUID.UUID.String(), err
}

func (s *sqlserver) ReadAll() ([]types.Macro, error) {

	var dbContext = context.Background()
	err := s.Db.PingContext(dbContext)
	if err != nil {
		return nil, err
	}

	stmt, err := s.Db.Prepare(readall)
	if err != nil {
		fmt.Printf("--> Prrepare query Error occured ")

		return nil, err
	}

	defer stmt.Close()

	m := make([]types.Macro, 0)

	data, queryErr := s.Db.QueryContext(dbContext, readall)
	if queryErr != nil {
		return nil, queryErr
	}

	for data.Next() {
		t := types.Macro{}

		nErr := data.Scan(&t.ID, &t.Carbs, &t.Protein, &t.Fat, &t.Alcohol, &t.Date)
		if nErr != nil {
			return nil, nErr
		}

		m = append(m, t)
	}

	return m, err
}

func (s *sqlserver) Read(_uuid uuid.UUID) (*types.Macro, error) {

	var dbContext = context.Background()
	err := s.Db.PingContext(dbContext)
	if err != nil {
		return nil, err
	}

	retMacro := types.Macro{}

	stmt, err := s.Db.Prepare(read)
	if err != nil {
		fmt.Printf("--> Prrepare query Error occured ")

		return nil, err
	}

	defer stmt.Close()

	//res, err := stmt.Exec(m.Carbs, m.Protein, m.Fat, m.Alcohol)

	newRecord := stmt.QueryRowContext(dbContext,
		sql.Named("id", _uuid),
	)

	//u := uuid.UUID{}
	err = newRecord.Scan(&retMacro.ID, &retMacro.Carbs, &retMacro.Protein, &retMacro.Fat, &retMacro.Alcohol, &retMacro.Date)
	//err = newRecord.Scan(&retID)

	if err != nil {
		fmt.Printf("--> Prrepare query Error occured ")

		return nil, err
	}

	return &retMacro, err
}

func (s *sqlserver) Update(_uuid uuid.UUID, m types.Macro) error {

	stmt, err := s.Db.Prepare(update)
	if err != nil {
		fmt.Printf("--> Prrepare query Error occured ")

		return err
	}

	defer stmt.Close()
	var dbContext = context.Background()

	_, err = s.Db.ExecContext(dbContext,
		update,
		sql.Named("id", _uuid),
		sql.Named("carbs", m.Carbs),
		sql.Named("protein", m.Protein),
		sql.Named("fat", m.Fat),
		sql.Named("alcohol", m.Alcohol),
	)

	return err
}
func (s *sqlserver) Delete(u uuid.UUID) (int64, error) {

	var dbContext = context.Background()

	result, err := s.Db.ExecContext(dbContext, delete, sql.Named("id", u))
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
