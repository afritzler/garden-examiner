package output

import (
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mandelsoft/cmdint/pkg/cmdint"

	"github.com/afritzler/garden-examiner/cmd/gex/context"
	"github.com/afritzler/garden-examiner/cmd/gex/util"
	"github.com/afritzler/garden-examiner/pkg"
)

type ShellOutput struct {
	*KubectlOutput
}

var _ Output = &ShellOutput{}

func NewShellOutput(node *string, pod *string, mapper ElementMapper) Output {
	args := []string{util.StringValue(node), util.StringValue(pod)}
	return &ShellOutput{NewKubectlOutput(args, mapper)}
}

func (this *ShellOutput) Out(ctx *context.Context) error {
	cluster := this.Elem.(gube.Cluster)
	nodes, err := cluster.GetNodes()
	if err != nil {
		return err
	}

	hostnames := map[string]string{}
	nodenames := []string{}

	podnodes := map[string]string{}
	podnames := []string{}

	name := ""
	lookupNode := ""
	lookupPod := ""

	if len(this.GetArgs()) > 0 {
		lookupNode = this.GetArgs()[0]
	}

	for _, n := range nodes {
		host := get_label(n.GetObjectMeta(), "kubernetes.io/hostname")
		nodenames = append(nodenames, n.GetName())
		if len(lookupNode) >= 5 && strings.HasSuffix(n.GetName(), lookupNode) {
			name = n.GetName()
		}
		hostnames[n.GetName()] = host
	}

	if len(this.GetArgs()) > 1 {
		lookupPod = this.GetArgs()[1]
		if lookupPod != "" {
			var pods map[string]corev1.Pod
			if i := strings.Index(lookupPod, "/"); i > 0 {
				ns := lookupPod[0:i]
				lookupPod = lookupPod[i+1:]
				pods, err = cluster.GetPods(ns)
			} else {
				pods, err = cluster.GetPods("")
			}
			if err != nil {
				return err
			}
			for _, n := range pods {
				podnames = append(podnames, n.GetName())
				podnodes[n.GetName()] = n.Spec.NodeName
			}
		}
	}

	if name == "" && len(podnames)+len(nodenames) > 0 {
		pname, c := cmdint.SelectBest(lookupPod, podnames...)
		name, _ = podnodes[pname]
		nname, nc := cmdint.SelectBest(lookupNode, nodenames...)
		if nc < c {
			name = nname
		}
	}
	if name == "" {
		fmt.Printf("select one of:\n")
		for _, n := range nodenames {
			fmt.Printf("- %s (%s)\n", n, hostnames[n])
		}
		if lookupNode == "" {
			return nil
		}
		return fmt.Errorf("node '%s' not found", lookupNode)
	}
	client, err := cluster.GetClientset()
	if err != nil {
		return fmt.Errorf("cannot access cluster: %s", err)
	}
	fmt.Printf("running shell on node '%s'\n", name)
	_, err = client.CoreV1().Pods("default").Get("rootpod", metav1.GetOptions{})
	if err == nil {
		this.Kubectl(nil, "delete", "pod", "rootpod")
	}
	manifest := strings.Replace(shell_manifest, "HOSTNAME", hostnames[name], -1)
	this.Kubectl([]byte(manifest), "apply", "-f", "-")

	for true {
		pod, err := client.CoreV1().Pods("default").Get("rootpod", metav1.GetOptions{})
		if err != nil {
			fmt.Printf("pod not found: %s\n", err)
		} else {
			ip := pod.Status.HostIP
			if ip != "" && pod.Status.Phase == "Running" {
				fmt.Printf("host ip found: %s\n", ip)
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
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
    securityContext:
      privileged: true
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
