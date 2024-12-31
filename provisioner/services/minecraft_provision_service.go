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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/ssh"
)

const (
	default_linode_image  = "linode/ubuntu22.04"
	minecraft_firewall_id = 1157654
)

type MinecraftProvisionService interface {
	Provision(provisionRequest *m.ProvisionMcServerRequest) (string, error)
	ListServersByOwner(owner string) ([]*m.MinecraftServer, error)
	DeleteServer(id string) error
	AnnounceMessage(id string, message string) error
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

func (s *MinecraftLinodeProvisionService) Provision(req *m.ProvisionMcServerRequest) (string, error) {

	if req.Instance == m.MINECRAFT_INSTANCE_INVALID {
		return "", fmt.Errorf("Invalid minecraft instance type provided")
	}
	if req.Region == m.INVALID {
		return "", fmt.Errorf("Invalid minecraft region type provided")
	}
	if req.Username == "" {
		return "", fmt.Errorf("Invalid minecraft region type provided")
	}
	if req.Owner == "" {
		return "", fmt.Errorf("Invalid server owner provided")
	}

	linodeReq, err := s.genLinodeRequest(req.Instance, req.Region, req.Username)
	if err != nil {
		return "", fmt.Errorf("error generating linode req %s", err.Error())
	}

	resp, err := s.linodeClient.CreateLinode(linodeReq)

	if err != nil {
		return "", fmt.Errorf("error creating mc server %s", err.Error())
	}
	if resp.Ipv4 == nil || len(resp.Ipv4) == 0 {
		return "", fmt.Errorf("no ip found from creating linode")
	}

	now := time.Now()

	server := &m.MinecraftServer{
		ID:           primitive.NewObjectID(),
		IP:           resp.Ipv4[0],
		LinodeId:     resp.Id,
		Username:     req.Username,
		InstanceType: req.Instance.String(),
		Region:       req.Region.String(),
		Label:        linodeReq.Label,
		Owner:        req.Owner,
		CreatedAt:    now,
		UpdatedAt:    now,
		Status:       "pending",
	}

	_, err = s.provisionerDb.SaveServer(server)
	if err != nil {
		return resp.Ipv4[0], fmt.Errorf("server created but failed to save to database: %v", err)
	}

	go helpers.PingMcServer(resp.Ipv4[0], func() {
		s.provisionerDb.UpdateServerStatus(server.ID, "running")
	})

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

func (s *MinecraftLinodeProvisionService) ListServersByOwner(owner string) ([]*m.MinecraftServer, error) {
	if owner == "" {
		return nil, fmt.Errorf("no owner provided")
	}

	servers, err := s.provisionerDb.ListMcServerByOwner(owner)

	if err != nil {
		return nil, err
	}

	if servers == nil || len(servers) == 0 {
		return []*m.MinecraftServer{}, nil
	}

	return servers, nil
}

func (s *MinecraftLinodeProvisionService) DeleteServer(id string) error {
	if id == "" {
		return fmt.Errorf("no server id provided")
	}

	server, err := s.provisionerDb.FindMCServer(id)

	if err != nil {
		return err
	}
	if server == nil {
		return nil
	}

	err = s.provisionerDb.UpdateServerStatus(server.ID, "deleting")
	if err != nil {
		return err
	}

	err = s.linodeClient.DeleteLinode(server.LinodeId)
	if err != nil {
		return err
	}

	err = s.provisionerDb.DeleteMCServer(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *MinecraftLinodeProvisionService) AnnounceMessage(id string, message string) error {
	if id == "" {
		return fmt.Errorf("no server id provided")
	}
	if message == "" {
		return fmt.Errorf("no message provided")
	}

	// Find the server
	server, err := s.provisionerDb.FindMCServer(id)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found")
	}

	// Create SSH config
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(s.serverRootPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Connect to the server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", server.IP), config)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create session
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Escape any quotes in the message to prevent command injection
	escapedMessage := strings.Replace(message, "'", "\\'", -1)
	escapedMessage = strings.Replace(escapedMessage, "\"", "\\\"", -1)

	// Execute the minecraft server command
	cmd := fmt.Sprintf("su - mcserver -c '/home/mcserver/mcserver send \"say %s\"'", escapedMessage)
	err = session.Run(cmd)
	if err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	return nil
}
