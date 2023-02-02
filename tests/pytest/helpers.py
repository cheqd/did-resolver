import copy
import sys
import pexpect
import json


RESOLVER_URL = "http://localhost:8080"
PATH = "/1.0/identifiers/"
METADATA_PATH = "/metadata"

FAKE_KEY = "#fake_key"

MAINNET_DID = "did:cheqd:mainnet:76471e8c-0d1c-4b97-9b11-17b65e024335"
TESTNET_DID = "did:cheqd:testnet:c1685ca0-1f5b-439c-8eb8-5c0e85ab7cd0"
FAKE_MAINNET_DID = "did:cheqd:mainnet:76471e8c-0d1c-4b27-9b11-17b65e024133"
FAKE_TESTNET_DID = "did:cheqd:testnet:76471e8c-0d1c-4b97-9b11-17b65e024133"

INDY_TESTNET_DID_1 = "did:cheqd:testnet:zHqbcXb3irKRCMst"
MIGRATED_INDY_TESTNET_DID_1 = "did:cheqd:testnet:CpeMubv5yw63jXyrgRRsxR"
INDY_TESTNET_DID_2 = "did:cheqd:testnet:zEv9FXHwp8eFeHbeTXamwda8YoPfgU12"
MIGRATED_INDY_TESTNET_DID_2 = "did:cheqd:testnet:3KpiDD6Hxs4i2G7FtpiGhu"

TESTNET_DID_VERSION_ID = "e5615fc2-6f13-42b1-989c-49576a574cef"
TESTNET_DID_VERSION = TESTNET_DID + "/version/" + TESTNET_DID_VERSION_ID
TESTNET_DID_VERSIONS = TESTNET_DID + "/versions"

MAINNET_FRAGMENT_1 = MAINNET_DID + "#key1"
TESTNET_FRAGMENT_1 = TESTNET_DID + "#key-1"
TESTNET_FRAGMENT_2 = INDY_TESTNET_DID_2 + "#key-1"
MIGRATED_TESTNET_FRAGMENT_2 = MIGRATED_INDY_TESTNET_DID_2 + "#key-1"
FAKE_MAINNET_FRAGMENT_1 = MAINNET_DID + FAKE_KEY
FAKE_TESTNET_FRAGMENT_1 = TESTNET_DID + FAKE_KEY
FAKE_TESTNET_FRAGMENT_2 = INDY_TESTNET_DID_2 + FAKE_KEY

FAKE_TESTNET_VERSION_ID = "e5615fc2-6f13-42b1-989c-49576a574ced"
FAKE_TESTNET_VERSION = FAKE_TESTNET_DID + "/version/" + FAKE_TESTNET_VERSION_ID
FAKE_TESTNET_DID_VERSIONS = FAKE_TESTNET_DID + "/versions"

TESTNET_RESOURCE_1 = TESTNET_DID + \
    "/resources/9ba3922e-d5f5-4f53-b265-fc0d4e988c77"
TESTNET_RESOURCE_2 = INDY_TESTNET_DID_2 + \
    "/resources/214b8b61-a861-416b-a7e4-45533af40ada"
TESTNET_RESOURCE_METADATA_1 = TESTNET_RESOURCE_1 + METADATA_PATH
TESTNET_RESOURCE_METADATA_2 = TESTNET_RESOURCE_2 + METADATA_PATH
TESTNET_RESOURCE_LIST = TESTNET_DID + METADATA_PATH
TESTNET_RESOURCE_NAME = "Demo Resource"
RESOURCE_DATA_1 = "{ \r\n    \"content\": \"test data\"\r\n}"
FAKE_TESTNET_FRAGMENT = TESTNET_DID + FAKE_KEY
FAKE_TESTNET_RESOURCE = TESTNET_DID + \
    "/resources/76471e8c-0d1c-4b97-9b11-17b65e024334"

DIDJSON = "application/did+json"
DIDLDJSON = "application/did+ld+json"
LDJSON = "application/ld+json"
JSON = "application/json"
HTML = "text/html"

IMPLICIT_TIMEOUT = 40
ENCODING = "utf-8"
READ_BUFFER = 60000


def run(command, params, expected_output):
    cli = pexpect.spawn(f"{command} {params}", encoding=ENCODING,
                        timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)
    return cli


def json_loads(s_to_load: str) -> dict:
    s = copy.copy(s_to_load)
    s = s.replace("\\", "")
    s = s.replace("\"[", "[")
    s = s.replace("]\"", "]")
    return json.loads(s)
