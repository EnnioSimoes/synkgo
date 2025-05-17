package prompts

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func RunPrompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt failed %v", err)
	}

	return result, nil
}

func RunMultiSelectTablesPrompt(items []string) ([]string, error) {
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
