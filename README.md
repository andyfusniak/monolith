# monolith (Universal monolith web app written in Go)

This repository provides a skeleton Go HTTP service that uses Sqlite 3 as
its underlying relational database. The service is designed to be a monolith
web app that can be used as a starting point for a variety of web
applications.

To compile, run:

```shell
$ make
```


## Environment Variables

| Env Var           | Required | Default | Description                            |
| ----------------- | -------- | ------- | -------------------------------------- |
| **`PORT`**        | Optional | 8080    | Port for the app service to listen on. |
| **`DB_FILEPATH`** | Required |         | Fullpath to the sqlite3 database file. |

```shell
$ export DB_FILEPATH='./monolith.db'
```
