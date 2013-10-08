mantra
======

A (very) simple library for persisting Go objects to MySQL with as little boilerplate as possible


PACKAGE DOCUMENTATION

    import "github.com/octoberxp/mantra"

Mantra is a simple ORM framework designed to work with MySQL.


TYPES

  type ModelSQL struct {
    TableName, PrimaryKeyField string

    Create, Retrieve, Update, Delete string

    DatabaseToStructMapping, StructToDatabaseMapping map[string]string
  }

    ModelSQL is a type that holds basic SQL statements and related data for
    a single type.


  func ModelSQLForObject(object interface{}, tableName string) (sqlObject ModelSQL, err error)

ModelSQLForObject will use reflection to generate a ModelSQL struct
based on the given data type. By default, the generated SQL will use all
struct fields as table names, converted into underscore_case.

If there is a struct field named "Id", Mantra will use it as the primary
key by default. If there is a struct field named "Deleted", Mantra will
use it as a deletion flag by default, and generate "soft-delete" SQL. If
there is no Deleted field, and no mantraDeletionFlag-tagged field,
Mantra will generate "hard-delete" SQL (an actual DELETE statement).

The following tags can be used to control the output:

- mantraIgnore:"true" applied to a struct field will cause that field to be completely ignored by Mantra.
- mantraColumn:"real_col_name" can be used to directly specify a column name for columns that don't fit Mantra's expectations.
- mantraPrimaryKey:"true" applied to a field will cause Mantra to treat that field as the primary key.
- mantraDeletionFlag:"true" applied to a field will cause Mantra to treat that field as the deletion flag.
