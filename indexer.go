package packman

import "errors"

// Indexer is the interface our actual implementations must satisfy to be
// considered package indexers.
type Indexer interface {
	Index(p Package) error
	Query(name string) (Package, error)
}

// Package represents a package to be indexed, removed and querried.
type Package struct {
	Name string
	Deps []Package
}

// PackageIndexer implements the Indexer capabilities of adding, removing and
// querying for a package.
type PackageIndexer struct {
	Dict map[string]*Package
}

// NewIndexer returns an implementation of our indexer.
func NewIndexer(initialDict map[string]*Package) PackageIndexer {
	pi := PackageIndexer{Dict: initialDict}
	if pi.Dict == nil {
		pi.Dict = make(map[string]*Package, 0)
	}
	return pi
}

// Index indexes a package. If the package is already indexed, nothing is changed.
func (pi *PackageIndexer) Index(p *Package) error {
	if _, ok := pi.Dict[p.Name]; !ok {
		pi.Dict[p.Name] = p
	}
	return nil
}

// Remove removes a package from the index.
func (pi *PackageIndexer) Remove(name string) error {
	delete(pi.Dict, name)
	return nil
}

// Query searches index for a package.
func (pi *PackageIndexer) Query(name string) (*Package, error) {
	if p, ok := pi.Dict[name]; ok {
		return p, nil
	}
	return nil, errors.New("Package not found")
}
