package main

const (
	scansText = `{{define "scans"}}// DON'T EDIT *** generated by scaneo *** DON'T EDIT //

package {{.PackageName}}

import "database/sql"

{{range .Tokens}}func {{$.Visibility}}can{{title .Name}}(r *sql.Row) ({{.Name}}, error) {
	var s {{.Name}}
	if err := r.Scan({{range .Fields}}
		&s.{{.Name}},{{end}}
	); err != nil {
		return {{.Name}}{}, err
	}
	return s, nil
}

type {{.Name}}s map[int64]{{.Name}}

func (m {{.Name}}s) Ids() pq.Int64Array {
	out := make(pq.Int64Array, 0, len(m))
	for i := range m {
		out = append(out, i)
	}
	return out
}

{{ if hasid .Fields }}
func {{$.Visibility}}can{{title .Name}}s(rs *sql.Rows) ({{.Name}}s, error) {
	structs := make({{.Name}}s,  16)
	var err error
	for rs.Next() {
		var s {{.Name}}
		if err = rs.Scan({{range .Fields}}
			&s.{{.Name}},{{end}}
		); err != nil {
			return nil, err
		}
		structs[s.Id] = s
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}
{{ else }}
func {{$.Visibility}}can{{title .Name}}s(rs *sql.Rows) ([]{{.Name}}, error) {
	structs := make([]{{.Name}}, 0, 16)
	var err error
	for rs.Next() {
		var s {{.Name}}
		if err = rs.Scan({{range .Fields}}
			&s.{{.Name}},{{end}}
		); err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}
{{ end }}

// Insert {{title .Name}} in the database and returns the item with id filled.
func (item {{title .Name}}) Insert(tx *sql.Tx) (out {{.Name}}, err error) {
	r := tx.QueryRow(` + "`" + `INSERT INTO {{snake .Name}}s (
		{{range $i, $e := noid .Fields}}{{if $i}},{{end}}{{snake $e.Name}}{{end}}
		) VALUES (
		{{range $i, $e := noid .Fields}}{{if $i}},{{end}}${{inc $i}}{{end}}
		) RETURNING 
		{{range $i, $e := .Fields}}{{if $i}},{{end}}{{snake $e.Name}}{{end}};
		` + "`" + `{{range noid .Fields}},item.{{.Name}}{{end}})
	return {{$.Visibility}}can{{title .Name}}(r)
}

{{ if hasid .Fields }}
// Update {{title .Name}} in the database and returns the new version.
func (item {{title .Name}}) Update(tx *sql.Tx) (out {{.Name}}, err error) {
	r := tx.QueryRow(` + "`" + `UPDATE {{snake .Name}}s SET (
		{{range $i, $e := noid .Fields}}{{if $i}},{{end}}{{snake $e.Name}}{{end}}
		) = (
		{{range $i, $e := noid .Fields}}{{if $i}},{{end}}${{inc (inc $i)}}{{end}}
		) WHERE id = $1 RETURNING 
		{{range $i, $e := .Fields}}{{if $i}},{{end}}{{snake $e.Name}}{{end}};
		` + "`" + `{{range .Fields}},item.{{.Name}}{{end}})
	return {{$.Visibility}}can{{title .Name}}(r)
}

// Delete {{title .Name}} in the database and the return the id.
// Only the field 'Id' is used.
func (item {{title .Name}}) Delete(tx *sql.Tx) (int64, error) {
	var deleted_id int64
	r := tx.QueryRow("DELETE FROM {{snake .Name}}s WHERE id = $1 RETURNING id;", item.Id)
	err := r.Scan(&deleted_id)
	return deleted_id, err
}
{{ end }}

{{end}}
{{end}}`

	scansTextTest = `{{define "scansTest"}}// DON'T EDIT *** generated by scaneo *** DON'T EDIT //

package {{.PackageName}}

import (
	"database/sql"
	"math/rand"
)

{{range .Tokens}}
func rand{{title .Name}}() {{.Name}} {
	return {{.Name}}{ {{range .Fields}}
		{{.Name}}: {{rand .Type}},{{end}}
	}
}

func queries{{.Name}}(tx *sql.Tx, item {{.Name}}) ({{.Name}}, error) {
	item, err := item.Insert(tx)
	{{ if hasid .Fields }}
	if err != nil {
		return item, err
	}
	return item.Update(tx)
	{{ else }} return item, err {{ end }}
}

{{end}}
{{end}}`
)
