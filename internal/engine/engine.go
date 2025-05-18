package engine

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/EnnioSimoes/synkgo/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbSource, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/database_source")
	if err != nil {
		panic(err)
	}
	defer dbSource.Close()

	dbDestination, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/database_destination")
	if err != nil {
		panic(err)
	}
	defer dbDestination.Close()
	// createConfigFile()
	// configSource, err := getConfig("destination")
	// if err != nil {
	// 	println("Erro ao obter a configuração de origem")
	// 	return
	// }
	// println("Configuração de origem:")
	// fmt.Println(configSource)

	// getSourceTables(dbSource)
	// getSourceTables(dbDestination)
}

// type Database struct {
// 	Host     string `json:"host"`
// 	Port     int    `json:"port"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	Database string `json:"database"`
// }

// type Config struct {
// 	DatabaseSource      Database `json:"database_source"`
// 	DatabaseDestination Database `json:"database_destination"`
// }

// func createConfigFile() {
// 	file, err := os.Create("synkgo.json")
// 	if err != nil {
// 		println("Erro ao criar o arquivo")
// 	}
// 	defer file.Close()

// 	config := dto.Config{
// 		DatabaseSource: dto.Database{
// 			Host:     "localhost",
// 			Port:     "5432",
// 			Username: "user",
// 			Password: "password",
// 			Database: "source_db",
// 		},
// 		DatabaseDestination: dto.Database{
// 			Host:     "localhost",
// 			Port:     "5432",
// 			Username: "user",
// 			Password: "password",
// 			Database: "destination_db",
// 		},
// 	}

// 	data, err := json.MarshalIndent(config, "", "  ")
// 	if err != nil {
// 		println("Erro ao formatar o JSON")
// 		return
// 	}
// 	_, err = file.WriteString(string(data))
// 	if err != nil {
// 		println("Erro ao escrever no arquivo")
// 	}
// 	defer file.Close()
// 	println("Arquivo criado com sucesso")
// }

// func getConfig(source string) (dto.Database, error) {
// 	file, err := os.Open("./synkgo.json")
// 	if err != nil {
// 		println("Erro ao abrir o arquivo")
// 		return dto.Database{}, err
// 	}
// 	defer file.Close()

// 	var config dto.Config
// 	decoder := json.NewDecoder(file)
// 	err = decoder.Decode(&config)
// 	if err != nil {
// 		println("Erro ao decodificar o JSON")
// 		return dto.Database{}, err
// 	}

// 	if source == "source" {
// 		return config.DatabaseSource, nil
// 	}

// 	if source == "destination" {
// 		return config.DatabaseDestination, nil
// 	}

// 	return dto.Database{}, nil
// }

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
	// check if the file not exists
	if _, err := os.Stat("synkgo.json"); os.IsNotExist(err) {
		fmt.Println("Config file not found")
		return nil, fmt.Errorf("config file not found")
	}

	config, err := config.GetConfigFromFile()
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
	// dbSource, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/database_source")

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
	dbDestination, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/database_destination")
	if err != nil {
		panic(err)
	}
	defer dbDestination.Close()
	result, err := getTables(dbDestination)
	if err != nil {
		panic(err)
	}
	return result, nil
}
