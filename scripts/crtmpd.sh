#!/bin/sh
echo "*** Running cRTMPserver"
exec /etc/crtmpserver/crtmpserver "/etc/crtmpserver/configs/pns.lua"  >> /var/log/crtmpserver/crtmpd.log 2>&1
