const axios = require('axios');

const SAHAMATI_ROUTER_BASE_URL = "https://api.dev.sahamati.org.in/router";
const DISCOVER_URL = `${SAHAMATI_ROUTER_BASE_URL}/v2/Accounts/discover`;

const DEFAULTS = {
  RECIPIENT_ENTITY_ID: "FIP-SIMULATOR",
  TXN_ID: "f35761ac-4a18-11e8-96ff-0277a9fbfedc2",
  JWS_SIGNATURE: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  SIMULATE_RES: "Ok",
  AUTH_TOKEN: "token"
};

const encodeRequestMetadata = (entityId) => btoa(JSON.stringify({ "recipient-id": entityId }));

function createDiscoverRequest() {
  return {
    ver: "2.0.0",
    timestamp: new Date().toISOString(),
    txnid: DEFAULTS.TXN_ID,
    Customer: {
      id: "customer_identifier@AA_identifier",
      Identifiers: [
        {
          category: "STRONG",
          type: "AADHAAR",
          value: "XXXXXXXXXXXXXXXX"
        }
      ]
    },
    FITypes: ["DEPOSIT"]
  };
}

async function discoverAccounts() {
  const requestBody = createDiscoverRequest();
  
  const config = {
    method: 'post',
    url: DISCOVER_URL,
    headers: {
      'Authorization': `Bearer ${DEFAULTS.AUTH_TOKEN}`,
      'Content-Type': 'application/json',
      'x-jws-signature': DEFAULTS.JWS_SIGNATURE,
      'x-simulate-res': DEFAULTS.SIMULATE_RES,
      'x-request-meta': encodeRequestMetadata(DEFAULTS.RECIPIENT_ENTITY_ID)
    },
    data: requestBody
  };

  return axios(config);
}

async function main() {
  try {
    const response = await discoverAccounts();
    console.log(response.data);
  } catch (error) {
    console.error('Error:', error.message);
  }
}

main();