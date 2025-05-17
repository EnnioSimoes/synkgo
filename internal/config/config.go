package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/EnnioSimoes/synkgo/internal/config/prompts"
	"github.com/EnnioSimoes/synkgo/internal/dto"
)

func InitializeConfig() {
	// Initialize the configuration file
	// CreateConfigFile()

	// Prompt for database destination configuration
	fmt.Print("Enter database destination configuration\n")
	host, _ := prompts.RunPrompt("Host")
	port, _ := prompts.RunPrompt("Port")
	username, _ := prompts.RunPrompt("Username")
	password, _ := prompts.RunPrompt("Password")
	database, _ := prompts.RunPrompt("Database")

	// Prompt for tables to sync
	fmt.Print("Enter tables to sync\n")
	tables, _ := prompts.RunMultiSelectTablesPrompt([]string{"clientes", "fornecedores", "pedidos", "produtos", "vendas"})

	// Print the results
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Database: %s\n", database)
	fmt.Printf("Tables: %v\n", tables)
}

func CreateConfigFile() {
	// Initialize the configuration file
	// Check if the file exists
	if _, err := os.Stat("synkgo.json"); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create("synkgo.json")
		if err != nil {
			log.Fatalf("Error creating config file: %v", err)
		}
		defer file.Close()

		// Write default configuration to the file
		defaultConfig := dto.Config{
			DatabaseSource: dto.Database{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "password",
				Database: "source_db",
			},
			DatabaseDestination: dto.Database{
				Host:     "localhost",
				Port:     5432,
				Username: "user",
				Password: "password",
				Database: "destination_db",
			},
			Tables: []string{},
		}

		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			log.Fatalf("Error formatting JSON: %v", err)
		}
		_, err = file.WriteString(string(data))
		if err != nil {
			log.Fatalf("Error writing to config file: %v", err)
		}
		defer file.Close()

	} else {
		// If the file exists, you can read it or perform other actions
		log.Println("Config file already exists.")
	}

}
