package splash

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/qdm12/REPONAME_GITHUB/internal/constants"
)

// Splash returns the welcome spash message
func Splash(version, vcsRef, buildDate string) string {
	lines := title()
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Running version %s built on %s (commit %s)", version, buildDate, vcsRef))
	lines = append(lines, "")
	lines = append(lines, annoucement()...)
	lines = append(lines, "")
	lines = append(lines, links()...)
	return strings.Join(lines, "\n")
}

func title() []string {
	return []string{
		"=========================================",
		"=========================================",
		"========= REPONAME_GITHUB ========",
		"=========================================",
		"=== Made with " + emoji.Sprint(":heart:") + " by github.com/qdm12 ====",
		"=========================================",
	}
}

func annoucement() []string {
	if len(constants.Annoucement) == 0 {
		return nil
	}
	expirationDate, _ := time.Parse("2006-01-02", constants.AnnoucementExpiration) // error covered by a unit test
	if time.Now().After(expirationDate) {
		return nil
	}
	return []string{emoji.Sprint(":mega: ") + constants.Annoucement}
}

func links() []string {
	return []string{
		emoji.Sprint(":wrench: ") + "Need help? " + constants.IssueLink,
		emoji.Sprint(":computer: ") + "Email? quentin.mcgaw@gmail.com",
		emoji.Sprint(":coffee: ") + "Slack? Join from the Slack button on Github",
		emoji.Sprint(":money_with_wings: ") + "Help me? https://github.com/sponsors/qdm12",
	}
}
