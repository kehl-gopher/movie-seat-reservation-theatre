package env

type AppEnv struct {
	ENV        string `mapstructure:"env"`
	VERSION    string `mapstructure:"app_version"`
	BASEURL    string `mapstructure:"base_url"`
	APP_PORT   string `mapstructure:"app_port"`
	APP_HOST   string `mapstructure:"app_host"`
	APP_NAME   string `mapstructure:"app_name"`
	SECRET_KEY string `mapstructure:"jwt_secret_key"`
	EXPIRES_IN string `mapstructure:"expiration_time_min"`
}

type Redis struct {
	RED_HOST     string `mapstructure:"red_host"`
	RED_PORT     string `mapstructure:"red_port"`
	RED_PASSWORD string `mapstructure:"red_password"`
	RED_DB       string `mapstructure:"red_db"`
	RED_URL      string `mapstructure:"red_url"`
	RED_SSLMODE  string `mapstructure:"red_sslmode"`
}

type Minio struct {
	MINIO_ENDPOINT  string `mapstructure:"minio_endpoint"`
	MINIO_ACCESSKEY string `mapstructure:"minio_access_key"`
	MINIO_BUCKET    string `mapstructure:"minio_bucket"`
	MINIO_SECRETKEY string `mapstructure:"minio_secret_key"`
	MINIO_USESSL    string `mapstructure:"minio_usessl"`
	MINIO_REGION    string `mapstructure:"minio_region"`
}
type Database struct {
	DBNAME      string `mapstructure:"db_name"`
	DBPASSWORD  string `mapstructure:"db_password"`
	DBUSER      string `mapstructure:"db_user"`
	DBHOST      string `mapstructure:"db_host"`
	DBPORT      string `mapstructure:"db_port"`
	DB_SSLMODE  string `mapstructure:"sslmode"`
	DB_TIMEZONE string `mapstructure:"db_timezone"`
}

type Config struct {
	Database `mapstructure:"database,squash"`
	Redis    `mapstructure:"redis,squash"`
	AppEnv   `mapstructure:"app_env,squash"`
	Minio    `mapstructure:"minio,squash"`
}
