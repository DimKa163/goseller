package configuration

type CategoryConfiguration struct {
	Addr          string `env:"ADDR" envDefault:":8080"`
	AppEnv        string `env:"APP_ENV" envDefault:"production"`
	Database      string `env:"DATABASE,required"`
	MigrationPath string `env:"MIGRATION_PATH" envDefault:"./internal/user/migrations"`
}
