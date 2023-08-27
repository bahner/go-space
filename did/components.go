package did

import (
	"strconv"
	"strings"
)

// extractComponents splits the identifier into its components
func extractComponents(identifier string) (string, string, string, string) {
	components := strings.Split(identifier, ":")
	if len(components) == 3 {
		return components[0], components[1], "1", components[2]
	}
	return components[0], components[1], components[2], components[3]
}

// isValidComponents validates that the components meet certain criteria
func isValidComponents(scheme, method, version, multibaseValue string) bool {
	return scheme == "did" && method == "key" && isValidVersion(version) && multibaseValue[0] == 'z'
}

// isValidVersion validates the version string
func isValidVersion(version string) bool {
	v, err := strconv.Atoi(version)
	return err == nil && v > 0
}
