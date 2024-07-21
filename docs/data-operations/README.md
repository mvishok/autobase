# Data Operations

Autobase provides a set of operations that can be performed on the data sources. The latest version of Autobase supports the following operations:

- `SELECT`
- `UPDATE`
- `WHERE`

These operations can be performed on the data sources by specifying the operation in the query string of the request.

## SELECT

The `SELECT` operation is used to retrieve specific fields from the data source. The `SELECT` operation is similar to the `SELECT` statement in SQL.

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