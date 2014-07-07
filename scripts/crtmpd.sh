#!/bin/sh
echo "*** Running cRTMPserver"
exec /etc/crtmpserver/crtmpserver "/etc/crtmpserver/configs/flvplayback.lua"  >> /var/log/crtmpserver/crtmpd.log 2>&1
