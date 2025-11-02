mkdir -p apps/oss/cmd/server
cat > apps/oss/cmd/server/main.go << 'EOF'
package main

import (
	"log"
	"neftac/storage/internal/api"
	"neftac/storage/internal/db"
	"neftac/storage/internal/lifecycle"
	"neftac/storage/internal/policy"
	"os"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    struct{ Host, Port string } `yaml:"server"`
	Database  struct{ File string }       `yaml:"database"`
	Storage   struct{ Root string }      `yaml:"storage"`
	JWT       struct{ Secret string; ExpiryHours int } `yaml:"jwt"`
	S3        struct{ Region, AccessKey, SecretKey string } `yaml:"s3"`
	Lifecycle struct{ ExpireDays int }  `yaml:"lifecycle"`
}

var AppConfig Config

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil { log.Fatal("Config not found:", err) }
	if err := yaml.Unmarshal(data, &AppConfig); err != nil { log.Fatal("Invalid config:", err) }

	db.InitDB(AppConfig.Database.File)
	os.MkdirAll(AppConfig.Storage.Root, 0755)
	policy.Init()

	go lifecycle.StartCleanup(AppConfig.Lifecycle.ExpireDays)

	app := fiber.New(fiber.Config{BodyLimit: 500 * 1024 * 1024})
	api.SetupRoutes(app)
	api.SetupS3Routes(app)

	log.Println("Neftac S3 LIVE at s3.airsoko.com")
	log.Fatal(app.Listen(AppConfig.Server.Host + ":" + AppConfig.Server.Port))
}
EOF
