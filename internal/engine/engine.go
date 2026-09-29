package engine

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/EnnioSimoes/synkgo/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// isValidTableName verifica nomes de tabela válidos (proteção contra SQL injection).
func isValidTableName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(name)
}

// Sync copia os dados das tabelas configuradas do banco de origem para o
// banco de destino, usando as credenciais definidas em synkgo.json.
func Sync() error {
	// 1. Carregar configuração (synkgo.json é a única fonte de credenciais)
	cfg, err := config.GetConfigFromFile()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// 2. Conectar aos bancos usando as credenciais configuradas
	dbSource, err := config.OpenDB(cfg.DatabaseSource)
	if err != nil {
		return fmt.Errorf("opening source db: %w", err)
	}
	defer dbSource.Close()

	dbDestination, err := config.OpenDB(cfg.DatabaseDestination)
	if err != nil {
		return fmt.Errorf("opening destination db: %w", err)
	}
	defer dbDestination.Close()

	// 3. Testar conexões
	if err := dbSource.Ping(); err != nil {
		return fmt.Errorf("source db unreachable: %w", err)
	}
	if err := dbDestination.Ping(); err != nil {
		return fmt.Errorf("destination db unreachable: %w", err)
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
	return nil
}
