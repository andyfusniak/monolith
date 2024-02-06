# monolith (Universal monolith web app writtin in Go)

This repository provides a skeleton web application written in Go.

To compile, run:

```shell
$ make
```


## Environment Variables

| Env Var           | Required | Default | Description |
| ----------------- | -------- | ------- | ------------|
| **`PORT`**        | Optional | 8080    | Port for the app service to listen on. |
| **`DB_FILEPATH`** | Required |         | Fullpath to the sqlite3 database file. |

```shell
$ export DB_FILEPATH='./monolith.db'
```
