package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// name of json config file
const filename = "config.json"

// name of the .env for development/test environment configs
const envFile = ".env"

// environment variable indicating whether it's running in dev/test/prod mode
const modeKey = "CAMERAROLL_MODE"

const (
	DevMode  = "dev"
	TestMode = "test"
	ProdMode = "prod"
)

const (
	addressKey   = "DB_ADDR"
	userKey      = "DB_USER"
	passwordKey  = "DB_PASS"
	nameKey      = "DB_NAME"
	jwtSecretKey = "JWT_SECRET"
)

// DatabaseSettings contains the configs of the MySQL database that the server connects to
type DatabaseSettings struct {
	Address  string `json:"address"` // database address in ip:port format
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// GoogleOAuthSettings contains information needed for Google OAuth2
type GoogleOAuthSettings struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Config contains all the configs this server requires
type Config struct {
	Mode        string
	Port        uint                 `json:"port"`
	RootURL     string               `json:"root_url"`
	JWTSecret   string               `json:"jwt_secret"`
	AdminID     string               `json:"admin_account"`
	GoogleOAuth *GoogleOAuthSettings `json:"google_oauth"`
	Database    *DatabaseSettings    `json:"database"`
}

func (config *Config) loadFromFile() {
	// find the path to current executable
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	// open config.json file
	file, err := os.Open(filepath.Join(exPath, filename))
	if err != nil {
		log.Printf("Failed to open config file [%s]. Error[%v]\n", filename, err)
		panic(err)
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Printf("Failed to parse config file [%s]. Error[%v]\n", filename, err)
		panic(err)
	}

	// find the path to .env file
	envPath := filepath.Join(exPath, envFile)

	// load .env file
	godotenv.Load(envPath)
}

func (config *Config) loadFromEnv(suffix string) {
	address := os.Getenv(addressKey + suffix)
	if len(address) > 0 {
		config.Database.Address = address
	}

	user := os.Getenv(userKey + suffix)
	if len(user) > 0 {
		config.Database.User = user
	}

	password := os.Getenv(passwordKey + suffix)
	if len(password) > 0 {
		config.Database.Password = password
	}

	name := os.Getenv(nameKey + suffix)
	if len(name) > 0 {
		config.Database.Name = name
	}

	jwtSecret := os.Getenv(jwtSecretKey + suffix)
	if len(jwtSecret) > 0 {
		config.JWTSecret = jwtSecret
	}
}

func NewConfig() *Config {
	// create a new siteOptions object
	config := Config{}

	// read config.json first
	config.loadFromFile()

	// check if it's dev or test mode
	mode := os.Getenv(modeKey)
	suffix := ""

	if mode == DevMode {
		suffix = "_DEV"
	} else if mode == TestMode {
		suffix = "_TEST"
	} else {
		// default to ProdMode
		mode = ProdMode
	}

	// override configs with values from environment variables
	config.loadFromEnv(suffix)

	config.Mode = mode

	return &config
}
