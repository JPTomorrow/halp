package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getSqlTableNameString(s interface{}) string {
	name := reflect.TypeOf(s).Name()
	return toSnakeCase(name + "s")
}

func getSqlFieldNames(s interface{}, exclude_props ...string) string {
	t := reflect.TypeOf(s)
	sb := strings.Builder{}

OUTER:
	for i := 0; i < t.NumField(); i++ {
		name, ok := t.Field(i).Tag.Lookup("sql_name")
		props, ok2 := t.Field(i).Tag.Lookup("sql_props")

		if !ok || !ok2 {
			continue
		}

		for _, val := range exclude_props {
			if strings.Contains(props, val) {
				continue OUTER
			}
		}

		sb.WriteString(name)
		if i < t.NumField()-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func getSqlFieldValues(s interface{}, exclude_props ...string) string {
	v := reflect.ValueOf(s)
	sb := strings.Builder{}

OUTER:
	for i := 0; i < v.NumField(); i++ {
		props, ok := v.Type().Field(i).Tag.Lookup("sql_props")

		if !ok {
			continue
		}

		for _, val := range exclude_props {
			if strings.Contains(props, val) {
				continue OUTER
			}
		}

		sb.WriteString(fmt.Sprintf("'%+v'", v.Field(i)))
		if i < v.NumField()-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func getSqlFieldValue(s interface{}, val string) (string, error) {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		if v.Type().Field(i).Tag.Get("sql_name") == val {
			return fmt.Sprintf("'%+v'", v.Field(i)), nil
		}
	}
	return "", fmt.Errorf("field not found")
}

func PushSchema(structs ...interface{}) (string, error) {
	sb := strings.Builder{}
	final := strings.Builder{}

	for _, s := range structs {
		t := reflect.TypeOf(s)
		fcount := t.NumField()
		struct_name := getSqlTableNameString(s)

		// create table statement
		sb.WriteString("CREATE TABLE IF NOT EXISTS " + struct_name + " (\n")

		for i := range fcount {
			v, ok := t.Field(i).Tag.Lookup("sql_name")
			if !ok {
				continue
			}
			sb.WriteString("\t" + v)
			end, ok := t.Field(i).Tag.Lookup("sql_props")
			if !ok {
				continue
			}
			sb.WriteString("\t\t\t" + end)
			if i < fcount-1 {
				end = ",\n"
				sb.WriteString(end)
			}
		}

		sb.WriteString(")\n\n")
		_, err := DbInstance.Exec(sb.String())
		if err != nil {
			return "", err
		}
		final.WriteString(sb.String())
		sb.Reset()
	}
	return final.String(), nil
}

func Insert(s interface{}) (sql.Result, error) {
	query := "INSERT INTO " + getSqlTableNameString(s) + " (" + getSqlFieldNames(s, "AUTOINCREMENT") + ") VALUES (" + getSqlFieldValues(s, "AUTOINCREMENT") + ")"
	fmt.Printf("\nInsert Query: %s\n\n\n", query)
	return DbInstance.Exec(query)
}

func Remove(s interface{}) (sql.Result, error) {
	id, err := getSqlFieldValue(s, "id")
	if err != nil {
		return nil, err
	}
	query := "DELETE FROM " + getSqlTableNameString(s) + " WHERE id = " + id
	fmt.Printf("\nRemove Query: %s\n\n\n", query)
	return DbInstance.Exec(query)
}

func Update() {

}

func Query(entry interface{}, values ...string) (*sql.Rows, error) {
	table_name := getSqlTableNameString(entry)
	query := strings.Builder{}
	query.WriteString("SELECT * FROM " + table_name + " WHERE ")
	for i, val := range values {
		first_letter := string(val[0])
		rest := val[1:]
		final := strings.ToUpper(first_letter) + rest
		query.WriteString(table_name + "." + val + " = " + fmt.Sprintf("%+v", reflect.ValueOf(entry).FieldByName(final)))
		if i < len(values)-1 {
			query.WriteString(" AND ")
		}
	}
	fmt.Printf("\nQuery: %s\n\n\n", query.String())
	return DbInstance.Query(query.String())
}
