/*
Copyright © 2024 Dustin Dortch
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"tfwriter/schema"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

func resourceDefinition(resource string) schema.Resource {
	var configSchema = schema.ReadSchema()
	providers := configSchema.Providers

	var resourceDefinition schema.Resource

	for _, provider := range providers {
		for k, v := range provider.Resources {
			if resource == k {
				resourceDefinition = v
			}
		}
	}

	return resourceDefinition
}

var resourceCmd = &cobra.Command{
	Use:   "resource <resource_type>",
	Short: "Generate boilerplate code for a specific Terraform resource type",
	Long: `Generate boilerplate code for a specific Terraform resource type. For example:
	
	tfwriter resource aws_instance
	`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Resource type is required")
			cmd.Help()
			return
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		f := hclwrite.NewEmptyFile()
		rootBody := f.Body()

		for i := 0; i < len(args); i++ {
			resourceType, resourceLabel, label := strings.Cut(args[i], ".")
			if !label {
				resourceLabel = strconv.Itoa(i)
			}
			resourceDefinition := resourceDefinition(resourceType)
			resourceBlock := rootBody.AppendNewBlock("resource", []string{resourceType, resourceLabel})
			resourceBody := resourceBlock.Body()

			for k, v := range resourceDefinition.Block.Attributes {
				if !v.Computed && v.Required {
					resourceBody.SetAttributeValue(k, cty.StringVal(fmt.Sprintf("Type: %s, Required: %s", v.TypeString(), strconv.FormatBool(v.Required))))
				}
			}

			for k, v := range resourceDefinition.Block.Attributes {
				if !v.Computed && !v.Required {
					resourceBody.SetAttributeValue(k, cty.StringVal(fmt.Sprintf("Type: %s, Required: %s", v.TypeString(), strconv.FormatBool(v.Required))))
				}
			}

			rootBody.AppendNewline()
		}

		fmt.Printf("%s", f.Bytes())
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
}
