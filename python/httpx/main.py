import json
import base64
import httpx
from typing import Dict, Any, Optional, Union
from datetime import datetime
import uuid


SAHAMATI_ROUTER_BASE_URL = "https://api.dev.sahamati.org.in/router"
DISCOVER_URL = f"{SAHAMATI_ROUTER_BASE_URL}/v2/Accounts/discover"


class SahamatiClient:
    def __init__(self):
        self.http_client = httpx.Client()

    def create_discover_request(self) -> Dict[str, Any]:
        return {
            "ver": "2.0.0",
            "timestamp": datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%S.%fZ"),
            "txnid": str(uuid.uuid4()),
            "Customer": {
                "id": "customer_identifier@AA_identifier",
                "Identifiers": [
                    {
                        "category": "STRONG",
                        "type": "AADHAAR",
                        "value": "XXXXXXXXXXXXXXXX"
                    }
                ]
            },
            "FITypes": ["DEPOSIT"]
        }

    def set_headers(self, request: httpx.Request) -> None:
        metadata_header = self.encode_request_metadata("FIP-SIMULATOR")
        
        request.headers.update({
            "x-jws-signature": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
            "Content-Type": "application/json",
            "x-simulate-res": "Ok",
            "Authorization": "Bearer token",
            "x-request-meta": metadata_header
        })

    @staticmethod
    def encode_request_metadata(entity_id: str) -> str:
        metadata = {"recipient-id": entity_id}
        json_bytes = json.dumps(metadata).encode('utf-8')
        return base64.b64encode(json_bytes).decode('utf-8')

    def discover_accounts(self) -> Dict[str, Any]:
        request_body = self.create_discover_request()
        
        request = httpx.Request(
            method="POST",
            url=DISCOVER_URL,
            json=request_body
        )
        
        self.set_headers(request)
        
        response = self.http_client.send(request)
        response.raise_for_status()
        
        return response.json()

    def __del__(self):
        self.http_client.close()


def main():
    client = SahamatiClient()
    try:
        response = client.discover_accounts()
        print(json.dumps(response, indent=2))
    except Exception as e:
        print(f"Error: {e}")


if __name__ == "__main__":
    main()