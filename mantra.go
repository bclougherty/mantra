// Mantra is a simple ORM framework designed to work with MySQL.
package mantra

import (
	"errors"
	"fmt"
	"github.com/octoberxp/stringmaps"
	"reflect"
)

// ModelSQL is a type that holds basic SQL statements and related data for a single type.
type ModelSQL struct {
	TableName, PrimaryKeyField string

	Create, Retrieve, Update, Delete string

	DatabaseToStructMapping, StructToDatabaseMapping map[string]string
}

// ModelSQLForObject will use reflection to generate a ModelSQL struct based on the given data type.
// By default, the generated SQL will use all struct fields as table names, converted into underscore_case.
//
// If there is a struct field named "Id", Mantra will use it as the primary key by default.
// If there is a struct field named "Deleted", Mantra will use it as a deletion flag by default, and generate "soft-delete" SQL.
// If there is no Deleted field, and no mantraDeletionFlag-tagged field, Mantra will generate "hard-delete" SQL (an actual DELETE statement).
//
// The following tags can be used to control the output:
//
//   - mantraIgnore:"true" applied to a struct field will cause that field to be completely ignored by Mantra.
//   - mantraColumn:"real_col_name" can be used to directly specify a column name for columns that don't fit Mantra's expectations.
//   - mantraPrimaryKey:"true" applied to a field will cause Mantra to treat that field as the primary key.
//   - mantraDeletionFlag:"true" applied to a field will cause Mantra to treat that field as the deletion flag.
func ModelSQLForObject(object interface{}, tableName string) (sqlObject ModelSQL, err error) {
	sqlObject = ModelSQL{TableName: tableName}

	if primaryKey := primaryKeyField(object); len(primaryKey) > 0 {
		sqlObject.PrimaryKeyField = primaryKey
	} else {
		err = errors.New(fmt.Sprintf("mantra: No primary key found for type %v!", reflect.TypeOf(object).Name))

		return sqlObject, err
	}

	sqlObject.StructToDatabaseMapping = structToDatabaseFieldMap(object)
	sqlObject.DatabaseToStructMapping = stringmaps.Reverse(sqlObject.StructToDatabaseMapping)

	sqlObject.Create = fmt.Sprintf("INSERT INTO `%v` (%v) VALUES (%v)", tableName, insertFieldList(object), valuePlaceholders(object))
	sqlObject.Update = fmt.Sprintf("UPDATE `%v` SET %v WHERE `%v` = ?", tableName, updateFieldList(object), sqlObject.PrimaryKeyField)
	sqlObject.Retrieve = fmt.Sprintf("SELECT %v FROM `%v` WHERE `%v` = ?", insertFieldList(object), tableName, sqlObject.PrimaryKeyField)

	// Look for a deletion flag
	if flag := deletionFlag(object); len(flag) > 0 {
		sqlObject.Delete = fmt.Sprintf("UPDATE `%v` SET `%v` = 0 WHERE `%v` = ?", tableName, flag, sqlObject.PrimaryKeyField)
	} else {
		sqlObject.Delete = fmt.Sprintf("DELETE FROM  `%v` WHERE `%v` = ?", tableName, sqlObject.PrimaryKeyField)
	}

	return sqlObject, nil
}
