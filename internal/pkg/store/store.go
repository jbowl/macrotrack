//////go:build prod
////// +build prod

package store

import (
	"fmt"

	"macrotrack/internal/pkg/types"

	"github.com/google/uuid"
)

type Storage interface {
	Init() error
	Open() error
	//Create(types.Macro) (uuid.NullUUID, error)
	Create(types.Macro) (string, error)
	Read(uuid.UUID) (*types.Macro, error)
	ReadAll() ([]types.Macro, error)
	Update(uuid.UUID, types.Macro) error
	Delete(uuid.UUID) (int64, error)
}

func GetStorage(storeagetype string, dsn string) Storage {

	var db Storage

	switch storeagetype {

	case "mongo":
		//db = &mongoDB

	case "sqlserver":
		db = &sqlserver{DSN: dsn}

	case "postgres":
		db = &postgres{DSN: dsn}

	}

	fmt.Println(storeagetype)

	//db := &sqlserver{DSN: dsn}

	//db := &mongoDB{}
	//db := &MysqlDB{}
	return db

}
