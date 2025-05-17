/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/EnnioSimoes/synkgo/internal/engine"
	"github.com/manifoldco/promptui"
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

	result, _ := runMultiSelectTablesPrompt(tablesMatched)

	if len(result) == 0 {
		println("No tables selected")
	}

}

func runMultiSelectTablesPrompt(items []string) ([]string, error) {
	// items := []string{"clientes", "fornecedores", "pedidos", "produtos", "vendas"}
	finalSelection := []string{}
	selectedItemsMap := make(map[string]bool)

	if len(items) == 0 {
		fmt.Println("Any items to select")
		return nil, fmt.Errorf("no items to select")
	}

	for _, item := range items {
		selectedItemsMap[item] = false
	}

	for {
		displayItems := []string{}
		// Show the selected items with a check mark using the symbol "✓"
		itemMap := make(map[string]string)

		for _, item := range items {
			prefix := "[ ]"
			if selectedItemsMap[item] {
				prefix = "[✅]"
			}
			displayItem := fmt.Sprintf("%s %s", prefix, item)
			displayItems = append(displayItems, displayItem)
			itemMap[displayItem] = item // Mapping the string to the original item
		}

		// Adiciona a opção de finalizar
		doneOption := "--- Finish ---"
		displayItems = append(displayItems, doneOption)
		itemMap[doneOption] = "__DONE__" // A special value to identify the finish action

		prompt := promptui.Select{
			Label: "Select items (use arrows, Enter to toggle, select Finish to Save Selection when done)",
			Items: displayItems,
			Size:  len(displayItems), // Show all items + finish option
			// Customizing templates for better visual feedback (optional)
			Templates: &promptui.SelectTemplates{
				Selected: fmt.Sprintf("%s {{ . | green | bold }}", promptui.IconGood),
				Active:   fmt.Sprintf("%s {{ . | cyan | underline }}", promptui.IconSelect),
				Inactive: "  {{ . }}",
			},
		}

		_, result, err := prompt.Run()

		if err != nil {
			// Treat Ctrl+C interruption
			if err == promptui.ErrInterrupt {
				fmt.Println("User interrupted selection.")
				os.Exit(0) // Graceful exit
			}
			fmt.Printf("Erro no prompt: %v\n", err)
			os.Exit(1) // Exit with error
			return nil, err
		}

		selectedKey := itemMap[result] // Obtain the original item or the special action

		if selectedKey == "__DONE__" {
			break // Exit the loop if "Finish Selection" was chosen
		}

		// Alternates the selection state of the chosen item
		selectedItemsMap[selectedKey] = !selectedItemsMap[selectedKey]
	}

	// Collects the items that were actually selected
	for item, isSelected := range selectedItemsMap {
		if isSelected {
			finalSelection = append(finalSelection, item)
		}
	}

	if len(finalSelection) > 0 {
		fmt.Println("\nSelected items:")
		for _, item := range finalSelection {
			fmt.Printf("- %s\n", item)
		}

		return finalSelection, nil
	}

	fmt.Println("\nNo item was selected.")
	return finalSelection, nil

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
