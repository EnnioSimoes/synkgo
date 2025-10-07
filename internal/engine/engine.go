package engine

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/EnnioSimoes/synkgo/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// Verifica nomes de tabela válidos
func isValidTableName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(name)
}

func Sync() {
	// 1. Conectar aos bancos
	dbSource, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/database_source")
	if err != nil {
		log.Fatalf("Error opening source DB: %v", err)
	}
	defer dbSource.Close()

	dbDestination, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/database_destination")
	if err != nil {
		log.Fatalf("Error opening destination DB: %v", err)
	}
	defer dbDestination.Close()

	// 2. Testar conexões
	if err := dbSource.Ping(); err != nil {
		log.Fatalf("Source DB unreachable: %v", err)
	}
	if err := dbDestination.Ping(); err != nil {
		log.Fatalf("Destination DB unreachable: %v", err)
	}

	// 3. Carregar configuração
	cfg, err := config.GetConfigFromFile()
	if err != nil {
		log.Printf("Error getting configuration: %v", err)
		return
	}

	tablesToSync := cfg.Tables
	fmt.Println("Starting synchronization for tables:", tablesToSync)

	// 4. Processar cada tabela
	for _, tableName := range tablesToSync {

		if !isValidTableName(tableName) {
			log.Printf("⚠️ Skipping invalid table name: %s\n", tableName)
			continue
		}

		fmt.Printf("\n🔄 Syncing table: %s ...\n", tableName)

		rows, err := dbSource.Query(fmt.Sprintf("SELECT * FROM `%s`", tableName))
		if err != nil {
			log.Printf("Error selecting data from %s: %v\n", tableName, err)
			continue
		}

		columns, err := rows.Columns()
		if err != nil {
			log.Printf("Error getting columns from %s: %v\n", tableName, err)
			rows.Close()
			continue
		}

		// Preparar placeholders
		placeholders := strings.TrimRight(strings.Repeat("?,", len(columns)), ",")

		// Desabilitar foreign keys
		_, _ = dbDestination.Exec("SET FOREIGN_KEY_CHECKS = 0")

		// Truncar a tabela
		fmt.Printf("Truncating table %s...\n", tableName)
		if _, err = dbDestination.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)); err != nil {
			log.Printf("Error truncating table %s: %v\n", tableName, err)
			rows.Close()
			continue
		}

		// Preparar INSERT
		query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
			tableName, strings.Join(columns, ","), placeholders)

		insertStmt, err := dbDestination.Prepare(query)
		if err != nil {
			log.Printf("Error preparing insert statement for %s: %v\n", tableName, err)
			rows.Close()
			continue
		}

		count := 0
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				log.Printf("Error scanning row in %s: %v\n", tableName, err)
				continue
			}

			if _, err := insertStmt.Exec(values...); err != nil {
				log.Printf("Error inserting row into %s: %v\n", tableName, err)
				continue
			}

			count++
		}

		rows.Close()
		insertStmt.Close()

		if err := rows.Err(); err != nil {
			log.Printf("Error iterating rows for %s: %v\n", tableName, err)
		}

		// Reativar foreign keys
		_, _ = dbDestination.Exec("SET FOREIGN_KEY_CHECKS = 1")

		fmt.Printf("✅ Successfully synced %d rows in table %s\n", count, tableName)
	}

	fmt.Println("\n🎉 Sync completed successfully!")
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
