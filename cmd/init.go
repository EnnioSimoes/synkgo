/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/EnnioSimoes/synkgo/internal/config"
	"github.com/spf13/cobra"
)

func InitializeConfig() {
	config.InitializeConfig()
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize synkgo configuration file (synkgo.json)",
	Long: `Initialize synkgo configuration file (synkgo.json).
If the file already exists, it will not be overwritten.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
