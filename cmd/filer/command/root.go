package command

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "filer service",
	Short: "A CLI for filer Service",
	Long:  `Filer CLI is a tool to handle upload and download files.`,
}

func init() {
	RootCmd.AddCommand(serveFilerCmd, serveUploadAppCmd)
}
