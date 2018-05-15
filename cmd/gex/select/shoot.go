package cmd_select

import (
	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/const"
	elem "github.com/afritzler/garden-examiner/cmd/gex/shoot"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
)

func init() {
	GetCmdTab().SimpleCommand("shoot", shoot).CmdDescription("select shoot cluster").CmdArgDescription("[<project>/]<shootname>").
		ArgOption(constants.O_SEL_PROJECT).
		ArgOption(constants.O_SEL_SEED)
}

func shoot(opts *cmdint.Options) error {
	h, err := NewShootHandler(opts)
	if err != nil {
		return err
	}
	return util.Doit(opts, h)
}

func NewShootHandler(opts *cmdint.Options) (util.Handler, error) {
	return &elem.GetHandler{elem.NewHandler(NewSelectOutput())}, nil
}
