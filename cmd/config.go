/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/EnnioSimoes/synkgo/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the current configuration",
	Long: `Show the current configuration saved in the synkgo.json file.
This command will display the database source and destination configuration,
as well as the tables to sync.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("config called")
		createFlag, err := cmd.Flags().GetBool("create")
		if err != nil {
			println("Error to get create flag")
			return
		}

		if _, err := os.Stat("synkgo.json"); err == nil {
			if createFlag {
				fmt.Println("synkgo.json already exists. Use 'synkgo init' to reinitialize.")
				return
			}
		}

		if createFlag {
			config.InitializeConfig()
		}

		showConfig()
	},
}

func showConfig() {
	config, err := config.GetConfigFromFile()
	if err != nil {
		println("Error to get config file")
		return
	}

	fmt.Print("Database source configuration:\n")
	fmt.Printf("Host: %s\n", config.DatabaseSource.Host)
	fmt.Printf("Port: %s\n", config.DatabaseSource.Port)
	fmt.Printf("Username: %s\n", config.DatabaseSource.Username)
	fmt.Printf("Password: %s\n", config.DatabaseSource.Password)
	fmt.Printf("Database: %s\n", config.DatabaseSource.Database)

	fmt.Print("\nDatabase destination configuration:\n")
	fmt.Printf("Host: %s\n", config.DatabaseDestination.Host)
	fmt.Printf("Port: %s\n", config.DatabaseDestination.Port)
	fmt.Printf("Username: %s\n", config.DatabaseDestination.Username)
	fmt.Printf("Password: %s\n", config.DatabaseDestination.Password)
	fmt.Printf("Database: %s\n", config.DatabaseDestination.Database)

	fmt.Print("\nTables configuration to sync:\n")
	for _, table := range config.Tables {
		fmt.Printf("Table: %s\n", table)
	}
	fmt.Print("\n")
	fmt.Print("Config file loaded successfully\n")
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configCmd.Flags().BoolP("create", "c", false, "Create a new synkgo.json file if it doesn't exist")
}
