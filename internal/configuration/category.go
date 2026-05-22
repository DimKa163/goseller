package configuration

type CategoryConfiguration struct {
	Addr          string `env:"CATEGORY_ADDR" envDefault:":8080"`
	AppEnv        string `env:"APP_ENV" envDefault:"production"`
	Database      string `env:"CATEGORY_DATABASE,required"`
	MigrationPath string `env:"MIGRATION_PATH" envDefault:"./internal/category/migrations"`
}
