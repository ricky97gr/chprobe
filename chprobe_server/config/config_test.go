package conf

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	configPath = "config.yaml"
	fmt.Println(GetConfig())
}
