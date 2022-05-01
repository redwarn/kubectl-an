package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kubectl-an",
	Short:   "kubectl-an is a k8s resource annotations display",
	Long:    ``,
	Version: "0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	KubernetesConfigFlags = genericclioptions.NewConfigFlags(true)
	annotationCmd.Flags().BoolP("deployments", "d", false, "show deployments annotations")
	annotationCmd.Flags().BoolP("daemonSets", "e", false, "show daemonSets annotations")
	annotationCmd.Flags().BoolP("statefulSets", "f", false, "show statefulSets annotations")
	annotationCmd.Flags().BoolP("service", "v", false, "show service annotations")
	annotationCmd.Flags().BoolP("ingress", "i", false, "show ingress annotations")
	annotationCmd.Flags().BoolP("pods", "p", false, "show pod annotations")
	KubernetesConfigFlags.AddFlags(rootCmd.PersistentFlags())

}
