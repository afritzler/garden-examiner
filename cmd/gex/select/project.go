package cmd_select

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	elem "github.com/afritzler/garden-examiner/cmd/gex/project"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	GetCmdTab().SimpleCommand("project", project).CmdDescription("select project").CmdArgDescription("<project>")
}

func project(opts *cmdint.Options) error {
	h, err := NewProjectHandler(opts)
	if err != nil {
		return err
	}
	return util.Doit(opts, h)
}

func NewProjectHandler(opts *cmdint.Options) (util.Handler, error) {
	return &elem.GetHandler{elem.NewHandler(NewSelectOutput())}, nil
}
