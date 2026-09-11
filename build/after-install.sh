#!/bin/bash

# Fix the config to contain proper values. No-ops on upgrades: the package
# manager keeps the existing config, so the placeholders are already gone.
NEW_SECRET=$(head -c 512 /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 32)
sed -i "s/<jwt-secret>/$NEW_SECRET/g" /etc/task64/config.yml
sed -i "s/<rootpath>/\/opt\/task64\//g" /etc/task64/config.yml
sed -i "s/path: \"\.\/task64.db\"/path: \"\\/opt\/task64\/task64.db\"/g" /etc/task64/config.yml

systemctl enable task64.service

# Nothing to reload or restart in chroots and containers without systemd.
if [ -d /run/systemd/system ]; then
	# Pick up changes to task64.service itself.
	systemctl daemon-reload || true
	# try-restart is a no-op while the unit is stopped, so fresh installs stay
	# stopped and upgrades pick up the new binary.
	systemctl try-restart task64.service || true
fi
