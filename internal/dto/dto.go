package dto

type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Config struct {
	DatabaseSource      Database `json:"database_source"`
	DatabaseDestination Database `json:"database_destination"`
	Tables              []string `json:"tables"`
}

func (c *Config) SetDatabaseSource(db Database) {
	c.DatabaseSource = db
}

func (c *Config) SetDatabaseDestination(db Database) {
	c.DatabaseDestination = db
}

func (c *Config) SetTables(tables []string) {
	c.Tables = tables
}

func (c *Config) GetDatabaseSource() Database {
	return c.DatabaseSource
}

func (c *Config) GetDatabaseDestination() Database {
	return c.DatabaseDestination
}

func (c *Config) GetTables() []string {
	return c.Tables
}

func (c *Config) AddTable(table string) {
	c.Tables = append(c.Tables, table)
}

func (c *Config) RemoveTable(table string) {
	for i, t := range c.Tables {
		if t == table {
			c.Tables = append(c.Tables[:i], c.Tables[i+1:]...)
			break
		}
	}
}

func (c *Config) ClearTables() {
	c.Tables = []string{}
}
