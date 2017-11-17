# Auto-make Sql function

```json
{
  "package": "sqlFunc",
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
          "name": "name",
          "type": "int"
        },
        {
          "name": "age",
          "type": "int"
        }
      ],
      "index":"id",
	  "unique":["id","name"],
      "model": [1, 2, 3]
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


`model`:
	1. insert function
	2. delete function
	3. update function
	4. select function
	5. insert into with on duplicate key update