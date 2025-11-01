#!/bin/bash
set -eu

###########################################################
# VARIABLES
###########################################################
TIMEZONE=Asia/Kolkata
USERNAME=twitbox

# Secure password prompt
read -s -p "Enter password for DB user '${USERNAME}': " DB_PASSWORD
echo

###########################################################
# ENVIRONMENT
###########################################################
export LC_ALL=en_US.UTF-8

###########################################################
# SCRIPT LOGIC
###########################################################

echo ">>> Updating system..."
apt-add-repository --yes universe
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo ">>> Setting timezone..."
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

echo ">>> Creating deploy user..."
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

echo ">>> Copying SSH key for deploy user..."
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME} || echo "No root SSH keys found — skipping copy"

###########################################################
# Firewall
###########################################################
echo ">>> Setting up UFW firewall..."
apt --yes install ufw
ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

###########################################################
# Fail2ban
###########################################################
echo ">>> Installing fail2ban..."
apt --yes install fail2ban

###########################################################
# MySQL
###########################################################
echo ">>> Installing MySQL..."
apt --yes install mysql-server

echo ">>> Securing MySQL..."
mysql --execute="ALTER USER 'root'@'localhost' IDENTIFIED BY 'rootpassword'; FLUSH PRIVILEGES;"


echo ">>> Creating MySQL database and user..."
mysql --user=root --password=rootpassword <<EOF
CREATE DATABASE IF NOT EXISTS ${USERNAME}
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS '${USERNAME}'@'localhost' IDENTIFIED BY '${DB_PASSWORD}';
GRANT ALL PRIVILEGES ON ${USERNAME}.* TO '${USERNAME}'@'localhost';
FLUSH PRIVILEGES;
EOF

echo ">>> Writing TWITBOX_DB_DSN to /etc/environment"
echo "TWITBOX_DB_DSN='${USERNAME}:${DB_PASSWORD}@tcp(localhost:3306)/${USERNAME}?parseTime=true&multiStatements=true'" >> /etc/environment

###########################################################
# Install Migrate (MySQL)
###########################################################
echo ">>> Installing migration tool..."
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz \
  | tar xvz && mv migrate.linux-amd64 /usr/local/bin/migrate

###########################################################
# Install Caddy (HTTPS reverse proxy)
###########################################################
echo ">>> Installing Caddy..."
apt --yes install -y debian-keyring debian-archive-keyring apt-transport-https curl gpg

curl -1sLf https://dl.cloudsmith.io/public/caddy/stable/gpg.key \
  | gpg --dearmor -o /usr/share/keyrings/caddy.gpg

curl -1sLf https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt \
  | sed 's#deb #deb [signed-by=/usr/share/keyrings/caddy.gpg] #' \
  | tee /etc/apt/sources.list.d/caddy-stable.list

apt update
apt --yes install caddy

###########################################################
# Finished
###########################################################
echo "✅ Server provisioning complete — rebooting now..."
reboot
