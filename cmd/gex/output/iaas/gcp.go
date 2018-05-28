package iaas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	. "github.com/afritzler/garden-examiner/cmd/gex/cleanup"
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
	data := map[string]interface{}{}
	tmpAccount := util.ExecCmdReturnOutput("bash", "-c", "gcloud config list account --format json")
	dec := json.NewDecoder(strings.NewReader(tmpAccount))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	tmpAccount, err := jq.String("core", "account")
	if err != nil {
		return fmt.Errorf("cannot list gcloud accounts: %s", err)
	}
	dec = json.NewDecoder(strings.NewReader(string(serviceaccount)))
	dec.Decode(&data)
	jq = jsonq.NewQuery(data)
	account, err := jq.String("client_email")
	if err != nil {
		return fmt.Errorf("cannot list gcloud client emails: %s", err)
	}
	project, err := jq.String("project_id")
	if err != nil {
		return fmt.Errorf("cannot find project id in account list: %s", err)
	}

	tmpfile, err := ioutil.TempFile("/tmp", "serviceaccount")
	if err != nil {
		log.Fatal(err)
	}
	defer Cleanup(func() { os.Remove(tmpfile.Name()) })()

	if _, err := tmpfile.Write(serviceaccount); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	defer Cleanup(func() {
		util.ExecCmd("gcloud config set account " + tmpAccount)
	})()

	err = util.ExecCmd("gcloud auth activate-service-account --key-file=" + tmpfile.Name())
	if err != nil {
		return fmt.Errorf("cannot activate service account: %s", err)
	}
	err = util.ExecCmd("gcloud " + strings.Join(args, " ") + " " + "--account=" + account + " --project=" + project)
	if err != nil {
		return fmt.Errorf("cannot execute 'gcloud': %s", err)
	}
	return nil
}
