#!/bin/sh

# Fix the config to contain proper values. No-ops on upgrades: the package
# manager keeps the existing config, so the placeholders are already gone.
NEW_SECRET=$(head -c 512 /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 32)
sed -i "s/<jwt-secret>/$NEW_SECRET/g" /etc/task64/config.yml
sed -i "s/<rootpath>/\/opt\/task64\//g" /etc/task64/config.yml
sed -i "s/path: \"\.\/task64.db\"/path: \"\\/opt\/task64\/task64.db\"/g" /etc/task64/config.yml

rc-update add task64 default

# Only restart when already running, so fresh installs stay stopped and
# upgrades pick up the new binary.
if rc-service task64 status >/dev/null 2>&1; then
	rc-service task64 restart || true
fi
