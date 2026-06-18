package detectors

import (
	"strings"

	"github.com/RevIneX/context/internal/data"
)

var databaseFileNames = map[string]string{
	"mysql.cnf":       "MySQL",
	"mongod.conf":     "MongoDB",
	"cassandra.yaml":  "Cassandra",
	"clickhouse.xml":  "ClickHouse",
	"redis.conf":      "Redis",
	"memcached.conf":  "Memcached",
	"elasticsearch.yml": "Elasticsearch",
	"neo4j.conf":      "Neo4j",
}

var databaseSignatures = map[string]string{
	"postgres":     "PostgreSQL",
	"postgresql":   "PostgreSQL",
	"mysql":        "MySQL",
	"mariadb":      "MariaDB",
	"oracle":       "Oracle",
	"sqlite":       "SQLite",
	"mssql":        "SQL Server",
	"sqlserver":    "SQL Server",
	"cassandra":    "Cassandra",
	"clickhouse":   "ClickHouse",
	"mongodb":      "MongoDB",
	"mongo":        "MongoDB",
	"couchdb":      "CouchDB",
	"couchbase":    "Couchbase",
	"neo4j":        "Neo4j",
	"dynamodb":     "DynamoDB",
	"dynamo":       "DynamoDB",
	"firebase":     "Firebase",
	"firestore":    "Firestore",
	"supabase":     "Supabase",
	"prisma":       "Prisma",
	"liquibase":    "Liquibase",
	"flyway":       "Flyway",
}

type DatabaseDetector struct{}

func (d *DatabaseDetector) Name() string {
	return "database"
}

func (d *DatabaseDetector) Detect(filePath string) []data.Finding {
	lastSlash := strings.LastIndexByte(filePath, '/')
	if lastSlash == -1 {
		lastSlash = strings.LastIndexByte(filePath, '\\')
	}
	fileName := filePath[lastSlash+1:]
	lowName := strings.ToLower(fileName)
	if tech, ok := databaseFileNames[lowName]; ok {
		return []data.Finding{
			data.NewFinding(data.CatDatabase, tech, "", filePath, 1, nil),
		}
	}

	extSep := strings.LastIndexByte(lowName, '.')
	var nameWithoutExt string
	if extSep == -1 {
		nameWithoutExt = lowName
	} else {
		nameWithoutExt = lowName[:extSep]
		if lowName[extSep:] == ".sql" {
			if tech, ok := databaseSignatures[nameWithoutExt]; ok {
				return []data.Finding{
					data.NewFinding(data.CatDatabase, tech, "", filePath, 1, nil),
				}
			}
			// Обычный SQL-файл
			return []data.Finding{
				data.NewFinding(data.CatDatabase, "SQL", "", filePath, 1, nil),
			}
		}
	}

	if tech, ok := databaseSignatures[nameWithoutExt]; ok {
		return []data.Finding{
			data.NewFinding(data.CatDatabase, tech, "", filePath, 1, nil),
		}
	}

	return []data.Finding{}
}