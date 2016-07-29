package packman

// Indexer is the interface our actual implementations must satisfy to be
// considered package indexers.
type Indexer interface {
	Index(p Package) error
	Remove(p Package) error
	Query(name string) (Package, error)
}

// Package represents a package to be indexed, removed and querried.
type Package struct {
	Name string
	Deps []*Package
}

// PackageIndexer implements the Indexer capabilities of adding, removing and
// querying for a package.
type PackageIndexer struct {
	dict map[string]*Package
}

// NewIndexer returns an implementation of our indexer.
func NewIndexer(initialDict map[string]*Package) PackageIndexer {
	pi := PackageIndexer{dict: initialDict}
	if pi.dict == nil {
		pi.dict = make(map[string]*Package, 0)
	}
	return pi
}

// Index indexes a package. If the package is already indexed, we update its dependencies.
// If any of the package's dependencies are unknown, we return an error.
func (pi *PackageIndexer) Index(p *Package) error {
	// check dependency exists
	for _, d := range p.Deps {
		_, err := pi.Query(d.Name)
		if _, ok := err.(PackageNotFoundError); ok {
			return UnknownDependentError{p.Name, d.Name}
		}
	}
	// add/update package
	pi.dict[p.Name] = p

	return nil
}

// Remove removes a package from the index if no other package depends on it.
// Returns a PackageHasDependentsError if any other package depends on it.
func (pi *PackageIndexer) Remove(p *Package) error {
	// ensure package is not depended upon by another
	for _, pkg := range pi.dict {
		for _, dep := range pkg.Deps {
			if p.Name == dep.Name {
				return PackageHasDependentsError{p.Name}
			}
		}
	}
	// we're free to remove it at this point
	delete(pi.dict, p.Name)

	return nil
}

// Query searches index for a package.
func (pi *PackageIndexer) Query(name string) (*Package, error) {
	if p, ok := pi.dict[name]; ok {
		return p, nil
	}
	return nil, PackageNotFoundError{name}
}
