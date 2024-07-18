package main

import (
	"regexp"

	"github.com/kwilteam/kwil-db/parse"
	"github.com/sourcegraph/go-lsp"
)

// Responsible for collecting required information from the parsed schema
const (
	actionOffset    = 7  // len(action) + 1
	procedureOffset = 10 // len(procedure) + 1
)

type kfDocs struct {
	rawKf        string
	parsedSchema *parse.SchemaParseResult
}

// Action and Procedure
type methodsLocation struct {
	location
	name   string
	params []string
}

type location struct {
	start, end int
}

func getDiagnostics(r *parse.SchemaParseResult) []lsp.Diagnostic {
	if r == nil {
		return []lsp.Diagnostic{}
	}

	diagnosis := make([]lsp.Diagnostic, 0)
	for _, err := range r.ParseErrs.Errors() {
		d := lsp.Diagnostic{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      err.Position.StartLine - 1,
					Character: err.Position.StartCol,
				},
				End: lsp.Position{
					Line:      err.Position.EndLine - 1,
					Character: err.Position.EndCol,
				},
			},
			Severity: lsp.Error,
			Message:  err.Err.Error() + ": " + err.Message,
		}
		diagnosis = append(diagnosis, d)
	}
	return diagnosis
}

func getTables(r *parse.SchemaParseResult) []string {
	items := make([]string, 0)
	if r == nil || r.Schema == nil {
		return items
	}

	for _, table := range r.Schema.Tables {
		items = append(items, table.Name)
	}
	return items
}

func getActions(r *parse.SchemaParseResult) []string {
	items := make([]string, 0)
	if r == nil || r.Schema == nil || r.Schema.Actions == nil {
		return items
	}

	for _, action := range r.Schema.Actions {
		items = append(items, action.Name)
	}
	return items
}

func getProcedures(r *parse.SchemaParseResult) []string {
	items := make([]string, 0)
	if r == nil || r.Schema == nil || r.Schema.Procedures == nil {
		return items
	}
	for _, procedure := range r.Schema.Procedures {
		items = append(items, procedure.Name)
	}
	return items
}

func getDatabaseName(r *parse.SchemaParseResult) string {
	if r == nil || r.Schema == nil {
		return ""
	}
	return r.Schema.Name
}

func getExtensions(r *parse.SchemaParseResult) []string {
	items := make([]string, 0)
	if r == nil || r.Schema == nil || r.Schema.Extensions == nil {
		return items
	}
	for _, extension := range r.Schema.Extensions {
		items = append(items, extension.Alias)
	}
	return items
}

func getActionLocations(r *parse.SchemaParseResult) []methodsLocation {
	locs := make([]methodsLocation, 0)
	if r == nil || r.Schema == nil || r.Schema.Actions == nil || r.SchemaInfo == nil {
		return locs
	}

	for _, action := range r.Schema.Actions {
		block, ok := r.SchemaInfo.Blocks[action.Name]
		if !ok {
			continue
		}
		params := extractAndDeduplicateMethodParams(action.Parameters, action.Body)
		locs = append(locs, methodsLocation{
			location: location{
				start: block.AbsStart,
				end:   block.AbsEnd,
			},
			name:   action.Name,
			params: params,
		})
	}
	return locs
}

func getProcedureLocations(r *parse.SchemaParseResult) []methodsLocation {
	locs := make([]methodsLocation, 0)
	if r == nil || r.Schema == nil || r.Schema.Procedures == nil || r.SchemaInfo == nil {
		return locs
	}

	for _, procedure := range r.Schema.Procedures {
		block, ok := r.SchemaInfo.Blocks[procedure.Name]
		if !ok {
			continue
		}

		var procedureParams []string
		for _, param := range procedure.Parameters {
			procedureParams = append(procedureParams, param.Name)
		}

		params := extractAndDeduplicateMethodParams(procedureParams, procedure.Body)
		locs = append(locs, methodsLocation{
			location: location{
				start: block.AbsStart,
				end:   block.AbsEnd,
			},
			name:   procedure.Name,
			params: params,
		})
	}
	return locs

}

func getTableLocations(r *parse.SchemaParseResult) []methodsLocation {
	locs := make([]methodsLocation, 0)
	if r == nil || r.Schema == nil || r.Schema.Tables == nil || r.SchemaInfo == nil {
		return locs
	}

	for _, table := range r.Schema.Tables {
		block, ok := r.SchemaInfo.Blocks[table.Name]
		if !ok {
			continue
		}

		// params are the columns of the table
		var params []string
		for _, column := range table.Columns {
			params = append(params, column.Name)
		}

		locs = append(locs, methodsLocation{
			location: location{
				start: block.AbsStart,
				end:   block.AbsEnd,
			},
			name:   table.Name,
			params: params,
		})
	}
	return locs
}

func getForeignProcedureLocations(r *parse.SchemaParseResult) []methodsLocation {
	locs := make([]methodsLocation, 0)
	if r == nil || r.Schema == nil || r.Schema.ForeignProcedures == nil || r.SchemaInfo == nil {
		return locs
	}

	for _, procedure := range r.Schema.ForeignProcedures {
		block, ok := r.SchemaInfo.Blocks[procedure.Name]
		if !ok {
			continue
		}

		var procedureParams []string
		for _, param := range procedure.Parameters {
			procedureParams = append(procedureParams, param.Name)
		}

		locs = append(locs, methodsLocation{
			location: location{
				start: block.AbsStart,
				end:   block.AbsEnd,
			},
			name:   procedure.Name,
			params: procedureParams,
		})
	}
	return locs
}

func extractAndDeduplicateMethodParams(params []string, body string) []string {
	pattern := `\$[a-zA-Z_]\w*`

	re := regexp.MustCompile(pattern)
	// extract the params in the action or procedure body
	matches := re.FindAllString(body, -1)

	// deduplicate the matches
	uniqueParams := make(map[string]struct{})
	for _, param := range params {
		uniqueParams[param] = struct{}{}
	}

	// remove the params that are already in the params list
	for _, match := range matches {
		if _, ok := uniqueParams[match]; !ok {
			params = append(params, match)
		}
	}

	return params
}

func getTokenPosition(uri lsp.DocumentURI, r *parse.SchemaParseResult, token string) []lsp.Location {
	var locations []lsp.Location
	if r == nil || r.Schema == nil || r.SchemaInfo == nil {
		return locations
	}

	// Check if token exists
	block, ok := r.SchemaInfo.Blocks[token]
	if !ok {
		return locations
	}

	// action or procedure?
	offset := procedureOffset
	if _, ok := r.ParsedActions[token]; ok {
		offset = actionOffset
	}

	locations = append(locations, lsp.Location{
		URI: uri,
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      block.StartLine - 1,
				Character: block.StartCol + offset,
			},
			End: lsp.Position{
				Line:      block.EndLine - 1,
				Character: block.EndCol,
			},
		},
	})

	return locations
}
