# Data Operations

Autobase provides a set of operations that can be performed on the data sources. The latest version of Autobase supports the following operations:

- `SELECT`
- `UPDATE`
- `WHERE`

These operations can be performed on the data sources by specifying the operation in the query string of the request.

## SELECT

The `SELECT` operation is used to retrieve specific fields from the data source. The `SELECT` operation is similar to the `SELECT` statement in SQL.

### Wildcard

The `*` wildcard can be used to select all columns from the data source.

### Syntax

```http
select=field1,field2,...
```

## UPDATE

The `UPDATE` operation is used to update the data in the data source. In the `UPDATE` operation, you can specify the fields that you want to update, followed by a `:` separator and the new value.

### Syntax

```http
update=field1:value1,field2:value2,...
```

## INSERT

The `INSERT` operation is used to insert new data into the data source. In the `INSERT` operation, you have to specify the values for all the fields in the data source, as the `INSERT` operation does not support partial insertion.

The order of the values should be the same as the order of the fields in the data source.

### Syntax

```http
insert=value1,value2,...
```

## DELETE

The `DELETE` operation is used to delete data from the data source. To delete data from the data source, you have to set the `DELETE` parameter to `true` in the query string.

### Syntax

```http
delete=true&where=field1:eq:value1,...
```

## WHERE

The `WHERE` operation is used to filter the data in the data source. The `WHERE` operation is similar to the `WHERE` clause in SQL. In the `WHERE` operation, you can specify the field, followed by a `:` separator, the  operator, and the value.

### Operators
- `eq`: Equal to
- `ne`: Not equal to
- `gt`: Greater than
- `lt`: Less than
- `ge`: Greater than or equal to
- `le`: Less than or equal to

### Syntax

```http
where=field1:eq:value1,field2:ne:value2,...
```