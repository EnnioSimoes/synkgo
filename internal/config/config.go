package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/EnnioSimoes/synkgo/internal/config/prompts"
	"github.com/EnnioSimoes/synkgo/internal/dto"
)

func InitializeConfig() {
	// Initialize the configuration file
	var config dto.Config

	if _, err := os.Stat("synkgo.json"); os.IsNotExist(err) {
		// Prompt for database destination configuration
		fmt.Print("Enter database source configuration\n")
		config.DatabaseSource.Host, _ = prompts.RunPrompt("Host")
		config.DatabaseSource.Port, _ = prompts.RunPrompt("Port")
		config.DatabaseSource.Username, _ = prompts.RunPrompt("Username")
		config.DatabaseSource.Password, _ = prompts.RunPrompt("Password")
		config.DatabaseSource.Database, _ = prompts.RunPrompt("Database")

		fmt.Print("\nEnter database destination configuration\n")
		config.DatabaseDestination.Host, _ = prompts.RunPrompt("Host")
		config.DatabaseDestination.Port, _ = prompts.RunPrompt("Port")
		config.DatabaseDestination.Username, _ = prompts.RunPrompt("Username")
		config.DatabaseDestination.Password, _ = prompts.RunPrompt("Password")
		config.DatabaseDestination.Database, _ = prompts.RunPrompt("Database")

		// Print the results
		fmt.Printf("Host: %s\n", config.DatabaseSource.Host)
		fmt.Printf("Port: %s\n", config.DatabaseSource.Port)
		fmt.Printf("Username: %s\n", config.DatabaseSource.Username)
		fmt.Printf("Password: %s\n", config.DatabaseSource.Password)
		fmt.Printf("Database: %s\n", config.DatabaseSource.Database)
		fmt.Printf("Host: %s\n", config.DatabaseDestination.Host)
		fmt.Printf("Port: %s\n", config.DatabaseDestination.Port)
		fmt.Printf("Username: %s\n", config.DatabaseDestination.Username)
		fmt.Printf("Password: %s\n", config.DatabaseDestination.Password)
		fmt.Printf("Database: %s\n", config.DatabaseDestination.Database)
		fmt.Printf("Tables: %v\n", config.Tables)
		// Save the configuration to a file
		CreateConfigFile(config)
	}

	matchedTables, err := CompareTables()
	if err != nil {
		println("Error to compare tables")
	}
	// fmt.Printf("Matched tables: %v\n", matchedTables)
	config, err = GetConfigFromFile()
	if err != nil {
		println("Error to get config file")
		return
	}

	// Prompt for tables to sync
	fmt.Print("Enter tables to sync\n")
	config.Tables, _ = prompts.RunMultiSelectTablesPrompt(matchedTables)
	CreateConfigFile(config)

}

func CreateConfigFile(config dto.Config) {
	// Initialize the configuration file
	// Check if the file exists
	// Create the file if it doesn't exist
	file, err := os.Create("synkgo.json")
	if err != nil {
		log.Fatalf("Error creating config file: %v", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Error formatting JSON: %v", err)
	}
	_, err = file.WriteString(string(data))
	if err != nil {
		log.Fatalf("Error writing to config file: %v", err)
	}
	defer file.Close()

	// } else {
	// 	// If the file exists, you can read it or perform other actions
	// 	log.Println("Config file already exists.")
	// }

}

func GetConfigFromFile() (dto.Config, error) {
	// Read the configuration file
	file, err := os.Open("synkgo.json")
	if err != nil {
		return dto.Config{}, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config dto.Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return dto.Config{}, fmt.Errorf("error decoding config file: %v", err)
	}

	return config, nil
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

func GetSourceTables() ([]string, error) {
	config, err := GetConfigFromFile()
	if err != nil {
		println("Error to get config file")
		return nil, fmt.Errorf("error to get config file")
	}

	stringConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DatabaseSource.Username,
		config.DatabaseSource.Password,
		config.DatabaseSource.Host,
		config.DatabaseSource.Port,
		config.DatabaseSource.Database,
	)
	dbSource, err := sql.Open("mysql", stringConnection)

	if err != nil {
		panic(err)
	}
	defer dbSource.Close()
	result, err := getTables(dbSource)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func GetDestinationTables() ([]string, error) {
	config, err := GetConfigFromFile()
	if err != nil {
		println("Error to get config file")
		return nil, fmt.Errorf("error to get config file")
	}

	stringConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DatabaseDestination.Username,
		config.DatabaseDestination.Password,
		config.DatabaseDestination.Host,
		config.DatabaseDestination.Port,
		config.DatabaseDestination.Database,
	)
	dbDest, err := sql.Open("mysql", stringConnection)

	if err != nil {
		panic(err)
	}
	defer dbDest.Close()
	result, err := getTables(dbDest)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func CompareTables() ([]string, error) {
	sourceTables, err := GetSourceTables()
	if err != nil {
		println("Error to get source tables")
		return nil, fmt.Errorf("error to get source tables")
	}

	destinationTables, err := GetDestinationTables()
	if err != nil {
		println("Error to get destination tables")
		return nil, fmt.Errorf("error to get destination tables")
	}

	var tablesMatched []string

	for _, table := range sourceTables {
		for _, tablesDest := range destinationTables {
			if table == tablesDest {
				tablesMatched = append(tablesMatched, table)
			}
		}
	}

	return tablesMatched, nil
}
