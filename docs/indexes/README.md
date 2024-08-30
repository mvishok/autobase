# Indexing

Autobase supports indexing on columns in the data source. Indexing can be used to speed up the data retrieval process by creating an index on the specified columns.

Currently, Autobase supports the following types of indexes:
- `PRIMARY`

## PRIMARY

In Autobase, the `PRIMARY` index can be created on a single column in the data source. The `PRIMARY` index is used to uniquely identify each row in the data source. The `PRIMARY` index is similar to the `PRIMARY KEY` in SQL.

A `PRIMARY` index can be created on a column by adding `$` prefix to the column name in the data source.

### Syntax

```csv
$id,roll,name,age,score,height
1,12,John Doe,20,85,175
2,31,Jane Smith,22,90,165
3,18,Bob Johnson,21,78,180
```

In the above example, the `PRIMARY` index is created on the `id` column.