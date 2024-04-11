#!/bin/sh

set -exo pipefail

if [ "$#" -eq 0 ] || [ "${1#-}" != "$1" ]; then
	set -- nats-server "$@"
fi

# If we're running in Fly, enable clustering and ensure the hostname is set
if ! [ -z "$FLY_APP_NAME" ]; then
	export HOSTNAME=$FLY_MACHINE_ID
fi
# 	cat <<EOF >>/etc/nats/nats-server.conf
# # Advertise our machine's DNS name
# client_advertise = ${FLY_ALLOC_ID}.vm.${FLY_APP_NAME}.internal:4222
#
# # Enable clustering
# cluster {
#   name = ${FLY_APP_NAME}
#
#   # Route connections to be received on any interface on port 6222
#   host = ::
#   port = 6222
#
#   # Advertise our machine's DNS name
#   cluster_advertise = ${FLY_ALLOC_ID}.vm.${FLY_APP_NAME}.internal:6222
#
#   # Routes are actively solicited and connected to from this server
#   routes = [
# EOF
#
# 	local_ip=$(dig +short AAAA $FLY_ALLOC_ID.vm.$FLY_APP_NAME.internal)
#
# 	# Add all routes except our own
# 	for route in $(dig +short AAAA top3.nearest.of.$FLY_APP_NAME.internal); do
# 		if [ "$route" != "$local_ip" ]; then
# 			echo "    \"nats://[$route]:6222\"" >>/etc/nats/nats-server.conf
# 		fi
# 	done
#
# 	echo -e "  ]\n}" >>/etc/nats/nats-server.conf
# fi

exec "$@"
