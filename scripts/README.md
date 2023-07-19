MEMO='{"forward":{"receiver":"osmo12g76fl4shtgca30mn25vakwdptxr027vq8q4h3","port":"transfer","channel":"channel-1","timeout":600000000000,"retries":0}}'
{"forward":{"receiver":"cosmos1yp0sn7r00z5x4hvkx7c965yqqlzplw7fru0q26","port":"transfer","channel":"channel-1","timeout":600000000000,"retries":0}}
hermes --config scripts/relayer_hermes/config_compose_gaia.toml create channel --a-chain centaurid-t1 --b-chain gaiad-t1 --a-port transfer --b-port transfer --new-client-connection --yes

hermes --config scripts/relayer_hermes/config_compose_osmosis.toml create channel --a-chain centaurid-t1 --b-chain testing --a-port transfer --b-port transfer --new-client-connection --yes

centaurid tx gov submit-proposal "scripts/proposalAddToken.json" --from myaccount --keyring-backend test --chain-id centaurid-t1 --yes
centaurid tx gov vote 1 yes --from test --keyring-backend test --chain-id centaurid-t1 --yes

centaurid tx gov submit-proposal "scripts/proposalRateLimit.json" --from myaccount --keyring-backend test --chain-id centaurid-t1 --yes
centaurid tx gov vote 2 yes --from test --keyring-backend test --chain-id centaurid-t1 --yes

gaiad tx ibc-transfer transfer transfer channel-0 centauri1efd63aw40lxf3n4mhf7dzhjkr453axur7jv2pe 100000000stake --from gnad --keyring-backend test --chain-id gaiad-t1 --yes --fees 5000stake --memo "$MEMO"

gaiad tx ibc-transfer transfer transfer channel-0 centauri1efd63aw40lxf3n4mhf7dzhjkr453axur7jv2pe 5000000stake --from gnad --keyring-backend test --chain-id gaiad-t1 --yes --fees 5000stake