// Copyright © 2019 morgulbrut
// This work is free. You can redistribute it and/or modify it under the
// terms of the Do What The Fuck You Want To Public License, Version 2,
// as published by Sam Hocevar. See the LICENSE file for more details.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// alternativesCmd represents the alternatives command
var alternativesCmd = &cobra.Command{
	Use:   "alternatives",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("alternatives called")
	},
}

func init() {
	rootCmd.AddCommand(alternativesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// alternativesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alternativesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
