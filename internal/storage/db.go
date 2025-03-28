package storage

type Database struct {
	Dialect    string           `yaml:"dialect"`
	PostgreSQL PostgreSQLConfig `yaml:"postgresql"`
	DynamoDB   DynamoDBConfig   `yaml:"dynamodb"`
}

type PostgreSQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslmode"`
}

type DynamoDBConfig struct {
	Region          string `yaml:"region"`
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
}
