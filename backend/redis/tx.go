package redis

import "github.com/matobet/verdi/env"

type Tx struct {
	Conn
}

var _ env.RedisTx = (*Tx)(nil)

func (tx *Tx) Begin() env.RedisTx {
	tx.Conn.Send("MULTI")
	return tx
}

func (tx *Tx) Commit() error {
	_, err := tx.Conn.Do("EXEC")
	return err
}
