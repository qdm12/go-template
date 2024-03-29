package settings

import (
	"fmt"

	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

type Settings struct {
	HTTP     HTTP
	Metrics  Metrics
	Log      Log
	Database Database
	Health   Health
}

func (s *Settings) SetDefaults() {
	s.HTTP.setDefaults()
	s.Metrics.setDefaults()
	s.Log.setDefaults()
	s.Database.setDefaults()
	s.Health.SetDefaults()
}

func (s *Settings) Validate() (err error) {
	nameToValidation := map[string]func() error{
		"http server": s.HTTP.validate,
		"metrics":     s.Metrics.validate,
		"logging":     s.Log.validate,
		"database":    s.Database.validate,
		"health":      s.Health.Validate,
	}

	for name, validation := range nameToValidation {
		err = validation()
		if err != nil {
			return fmt.Errorf("%s settings: %w", name, err)
		}
	}

	return nil
}

func (s *Settings) String() string {
	return s.toLinesNode().String()
}

func (s *Settings) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Settings summary:")
	node.AppendNode(s.HTTP.toLinesNode())
	node.AppendNode(s.Metrics.toLinesNode())
	node.AppendNode(s.Log.toLinesNode())
	node.AppendNode(s.Database.toLinesNode())
	node.AppendNode(s.Health.toLinesNode())
	return node
}

func (s *Settings) Copy() (copied Settings) {
	return Settings{
		HTTP:     s.HTTP.copy(),
		Metrics:  s.Metrics.copy(),
		Log:      s.Log.copy(),
		Database: s.Database.copy(),
		Health:   s.Health.copy(),
	}
}

func (s *Settings) OverrideWith(other Settings) {
	s.HTTP.overrideWith(other.HTTP)
	s.Metrics.overrideWith(other.Metrics)
	s.Log.overrideWith(other.Log)
	s.Database.overrideWith(other.Database)
	s.Health.overrideWith(other.Health)
}

func (s *Settings) Read(r *reader.Reader) (err error) {
	err = s.HTTP.read(r)
	if err != nil {
		return fmt.Errorf("HTTP server settings: %w", err)
	}

	s.Metrics.read(r)
	s.Log.read(r)
	s.Database.read(r)
	s.Health.Read(r)

	return nil
}
