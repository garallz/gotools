package test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/garallz/Go/reflect/table"
	_ "github.com/go-sql-driver/mysql"
)

type DataStruct struct {
	Name      string    `db:"name"`
	Phone     int       `db:"phone"`
	Addr      string    `db:"addr"`
	Timestamp time.Time `db:"timestamp"`
}

func OpenMysql(user, passwd, database string) *sql.DB {
	// username:password@protocol(address)/dbname?param=value
	var str = fmt.Sprintf("%s:%s@/%s", user, passwd, database)
	db, err := sql.Open("mysql", str)
	if err != nil {
		panic(err)
	} else if db.Ping() != nil {
		panic("sql can not connection")
	}
	return db
}

func (d *DataStruct) DB() *table.TableStruct {
	db := OpenMysql("root", "Gara123", "lin")
	table := table.NewTable(d, "tabletest", "db", db)
	//	table.SetDBName(table.MYSQL_TYPE)
	table.SetIndex("Name")
	return table
}

func TestExecUse(t *testing.T) {
	var data = &DataStruct{
		Name:      "Gara",
		Phone:     12356788765,
		Addr:      "China",
		Timestamp: time.Now(),
	}

	if err := data.DB().Insert(); err != nil {
		t.Error(err)
	}

	var where = "name = 'Gara'"
	data = &DataStruct{}

	if err := data.DB().Query(where); err != nil {
		log.Panicln(err)
	} else {
		fmt.Println(data)
	}

	if rows, err := data.DB().QueryArray(where); err != nil {
		log.Panicln(err)
	} else {
		fmt.Println(rows)
	}
}

func TestQueryNullValues(t *testing.T) {
	var where = "name = 'Gara'"
	var data = &DataStruct{}

	if err := data.DB().QueryWithNull(where); err != nil {
		log.Panicln(err)
	} else {
		fmt.Println(data)
	}

	if rows, err := data.DB().QueryArrayWithNull(where); err != nil {
		log.Panicln(err)
	} else {
		fmt.Println(rows)
	}
}
