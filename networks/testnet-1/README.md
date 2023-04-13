# Composable Testnet 1


## Hardware Requirements

* TBD

## Installation Steps


#### Prerequisites

```shell
# Install Updates
sudo apt update && sudo apt upgrade -y
sudo apt install make build-essential gcc git jq chrony -y

# Install go
wget https://golang.org/dl/go1.19.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
rm go1.19.3.linux-amd64.tar.gz

# Add go to path
echo "export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GO111MODULE=on
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin" >> /home/ubuntu/.bashrc
```

#### Clone git repository

```shell
git clone https://github.com/notional-labs/composable-testnet.git
```

#### Install

```shell
cd composable-testnet/
make install
```

#### Initialize chain

```shell
banksyd init [node_name] --chain-id banksy-testnet-1
```

#### Generate keys

```shell
banksyd keys add [key_name]
```

#### Create gentx

```shell
# Add genesis account 
banksyd add-genesis-account [key_name] 1000000000ubanksy

# Create a validator at genesis
banksyd gentx [key_name] 1000000000ubanksy --moniker [node_name] --chain-id banksy-testnet-1 \
  --commission-max-change-rate 0.1  \
  --commission-max-rate 0.2   \
  --commission-rate 0.05   \
  --min-self-delegation "1"   \
  --website "" \
  --security-contact=""   \
  --identity=   \
  --keyring-backend os \
  --details=""
```

#### Create gentx
Create a new file in `composable-testnet/networks/testnet-1/gentxs/` with your gentx.
