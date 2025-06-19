#!/bin/bash
apt-get -y update

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

crontab -l > gamecron
echo "*/5 * * * * su - mcserver -c '/home/mcserver/mcserver monitor' > /dev/null 2>&1" >> gamecron
echo "*/30 * * * * su - mcserver -c '/home/mcserver/mcserver update' > /dev/null 2>&1" >> gamecron
echo "0 0 * * 0 su - mcserver -c '/home/mcserver/mcserver update-lgsm' > /dev/null 2>&1" >> gamecron
crontab gamecron
rm gamecron
