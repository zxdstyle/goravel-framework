package schema

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	ormcontract "github.com/goravel/framework/contracts/database/orm"
	schemacontract "github.com/goravel/framework/contracts/database/schema"
	testingcontract "github.com/goravel/framework/contracts/testing"
	"github.com/goravel/framework/database/gorm"
	configmock "github.com/goravel/framework/mocks/config"
	ormmock "github.com/goravel/framework/mocks/database/orm"
	logmock "github.com/goravel/framework/mocks/log"
	"github.com/goravel/framework/support/env"
)

type TestSchema struct {
	dbConfig   testingcontract.DatabaseConfig
	driver     ormcontract.Driver
	mockConfig *configmock.Config
	mockOrm    *ormmock.Orm
	mockLog    *logmock.Log
	query      ormcontract.Query
	schema     *Schema
}

type SchemaSuite struct {
	suite.Suite
	schemas []TestSchema
	//mysqlQuery      ormcontract.Query
	postgresQuery ormcontract.Query
	//sqliteQuery     ormcontract.Query
	//sqlserverDB     ormcontract.Query
}

func TestSchemaSuite(t *testing.T) {
	if env.IsWindows() {
		t.Skip("Skipping tests of using docker")
	}

	if err := testDatabaseDocker.Fresh(); err != nil {
		t.Fatal(err)
	}

	//mysqlDocker := gorm.NewMysqlDocker(testDatabaseDocker)
	//mysqlQuery, err := mysqlDocker.New()
	//if err != nil {
	//	log.Fatalf("Init mysql docker error: %v", err)
	//}

	postgresqlDocker := gorm.NewPostgresqlDocker(testDatabaseDocker)
	postgresqlQuery, err := postgresqlDocker.New()
	if err != nil {
		log.Fatalf("Init postgresql docker error: %v", err)
	}

	//sqliteDocker := gorm.NewSqliteDocker("goravel")
	//sqliteQuery, err := sqliteDocker.New()
	//if err != nil {
	//	log.Fatalf("Get sqlite error: %s", err)
	//}
	//
	//sqlserverDocker := gorm.NewSqlserverDocker(testDatabaseDocker)
	//sqlserverQuery, err := sqlserverDocker.New()
	//if err != nil {
	//	log.Fatalf("Init sqlserver docker error: %v", err)
	//}

	suite.Run(t, &SchemaSuite{
		//mysqlQuery:      mysqlQuery,
		postgresQuery: postgresqlQuery,
		//sqliteQuery:     sqliteQuery,
		//sqlserverDB:     sqlserverQuery,
	})

	//assert.Nil(t, file.Remove("goravel"))
}

func (s *SchemaSuite) SetupTest() {
	mockConfig := &configmock.Config{}
	mockOrm := &ormmock.Orm{}
	//mockOrmOfConnection := &ormmock.Orm{}
	mockLog := &logmock.Log{}
	mockConfig.On("GetString", "database.default").Return("mysql").Once()
	mockOrm.On("Connection", "mysql").Return(mockOrm).Once()
	mockOrm.On("Query").Return(s.postgresQuery).Once()
	postgresSchema, err := NewSchema("", mockConfig, mockOrm, mockLog)
	s.Nil(err)
	s.schemas = append(s.schemas, TestSchema{
		dbConfig:   testDatabaseDocker.Postgres.Config(),
		driver:     ormcontract.DriverPostgres,
		mockConfig: mockConfig,
		mockOrm:    mockOrm,
		mockLog:    mockLog,
		query:      s.postgresQuery,
		schema:     postgresSchema,
	})
}

func (s *SchemaSuite) TestConnection() {
	for _, schema := range s.schemas {
		schema.mockOrm.On("Connection", "postgres").Return(schema.mockOrm).Once()
		schema.mockOrm.On("Query").Return(schema.query).Once()
		s.NotNil(schema.schema.Connection("postgres"))

		schema.mockOrm.AssertExpectations(s.T())
	}
}

func (s *SchemaSuite) TestCreate() {
	for _, schema := range s.schemas {
		s.Run(schema.driver.String(), func() {
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.prefix", schema.schema.connection)).
				Return("goravel_").Once()

			err := schema.schema.Create("create", func(table schemacontract.Blueprint) {
				table.String("name")
			})
			s.Nil(err)
			s.True(schema.schema.HasTable("goravel_create"))
		})
	}
}

func (s *SchemaSuite) TestGetColumns() {
	for _, schema := range s.schemas {
		s.Run(schema.driver.String(), func() {
			table := "get_columns"
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.database", schema.schema.connection)).
				Return(schema.dbConfig.Database).Once()
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.schema", schema.schema.connection)).
				Return("").Once()
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.prefix", schema.schema.connection)).
				Return("").Twice()

			err := schema.schema.Create(table, func(table schemacontract.Blueprint) {
				table.Char("char")
				table.String("string")
			})

			s.Nil(err)
			s.True(schema.schema.HasTable(table))

			columns, err := schema.schema.GetColumns(table)

			s.Nil(err)
			s.Equal(2, len(columns))
			for _, column := range columns {
				if column.Name == "char" {
					s.False(column.AutoIncrement)
					s.Empty(column.Collation)
					s.Empty(column.Comment)
					s.Nil(column.Default)
					s.True(column.Nullable)
					s.Equal("character(255)", column.Type)
					s.Equal("bpchar", column.TypeName)
				}
				if column.Name == "string" {
					s.False(column.AutoIncrement)
					s.Empty(column.Collation)
					s.Empty(column.Comment)
					s.Nil(column.Default)
					s.True(column.Nullable)
					s.Equal("character varying(255)", column.Type)
					s.Equal("varchar", column.TypeName)
				}
			}

			columnListing := schema.schema.GetColumnListing(table)

			s.Equal(2, len(columnListing))
			s.Contains(columnListing, "char")
			s.Contains(columnListing, "string")

			s.True(schema.schema.HasColumn(table, "char"))
			s.True(schema.schema.HasColumns(table, []string{"char", "string"}))

			schema.mockConfig.AssertExpectations(s.T())
		})
	}
}

func (s *SchemaSuite) TestGetTables() {
	for _, schema := range s.schemas {
		s.Run(schema.driver.String(), func() {
			tables, err := schema.schema.GetTables()
			s.Greater(len(tables), 0)
			s.Nil(err)
		})
	}
}

func (s *SchemaSuite) TestHasTable() {
	for _, schema := range s.schemas {
		s.Run(schema.driver.String(), func() {
			s.True(schema.schema.HasTable("users"))
			s.False(schema.schema.HasTable("unknow"))
		})
	}
}

func (s *SchemaSuite) TestInitGrammarAndProcess() {
	for _, schema := range s.schemas {
		s.Run(schema.driver.String(), func() {
			s.Nil(schema.schema.initGrammarAndProcess())
			grammarType := reflect.TypeOf(schema.schema.grammar)
			grammarName := grammarType.Elem().Name()
			processorType := reflect.TypeOf(schema.schema.processor)
			processorName := processorType.Elem().Name()

			switch schema.driver {
			case ormcontract.DriverMysql:
				s.Equal("Mysql", grammarName)
				s.Equal("Mysql", processorName)
			case ormcontract.DriverPostgres:
				s.Equal("Postgres", grammarName)
				s.Equal("Postgres", processorName)
			case ormcontract.DriverSqlserver:
				s.Equal("Sqlserver", grammarName)
				s.Equal("Sqlserver", processorName)
			case ormcontract.DriverSqlite:
				s.Equal("Sqlite", grammarName)
				s.Equal("Sqlite", processorName)
			default:
				s.Fail("unsupported database driver")
			}
		})
	}
}

func (s *SchemaSuite) TestParseDatabaseAndSchemaAndTable() {
	for _, schema := range s.schemas {
		s.Run(schema.driver.String(), func() {
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.database", schema.schema.connection)).
				Return(schema.dbConfig.Database).Once()
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.schema", schema.schema.connection)).
				Return("").Once()
			database, schemaName, table := schema.schema.parseDatabaseAndSchemaAndTable("users")
			s.Equal(schema.dbConfig.Database, database)
			s.Equal("public", schemaName)
			s.Equal("users", table)

			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.database", schema.schema.connection)).
				Return(schema.dbConfig.Database).Once()
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.schema", schema.schema.connection)).
				Return("").Once()
			database, schemaName, table = schema.schema.parseDatabaseAndSchemaAndTable("goravel.users")
			s.Equal(schema.dbConfig.Database, database)
			s.Equal("goravel", schemaName)
			s.Equal("users", table)

			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.database", schema.schema.connection)).
				Return(schema.dbConfig.Database).Once()
			schema.mockConfig.On("GetString", fmt.Sprintf("database.connections.%s.schema", schema.schema.connection)).
				Return("hello").Once()
			database, schemaName, table = schema.schema.parseDatabaseAndSchemaAndTable("goravel.users")
			s.Equal(schema.dbConfig.Database, database)
			s.Equal("goravel", schemaName)
			s.Equal("users", table)

			schema.mockConfig.AssertExpectations(s.T())
		})
	}
}
