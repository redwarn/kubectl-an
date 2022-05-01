package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	kv1 "k8s.io/api/apps/v1"
	cv1 "k8s.io/api/core/v1"
	nv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"kubectl-an/pkg/kube"
	"kubectl-an/pkg/table"
	"os"
)

// annotationCmd represents the annotation command
var annotationCmd = &cobra.Command{
	Use:   "annotation",
	Short: "show resource annotations",
	Long:  `show k8s resource annotations`,
	RunE:  annotations,
}

func init() {
	rootCmd.AddCommand(annotationCmd)
}

func getResources(cmd *cobra.Command, clientSet *kubernetes.Clientset, ns, name string) []interface{} {
	var rList []interface{}
	listOptions := v1.ListOptions{}
	if name != "" {
		listOptions.FieldSelector = fields.OneTermEqualSelector("metadata.name", name).String()
	}
	if flag, _ := cmd.Flags().GetBool("deployments"); flag {
		deploymentList, err := clientSet.AppsV1().Deployments(ns).List(context.Background(), listOptions)
		if err != nil {
			fmt.Printf("list deployments error: %s", err.Error())
		}
		rList = append(rList, deploymentList)
	}

	if flag, _ := cmd.Flags().GetBool("daemonSets"); flag {
		daemonSetList, err := clientSet.AppsV1().DaemonSets(ns).List(context.Background(), listOptions)
		if err != nil {
			fmt.Printf("list daemonsets error: %s", err.Error())
		}
		rList = append(rList, daemonSetList)
	}
	if flag, _ := cmd.Flags().GetBool("statefulSets"); flag {
		statefulSetList, err := clientSet.AppsV1().StatefulSets(ns).List(context.Background(), listOptions)
		if err != nil {
			fmt.Printf("list statefulsets error: %s", err.Error())
		}
		rList = append(rList, statefulSetList)
	}
	if flag, _ := cmd.Flags().GetBool("ingress"); flag {
		ingressList, err := clientSet.NetworkingV1().Ingresses(ns).List(context.Background(), listOptions)
		if err != nil {
			fmt.Printf("list ingress error: %s", err.Error())
		}
		rList = append(rList, ingressList)
	}
	if flag, _ := cmd.Flags().GetBool("service"); flag {
		serviceList, err := clientSet.CoreV1().Services(ns).List(context.Background(), listOptions)
		if err != nil {
			fmt.Printf("list service error: %s", err.Error())
		}
		rList = append(rList, serviceList)
	}
	if flag, _ := cmd.Flags().GetBool("pods"); flag {
		podsList, err := clientSet.CoreV1().Pods(ns).List(context.Background(), listOptions)
		if err != nil {
			fmt.Printf("list pod error: %s", err.Error())
		}
		rList = append(rList, podsList)
	}
	return rList
}
func annotations(cmd *cobra.Command, args []string) error {
	clientSet := kube.ClientSet(KubernetesConfigFlags)
	ns, _ := rootCmd.Flags().GetString("namespace")
	var rList []interface{}
	if len(args) == 0 {
		rList = getResources(cmd, clientSet, ns, "")
	}
	if len(args) > 0 {
		for i := range args {
			rs := getResources(cmd, clientSet, ns, args[i])
			rList = append(rList, rs...)
		}
	}
	resources := make([]table.Resource, 0)
	for i := 0; i < len(rList); i++ {
		switch t := rList[i].(type) {
		case *kv1.DeploymentList:
			for k := 0; k < len(t.Items); k++ {
				resource := table.Resource{
					Namespace:   t.Items[k].Namespace,
					Name:        t.Items[k].GetName(),
					Type:        "Deployment",
					Annotations: t.Items[k].Annotations,
				}
				delete(resource.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
				resources = append(resources, resource)
			}
		case *kv1.StatefulSetList:
			for k := 0; k < len(t.Items); k++ {
				resource := table.Resource{
					Namespace:   t.Items[k].Namespace,
					Name:        t.Items[k].GetName(),
					Type:        "StatefulSet",
					Annotations: t.Items[k].Annotations,
				}
				delete(resource.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
				resources = append(resources, resource)
			}
		case *kv1.DaemonSetList:
			for k := 0; k < len(t.Items); k++ {
				resource := table.Resource{
					Namespace:   t.Items[k].Namespace,
					Name:        t.Items[k].GetName(),
					Type:        "DaemonSet",
					Annotations: t.Items[k].Annotations,
				}
				delete(resource.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
				resources = append(resources, resource)
			}
		case *cv1.ServiceList:
			for k := 0; k < len(t.Items); k++ {
				resource := table.Resource{
					Namespace:   t.Items[k].Namespace,
					Name:        t.Items[k].GetName(),
					Type:        "Service",
					Annotations: t.Items[k].Annotations,
				}
				delete(resource.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
				delete(resource.Annotations, "workflows.argoproj.io/description")
				resources = append(resources, resource)
			}
		case *nv1.IngressList:
			for k := 0; k < len(t.Items); k++ {
				resource := table.Resource{
					Namespace:   t.Items[k].Namespace,
					Name:        t.Items[k].GetName(),
					Type:        "Ingress",
					Annotations: t.Items[k].Annotations,
				}
				resources = append(resources, resource)
			}
		case *cv1.PodList:
			for k := 0; k < len(t.Items); k++ {
				resource := table.Resource{
					Namespace:   t.Items[k].Namespace,
					Name:        t.Items[k].GetName(),
					Type:        "Pod",
					Annotations: t.Items[k].Annotations,
				}
				delete(resource.Annotations, "kubernetes.io/limit-ranger")
				resources = append(resources, resource)
			}
		}
	}
	if len(resources) == 0 {
		fmt.Println("no resource")
	}

	t := table.GenTable(resources)
	t.SetOutputMirror(os.Stdout)
	t.Render()
	return nil
}
