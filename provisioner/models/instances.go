package models

type LinodeInstance uint8

const (
	LINODE_INSTANCE_INVALID LinodeInstance = 0
	G6_NANODE_1             LinodeInstance = 1 // Nanode 1gb
	G6_STANDARD_1           LinodeInstance = 2 // Linode 2gb
	G6_STANDARD_2           LinodeInstance = 3 // Linode 4gb
	G6_STANDARD_4           LinodeInstance = 4 // Linode 8gb
	G6_STANDARD_6           LinodeInstance = 5 // Linode 16gb
	G6_STANDARD_8           LinodeInstance = 5 // Linode 32gb
)

func (i LinodeInstance) String() string {
	return [...]string{
		"invalid",
		"g6-nanode",
		"g6-standard-1",
		"g6-standard-2",
		"g6-standard-4",
		"g6-standard-6",
		"g6-standard-8",
	}[i]
}

type MinecraftInstance uint8

const (
	MINECRAFT_INSTANCE_INVALID    MinecraftInstance = 0
	MINECRAFT_INSTANCE_BASIC_1    MinecraftInstance = 1
	MINECRAFT_INSTANCE_STANDARD_1 MinecraftInstance = 2
	MINECRAFT_INSTANCE_PREMIUM_1  MinecraftInstance = 3
	MINECRAFT_INSTANCE_SUPER_1    MinecraftInstance = 4
	MINECRAFT_INSTANCE_ULTIMATE_1 MinecraftInstance = 5
)

func (i MinecraftInstance) String() string {
	return [...]string{
		"MINECRAFT_INSTANCE_INVALID",
		"MINECRAFT_INSTANCE_BASIC_1",
		"MINECRAFT_INSTANCE_STANDARD_1",
		"MINECRAFT_INSTANCE_PREMIUM_1",
		"MINECRAFT_INSTANCE_SUPER_1",
		"MINECRAFT_INSTANCE_ULTIMATE_1",
	}[i]
}
