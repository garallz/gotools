# Auto-make Sql function

```json
{
  "package": "sqlFunc",
  "data":[
    {
      "table":"node",
      "fields":[
        {
          "name": "xid",
          "type": "int"
        }
      ],
      "index":"xid",
      "model": [1, 2, 3], 
      "message":"测试测试"
    }
  ]
}
```

- `package`: go file package name.
- `data`: table array.
    - `table`: database table name.
    - `fields`: table fields array.
        - `name`: field name.
        - `type`: field type (golang type eg: int, string, time.Time...)
    - `index`: table index field.
    - `model`: make function model.
    - `message`: function message.