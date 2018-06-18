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
	"github.com/mandelsoft/filepath/pkg/filepath"
)

func init() {
	RegisterIaasHandler(&gcp{}, "gcp")
}

type gcp struct {
}

func (this *gcp) Execute(shoot gube.Shoot, config map[string]string, args ...string) error {
	data := map[string]interface{}{}
	tmpAccount := util.ExecCmdReturnOutput("bash", "-c", "gcloud config list account --format json")
	dec := json.NewDecoder(strings.NewReader(tmpAccount))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)
	tmpAccount, err := jq.String("core", "account")
	if err != nil {
		return fmt.Errorf("cannot list gcloud accounts: %s", err)
	}

	serviceaccount := []byte(config["serviceaccount.json"])
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
		return fmt.Errorf("cannot get temporary key file name: %s", err)
	}
	defer Cleanup(func() { os.Remove(tmpfile.Name()) })()

	if _, err := tmpfile.Write(serviceaccount); err != nil {
		return fmt.Errorf("cannot write temporary key file '%s' for key file: %s", tmpfile.Name, err)
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

func (this *gcp) Export(shoot gube.Shoot, config map[string]string, cachedir string) error {
	serviceaccount := []byte(config["serviceaccount.json"])
	err := os.MkdirAll(cachedir, 0700)
	if err != nil {
		return fmt.Errorf("cannot create cache dir '%s' for key file: %s", cachedir, err)
	}
	keyfile := filepath.Join(cachedir, "gcp.serviceaccount")

	file, err := os.OpenFile(keyfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return fmt.Errorf("cannot create key file '%s' for key file: %s", keyfile, err)
	}
	if _, err := file.Write(serviceaccount); err != nil {
		return fmt.Errorf("cannot write key file '%s' for key file: %s", keyfile, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("cannot close key file '%s' for key file: %s", keyfile, err)
	}

	fmt.Printf("activating gcloud service account for %s\n", shoot.GetName())
	err = util.ExecCmd("gcloud auth activate-service-account --key-file=" + keyfile)
	if err != nil {
		return fmt.Errorf("cannot activate service account: %s", err)
	}
	return nil
}

func (this *gcp) Describe(shoot gube.Shoot) error {
	info, err := shoot.GetIaaSInfo()
	if err == nil {
		iaas := info.(*gube.GCPInfo)
		attrs := util.NewAttributeSet()
		fmt.Printf("GCP Information:\n")
		attrs.Attribute("Region", iaas.GetRegion())
		attrs.Attribute("VPC Name", iaas.GetVpcName())
		attrs.Attribute("Service Accout EMail", iaas.GetServiceAccountEMail())
		attrs.PrintAttributes()
	}
	return nil
}
