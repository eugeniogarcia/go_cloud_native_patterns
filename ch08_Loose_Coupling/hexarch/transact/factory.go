package transact

import (
	"fmt"

	"egsmartin.com/hexarch/core"
)

//Crea una factoria que retorna un TransactionLogger. Tenemos tres implementaciones que implementan esta interface
func NewTransactionLogger(s string) (core.TransactionLogger, error) {
	switch s {
	case "test":
		return NewTestTransactionLogger()

	case "file":
		return NewFileTransactionLogger("./transactions.txt")

	case "postgres":
		params := PostgresDbParams{
			host: "localhost", dbName: "kvs",
			user: "test", password: "hunter2",
		}
		return NewPostgresTransactionLogger(params)

	default:
		return nil, fmt.Errorf("no such transaction logger %s", s)
	}
}
