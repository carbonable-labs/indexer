default:
    just --list

feeder_gateway := "https://alpha-sepolia.starknet.io/feeder_gateway"

# start docker database
start_db:
    docker compose -f compose.yaml up -d

# stop docker database
stop_db:
    docker compose -f compose.yaml stop

# run synchronizer application
sync:
    FEEDER_GATEWAY={{feeder_gateway}} go run main.go
