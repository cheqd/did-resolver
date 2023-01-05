import copy
import sys
import pexpect
import json


RESOLVER_URL = "http://localhost:8080"
PATH = "/1.0/identifiers/"

TESTNET_DID = "did:cheqd:testnet:305b4368-87ee-4690-874e-52f628dbd902"
TESTNET_FRAGMENT = TESTNET_DID + "#key1"
FAKE_TESTNET_DID = "did:cheqd:testnet:76471e8c-0d1c-4b97-9b11-17b65e024133"
TESTNET_RESOURCE = TESTNET_DID + "/resources/44547089-170b-4f5a-bcbc-06e46e0089e4"
TESTNET_RESOURCE_METADATA = TESTNET_RESOURCE + "/metadata"
TESTNET_RESOURCE_LIST = TESTNET_DID + "/resources/all"
TESTNET_RESOURCE_LIST_REDIRECT = TESTNET_DID + "/resources/"
TESTNET_RESOURCE_NAME = "Demo Resource"
RESOURCE_DATA = "{ \"content\": \"test data\" }"
FAKE_TESTNET_FRAGMENT = TESTNET_DID + "#fake_key"
FAKE_TESTNET_RESOURCE = TESTNET_DID + "/resources/76471e8c-0d1c-4b97-9b11-17b65e024334"

MAINNET_DID = "did:cheqd:mainnet:76471e8c-0d1c-4b97-9b11-17b65e024335"
MAINNET_FRAGMENT = MAINNET_DID + "#key1"
FAKE_MAINNET_DID = "did:cheqd:mainnet:76471e8c-0d1c-4b27-9b11-17b65e024133"
FAKE_MAINNET_FRAGMENT = MAINNET_DID + "#fake_key"

DIDJSON = "application/did+json"
DIDLDJSON = "application/did+ld+json"
LDJSON = "application/ld+json"
JSON = "application/json"
HTML = "text/html"

IMPLICIT_TIMEOUT = 40
ENCODING = "utf-8"
READ_BUFFER = 60000


def run(command, params, expected_output):
    cli = pexpect.spawn(f"{command} {params}", encoding=ENCODING, timeout=IMPLICIT_TIMEOUT, maxread=READ_BUFFER)
    cli.logfile = sys.stdout
    cli.expect(expected_output)
    return cli


def json_loads(s_to_load: str) -> dict:
    s = copy.copy(s_to_load)
    s = s.replace("\\", "")
    s = s.replace("\"[", "[")
    s = s.replace("]\"", "]")
    return json.loads(s)
