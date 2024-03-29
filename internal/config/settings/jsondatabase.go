package settings

import (
	"errors"
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

type JSONDatabase struct {
	Filepath string
}

func (j *JSONDatabase) setDefaults() {
	j.Filepath = gosettings.DefaultComparable(j.Filepath, "data.json")
}

var (
	ErrJSONFilepathIsDirectory = errors.New("JSON filepath is a directory")
)

func (j *JSONDatabase) validate() (err error) {
	stats, err := os.Stat(j.Filepath)
	if err != nil {
		return fmt.Errorf("file path: %w", err)
	} else if stats.IsDir() {
		return fmt.Errorf("%w: %s", ErrJSONFilepathIsDirectory, j.Filepath)
	}
	return nil
}

func (j *JSONDatabase) toLinesNode() (node *gotree.Node) {
	node = gotree.New("JSON database settings:")
	node.Appendf("File path: %s", j.Filepath)
	return node
}

func (j *JSONDatabase) copy() (copied JSONDatabase) {
	return JSONDatabase{
		Filepath: j.Filepath,
	}
}

func (j *JSONDatabase) overrideWith(other JSONDatabase) {
	j.Filepath = gosettings.OverrideWithComparable(j.Filepath, other.Filepath)
}

func (j *JSONDatabase) read(r *reader.Reader) {
	j.Filepath = r.String("JSON_FILEPATH")
}
