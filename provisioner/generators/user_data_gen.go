package generators

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"strings"
)

type UserDataTemplateType uint16
const (
	LINODE_UBUNTU_22_04_MINECRAFT = iota
)

//go:embed minecraft_ubuntu_linode_template.sh
var linodeUbuntuMinecraftInit string


type UserDataGenerator struct {
	templateString string
}

func CreateUserDataGenerator(templateType UserDataTemplateType) *UserDataGenerator {
	switch templateType {
		case LINODE_UBUNTU_22_04_MINECRAFT:
			return &UserDataGenerator{linodeUbuntuMinecraftInit}
		default:
			return nil
	}
}

func (g *UserDataGenerator) Generate(params map[string]string) (string, error) {
	if g == nil || g.templateString == ""  {
		return "", fmt.Errorf("No generator or template found")
	}

	str := g.templateString
	for paramKey, replaceValue := range params {
		str = strings.ReplaceAll(str, paramKey, replaceValue)
	}

	return base64.StdEncoding.EncodeToString([]byte(str)), nil
}
