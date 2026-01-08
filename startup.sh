#!/bin/bash
exec > >(tee /var/log/tpot-install.log|logger -t user-data -s 2>/dev/console) 2>&1

echo "Starting T-Pot Cloud-Init..."


if [ -f /opt/tpot/etc/tpot.yml ]; then
    echo "T-Pot is already installed. Skipping provisioning."
    exit 0
fi


echo "Updating system..."
DEBIAN_FRONTEND=noninteractive apt-get update
DEBIAN_FRONTEND=noninteractive apt-get upgrade -y
DEBIAN_FRONTEND=noninteractive apt-get install -y git ansible net-tools

echo "Cloning T-Pot..."
git clone https://github.com/telekom-security/tpotce /opt/tpotce

cat <<EOF > /opt/tpotce/iso/installer/tpot.conf
myCONF_TPOT_FLAVOR='STANDARD'
myCONF_WEB_USER='admin'
myCONF_WEB_PW='RazzleDazzle2026!'
myCONF_REMOTE_LOG='n'
EOF


echo "Running Installer..."
cd /opt/tpotce/iso/installer/
./install.sh --type=auto --conf=/opt/tpotce/iso/installer/tpot.conf

echo "Installation complete. System will reboot shortly."