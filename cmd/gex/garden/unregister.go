package garden

import (
	"fmt"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/output"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/afritzler/garden-examiner/pkg/data"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "unregister", unregister).
		CmdDescription("run unregister for garden cluster").
		ArgOption("email").Short('e').Description("email address"))
}

func unregister(opts *cmdint.Options) error {
	ctx := context.Get(opts)
	email := opts.GetOptionValue("email")
	githubURL := ""
	if ctx.GardenSetConfig != nil {
		githubURL = ctx.GardenSetConfig.GetGithubURL()
	}
	if email == nil {
		if githubURL == "" {
			return fmt.Errorf("No email specified and no github url configured in gex config")
		}
		email := getEmail(githubURL)
		if email == "null" {
			return fmt.Errorf("Could not read github email address")
		}
		return cmdline.ExecuteOutput(opts, NewUnregisterOutput(githubURL, email), TypeHandler)
	} else {
		return cmdline.ExecuteOutput(opts, NewUnregisterOutput(githubURL, *email), TypeHandler)
	}
}

func unregister_garden(githubURL string, g gube.Garden, email string) (string, error) {
	kubeset, err := g.GetClientset()
	if err != nil {
		return "", fmt.Errorf("failed to get garden clientset: %s", err)
	}
	clusterRoleBinding, err := kubeset.RbacV1().ClusterRoleBindings().Get("garden-administrators", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	for k, subject := range clusterRoleBinding.Subjects {
		if subject.Kind == "User" && subject.Name == email {
			clusterRoleBinding.Subjects = append(clusterRoleBinding.Subjects[:k], clusterRoleBinding.Subjects[k+1:]...)
			_, err := kubeset.RbacV1().ClusterRoleBindings().Update(clusterRoleBinding)
			return fmt.Sprintf("user '%s' unregisted", email), err
		}
	}
	return fmt.Sprintf("user '%s' already unregisted", email), nil
}

func createUnregisterMapper(githubURL, email string) data.MappingFunction {
	return func(e interface{}) interface{} {
		cfg := e.(gube.GardenConfig)
		g, err := cfg.GetGarden()
		if err != nil {
			return err
		}
		s, err := unregister_garden(githubURL, g, email)
		if err != nil {
			return fmt.Errorf("%s: %s", cfg.GetName(), err)
		}
		return fmt.Sprintf("%s: %s", cfg.GetName(), s)
	}
}

func NewUnregisterOutput(githubURL, email string) output.Output {
	return output.NewStringOutput(createUnregisterMapper(githubURL, email), "")
}
