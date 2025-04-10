package env

type AppEnv struct {
	ENV      string `mapstructure:"env"`
	VERSION  string `mapstructure:"app_version"`
	BASEURL  string `mapstructure:"base_url"`
	APP_PORT string `mapstructure:"app_port"`
	APP_NAME string `mapstructure:"app_name"`
}

type Redis struct {
	RED_CONNSTR string `mapstructure:"red_connstr"`
}
type Database struct {
	DBNAME     string `mapstructure:"db_name"`
	DBPASSWORD string `mapstructure:"db_password"`
	DBUSER     string `mapstructure:"db_user"`
	DBHOST     string `mapstructure:"db_host"`
	DBPORT     string `mapstructure:"db_port"`
}

type Config struct {
	Database `mapstructure:"database,squash"`
	Redis    `mapstructure:"redis,squash"`
	AppEnv   `mapstructure:"app_env,squash"`
}
