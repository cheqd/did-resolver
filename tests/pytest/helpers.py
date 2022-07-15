import copy
import sys
import pexpect
import json


RESOLVER_URL = "http://localhost:1313"
PATH = "/1.0/identifiers/"

TESTNET_DID = "did:cheqd:testnet:zFWM1mKVGGU2gHYuLAQcTJfZBebqBpGf"
TESTNET_FRAGMENT = TESTNET_DID + "#key1"
FAKE_TESTNET_DID = "did:cheqd:testnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY"
FAKE_TESTNET_FRAGMENT = TESTNET_DID + "#fake_key"
FAKE_TESTNET_RESOURCE = TESTNET_DID + "/resources/76471e8c-0d1c-4b97-9b11-17b65e024334"

MAINNET_DID = "did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY"
MAINNET_FRAGMENT = MAINNET_DID + "#key1"
FAKE_MAINNET_DID = "did:cheqd:mainnet:zFWM1mKVGGU2gHYuLAQcTJfZBebqBpGf"
FAKE_MAINNET_FRAGMENT = MAINNET_DID + "#fake_key"

DIDJSON = "application/did+json"
DIDLDJSON = "application/did+ld+json"
LDJSON = "application/ld+json"
LDJSONProfile = LDJSON + ";profile=\"https://w3id.org/did-resolution\""
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

