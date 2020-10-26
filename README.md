<h1 align="center">Welcome to gosqlexec 👋</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-1.0.1-blue.svg?cacheSeconds=2592000" />
  <img alt="documentation: yes" src="https://img.shields.io/badge/Documentation-Yes-green.svg" />
  <img alt="maintained: yes" src="https://img.shields.io/badge/Maintained-Yes-green.svg" />
</p>


>This package used for executing SQL file, through command line tools, you can attach command template from this package, into your command line app, and you can easily execute SQL for schema creations, dropping tables, altering tables, or you can write your own SQL syntax inside custom_query.sql file and you can apply the SQL execution into your database.



### How to install

```bash
go get -u github.com/muhammadisa/gosqlexec
```



### Create database SQL files structure like this

```
+---db
|   +---alter
|   |       alter_tables.sql
|   |
|   +---drop
|   |       drop_tables.sql
|   |
|   +---query
|   |       custom_query.sql
|   |
|   \---schemas
|           articles.sql
|           authors.sql
|           likes.sql
|           users.sql
```



### Create query executor definition

```go
qe := gosqlexec.GoSQLExec{
    AlterQuery:  "db/alter/alter_tables.sql",
    DropQuery:   "db/drop/drop_tables.sql",
    CustomQuery: "db/query/custom_query.sql",
    Schemas: []string{
        "db/schemas/users.sql",
        "db/schemas/authors.sql",
        "db/schemas/articles.sql",
        "db/schemas/likes.sql",
    },
}
```



### Create CLI for your project

This project is using CLI from [https://github.com/urfave/cli]()

```go
app := cli.NewApp()
app.Name = "Your Project Name"
app.Usage = "your project usages"

app.Commands = []cli.Command{
    gosqlexec.MigrateCommand(qe),
    gosqlexec.DropTablesCommand(qe),
    gosqlexec.AlterTablesCommand(qe),
    gosqlexec.CustomQueryExecCommand(qe),
    {
        Name:  "run-server",
        Usage: "Start Server",
        Action: func(c *cli.Context) error {
            // Calls start server function here
            return nil
        },
    },
}

err = app.Run(os.Args)
if err != nil {
    log.Fatal(err)
}
```

