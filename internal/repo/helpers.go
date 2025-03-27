package repo

import (
	"os"
	"regexp"

	"github.com/not-for-prod/implgen/pkg/clog"
	"golang.org/x/mod/modfile"
)

func getModuleName() string {
	// Read go.mod file
	data, err := os.ReadFile("./go.mod")
	if err != nil {
		clog.Fatal(err.Error())
	}

	// Parse go.mod content
	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		clog.Fatal(err.Error())
	}

	// Return module name
	return modFile.Module.Mod.Path
}

// parseSQLXComment extracts the sqlx directive from a comment
func parseSQLXComment(comment string) string {
	re := regexp.MustCompile(`sqlx:\s*(\w+)`)
	match := re.FindStringSubmatch(comment)
	if len(match) > 1 {
		return match[1] // Extract the actual value (e.g., "ExecContext")
	}
	return "" // No match found
}
