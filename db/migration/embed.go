// db/migration/embed.go
package migration

import _ "embed"

//go:embed 000_init.sql
var InitSQL string
