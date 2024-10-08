package docker

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/goravel/framework/contracts/database/orm"
	contractstesting "github.com/goravel/framework/contracts/testing"
	configmocks "github.com/goravel/framework/mocks/config"
	"github.com/goravel/framework/support/env"
)

type PostgresTestSuite struct {
	suite.Suite
	mockConfig *configmocks.Config
	postgres   *PostgresImpl
}

func TestPostgresTestSuite(t *testing.T) {
	if env.IsWindows() {
		t.Skip("Skipping tests of using docker")
	}

	suite.Run(t, new(PostgresTestSuite))
}

func (s *PostgresTestSuite) SetupTest() {
	s.mockConfig = &configmocks.Config{}
	s.postgres = NewPostgresImpl("goravel", "goravel", "goravel")
}

func (s *PostgresTestSuite) TestBuild() {
	s.Nil(s.postgres.Build())
	instance, err := s.postgres.connect()
	s.Nil(err)
	s.NotNil(instance)

	s.Equal("127.0.0.1", s.postgres.Config().Host)
	s.Equal("goravel", s.postgres.Config().Database)
	s.Equal("goravel", s.postgres.Config().Username)
	s.Equal("goravel", s.postgres.Config().Password)
	s.True(s.postgres.Config().Port > 0)

	res := instance.Exec(`
	CREATE TABLE users (
	 id SERIAL PRIMARY KEY NOT NULL,
	 name varchar(255) NOT NULL
	);
	`)
	s.Nil(res.Error)

	res = instance.Exec(`
	INSERT INTO users (name) VALUES ('goravel');
	`)
	s.Nil(res.Error)
	s.Equal(int64(1), res.RowsAffected)

	var count int64
	res = instance.Raw(`
	SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' and table_name = 'users';
	`).Scan(&count)
	s.Nil(res.Error)
	s.Equal(int64(1), count)

	s.Nil(s.postgres.Fresh())

	res = instance.Raw(`
		SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' and table_name = 'users';
		`).Scan(&count)
	s.Nil(res.Error)
	s.Equal(int64(0), count)

	s.Nil(s.postgres.Stop())
}

func (s *PostgresTestSuite) TestImage() {
	image := contractstesting.Image{
		Repository: "postgres",
	}
	s.postgres.Image(image)
	s.Equal(&image, s.postgres.image)
}

func (s *PostgresTestSuite) TestName() {
	s.Equal(orm.DriverPostgres, s.postgres.Name())
}
