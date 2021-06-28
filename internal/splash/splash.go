// Package splash contains functions to log out splash initial information.
package splash

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/go-template/internal/constants"
	"github.com/qdm12/go-template/internal/models"
)

// Splash returns the welcome spash message.
func Splash(buildInfo models.BuildInformation) string {
	lines := title()
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Running version %s built on %s (commit %s)",
		buildInfo.Version, buildInfo.BuildDate, buildInfo.Commit))
	lines = append(lines, "")
	lines = append(lines, announcement()...)
	lines = append(lines, "")
	lines = append(lines, links()...)
	return strings.Join(lines, "\n")
}

func title() []string {
	return []string{
		"=========================================",
		"=========================================",
		"============== go-template ==============",
		"=========================================",
		"=== Made with ❤️  by github.com/qdm12 ====",
		"=========================================",
	}
}

func announcement() []string {
	if len(constants.Annoucement) == 0 {
		return nil
	}
	expirationDate, _ := time.Parse("2006-01-02", constants.AnnoucementExpiration) // error covered by a unit test
	if time.Now().After(expirationDate) {
		return nil
	}
	return []string{"📣" + constants.Annoucement}
}

func links() []string {
	return []string{
		"🔧 Need help? " + constants.IssueLink,
		"💻 Email? quentin.mcgaw@gmail.com",
		"☕ Slack? Join from the Slack button on Github",
		"💰 Help me? https://github.com/sponsors/qdm12",
	}
}
