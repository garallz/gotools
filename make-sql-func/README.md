# Auto-make MySql function

## Use Ways (two ways)
- Test run
1. go get package
```shell
$ go get -u github.com/garallz/Go/make-sql-func
```

2. make a test function to generate code.
```go
import (
	"fmt"
	"testing"

	"github.com/garallz/Go/make-sql-func"
)

func TestMakeSqlFunction(t *testing.T) {
	err := MakeSqlFunction(profile, "")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("Make sql function success!")
	}
}
```

- Build run
1. go get package.
```shell
$ go get -u github.com/garallz/Go/make-sql-func
```

2. build or install make-sql-func
```shell
$ cd $GOPATH/src/github.com/garallz/Go/make-sql-func

$ go install	// make in $GOPATH/bin
$ go build		// make in current path
```

3. run executable file to make the sql function.
```shell
$ ./make-sql-func profile
```

`default: function file generate under the current path`

-----------------------------------------------------------------------------
// Profile set.
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