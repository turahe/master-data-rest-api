# Database Logging

This application includes comprehensive database logging capabilities using pgx's tracelog functionality integrated with our custom logger.

## Features

### 1. **Query Logging**
- Logs all SQL queries with execution time
- Includes query parameters in debug mode
- Highlights slow queries with different log levels

### 2. **Connection Logging**
- Logs database connection events
- Tracks connection pool metrics
- Monitors connection health

### 3. **Configurable Log Levels**
- `trace` - Most verbose, logs everything
- `debug` - Includes query details and parameters
- `info` - Standard operations and slow queries
- `warn` - Warnings and slow queries only
- `error` - Errors only

## Configuration

Configure database logging through environment variables:

```bash
# Database Logging Configuration
DB_LOG_LEVEL=debug          # trace, debug, info, warn, error
DB_LOG_QUERIES=true         # Enable/disable query logging
DB_LOG_SLOW_QUERY=100ms     # Threshold for slow query warnings
```

## Log Examples

### Connection Events
```
time="2025-08-09 15:56:17" level=info msg=Connect component=database database=master_data fields.time=30.1085ms host=localhost pid=108 port=5432
```

### Query Execution (Debug Mode)
```
time="2025-08-09 15:56:20" level=debug msg="Query executed" component=database query="SELECT id, name FROM tm_banks WHERE code = $1" args_count=1 duration_ms=2 duration_str="2.145ms"
```

### Slow Query Detection
```
time="2025-08-09 15:56:22" level=warn msg="Slow database query detected" component=database query="SELECT * FROM tm_geodirectories ORDER BY name" duration_ms=156 duration_str="156.234ms"
```

## Performance Monitoring

The logging system automatically categorizes queries by execution time:

- **< 100ms**: Debug level logging
- **100ms - 1s**: Info level logging  
- **> 1s**: Warning level (slow query detected)

## Log Levels

### Trace Level
- All connection events
- All queries with full details
- Connection pool events
- Most verbose output

### Debug Level
- Query execution details
- Query parameters
- Connection events
- Performance metrics

### Info Level
- Successful operations
- Slow query warnings
- Connection status
- Standard operational logs

### Warn Level
- Slow queries only
- Connection issues
- Performance warnings

### Error Level
- Failed queries only
- Connection errors
- Critical issues

## Integration

The database logger integrates seamlessly with the application's main logger:

```go
// Logger is automatically configured when creating connections
dbConnection := database.NewPgxConnectionWithLogger(config.Database, log)
```

## Best Practices

1. **Development**: Use `debug` level for detailed troubleshooting
2. **Staging**: Use `info` level for performance monitoring
3. **Production**: Use `warn` or `error` level to minimize overhead
4. **Set appropriate slow query thresholds** based on your performance requirements
5. **Monitor logs for patterns** in slow queries or connection issues

## Disabling Logging

To disable database logging entirely:

```bash
DB_LOG_QUERIES=false
```

Or set the log level to a higher threshold:

```bash
DB_LOG_LEVEL=error  # Only log errors
```
