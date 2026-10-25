# Protocol: wfc-dls

Internal HTTP defaults to WFC_HTTP_ADDR=:20206. GET /healthz is unauthenticated;
protocol/internal paths require WFC_INTERNAL_KEY. Console HTTP reaches this owner
through wfc-nas; TCP/UDP ports belong directly to their protocol's container.
Runtime frontend/backend net/rpc endpoints stay on loopback inside each TCP
service container. Restarting another protocol does not restart these sockets.
