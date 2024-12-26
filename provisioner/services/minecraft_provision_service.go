package services

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/alancuriel/game-hosting-sass/provisioner/clients"
	gen "github.com/alancuriel/game-hosting-sass/provisioner/generators"
	"github.com/alancuriel/game-hosting-sass/provisioner/helpers"
	m "github.com/alancuriel/game-hosting-sass/provisioner/models"
)

const (
	default_linode_image  = "linode/ubuntu22.04"
	minecraft_firewall_id = 1157654
)

type MinecraftProvisionService interface {
	Provision(instance m.MinecraftInstance, region m.Region, minecraftUser string) (string, error)
}

type MinecraftLinodeProvisionService struct {
	linodeClient   clients.Linode
	serverRootPass string
}

func NewMinecraftLinodeProvisionService() (MinecraftProvisionService, error) {
	linodeApiKey := strings.TrimSpace(os.Getenv("LINODE_API_KEY"))

	if linodeApiKey == "" {
		return nil, fmt.Errorf("LINODE_API_KEY not found")
	}

	root_pass := strings.TrimSpace(os.Getenv("MC_ROOT_PASS"))

	if root_pass == "" {
		return nil, fmt.Errorf("MC_ROOT_PASS not found")
	}

	return &MinecraftLinodeProvisionService{
		linodeClient: clients.Linode{
			HttpClient: &http.Client{},
			ApiKey:     linodeApiKey,
		},
		serverRootPass: root_pass,
	}, nil
}

func (s *MinecraftLinodeProvisionService) Provision(
	instance m.MinecraftInstance, region m.Region, minecraftUser string) (string, error) {

	if instance == m.MINECRAFT_INSTANCE_INVALID {
		return "", fmt.Errorf("Invalid minecraft instance type provided")
	}
	if region == m.INVALID {
		return "", fmt.Errorf("Invalid minecraft region type provided")
	}
	if minecraftUser == "" {
		return "", fmt.Errorf("Invalid minecraft region type provided")
	}

	g := gen.CreateUserDataGenerator(gen.LINODE_UBUNTU_22_04_MINECRAFT)
	user_data, err := g.Generate(map[string]string{"${{OPUSER}}": minecraftUser})

	if err != nil {
		return "", err
	}

	instanceType := s.mapMinecraftTypeToLinode(instance)

	if instanceType == m.LINODE_INSTANCE_INVALID {
		return "", fmt.Errorf("Could not find  instance type from %s", instance.String())
	}

	label := "mc_" + minecraftUser + "_" + helpers.GenRandAlphaNumeric(4)

	req := &m.CreateLinodeRequest{
		Image:        default_linode_image,
		Region:       region.String(),
		InstanceType: instanceType.String(),
		Label:        label,
		RootPass:     s.serverRootPass,
		FirewallId:   minecraft_firewall_id,
		Metadata: map[string]string{
			"user_data": user_data,
		},
	}

	resp, err := s.linodeClient.CreateLinode(req)

	if err != nil {
		return "", fmt.Errorf("error creating mc server %s", err.Error())
	}
	if resp.Ipv4 == nil || len(resp.Ipv4) == 0 {
		return "", fmt.Errorf("no ip found from creating linode")
	}

	return resp.Ipv4[0], nil
}

func (s *MinecraftLinodeProvisionService) mapMinecraftTypeToLinode(
	instanceType m.MinecraftInstance) m.LinodeInstance {

	switch instanceType {
	case m.MINECRAFT_INSTANCE_BASIC_1:
		return m.G6_NANODE_1
	case m.MINECRAFT_INSTANCE_STANDARD_1:
		return m.G6_STANDARD_1
	case m.MINECRAFT_INSTANCE_PREMIUM_1:
		return m.G6_STANDARD_2
	case m.MINECRAFT_INSTANCE_SUPER_1:
		return m.G6_STANDARD_4
	case m.MINECRAFT_INSTANCE_ULTIMATE_1:
		return m.G6_STANDARD_6
	default:
		return m.LINODE_INSTANCE_INVALID
	}
}
