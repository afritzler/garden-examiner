package iaas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
	"github.com/jmoiron/jsonq"
)

func init() {
	RegisterIaasHandler(&gcp{}, "gcp")
}

type gcp struct {
}

func (this *gcp) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	serviceaccount := []byte(config["serviceaccount.json"])
	tmpfile, err := ioutil.TempFile("", "serviceaccount")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(serviceaccount); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	data := map[string]interface{}{}
	tmpAccount := util.ExecCmdReturnOutput("bash", "-c", "gcloud config list account --format json")
	dec := json.NewDecoder(strings.NewReader(tmpAccount))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	tmpAccount, err = jq.String("core", "account")
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	dec = json.NewDecoder(strings.NewReader(string(serviceaccount)))
	dec.Decode(&data)
	jq = jsonq.NewQuery(data)
	account, err := jq.String("client_email")
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	project, err := jq.String("project_id")
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	err = util.ExecCmd("gcloud auth activate-service-account --key-file=" + tmpfile.Name())
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	err = util.ExecCmd(strings.Join(args, " ") + " " + "--account=" + account + " --project=" + project)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
	err = util.ExecCmd("gcloud config set account " + tmpAccount)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	return nil
}
