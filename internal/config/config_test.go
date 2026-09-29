package config

import (
	"testing"

	"github.com/EnnioSimoes/synkgo/internal/dto"
)

func TestBuildDSN(t *testing.T) {
	db := dto.Database{
		Host:     "localhost",
		Port:     "3306",
		Username: "root",
		Password: "secret",
		Database: "mydb",
	}
	want := "root:secret@tcp(localhost:3306)/mydb"
	got := buildDSN(db)
	if got != want {
		t.Errorf("buildDSN() = %q, want %q", got, want)
	}
}

func TestBuildDSN_EmptyFields(t *testing.T) {
	db := dto.Database{}
	want := ":@tcp(:)/"
	got := buildDSN(db)
	if got != want {
		t.Errorf("buildDSN() = %q, want %q", got, want)
	}
}
