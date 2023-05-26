package settings

import (
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

const (
	MemoryStoreType   = "memory"
	JSONStoreType     = "json"
	PostgresStoreType = "postgres"
)

type Database struct {
	Type     *string
	Memory   MemoryDatabase
	JSON     JSONDatabase
	Postgres PostgresDatabase
}

func (d *Database) setDefaults() {
	d.Type = ptrTo(MemoryStoreType)
	d.Memory.setDefaults()
	d.JSON.setDefaults()
	d.Postgres.setDefaults()
}

var (
	ErrDatabaseTypeUnknown = fmt.Errorf("database type is unknown")
)

func (d *Database) validate() (err error) {
	switch *d.Type {
	case MemoryStoreType:
		err = d.Memory.validate()
		if err != nil {
			return fmt.Errorf("memory database: %w", err)
		}
	case JSONStoreType:
		err = d.JSON.validate()
		if err != nil {
			return fmt.Errorf("json database: %w", err)
		}
	case PostgresStoreType:
		err = d.Postgres.validate()
		if err != nil {
			return fmt.Errorf("postgres database: %w", err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrDatabaseTypeUnknown, *d.Type)
	}

	return nil
}

func (d *Database) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Database settings:")
	node.Appendf("Type: %s", *d.Type)
	switch *d.Type {
	case MemoryStoreType:
		node.AppendNode(d.Memory.toLinesNode())
	case JSONStoreType:
		node.AppendNode(d.JSON.toLinesNode())
	case PostgresStoreType:
		node.AppendNode(d.Postgres.toLinesNode())
	}
	return node
}

func (d *Database) copy() (copied Database) {
	return Database{
		Type:     gosettings.CopyPointer(d.Type),
		Memory:   d.Memory.copy(),
		JSON:     d.JSON.copy(),
		Postgres: d.Postgres.copy(),
	}
}

func (d *Database) mergeWith(other Database) {
	d.Type = gosettings.MergeWithPointer(d.Type, other.Type)
	d.Memory.mergeWith(other.Memory)
	d.JSON.mergeWith(other.JSON)
	d.Postgres.mergeWith(other.Postgres)
}

func (d *Database) overrideWith(other Database) {
	d.Type = gosettings.OverrideWithPointer(d.Type, other.Type)
	d.Memory.overrideWith(other.Memory)
	d.JSON.overrideWith(other.JSON)
	d.Postgres.overrideWith(other.Postgres)
}
