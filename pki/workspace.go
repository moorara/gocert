package pki

import (
	"github.com/moorara/go-box/util"
)

// NewWorkspace initialize a new workspace
func NewWorkspace(state *State, spec *Spec) error {
	// Make sub-directories
	_, err := util.MkDirs("", DirRoot, DirInterm, DirServer, DirClient, DirCSR)
	if err != nil {
		return err
	}

	// Write state file
	err = SaveState(state, FileState)
	if err != nil {
		return err
	}

	// Write spec file
	err = SaveSpec(spec, FileSpec)
	if err != nil {
		return err
	}

	return nil
}
