package services

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alancuriel/game-hosting-sass/provisioner/clients"
	"github.com/alancuriel/game-hosting-sass/provisioner/db"
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
	linodeClient   *clients.Linode
	serverRootPass string
	provisionerDb  *db.ProvisionerDB
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

	db, err := db.NewProvisioner()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize provisioner db: %v", err)
	}

	return &MinecraftLinodeProvisionService{
		linodeClient: &clients.Linode{
			HttpClient: &http.Client{},
			ApiKey:     linodeApiKey,
		},
		serverRootPass: root_pass,
		provisionerDb:  db,
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

	req, err := s.genLinodeRequest(instance, region, minecraftUser)
	if err != nil {
		return "", fmt.Errorf("error generating linode req %s", err.Error())
	}

	resp, err := s.linodeClient.CreateLinode(req)

	if err != nil {
		return "", fmt.Errorf("error creating mc server %s", err.Error())
	}
	if resp.Ipv4 == nil || len(resp.Ipv4) == 0 {
		return "", fmt.Errorf("no ip found from creating linode")
	}

	now := time.Now()

	server := &m.MinecraftServer{
		IP:           resp.Ipv4[0],
		Username:     minecraftUser,
		InstanceType: instance.String(),
		Region:       region.String(),
		Label:        req.Label,
		CreatedAt:    now,
		UpdatedAt:    now,
		Status:       "active",
	}

	err = s.provisionerDb.SaveServer(server)
	if err != nil {
		return resp.Ipv4[0], fmt.Errorf("server created but failed to save to database: %v", err)
	}

	return resp.Ipv4[0], nil
}

func (s *MinecraftLinodeProvisionService) genLinodeRequest(
	instance m.MinecraftInstance,
	region m.Region,
	minecraftUser string) (*m.CreateLinodeRequest, error) {
	g := gen.CreateUserDataGenerator(gen.LINODE_UBUNTU_22_04_MINECRAFT)
	user_data, err := g.Generate(map[string]string{"${{OPUSER}}": minecraftUser})

	if err != nil {
		return nil, err
	}

	instanceType := s.mapMinecraftTypeToLinode(instance)

	if instanceType == m.LINODE_INSTANCE_INVALID {
		return nil, fmt.Errorf("Could not find  instance type from %s", instance.String())
	}

	label := "mc_" + minecraftUser + "_" + helpers.GenRandAlphaNumeric(4)

	return &m.CreateLinodeRequest{
		Image:        default_linode_image,
		Region:       region.String(),
		InstanceType: instanceType.String(),
		Label:        label,
		RootPass:     s.serverRootPass,
		FirewallId:   minecraft_firewall_id,
		Metadata: map[string]string{
			"user_data": user_data,
		},
	}, nil
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
