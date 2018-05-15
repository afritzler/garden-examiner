package cmd_select

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	elem "github.com/afritzler/garden-examiner/cmd/gex/seed"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	GetCmdTab().SimpleCommand("seed", seed).CmdDescription("select shoot cluster").CmdArgDescription("[<project>/]<shootname>").
		ArgOption(constants.O_SEL_PROJECT).
		ArgOption(constants.O_SEL_SEED)
}

func seed(opts *cmdint.Options) error {
	h, err := NewSeedHandler(opts)
	if err != nil {
		return err
	}
	return util.Doit(opts, h)
}

func NewSeedHandler(opts *cmdint.Options) (util.Handler, error) {
	return &elem.GetHandler{elem.NewHandler(NewSelectOutput())}, nil
}
