/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yndd/ztp-dhcp/pkg/backend/k8s"
	_ "github.com/yndd/ztp-webserver/pkg/devices/all"
	"github.com/yndd/ztp-webserver/pkg/webserver"
)

var (
	port          int
	storageFolder string
	kubeconfig    string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the ZTP Web-Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ws := webserver.GetWebserverOperations()
		// # to utilize K8s apiserver use
		// backend := k8s.NewZtpK8sBackend(kubeconfig)
		// # or, to use static backend, use
		// backend := static.NewZtpStaticBackend()
		backend := k8s.NewZtpK8sBackend(kubeconfig)
		ws.SetBackend(backend)
		// execute
		ws.Run(port, storageFolder)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&port, "port", 8000, "The port to bind the webserver to")
	runCmd.Flags().StringVar(&storageFolder, "storagefolder", "/webserver", "Folder that contains content for the webserver to deliver")
	runCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Pointer to the kubeconfig file")
}
