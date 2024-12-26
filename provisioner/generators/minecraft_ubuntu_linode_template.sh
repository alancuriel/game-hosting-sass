#!/bin/bash

sudo apt update
sudo apt upgrade -y
sudo apt install -y lib32gcc-s1 lib32stdc++6 libsdl2-2.0-0

export DEBIAN_FRONTEND=noninteractive
adduser --disabled-password --gecos "" mcserver
echo "mcserver ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

wget https://linuxgsm.com/dl/linuxgsm.sh -P /home/mcserver/
chmod +x /home/mcserver/linuxgsm.sh
chown -R mcserver:mcserver /home/mcserver/*

echo running LinuxGSM script
su - mcserver -c "/home/mcserver/linuxgsm.sh mcserver"

su - mcserver -c "/home/mcserver/mcserver auto-install"
su - mcserver -c "sed -i \"s/motd=.*/motd=${{OPUSER}}'s server/\" /home/mcserver/serverfiles/server.properties"

echo Starting up the gameserver!
su - mcserver -c "/home/mcserver/mcserver start"

sleep 10
su - mcserver -c '/home/mcserver/mcserver send "op ${{OPUSER}}"'
