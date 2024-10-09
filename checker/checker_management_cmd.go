package checker

import (
	"fmt"
	"github.com/3rd_rec/air_api_tool/helper"

	"github.com/spf13/cobra"
)

type Args struct {
	// Required.
	// Schema's industry.
	// Support [saas_retail、saas_content].
	// Schema default fields are different in different industries.
	Industry string

	// Required.
	// Schema table name.
	// Support [user、product、user_event] for retail industry.
	// Support [user、content、user_event] for content industry.
	Table string

	// Optional.
	// If you have multiple schemas for the same industry and table,
	// you should specify a namespace to differentiate them.
	Namespace string

	// Ignore field checker validation.
	// Examples: []string{"custom_field1", "custom_field2", "custom_field3"}
	IgnoreFields []string

	// Recover ignored fields.
	// Examples: []string{"custom_field1", "custom_field2", "custom_field3"}
	RecoverFields []string

	// Recover default Schema configuration.
	Clear bool

	// Print the current Schema configuration after executing the command.
	Show bool
}

type ManagementCommand struct {
	command *cobra.Command
	args    *Args
}

func (c *ManagementCommand) GetCMD() *cobra.Command {
	c.args = &Args{}
	c.command = &cobra.Command{
		Use:   "checker",
		Short: "Manage the checker of the API. Use `api_tool checker --help` or see https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#Hsuudys5soj8uxxLVFjcAv2Qn9d for more details.",
		Long:  "Manage the checker of the API. Use `api_tool checker --help` or see https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#Hsuudys5soj8uxxLVFjcAv2Qn9d for more details.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := newService(c.args).Run(); err != nil {
				fmt.Printf("[Error] run `api_tool checker` occurred err: %s\n", err)
			}
		},
	}
	c.command.Flags().StringVarP(
		&c.args.Table,
		"table",
		"t",
		"",
		"Required. Supports: user, product, user_event in saas_retail industry. Supports: user, content, user_event in saas_content industry.",
	)
	c.command.Flags().StringVarP(
		&c.args.Industry,
		"industry",
		"i",
		"",
		"Required. Data's industry. Supports: saas_retail, saas_content.",
	)
	c.command.Flags().StringVarP(
		&c.args.Namespace,
		"namespace",
		"n",
		"default",
		"Optional. Specify a namespace to differentiate multiple schemas for the same industry and table.",
	)
	c.command.Flags().VarP(
		helper.NewStringSliceFlag(&c.args.IgnoreFields),
		"ignore-field",
		"g",
		"Ignore preset field checker logic.\nExamples: --ignore-field=user_id --ignore-field=registration_timestamp.",
	)
	c.command.Flags().VarP(
		helper.NewStringSliceFlag(&c.args.RecoverFields),
		"recover-field",
		"r",
		"Recover ignored fields.\nExamples: --recover-field=user_id --recover-field=registration_timestamp.",
	)
	c.command.Flags().BoolVarP(
		&c.args.Show,
		"show",
		"s",
		false,
		"Print the current checker configuration.",
	)
	c.command.Flags().BoolVarP(
		&c.args.Clear,
		"clear",
		"c",
		false,
		"Recover default checker configuration.",
	)
	return c.command
}
