#!/bin/bash

# publish dids
cheqd-noded tx cheqd create-did did.json --from cheqd1rnr5jrt4exl0samwj0yegv99jeskl0hsxmcz96 --keyring-backend test

# publish resources
cheqd-noded tx resource create resource.json --from cheqd1rnr5jrt4exl0samwj0yegv99jeskl0hsxmcz96 --keyring-backend test