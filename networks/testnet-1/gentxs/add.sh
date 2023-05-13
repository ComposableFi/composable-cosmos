
#!/bin/bash

# Set the path to the folder containing the JSON files
JSON_FOLDER="/home/duc/Notional/composable-testnet/networks/testnet-1/gentxs/"

# Loop through each JSON file in the folder and extract the "delegator_address" value
for file in ${JSON_FOLDER}/*.json; do
  echo "Processing ${file}..."
  address=$(jq -r '.body.messages[0].delegator_address' ${file})
  ~/go/bin/banksyd add-genesis-account $address 1000000000ubanksy
done