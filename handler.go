package packman

import (
	"log"
	"regexp"
	"strings"
)

const (
	// INDEX is the command to track a new package
	INDEX = "INDEX"

	// REMOVE is the command to remove an existing package
	REMOVE = "REMOVE"

	// QUERY is the command to lookup a package
	QUERY = "QUERY"
)

var msgRegex = regexp.MustCompile(`^(INDEX|REMOVE|QUERY){1}\|(.+){1}\|(.*)`)

// Command pulls together all the component parts of a message to be processed.
type Command struct {
	Verb    string
	Package Package
}

// CommandFromMessage decomposes an incoming message string into a Command if
// all the parts of a expected `<command>|<package>|<dependencies>\n` format
// are present. Returns an InvalidMessageError otherwise.
func CommandFromMessage(message string) (*Command, error) {
	matches := msgRegex.FindAllStringSubmatch(message, -1)
	if matches == nil || len(matches[0]) != 4 {
		return nil, InvalidMessageError{message}
	}

	verb, pkg, deps := matches[0][1], matches[0][2], strings.Split(matches[0][3], ",")
	if len(pkg) == 0 {
		return nil, InvalidMessageError{message}
	}

	p := Package{
		Name: pkg,
	}

	for _, name := range deps {
		if len(name) == 0 || name == p.Name {
			continue
		}
		p.Deps = append(p.Deps, &Package{Name: name})
	}

	c := Command{
		Verb:    verb,
		Package: p,
	}
	return &c, nil
}

// ProcessCommand processes a command obtained from a client message.
// The results will vary depending on the command that was issued but
// generally, only a QUERY command will return anything.
func ProcessCommand(cmd *Command, idxr *PackageIndexer) (interface{}, error) {
	// Ensure that processing of this command is synchronized.
	idxr.Lock()
	defer idxr.Unlock()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered | err=%v, cmd=%v", r, cmd)
		}
	}()

	switch cmd.Verb {
	case INDEX:
		return nil, idxr.Index(&cmd.Package)
	case REMOVE:
		return nil, idxr.Remove(&cmd.Package)
	case QUERY:
		return idxr.Query(cmd.Package.Name)
	}
	return nil, nil
}
