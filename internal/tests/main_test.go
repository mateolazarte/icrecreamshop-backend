package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"icecreamshop/internal/api"
	"icecreamshop/internal/storage"
	"icecreamshop/internal/types"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	envPath := "../../.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env:", err)
	}
	testMode := os.Getenv("TEST_MODE")
	apiEnv := os.Getenv("API_ENV")

	if apiEnv != "testing" {
		panic("env var API_ENV must be set to testing")
	}

	code := m.Run()
	if testMode == "integration" {
		println("Integration tests have finished")
	} else {
		println("Mocking tests have finished")
	}
	os.Exit(code)

}

func newStorage(flavors []types.Flavor, users []types.User, prices map[uint]uint) storage.Storage {
	testMode := os.Getenv("TEST_MODE")
	if testMode == "integration" {
		return storage.NewDBStorage(flavors, users, prices)
	} else {
		return storage.NewMemoryStorage(flavors, users, prices)
	}
}

var sv *api.Server
var router *gin.Engine

func setup() {
	sv = api.NewServer(newStorage(flavors, users, prices))
	router = sv.SetupRouter()
}

func clearAndCloseConnection(t *testing.T, store storage.Storage) {
	err := store.CleanDB()
	if err != nil {
		t.Error(err)
	}
	err = store.Close()
	if err != nil {
		t.Error(err)
	}
}
