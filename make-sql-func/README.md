# Auto-make MySql function

## Use Ways (two ways)
- Test run
1. go get package
```shell
$ go get github.com/garallz/Go/make-sql-func
```

2. make a test function to generate code.
```go
import (
	"fmt"
	"testing"

	"github.com/garallz/Go/make-sql-func"
)

func TestMakeSqlFunction(t *testing.T) {
	err := MakeSqlFunction("envFile", "")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("Make sql function success!")
	}
}
```

- Build run
1. go get package and install.
```shell
$ go get -d github.com/garallz/Go/make-sql-func
```

2. run executable file to make the sql function.
```shell
// default: generate under the current path

$ ./make-sql-func envFile
```

-----------------------------------------------------------------------------
// Env file set.
```json
{
  	"package": "main",
  	"data":[
	    {
	  		"message":"Testing",
	      	"table":"test",
	      	"fields":[
		        {
		          	"name": "id",
		          	"type": "int"
		        },
		        {
		          	"name": "user_name",
		          	"type": "string"
		        },
		        {
		          	"name": "user_age",
		          	"type": "int"
		        }
	      	],
	      	"index":"id",
			"autogrow": "id",
			"unique":["id", "user_name"],
	      	"model": []
	    }
  	]
}
```

- `package`: go file package name.
- `data`: table array.
    - `message`: table message.
    - `table`: database table name.
    - `fields`: table fields array.
        - `name`: field name.
        - `type`: field type (golang type eg: int, string, time.Time...)
    - `index`: table index field.
	- `unique`: table unique index.
    - `model`: make function model. (default null mean all)

```

`model`
	1. insert function
	2. delete function
	3. update function
	4. select function
	5. insert into with on duplicate key update
```