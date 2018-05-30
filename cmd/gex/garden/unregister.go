package garden

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "unregister", unregister).Raw().
		CmdDescription("run unregister for garden cluster").
		CmdArgDescription("unregister [email]").
		ArgOption("garden"))
}

func unregister(opts *cmdint.Options) error {
	ctx := context.Get(opts)
	githubURL := ""
	if ctx.GardenSetConfig != nil {
		githubURL = ctx.GardenSetConfig.GetGithubURL()
	}
	switch len(opts.Arguments) {
	case 0:
		return unregister_garden(githubURL, ctx.Garden, "")
	case 1:
		return unregister_garden(githubURL, ctx.Garden, opts.Arguments[0])
	default:
		return fmt.Errorf("One optional email argument required")
	}
}

func unregister_garden(githubURL string, g gube.Garden, email string) error {
	if email == "" {
		email = getEmail(githubURL)
		if email == "null" {
			return fmt.Errorf("Could not read github email address")
		}
	}

	kubeset, err := g.GetClientset()
	if err != nil {
		return fmt.Errorf("failed to get garden clientset: %s", err)
	}
	clusterRoleBinding, err := kubeset.RbacV1().ClusterRoleBindings().Get("garden-administrators", metav1.GetOptions{})
	if err != nil {
		return err
	}
	for k, subject := range clusterRoleBinding.Subjects {
		if subject.Kind == "User" && subject.Name == email {
			clusterRoleBinding.Subjects = append(clusterRoleBinding.Subjects[:k], clusterRoleBinding.Subjects[k+1:]...)
			_, err := kubeset.RbacV1().ClusterRoleBindings().Update(clusterRoleBinding)
			return err
		}
	}
	return nil
}
