package init

import (
	"app/model"
	"app/mqtt"
	"encoding/json"
	"fmt"
	"framework/database"
	utils "framework/utils/common"
	"framework/utils/cron"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	log.Println("Load environment variable and start services...")
}

// Load env variable
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

/*
Initialize database connection
*/

type DatabaseConfig struct {
	Driver      string `json:"driver"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Pass        string `json:"pass"`
	Schema      string `json:"schema"`
	Timezone    string `json:"timezone"`
	Ssl         string `json:"ssl"`
	AutoMigrate bool   `json:"auto-migrate"`
	Enable      bool   `json:"enable"`
}

type Databases struct {
	Default DatabaseConfig   `json:"default"`
	Other   []DatabaseConfig `json:"other"`
}

type MqttConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	Tls_path string `json:"tls"`
	Enable   bool   `json:"enable"`
}

type MqttStruct struct {
	Default MqttConfig   `json:"default"`
	Other   []MqttConfig `json:"other"`
}

// Load config/app.json
type configFile struct {
	Database Databases  `json:"database"`
	Mqtt     MqttStruct `json:"mqtt"`
	Cron     bool       `json:"cron"`
}

var config *configFile

const serviceJson = "./config/service.json"

func init() {
	file, err := os.Open(serviceJson)
	if err != nil {
		log.Fatalf("Fail to open config/app.json:\n %v", err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal([]byte(byteValue), &config)
}

func init() {
	// Database
	if config.Database.Default.Enable {
		err, _ := database.NewDB()
		if err != nil {
			os.Exit(1)
		}

		if config.Database.Default.AutoMigrate {
			// Migrate Model
			model.Execute()
		}
	}

	if len(config.Database.Other) > 0 {
		// Load database instance
		for _, service := range config.Database.Other {
			if service.Enable {
				err, _ := database.NewDB()
				if err != nil {
					os.Exit(1)
				}
				if service.AutoMigrate {
					// Migrate Model
					model.Execute()
				}
			}
		}
	}
}

/*
	Initialize cron service
*/
func init() {
	if config.Cron {
		// Enable cron service
		err, _ := cron.NewCron()
		if err != nil {
			log.Println("ERR_CRON_INIT", err)
		}
	}
}

/*
Initialize MQTT service
*/
func init() {
	// MQTT
	if config.Mqtt.Default.Enable {
		host := utils.GetEnv(config.Mqtt.Default.Host, "localhost")
		port := utils.GetEnv(config.Mqtt.Default.Port, "1883")
		mqtt.NewMQTTClient(fmt.Sprintf("%s:%s", host, port))
	}

	if len(config.Mqtt.Other) > 0 {
		// start other mqtt service
		for _, service := range config.Mqtt.Other {
			if service.Enable {
				host := utils.GetEnv(service.Host, "localhost")
				port := utils.GetEnv(service.Port, "1883")
				mqtt.NewMQTTClient(fmt.Sprintf("%s:%s", host, port))
			}
		}
	}

	// Other mqtt services
}

func Initialize() {
	fmt.Println("Go")
}
