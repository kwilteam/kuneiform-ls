package main

import "github.com/sourcegraph/go-lsp"

var (
	triggerKeywords = []string{"table, action, use, procedure, foreign, database"}

	// database completion items
	dbCompletionItems = []lsp.CompletionItem{
		{
			Label:            "database dbname;",
			Kind:             lsp.CIKClass,
			InsertText:       "database ${1:dbname};",
			InsertTextFormat: lsp.ITFSnippet,
		},
	}

	// base level kf completion items
	kfCompletionItems = []lsp.CompletionItem{
		{ // table declaration
			Label:            "table name {columns..}",
			Kind:             lsp.CIKClass,
			InsertText:       "table ${1:name} {\n\t${2:columns}\n}",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // action declaration
			Label:            "action name(params...) modifiers {body}",
			Kind:             lsp.CIKClass,
			InsertText:       "action ${1:name}(${2:params}) ${3:modifiers} {\n\t${4:body}\n}",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // use declaration
			Label:            "use ext {data} as alias;",
			Kind:             lsp.CIKClass,
			InsertText:       "use ${1:module} {\n\t${2:kvpairs}\n} as ${3:alias};",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // procedure declaration without return type
			Label:            "procedure name(params...) modifiers {body}",
			Kind:             lsp.CIKClass,
			InsertText:       "procedure ${1:name}(${2:params}) ${3:modifiers}  {\n\t${4:body}\n}",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // procedure declaration with return type
			Label:            "procedure name(params...) modifiers returns (return_vals...) {body}",
			Kind:             lsp.CIKClass,
			InsertText:       "procedure ${1:name}(${2:params}) ${3:modifiers} returns (${4:return_vals}) {\n\t${5:body}\n}",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // procedure declaration with table return type
			Label:            "procedure name(params...) modifiers returns table(return_vals...) {body}",
			Kind:             lsp.CIKClass,
			InsertText:       "procedure ${1:name}(${2:params}) ${3:modifiers} returns table(${4:return_vals}) {\n\t${5:body}\n}",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			// foreign procedure declaration without return type
			Label:            "foreign procedure name(params...) modifiers",
			Kind:             lsp.CIKClass,
			InsertText:       "foreign procedure ${1:name}(${2:params}) ${3:modifiers}",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // foreign procedure declaration
			Label:            "foreign procedure name(params...) returns (return_types...)",
			Kind:             lsp.CIKClass,
			InsertText:       "foreign procedure ${1:name}(${2:params}) returns (${3:return_type})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // foreign procedure declaration  with table return type
			Label:            "foreign procedure name(params...) returns table(return_types...)",
			Kind:             lsp.CIKClass,
			InsertText:       "foreign procedure ${1:name}(${2:params}) returns table(${3:return_type})",
			InsertTextFormat: lsp.ITFSnippet,
		},
	}

	// modifier and contextual completion items
	modifierAndContextualKeys = []string{"@caller", "@signer", "@txid", "@height", "@action", "@dataset", "public", "private", "view", "owner", "returns"}
	modifierCompletionItems   = getDefaultCompletionItems(modifierAndContextualKeys)

	// datatype completion items
	datatypes             = []string{"text", "int", "uuid", "blob", "bool", "uint256", "decimal"} // decimal(precision, scale)
	decimalCompletionItem = lsp.CompletionItem{
		Label:            "decimal(precision, scale)",
		Kind:             lsp.CIKProperty,
		InsertText:       "decimal(${1:precision}, ${2:scale})",
		InsertTextFormat: lsp.ITFSnippet,
	}
	datatypeCompletionItems = append(getDefaultCompletionItems(datatypes), decimalCompletionItem)

	// table completion items
	tableKeywords                = []string{"notnull", "primary key", "default", "unique", "on_delete", "on_update", "cascade", "restrict", "set_null", "set_default", "no_action", "references"}
	tableKeywordsCompletionItems = getDefaultCompletionItems(tableKeywords)
	tableCompletionItems         = append([]lsp.CompletionItem{
		{ // maxlen() attribute
			Label:            "maxlen(num)",
			Kind:             lsp.CIKFunction,
			InsertText:       "maxlen(${1:num})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // minlen() attribute
			Label:            "minlen(num)",
			Kind:             lsp.CIKFunction,
			InsertText:       "minlen(${1:num})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // max attribute
			Label:            "max(num)",
			Kind:             lsp.CIKFunction,
			InsertText:       "max(${1:num})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // min attribute
			Label:            "min(num)",
			Kind:             lsp.CIKFunction,
			InsertText:       "min(${1:num})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // foreign key declaration
			Label:            "foreign key column references table(column)",
			Kind:             lsp.CIKClass,
			InsertText:       "foreign key ${1:column} references ${2:table}(${3:column})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // foreign key declaration with on delete|update action
			Label:            "foreign key column references table(column) on delete|update action",
			Kind:             lsp.CIKClass,
			InsertText:       "foreign key ${1:column} references ${2:table}(${3:column}) on ${4:delete|update} ${5:action}",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // index declaration
			Label:            "#index_name index|unique|primary(columns...)",
			Kind:             lsp.CIKClass,
			InsertText:       "${1:index_name} ${2:index|primary|unqiue}(${3:columns})",
			InsertTextFormat: lsp.ITFSnippet,
		},
	}, tableKeywordsCompletionItems...)

	// SQL specific completions
	sqlFunctionsCompletionItems = []lsp.CompletionItem{
		{ // uuid generate function
			Label:            "uuid_generate_v5(uuid, text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "uuid_generate_v5(${1:uuid}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{ // abs function
			Label:            "abs(num)",
			Kind:             lsp.CIKFunction,
			InsertText:       "abs(${1:num})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		// Encoding Functions
		{
			Label:            "encode(blob, text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "encode(${1:blob}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "decode(text, text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "decode(${1:text}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		// Digest Functions
		{
			Label:            "digest(text|blob, text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "digest(${1:text|blob}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		// Array Functions
		{
			Label:            "array_append(any[], any)",
			Kind:             lsp.CIKFunction,
			InsertText:       "array_append(${1:any[]}, ${2:any})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "array_prepend(any, any[])",
			Kind:             lsp.CIKFunction,
			InsertText:       "array_prepend(${1:any}, ${2:any[]})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "array_cat(any[], any[])",
			Kind:             lsp.CIKFunction,
			InsertText:       "array_cat(${1:any[]}, ${2:any[]})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "array_length(any[])",
			Kind:             lsp.CIKFunction,
			InsertText:       "array_length(${1:any[]})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		// String Functions
		{
			Label:            "bit_length(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "bit_length(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "char_length(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "char_length(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "character_length(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "character_length(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "length(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "length(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "lower(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "lower(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "lpad(text, int, text?)",
			Kind:             lsp.CIKFunction,
			InsertText:       "lpad(${1:text}, ${2:int}, ${3:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "ltrim(text, text?)",
			Kind:             lsp.CIKFunction,
			InsertText:       "ltrim(${1:text}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "octet_length(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "octet_length(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "overlay(text, text, int, int?)",
			Kind:             lsp.CIKFunction,
			InsertText:       "overlay(${1:text}, ${2:text}, ${3:int}, ${4:int})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "position(text, text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "position(${1:text}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "rpad(text, int, text?)",
			Kind:             lsp.CIKFunction,
			InsertText:       "rpad(${1:text}, ${2:int}, ${3:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "rtrim(text, text?)",
			Kind:             lsp.CIKFunction,
			InsertText:       "rtrim(${1:text}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "substring(text, int, int?)",
			Kind:             lsp.CIKFunction,
			InsertText:       "substring(${1:text}, ${2:int}, ${3:int})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "trim(text, text?)",
			Kind:             lsp.CIKFunction,
			InsertText:       "trim(${1:text}, ${2:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "upper(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "upper(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "format(text, ...any)",
			Kind:             lsp.CIKFunction,
			InsertText:       "format(${1:text}, ${2:...any})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		// Aggregate Functions
		{
			Label:            "count(*)",
			Kind:             lsp.CIKFunction,
			InsertText:       "count(*)",
			InsertTextFormat: lsp.ITFSnippet,
		},
		{
			Label:            "sum(int|decimal|uint256)",
			Kind:             lsp.CIKFunction,
			InsertText:       "sum(${1:int|decimal|uint256})",
			InsertTextFormat: lsp.ITFSnippet,
		},
		// Misc Functions
		{
			Label:            "error(text)",
			Kind:             lsp.CIKFunction,
			InsertText:       "error(${1:text})",
			InsertTextFormat: lsp.ITFSnippet,
		},
	}

	sqlKeywords = []string{
		"ABORT", "ADD", "ALL", "AND", "AS", "ASC", "BETWEEN", "BY",
		"CASE", "COLLATE", "COMMIT", "CONFLICT", "CREATE", "CROSS",
		"DEFAULT", "DELETE", "DESC", "DISTINCT", "ELSE",
		"END", "ESCAPE", "EXCEPT", "EXISTS", "FAIL", "FROM",
		"FULL", "GLOB", "GROUP", "HAVING", "IGNORE", "IN",
		"INDEXED", "INNER", "INSERT", "INTERSECT", "INTO",
		"IS", "ISNULL", "JOIN", "LEFT", "LIKE", "LIMIT", "MATCH",
		"NATURAL", "NOT", "NULL", "OF", "OFFSET", "ON", "OR",
		"ORDER", "OUTER", "RAISE", "REGEXP", "REPLACE", "RETURNING",
		"RIGHT", "ROLLBACK", "SELECT", "SET", "THEN", "UNION",
		"UPDATE", "USING", "VALUES", "WHEN", "WHERE", "WITH", "TRUE",
		"FALSE", "NULLS", "FIRST", "LAST", "FILTER", "GROUPS", "DO", "NOTHING",
		"abort", "add", "all", "and", "as", "asc", "between", "by", "case", "collate",
		"commit", "conflict", "create", "cross", "default", "delete", "desc",
		"distinct", "else", "end", "escape", "except", "exists", "fail", "from",
		"full", "glob", "group", "having", "ignore", "in", "indexed", "inner",
		"insert", "intersect", "into", "is", "isnull", "join", "left", "like",
		"limit", "match", "natural", "not", "null", "of", "offset", "on", "or",
		"order", "outer", "raise", "regexp", "replace", "returning", "right",
		"rollback", "select", "set", "then", "union", "update", "using", "values",
		"when", "where", "with", "true", "false", "nulls", "first", "last",
		"filter", "groups", "do", "nothing", "delete", "update",
	}

	sqlKeywordsCompletionItems = append(getDefaultCompletionItems(sqlKeywords),
		[]lsp.CompletionItem{
			{ // insert statement
				Label:            "insert into table (columns...) values (values...)",
				Kind:             lsp.CIKClass,
				InsertText:       "insert into ${1:table} (${2:columns})\n\tvalues (${3:values});",
				InsertTextFormat: lsp.ITFSnippet,
			},
			{ // update statement
				Label:            "update table set (column = value) where {condition}",
				Kind:             lsp.CIKClass,
				InsertText:       "update ${1:table} set ${2:column} = ${3:value} where ${4:condition};",
				InsertTextFormat: lsp.ITFSnippet,
			},
			{ // select statement
				Label:            "select (columns...) from table where {condition}",
				Kind:             lsp.CIKClass,
				InsertText:       "select ${1:columns} from ${2:table} where ${3:condition}",
				InsertTextFormat: lsp.ITFSnippet,
			},
			{ // delete statement
				Label:            "delete from table where {condition}",
				Kind:             lsp.CIKClass,
				InsertText:       "delete from ${1:table} where ${2:condition};",
				InsertTextFormat: lsp.ITFSnippet,
			},
			{ // on conflict clause
				Label:            "on conflict (column) do {update|nothing}",
				Kind:             lsp.CIKClass,
				InsertText:       "on conflict (${1:column}) do ${2:update|nothing};",
				InsertTextFormat: lsp.ITFSnippet,
			},
			{ // conflict clause
				Label:            "conflict (column) do {update|nothing}",
				Kind:             lsp.CIKClass,
				InsertText:       "conflict (${1:column}) do ${2:update|nothing};",
				InsertTextFormat: lsp.ITFSnippet,
			},
			{ // do update clause
				Label:            "do update set {column} = {value}",
				Kind:             lsp.CIKClass,
				InsertText:       "do update set ${1:} = ${2:};",
				InsertTextFormat: lsp.ITFSnippet,
			},
		}...,
	)

	controlFlowCompletionItems = []lsp.CompletionItem{
		{ // if statement
			Label:      "if",
			Kind:       lsp.CIKKeyword,
			InsertText: "if",
		},
		{ // else if statement
			Label:      "elseif",
			Kind:       lsp.CIKKeyword,
			InsertText: "elseif",
		},
		{ // else if statement
			Label:      "for",
			Kind:       lsp.CIKKeyword,
			InsertText: "for",
		},
		{ // break statement
			Label:      "break",
			Kind:       lsp.CIKKeyword,
			InsertText: "break",
		},
		{ // return statement
			Label:      "return",
			Kind:       lsp.CIKKeyword,
			InsertText: "return",
		},
		{ // return next statement
			Label:      "return next",
			Kind:       lsp.CIKKeyword,
			InsertText: "return next",
		},
		{ // next
			Label:      "next",
			Kind:       lsp.CIKKeyword,
			InsertText: "next",
		},
	}
	// Method completion items
	methodCompletionItems = append(append(sqlFunctionsCompletionItems, sqlKeywordsCompletionItems...), controlFlowCompletionItems...)
)

func getDefaultCompletionItems(keys []string) []lsp.CompletionItem {
	var items []lsp.CompletionItem
	for _, kw := range keys {
		items = append(items, lsp.CompletionItem{
			Label:      kw,
			Kind:       lsp.CIKKeyword,
			InsertText: kw,
		})
	}

	return items
}
