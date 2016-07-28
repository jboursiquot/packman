package packman_test

import (
	"testing"

	"github.com/jboursiquot/packman"
)

func TestIndexerIndexWithoutDeps(t *testing.T) {
	idxr := packman.NewIndexer(nil)
	p := packman.Package{Name: "somepackage"}
	err := idxr.Index(&p)
	if err != nil {
		t.Errorf("Expected successful indexing of package: %v", err)
	}
}

func TestIndexerRemove(t *testing.T) {
	p := packman.Package{Name: "somepackage"}
	initialDict := make(map[string]*packman.Package, 1)
	initialDict[p.Name] = &p
	idxr := packman.NewIndexer(initialDict)

	err := idxr.Remove(p.Name)
	if err != nil {
		t.Errorf("Expected successful removal of package: %v", err)
	}

	qp, err := idxr.Query(p.Name)
	if qp != nil {
		t.Errorf("Expected querried package to not be present: %#v", err)
	}
	if err != nil && err.Error() != "Package not found" {
		t.Error(err)
	}
}

func TestIndexerQuery(t *testing.T) {
	p := packman.Package{Name: "somepackage"}
	initialDict := make(map[string]*packman.Package, 1)
	initialDict[p.Name] = &p
	idxr := packman.NewIndexer(initialDict)

	qp, err := idxr.Query(p.Name)
	if err != nil {
		t.Errorf("Expected successful querying for package: %v", err)
	}
	if p.Name != qp.Name {
		t.Errorf("Expected querried package to be the same as stored package: %#v != %#v", p, qp)
	}
}
