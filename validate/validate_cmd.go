package validate

import (
	"fmt"
	"github.com/3rd_rec/air_api_tool/validate/parser"
	"github.com/spf13/cobra"
	"strings"
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

	// Local file path that needs to be verified
	FilePath string

	// Preview Lines. By default, the entire file will be validated.
	Lines int

	ContentType string
}

type Command struct {
	command *cobra.Command
	args    *Args
}

func (c *Command) GetCMD() *cobra.Command {
	c.args = &Args{}
	c.command = &cobra.Command{
		Use:   "validate",
		Short: "After configuring the schema and checker, verify your local file and output the verification results.\nUse `api_tool validate --help` or see https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#XDzZdQeNYoxA7IxLTKuc1D1knTh for more details.",
		Long:  "After configuring the schema and checker, verify your local file and output the verification results.\nUse `api_tool validate --help` or see https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#XDzZdQeNYoxA7IxLTKuc1D1knTh for more details.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := NewService(c.args).Run(); err != nil {
				fmt.Printf("[Error] run `api_tool validate` occurred err: %s\n", err)
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
	c.command.Flags().StringVarP(
		&c.args.FilePath,
		"file-path",
		"f",
		"",
		fmt.Sprintf("Required. The file path that needs to be validated."),
	)
	c.command.Flags().StringVarP(
		&c.args.ContentType,
		"content-type",
		"c",
		"",
		fmt.Sprintf("Required. Content type of file. Supports: %s", strings.Join(parser.GetSupportedContentType(), ", ")),
	)
	c.command.Flags().IntVarP(
		&c.args.Lines,
		"lines",
		"l",
		0,
		"Optional. The number of data items to validate. If the file is too large, the validation speed will slow down. "+
			"You can specify the number of data items to validate. By default, the entire file will be validated.",
	)
	return c.command
}
