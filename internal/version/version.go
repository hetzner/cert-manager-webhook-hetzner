package version

var (
	// version is a semver version (https://semver.org).
	version = "0.6.2" // x-releaser-pleaser-version

	// versionPrerelease is a semver version pre-release identifier (https://semver.org).
	//
	// For final releases, we set this to an empty string.
	versionPrerelease = "dev"

	// Version of the application.
	Version = func() string {
		if versionPrerelease != "" {
			return version + "-" + versionPrerelease
		}
		return version
	}()
)
