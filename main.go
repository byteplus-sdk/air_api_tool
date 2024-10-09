package main

import (
	"github.com/3rd_rec/air_api_tool/checker"
	_ "github.com/3rd_rec/air_api_tool/checker/fieldchecker"
	"github.com/3rd_rec/air_api_tool/schema"
	"github.com/3rd_rec/air_api_tool/validate"
	"github.com/spf13/cobra"
)

type Command interface {
	GetCMD() *cobra.Command
}

var (
	rootCMD = &cobra.Command{
		Use:   "api_tool",
		Short: "API Tool",
	}
	supportedModules = []Command{
		&validate.Command{},
		&schema.ManagementCommand{},
		&checker.ManagementCommand{},
	}
)

func main() {
	for _, m := range supportedModules {
		rootCMD.AddCommand(m.GetCMD())
	}
	_ = rootCMD.Execute()
}
