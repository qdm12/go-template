package settings

import (
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid"
)

type PostgresDatabase struct {
	Address  string
	User     string
	Password string
	Database string
}

func (p *PostgresDatabase) setDefaults() {
	p.Address = gosettings.DefaultString(p.Address, "psql:5432")
	p.User = gosettings.DefaultString(p.User, "postgres")
	p.Password = gosettings.DefaultString(p.Password, "postgres")
	p.Database = gosettings.DefaultString(p.Database, "postgres")
}

func (p *PostgresDatabase) validate() (err error) {
	err = govalid.ValidateAddress(p.Address)
	if err != nil {
		return fmt.Errorf("connection address: %w", err)
	}

	return nil
}

func (p *PostgresDatabase) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Postgres database settings:")
	node.Appendf("Connection address: %s", p.Address)
	node.Appendf("User: %s", p.User)
	node.Appendf("Password: %s", obfuscatePassword(p.Password))
	node.Appendf("Database name: %s", p.Database)
	return node
}

func (p *PostgresDatabase) copy() (copied PostgresDatabase) {
	return PostgresDatabase{
		Address:  p.Address,
		User:     p.User,
		Password: p.Password,
		Database: p.Database,
	}
}

func (p *PostgresDatabase) mergeWith(other PostgresDatabase) {
	p.Address = gosettings.MergeWithString(p.Address, other.Address)
	p.User = gosettings.MergeWithString(p.User, other.User)
	p.Password = gosettings.MergeWithString(p.Password, other.Password)
	p.Database = gosettings.MergeWithString(p.Database, other.Database)
}

func (p *PostgresDatabase) overrideWith(other PostgresDatabase) {
	p.Address = gosettings.OverrideWithString(p.Address, other.Address)
	p.User = gosettings.OverrideWithString(p.User, other.User)
	p.Password = gosettings.OverrideWithString(p.Password, other.Password)
	p.Database = gosettings.OverrideWithString(p.Database, other.Database)
}
