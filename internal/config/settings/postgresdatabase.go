package settings

import (
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

type PostgresDatabase struct {
	Address  string
	User     string
	Password string
	Database string
}

func (p *PostgresDatabase) setDefaults() {
	p.Address = gosettings.DefaultComparable(p.Address, "psql:5432")
	p.User = gosettings.DefaultComparable(p.User, "postgres")
	p.Password = gosettings.DefaultComparable(p.Password, "postgres")
	p.Database = gosettings.DefaultComparable(p.Database, "postgres")
}

func (p *PostgresDatabase) validate() (err error) {
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

func (p *PostgresDatabase) overrideWith(other PostgresDatabase) {
	p.Address = gosettings.OverrideWithComparable(p.Address, other.Address)
	p.User = gosettings.OverrideWithComparable(p.User, other.User)
	p.Password = gosettings.OverrideWithComparable(p.Password, other.Password)
	p.Database = gosettings.OverrideWithComparable(p.Database, other.Database)
}

func (p *PostgresDatabase) read(r *reader.Reader) {
	p.Address = r.String("POSTGRES_ADDRESS")
	p.User = r.String("POSTGRES_USER", reader.ForceLowercase(false))
	p.Password = r.String("POSTGRES_PASSWORD", reader.ForceLowercase(false))
	p.Database = r.String("POSTGRES_DATABASE", reader.ForceLowercase(false))
}
