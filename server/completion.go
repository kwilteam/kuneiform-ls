package main

import (
	"log/slog"

	"github.com/kwilteam/kwil-db/parse"
	"github.com/sourcegraph/go-lsp"
)

// onCompletion handler support

// Defaults
func (l *lspHandler) getCompletionItems(r *parse.SchemaParseResult, pos int) []lsp.CompletionItem {

	inTable := isWithinTableBlock(r, pos)
	inProcedures := isWithinProcedureBlock(r, pos)
	inActions := isWithinActionBlock(r, pos)
	inForeginProcedures := isWithinForeignProcedureBlock(r, pos)
	dbDefined := getDatabaseName(r) != ""

	l.logger.Debug("Parse result: ", slog.Bool("dbDefined", dbDefined), slog.Bool("inTable", inTable), slog.Bool("inProcedures", inProcedures), slog.Bool("inActions", inActions), slog.Bool("inForeginProcedures", inForeginProcedures))

	tables := getTableCompletionItems(r)
	params := getParamsCompletionItems(r, pos)
	procedures := getProcedureCompletionItems(r, pos)
	actions := getActionCompletionItems(r, pos)
	extensions := getExtensionsCompletionItems(r)

	var items []lsp.CompletionItem

	if !dbDefined {
		// Database is not defined, so we can't proceed with any completions
		return dbCompletionItems

	} else if !inTable && !inProcedures && !inActions && !inForeginProcedures {
		// Not in any block, base level completions and modifiers, datatypes
		return append(append(kfCompletionItems, datatypeCompletionItems...), modifierCompletionItems...)

	} else if inTable {
		// datatypes, table completion items,
		return append(append(tables, datatypeCompletionItems...), tableCompletionItems...)

	} else if inProcedures {
		items = append(procedures, params...)             // procedures and params
		items = append(items, tables...)                  // tables
		items = append(items, methodCompletionItems...)   // methods
		items = append(items, datatypeCompletionItems...) // datatypes
		return append(items, modifierCompletionItems...)  // modifiers

	} else if inActions {
		items = append(actions, procedures...)            // actions and procedures
		items = append(items, extensions...)              // extensions
		items = append(items, params...)                  // params
		items = append(items, tables...)                  // tables
		items = append(items, methodCompletionItems...)   // methods
		items = append(items, datatypeCompletionItems...) // datatypes
		return append(items, modifierCompletionItems...)  // modifiers

	} else if inForeginProcedures {
		items := append(params, datatypeCompletionItems...) // datatypes, params
		return append(items, modifierCompletionItems...)    // modifiers
	}

	return items
}

func getTableCompletionItems(r *parse.SchemaParseResult) []lsp.CompletionItem {
	return getDefaultCompletionItems(getTables(r))
}

func getActionCompletionItems(r *parse.SchemaParseResult, pos int) []lsp.CompletionItem {
	items := []lsp.CompletionItem{}

	// Actions can only be called within an action block
	if isWithinActionBlock(r, pos) {
		actions := getActions(r)
		for _, action := range actions {
			items = append(items, lsp.CompletionItem{
				Label:            action,
				Kind:             lsp.CIKFunction,
				InsertText:       action + "(${1:params});",
				InsertTextFormat: lsp.ITFSnippet,
			})
		}
	}
	return items
}

func getProcedureCompletionItems(r *parse.SchemaParseResult, pos int) []lsp.CompletionItem {
	// Procedures can be called either from procedure block or action block
	items := []lsp.CompletionItem{}
	if isWithinProcedureBlock(r, pos) || isWithinActionBlock(r, pos) {
		procedures := getProcedures(r)
		for _, procedure := range procedures {
			items = append(items, lsp.CompletionItem{
				Label:            procedure,
				Kind:             lsp.CIKFunction,
				InsertText:       procedure + "(${1:params});",
				InsertTextFormat: lsp.ITFSnippet,
			})
		}
	}
	return items
}

func getExtensionsCompletionItems(r *parse.SchemaParseResult) []lsp.CompletionItem {
	aliases := getExtensions(r)
	var items []lsp.CompletionItem
	for _, alias := range aliases {
		items = append(items, lsp.CompletionItem{
			Label:            alias,
			Kind:             lsp.CIKFunction,
			InsertText:       alias + ".${1:ext_method}(${2:params});",
			InsertTextFormat: lsp.ITFSnippet,
		})
	}
	return items
}

func getParamsCompletionItems(r *parse.SchemaParseResult, pos int) []lsp.CompletionItem {
	items := []lsp.CompletionItem{}
	actionLocs := getActionLocations(r)
	params := make(map[string]struct{})
	// push only the params for the current action
	for _, loc := range actionLocs {
		if pos >= loc.start && pos <= loc.end {
			for _, param := range loc.params {
				params[param] = struct{}{}
			}
		}
	}

	procedureLocs := getProcedureLocations(r)
	for _, loc := range procedureLocs {
		if pos >= loc.start && pos <= loc.end {
			for _, param := range loc.params {
				params[param] = struct{}{}
			}
		}
	}

	tableLocs := getTableLocations(r)
	for _, loc := range tableLocs {
		for _, param := range loc.params {
			params[param] = struct{}{}
		}
	}

	for param := range params {
		items = append(items, lsp.CompletionItem{
			Label:            param,
			Kind:             lsp.CIKVariable,
			InsertText:       param,
			InsertTextFormat: lsp.ITFPlainText,
		})
	}
	return items
}

func isWithinActionBlock(r *parse.SchemaParseResult, pos int) bool {
	// Check if the position is within an action block
	locs := getActionLocations(r)
	for _, loc := range locs {
		if pos >= loc.start && pos <= loc.end {
			return true
		}
	}
	return false
}

func isWithinProcedureBlock(r *parse.SchemaParseResult, pos int) bool {
	// Check if the position is within a procedure block
	locs := getProcedureLocations(r)
	for _, loc := range locs {
		if pos >= loc.start && pos <= loc.end {
			return true
		}
	}
	return false
}

func isWithinTableBlock(r *parse.SchemaParseResult, pos int) bool {
	// Check if the position is within a table block
	locs := getTableLocations(r)
	for _, loc := range locs {
		if pos >= loc.start && pos <= loc.end {
			return true
		}
	}
	return false
}

func isWithinForeignProcedureBlock(r *parse.SchemaParseResult, pos int) bool {
	// Check if the position is within a foreign procedure block
	locs := getForeignProcedureLocations(r)
	for _, loc := range locs {
		if pos >= loc.start && pos <= loc.end {
			return true
		}
	}
	return false
}
