package configuration

type UserConfiguration struct {
	Addr          string `env:"USER_ADDR" envDefault:":8080"`
	AppEnv        string `env:"APP_ENV" envDefault:"production"`
	Database      string `env:"USER_DATABASE,required"`
	MigrationPath string `env:"MIGRATION_PATH" envDefault:"./internal/user/migrations"`
}
