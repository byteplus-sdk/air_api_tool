package schema

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

	// (FieldName, FieldType) tuple to be added in schema. Use ',' to separate each field_name and data_type.
	// Examples: []string{"custom_field1,int32", "custom_field2,string", "custom_field3,json_string"}
	AddFields []string

	// FieldName to be deleted in schema. Only supports deleting custom fields.
	// Examples: []string{"custom_field1", "custom_field2", "custom_field3"}
	DeleteFields []string

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
		Use:   "schema",
		Short: "Manage the schema of the API. Use `api_tool schema --help` or see https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#Vl0zdDrcloRV68xkVFXc6rvMn3g for more details.",
		Long:  "Manage the schema of the API. Use `api_tool schema --help` or see https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#Vl0zdDrcloRV68xkVFXc6rvMn3g for more details.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := newService(c.args).Run(); err != nil {
				fmt.Printf("[Error] run `api_tool schema` occurred err: %s\n", err)
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
		"Required. Schema's industry. Supports: saas_retail, saas_content.",
	)
	c.command.Flags().StringVarP(
		&c.args.Namespace,
		"namespace",
		"n",
		"default",
		"Optional. Specify a namespace to differentiate multiple schemas for the same industry and table.",
	)
	c.command.Flags().VarP(
		helper.NewStringSliceFlag(&c.args.AddFields),
		"add-field",
		"a",
		"Add custom fields. Field name and Data type need to be separated by `,`.\nExamples: --add-field=custom_field1,int32 --add-field=custom_field2,string\nCurrently supported data types are [\"int32\", \"int64\", \"float\", \"double\", \"string\", \"json_string\"].\nSee https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#Vl0zdDrcloRV68xkVFXc6rvMn3g for more details.",
	)
	c.command.Flags().VarP(
		helper.NewStringSliceFlag(&c.args.DeleteFields),
		"delete-field",
		"d",
		"Delete custom fields.\nExamples: --delete-field=custom_field1 --delete-field=custom_field2.",
	)
	c.command.Flags().BoolVarP(
		&c.args.Show,
		"show",
		"s",
		false,
		"Print the current schema configuration.",
	)
	c.command.Flags().BoolVarP(
		&c.args.Clear,
		"clear",
		"c",
		false,
		"Recover default schema configuration.",
	)
	return c.command
}
