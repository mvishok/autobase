# Configuration

The configuration file is a JSON file that contains the details of your Autobase project. The configuration file is mandatory to start the Autobase server. The configuration file contains the following fields:

# port

The `port` on which the Autobase server will run. The default port is `5000`.

```json
{
  "port": 5000
}
```

# dir

The `dir` property specifies the directory where the Autobase server will look for the data sources. This property is mandatory.

```json
{
  "dir": "data"
}
```

>[!NOTE]
> The `dir` property is relative to the path of the configuration file.

# env

The `env` property specifies the path to the environment file. The environment file contains the environment variables that will be used by the Autobase server. This property is optional.

```json
{
  "env": ".env"
}
```

# Authentication

Autobase supports basic authentication. You can specify the authentication token in your system environment variables. The environment variable should be named `AB_AUTH`. The `AB_AUTH` environment variable should contain A JSON object with `key` and `access_level` pairs.

### Example
```env
AB_AUTH={"MyKEY1":"read","MyTOKEN2":"write"}
```

If authentication is enabled, the client must send the `Authorization` header with the value `Bearer <key>` in the request. The key should be one of the keys specified in the `AB_AUTH` environment variable, and the requested operation should be allowed for that key.

### Access Levels

- `read`: Allows the client to read data from the server.
- `write`: Allows the client to write data to the server.

>[!ATTENTION]
> If the `AB_AUTH` environment variable is not set, the server will not require authentication.