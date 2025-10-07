package engine

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/EnnioSimoes/synkgo/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

func joinColumns(columns []string, separator string) string {
	result := ""
	for i, column := range columns {
		if i > 0 {
			result += separator
		}
		result += column
	}
	return result
}

func Sync() {
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

	cfg, err := config.GetConfigFromFile()
	if err != nil {
		fmt.Println("Error getting configuration:", err)
		return
	}

	tablesToSync := cfg.Tables

	// disable FOREIGN_KEY_CHECKS
	_, err = dbDestination.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		fmt.Printf("Error disabling FOREIGN_KEY_CHECKS in destination: %v\n", err)
		return
	}

	for _, tableName := range tablesToSync {
		// go func(tableName string) {
		fmt.Printf("\nProcessing table: %s ...\n", tableName)

		// Select all data from source database
		rows, err := dbSource.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
		if err != nil {
			fmt.Printf("Error selecting data from table %s in source: %v\n", tableName, err)
			return
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			fmt.Printf("Error getting columns for table %s: %v\n", tableName, err)
			return
		}

		placeholders := ""
		for i := 0; i < len(columns); i++ {
			if i > 0 {
				placeholders += ", "
			}
			placeholders += "?"
		}

		// TRUCATE TABLE
		_, err = dbDestination.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName))
		fmt.Printf("Truncating table %s in destination...\n", tableName)
		if err != nil {
			fmt.Printf("Error truncating table %s in destination: %v\n", tableName, err)
			return
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, joinColumns(columns, ","), placeholders)

		insertStmt, err := dbDestination.Prepare(query)
		if err != nil {
			fmt.Printf("Error preparing insert statement for table %s: %v\n", tableName, err)
			return
		}
		defer insertStmt.Close()

		// Loop through rows and insert into destination
		countRows := 0
		for rows.Next() {
			// Create a slice of interface{} to hold the column values
			values := make([]interface{}, len(columns))
			for i := range values {
				values[i] = new(interface{})
			}

			if err := rows.Scan(values...); err != nil {
				fmt.Printf("Error scanning row for table %s: %v\n", tableName, err)
				continue
			}

			// Insert the row into the destination database
			if _, err := insertStmt.Exec(values...); err != nil {
				fmt.Printf("Error inserting row into table %s in destination: %v\n", tableName, err)
				continue
			}
			countRows++
		}

		if err := rows.Err(); err != nil {
			fmt.Printf("Error iterating over rows for table %s: %v\n", tableName, err)
			return
		}

		fmt.Printf("Successfully synced a total: %d rows in table %s\n", countRows, tableName)
	}

	// enable FOREIGN_KEY_CHECKS
	_, err = dbDestination.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		fmt.Printf("Error enabling FOREIGN_KEY_CHECKS in destination: %v\n", err)
		return
	}

	fmt.Println("\nSync completed successfully!")
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
