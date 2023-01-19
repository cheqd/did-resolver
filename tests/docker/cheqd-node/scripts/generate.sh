#!/bin/bash

# publish dids
cheqd-noded tx cheqd create-did did.json --from base_account_1 --keyring-backend test --fees 5000000ncheq --chain-id cheqd -y

# publish resources
cheqd-noded tx resource create resource.json --from base_account_1 --keyring-backend test --fees 5000000ncheq --chain-id cheqd -y