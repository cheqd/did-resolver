#!/bin/bash

# publish dids
cheqd-noded tx cheqd create-did did.json --from base_account_1 --keyring-backend test --gas auto --gas-adjustment 1.2 --gas-prices 50ncheq --chain-id cheqd

# publish resources
cheqd-noded tx resource create resource.json --from base_account_1 --keyring-backend test --gas auto --gas-adjustment 1.2 --gas-prices 25ncheq --chain-id cheqd