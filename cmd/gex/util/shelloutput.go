package util

import (
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/pkg"
)

type ShellOutput struct {
	*KubectlOutput
}

var _ Output = &ShellOutput{}

func NewShellOutput(node *string, mapper ElementMapper) Output {
	args := []string{}
	if node != nil && *node != "" {
		args = []string{*node}
	}
	return &ShellOutput{NewKubectlOutput(args, mapper)}
}

func (this *ShellOutput) Out(ctx *context.Context) error {
	cluster := this.Elem.(gube.Cluster)
	nodes, err := cluster.GetNodes()
	if err != nil {
		return err
	}
	hostnames := map[string]string{}
	names := []string{}
	name := ""
	lookup := ""
	if len(this.GetArgs()) > 0 {
		lookup = this.GetArgs()[0]
	}
	for _, n := range nodes {
		host := get_label(n.GetObjectMeta(), "kubernetes.io/hostname")
		names = append(names, n.GetName())
		if len(lookup) >= 5 && strings.HasSuffix(n.GetName(), lookup) {
			name = n.GetName()
		}
		hostnames[n.GetName()] = host
	}
	if name == "" && lookup != "" {
		name, _ = cmdint.SelectBest(lookup)
	}
	if name == "" {
		fmt.Printf("select one of:\n")
		for _, n := range names {
			fmt.Printf("- %s (%s)\n", n, hostnames[n])
		}
		if lookup == "" {
			return nil
		}
		return fmt.Errorf("node '%s' not found", lookup)
	}
	fmt.Printf("running shell on node '%s'\n", name)
	manifest := strings.Replace(shell_manifest, "HOSTNAME", hostnames[name], -1)
	this.Kubectl([]byte(manifest), "apply", "-f", "-")
	time.Sleep(5 * time.Second)
	this.Kubectl(nil, "exec", "-it", "rootpod", "--", "/bin/sh", "-c", "chroot /hostroot")
	this.Kubectl(nil, "delete", "pod", "rootpod")
	return nil
}

func get_label(n metav1.Object, label string) string {
	return n.GetLabels()[label]
}

var shell_manifest = `
apiVersion: v1
kind: Pod
metadata:
  name: rootpod
spec:
  containers:
  - image: busybox
    name: root-container
    command:
    - sleep 
    - "10000000"
    volumeMounts:
    - mountPath: /hostroot
      name: root-volume
  hostNetwork: true
  hostPID: true
  nodeSelector:
    kubernetes.io/hostname: "HOSTNAME"
  tolerations:
  - key: node-role.kubernetes.io/master
    operator: Exists
    effect: NoSchedule
  volumes:
  - name: root-volume
    hostPath:
      path: /
`
