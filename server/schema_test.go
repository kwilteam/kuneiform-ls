package main

import (
	"fmt"
	"testing"

	"github.com/kwilteam/kwil-db/parse"
)

func Test_ParseSchema(t *testing.T) {
	var schema = `
	database glow;

table data {
    id uuid primary key,
    owner_id uuid notnull,
    foreign key (owner_id) references users(id) on update cascade
    // TODO: add other columns
}`

	res, err := parse.ParseAndValidate([]byte(schema))
	if err != nil {
		t.Error("Error parsing schema")
	}

	fmt.Println(getDiagnostics(res))
}
