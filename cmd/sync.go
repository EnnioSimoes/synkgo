/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/EnnioSimoes/synkgo/internal/engine"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Start the synchronization process",
	Long: `Start the synchronization process between source and destination databases.
This command will read the configuration from the synkgo.json file and perform
the data synchronization based on the specified tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("sync called")

		engine.Sync()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
