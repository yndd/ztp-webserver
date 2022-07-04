/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
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
	Short: "Execute the ZTP DHCP-Server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ws := webserver.GetWebserverOperations()
		ws.SetKubeConfig(kubeconfig)
		ws.Run(port, storageFolder)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// TODO: allow for viper environment variables here
	// viper.SetEnvPrefix("NDD")
	// viper.BindEnv("dhcpv4_port")

	// viper.BindPFlag("dhcpv4_port", runCmd.Flags().Lookup("dhcpv4-port"))

	runCmd.Flags().IntVar(&port, "port", 8000, "The port to bind the webserver to")
	//runCmd.Flags().IntVar(&dhcpv6_port, "dhcpv6-port", 567, "The port to bind the dhcpv6 server to.")
	runCmd.Flags().StringVar(&storageFolder, "storagefolder", "/webserver", "Folder that contains content for the webserver to deliver")
	runCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Pointer to the kubeconfig file")
}
