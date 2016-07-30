package packman_test

import (
	"reflect"
	"testing"

	"github.com/jboursiquot/packman"
)

var malformedMessageTests = []struct {
	msg      string
	expected error
}{
	{"", packman.InvalidMessageError{Message: ""}},
	{"||", packman.InvalidMessageError{Message: "||"}},
	{`badcommand|pkg|`, packman.InvalidMessageError{Message: `badcommand|pkg|`}},
	{"index|pkg", packman.InvalidMessageError{Message: "index|pkg"}},
	{"index||", packman.InvalidMessageError{Message: "index||"}},
	{`index||\t`, packman.InvalidMessageError{Message: `index||\t`}},
	{`index|pkg`, packman.InvalidMessageError{Message: `index|pkg`}},
	{`index\n|pkg`, packman.InvalidMessageError{Message: `index\n|pkg`}},
	{`INDEX|pkg`, packman.InvalidMessageError{Message: `INDEX|pkg`}},
	{`BLINDEX|pkg|dep`, packman.InvalidMessageError{Message: `BLINDEX|pkg|dep`}},
	{`QUERY||`, packman.InvalidMessageError{Message: `QUERY||`}},
	{`QUER|a|b`, packman.InvalidMessageError{Message: `QUER|a|b`}},
}

func TestHandlerHandlesInvalidMessages(t *testing.T) {
	for _, tt := range malformedMessageTests {
		_, err := packman.CommandFromMessage(tt.msg)
		if !reflect.DeepEqual(tt.expected, err) {
			t.Errorf("Expected %v but got %v", tt.expected, err)
		}
	}
}

var validMessageTests = []struct {
	msg      string
	expected *packman.Command
}{
	{
		`INDEX|cloog|gmp,isl,pkg-config`,
		&packman.Command{
			Verb: packman.INDEX,
			Package: packman.Package{
				Name: "cloog",
				Deps: []*packman.Package{
					&packman.Package{Name: "gmp"},
					&packman.Package{Name: "isl"},
					&packman.Package{Name: "pkg-config"},
				},
			},
		},
	},
	{
		`INDEX|ceylon|`,
		&packman.Command{
			Verb: packman.INDEX,
			Package: packman.Package{
				Name: "ceylon",
			},
		},
	},
	{
		`REMOVE|cloog3|`,
		&packman.Command{
			Verb: packman.REMOVE,
			Package: packman.Package{
				Name: "cloog3",
			},
		},
	},
	{
		`QUERY|cloog|`,
		&packman.Command{
			Verb: packman.QUERY,
			Package: packman.Package{
				Name: "cloog",
			},
		},
	},
}

func TestHandlerTranslatesMessageToCommand(t *testing.T) {
	for _, vmt := range validMessageTests {
		cmd, _ := packman.CommandFromMessage(vmt.msg)
		if !reflect.DeepEqual(vmt.expected, cmd) {
			t.Errorf("Expected command to be like %#v", vmt.expected)
		}
	}
}

var validCommandTests = []struct {
	cmd                *packman.Command
	initialIndexerDict map[string]*packman.Package
	expected           interface{}
}{
	{
		&packman.Command{
			Verb: packman.INDEX,
			Package: packman.Package{
				Name: "cloog",
				Deps: []*packman.Package{
					&packman.Package{Name: "gmp"},
					&packman.Package{Name: "isl"},
					&packman.Package{Name: "pkg-config"},
				},
			},
		},
		map[string]*packman.Package{
			"gmp":        &packman.Package{Name: "gmp"},
			"isl":        &packman.Package{Name: "isl"},
			"pkg-config": &packman.Package{Name: "pkg-config"},
		},
		nil,
	},
	{
		&packman.Command{
			Verb: packman.INDEX,
			Package: packman.Package{
				Name: "ceylon",
			},
		},
		nil,
		nil,
	},
	{
		&packman.Command{
			Verb: packman.REMOVE,
			Package: packman.Package{
				Name: "cloog",
			},
		},
		nil,
		nil,
	},
	{
		&packman.Command{
			Verb: packman.QUERY,
			Package: packman.Package{
				Name: "cloog",
			},
		},
		map[string]*packman.Package{
			"cloog": &packman.Package{Name: "cloog"},
		},
		nil,
	},
}

func TestHandlerProcessesCommandAccurately(t *testing.T) {
	for _, tt := range validCommandTests {
		idxr := packman.NewIndexer(tt.initialIndexerDict)
		_, err := packman.ProcessCommand(tt.cmd, &idxr)
		if err != nil {
			t.Errorf("Expected %v, got %v", tt.expected, err)
		}
	}
}
