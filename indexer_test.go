package packman_test

import (
	"testing"

	"github.com/jboursiquot/packman"
)

func TestIndexerIndexesPackageWithoutDeps(t *testing.T) {
	idxr := packman.NewIndexer(nil)
	p := packman.Package{Name: "somepackage"}
	err := idxr.Index(&p)
	if err != nil {
		t.Errorf("Expected successful indexing of package: %v", err)
	}
}

func TestIndexerIndexesPackageWithKnownDeps(t *testing.T) {
	d1 := packman.Package{Name: "Dep1"}
	d2 := packman.Package{Name: "Dep2"}
	initialDict := make(map[string]*packman.Package, 0)
	initialDict[d1.Name] = &d1
	initialDict[d2.Name] = &d2
	idxr := packman.NewIndexer(initialDict)
	p := packman.Package{Name: "somepackage", Deps: []*packman.Package{&d1, &d2}}
	err := idxr.Index(&p)
	if err != nil {
		t.Error(err)
	}
}

func TestIndexerFailsToIndexPackageWithUnknownDeps(t *testing.T) {
	idxr := packman.NewIndexer(nil)
	d := packman.Package{Name: "dep"}
	p := packman.Package{Name: "somepackage", Deps: []*packman.Package{&d}}
	err := idxr.Index(&p)
	if _, ok := err.(packman.UnknownDependentError); !ok {
		t.Error("Expected packman.UnknownDependentError")
	}
}

func TestIndexerRemovesPackageWithoutDeps(t *testing.T) {
	p := packman.Package{Name: "somepackage"}
	initialDict := make(map[string]*packman.Package, 0)
	initialDict[p.Name] = &p
	idxr := packman.NewIndexer(initialDict)

	err := idxr.Remove(&p)
	if err != nil {
		t.Errorf("Expected successful removal of package: %v", err)
	}

	qp, err := idxr.Query(p.Name)
	if qp != nil {
		t.Errorf("Expected querried package to not be present: %#v", err)
	}
	if _, ok := err.(packman.PackageNotFoundError); !ok {
		t.Error("Expected packman.PackageNotFoundError")
	}
}

func TestIndexerFailsToRemovePackageWithDeps(t *testing.T) {
	d1 := packman.Package{Name: "Dep1"}
	d2 := packman.Package{Name: "Dep2"}
	p := packman.Package{Name: "somepackage", Deps: []*packman.Package{&d1, &d2}}
	initialDict := make(map[string]*packman.Package, 0)
	initialDict[d1.Name] = &d1
	initialDict[d1.Name] = &d1
	initialDict[d2.Name] = &p
	idxr := packman.NewIndexer(initialDict)

	err := idxr.Remove(&d1)
	if _, ok := err.(packman.PackageHasDependentsError); !ok {
		t.Error("Expected packman.PackageHasDependentsError")
	}
}

func TestIndexerFindsIndexedPackage(t *testing.T) {
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

func TestIndexerFailsToFindUnindexedPackage(t *testing.T) {
	idxr := packman.NewIndexer(nil)
	_, err := idxr.Query("somepackage")
	if _, ok := err.(packman.PackageNotFoundError); !ok {
		t.Error("Expected packman.PackageNotFoundError")
	}
}
