package packman

import "fmt"

// PackageNotFoundError used to denote an unknown package.
type PackageNotFoundError struct {
	name string
}

func (e PackageNotFoundError) Error() string {
	return fmt.Sprintf("Package '%s' not found", e.name)
}

// UnknownDependentError used when a package specifies an unknown dependency.
type UnknownDependentError struct {
	pkg string
	dep string
}

func (e UnknownDependentError) Error() string {
	return fmt.Sprintf("Package '%s' has a dependency on unknown package '%s'", e.pkg, e.dep)
}

// PackageHasDependentsError used to identify cases where a package should
// not be removed because others depend on it.
type PackageHasDependentsError struct {
	name string
}

func (e PackageHasDependentsError) Error() string {
	return fmt.Sprintf("Package '%s' has dependents and cannot be removed", e.name)
}
