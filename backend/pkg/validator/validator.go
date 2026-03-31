package validator

import (
	"os"
	"sync"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type ValidationSchema struct {
	Rules map[string]map[string]string `yaml:"rules"`
}

var (
	validate *validator.Validate
	schemas  map[string]ValidationSchema
	mu       sync.RWMutex
)

func Init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	if validate == nil {
		Init()
	}
	return validate.Struct(s)
}

func RegisterTranslator(fn validator.RegisterTranslationsFunc, t validator.Translation) error {
	if validate == nil {
		Init()
	}
	return validate.RegisterTranslation("en", fn, t)
}

func LoadSchemas(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var s map[string]ValidationSchema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return err
	}

	mu.Lock()
	schemas = s
	mu.Unlock()

	return nil
}

func GetSchema(name string) ValidationSchema {
	mu.RLock()
	defer mu.RUnlock()
	return schemas[name]
}
