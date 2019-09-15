# Database Reflect Exec

```
Convenient for native database/sql use
（方便原生的 database/sql 使用）

database: ["mysql", "postgresql"]
（适用数据库：MySQL， POSTGRESQL)

Function: [Query, Insert, Update]
（功能如下： 查询， 插入， 更新）
```

## Only for Go1.13 Version and above

(因使用了Go1.13的部分功能，所以需要使用Go1.13及以上的版本))

## Example:

```go
type DataStruct struct {
    Data string `db:"data"`
    Name string `db:"name"`
}

func (d *DataStruct) DB() *TableStruct {
	table := NewTable(d, "data", "db", nil)
	table.SetDBName(MYSQL_TYPE)
	table.SetIndex("Name")
	return table
}

func main() {
    var data = &DataStruct{}
	if err := data.DB().Query("name = 'gara'"); err != nil {
		log.Panicln(err)
	} else {
		fmt.Println(data)
	}
}
```