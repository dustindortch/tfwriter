/*
Copyright © 2024 Dustin Dortch
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"tfwriter/schema"

	"github.com/spf13/cobra"
)

var filterProvider string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available resource types",
	Long:  `List available resource types`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			filterProvider = args[0]
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var configSchema = schema.ReadSchema()
		providers := configSchema.Providers

		if filterProvider != "" {
			for provider := range providers {
				parts := strings.Split(provider, "/")
				lastPart := parts[len(parts)-1]
				if lastPart != filterProvider {
					delete(providers, provider)
				}
			}
		}

		if len(providers) == 0 {
			fmt.Println("No providers found")
			os.Exit(1)
		}

		for _, provider := range providers {
			for resource := range provider.Resources {
				fmt.Printf("  Resource: %s\n", resource)
			}
		}
	},
}

func init() {
	listCmd.Flags().StringP("provider", "p", "", "Provider to list resources for")
	rootCmd.AddCommand(listCmd)
}
