/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/EnnioSimoes/synkgo/internal/config/prompts"
	"github.com/EnnioSimoes/synkgo/internal/engine"
	"github.com/spf13/cobra"
)

func selectTables() {
	tablesSource, err := engine.GetSourceTables()
	if err != nil {
		println("Error to get source tables")
	}

	tablesDestination, err := engine.GetDestinationTables()
	if err != nil {
		println("Error to get destination tables")
	}

	var tablesMatched []string

	for _, table := range tablesSource {
		for _, tablesDest := range tablesDestination {
			if table == tablesDest {
				tablesMatched = append(tablesMatched, table)
			}
		}
	}

	result, _ := prompts.RunMultiSelectTablesPrompt(tablesMatched)

	if len(result) == 0 {
		println("No tables selected")
	}

}

// tablesCmd represents the tables command
var tablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Show source and destination tables",
	Long:  `This command show source and destination tables based in configurations settings on synkgo.json`,
	Run: func(cmd *cobra.Command, args []string) {

		selectFlag, err := cmd.Flags().GetBool("config")
		if err != nil {
			println("Error to get select flag")
		}

		if selectFlag {
			println("Select tables")
			selectTables()

		}
		// result, err := engine.GetSourceTables()
		// if err != nil {
		// 	println("Error to get source tables")
		// }
		// println("Source tables:")
		// for _, table := range result {
		// 	println(table)
		// }

		// println("=====================================")

		// result, err = engine.GetDestinationTables()
		// if err != nil {
		// 	println("Error to get destination tables")
		// }
		// println("Destination tables:")
		// for _, table := range result {
		// 	println(table)
		// }
		// println("\nEnd of tables")
	},
}

func init() {
	rootCmd.AddCommand(tablesCmd)
	tablesCmd.Flags().BoolP("config", "c", false, "Select tables for copy")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tablesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tablesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
