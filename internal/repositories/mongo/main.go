package mongo

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Connect(uri string, dbName string) error {
	err := mgm.SetDefaultConfig(&mgm.Config{CtxTimeout: 12 * time.Second}, dbName, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	return nil
}
