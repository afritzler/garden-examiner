package garden

import (
	"fmt"
	"os"

	"github.com/afritzler/garden-examiner/cmd/gex/output"

	"github.com/afritzler/garden-examiner/pkg/data"

	"github.com/afritzler/garden-examiner/cmd/gex/cmdline"
	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/mandelsoft/cmdint/pkg/cmdint"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	filters.AddOptions(cmdline.AddAsVerb(GetCmdTab(), "register", register).
		CmdDescription("run register for garden cluster").
		CmdArgDescription("[check]")).
		FlagOption("check").Short('c').Description("check registration").
		ArgOption("email").Short('e').Description("email address")
}

func register(opts *cmdint.Options) error {
	ctx := context.Get(opts)
	check := opts.IsFlag("check")
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
		return cmdline.ExecuteOutput(opts, NewRegisterOutput(check, githubURL, email), TypeHandler)
	} else {
		return cmdline.ExecuteOutput(opts, NewRegisterOutput(check, githubURL, *email), TypeHandler)
	}
}

func getEmail(githubURL string) string {
	if githubURL == "" {
		return "null"
	}
	res := util.ExecCmdReturnOutput("bash", "-c", "curl -ks "+githubURL+"/api/v3/users/"+os.Getenv("USER")+" | jq -r .email")
	fmt.Printf("used github email: %s\n", res)
	return res
}

func register_garden(check bool, githubURL string, g gube.Garden, email string) (string, error) {
	kubeset, err := g.GetClientset()
	if err != nil {
		return "", fmt.Errorf("failed to get garden clientset: %s", err)
	}
	clusterRoleBinding, err := kubeset.RbacV1().ClusterRoleBindings().Get("garden-administrators", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	for _, subject := range clusterRoleBinding.Subjects {
		if subject.Kind == "User" && subject.Name == email {
			return fmt.Sprintf("user '%s' already registered", email), nil
		}
	}
	if check {
		return fmt.Sprintf("user '%s' not registed", email), nil
	}
	clusterRoleBinding.Subjects = append(clusterRoleBinding.Subjects, rbacv1.Subject{Kind: "User", Name: email})
	_, err = kubeset.RbacV1().ClusterRoleBindings().Update(clusterRoleBinding)
	return fmt.Sprintf("user '%s' registed", email), err
}

func createRegisterMapper(check bool, githubURL, email string) data.MappingFunction {
	return func(e interface{}) interface{} {
		cfg := e.(gube.GardenConfig)
		g, err := cfg.GetGarden()
		if err != nil {
			return err
		}
		s, err := register_garden(check, githubURL, g, email)
		if err != nil {
			return fmt.Errorf("%s: %s", cfg.GetName(), err)
		}
		return fmt.Sprintf("%s: %s", cfg.GetName(), s)
	}
}

func NewRegisterOutput(check bool, githubURL, email string) output.Output {
	return output.NewStringOutput(createRegisterMapper(check, githubURL, email), "")
}
