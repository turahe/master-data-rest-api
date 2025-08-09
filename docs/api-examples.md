# 🌐 API Examples & Usage Patterns

This guide provides comprehensive examples for using the Master Data REST API. All examples assume you have a valid API key and the server is running on `localhost:8080`.

## 🔐 Authentication

All API requests require an API key in the Authorization header:

```bash
# Set your API key as an environment variable
export API_KEY="your-api-key-here"

# Use in requests
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories"
```

## 🗺️ Geodirectories API

### Basic Operations

#### Get All Geodirectories
```bash
# Get all geodirectories with pagination
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories?limit=50&offset=0"
```

**Response:**
```json
{
  "success": true,
  "message": "Geodirectories retrieved successfully",
  "data": {
    "geodirectories": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Asia",
        "type": "CONTINENT",
        "code": "AS",
        "postal_code": null,
        "longitude": "139.6917",
        "latitude": "35.6895",
        "record_left": 1,
        "record_right": 100,
        "record_ordering": 1,
        "parent_id": null,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "limit": 50,
    "offset": 0
  }
}
```

#### Create a New Geodirectory
```bash
# Create a new country
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Indonesia",
       "type": "COUNTRY",
       "code": "ID",
       "postal_code": null,
       "longitude": "106.8456",
       "latitude": "-6.2088",
       "parent_id": "550e8400-e29b-41d4-a716-446655440000"
     }' \
     "http://localhost:8080/api/v1/geodirectories"
```

#### Get Geodirectory by ID
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/550e8400-e29b-41d4-a716-446655440000"
```

#### Update a Geodirectory
```bash
curl -X PUT \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Republic of Indonesia",
       "latitude": "-6.2088",
       "longitude": "106.8456"
     }' \
     "http://localhost:8080/api/v1/geodirectories/indonesia-id"
```

#### Delete a Geodirectory
```bash
curl -X DELETE \
     -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/indonesia-id"
```

### Hierarchical Operations

#### Get Countries (Type Filtering)
```bash
# Get all countries
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/type/COUNTRY"

# Get all provinces  
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/type/PROVINCE"
```

#### Get Children of a Geodirectory
```bash
# Get all provinces of Indonesia
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/indonesia-id/children?type=PROVINCE"

# Get all direct children (any type)
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/indonesia-id/children"
```

#### Get All Descendants (Nested Set Model)
```bash
# Get all descendants of Indonesia (provinces, cities, districts, villages)
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/indonesia-id/descendants?limit=100"
```

#### Get Ancestors
```bash
# Get all ancestors of a village (district -> city -> province -> country -> continent)
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/village-id/ancestors"
```

#### Get Hierarchy (Parent + Children)
```bash
# Get a geodirectory with its parent and children
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/jakarta-id/hierarchy"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "geodirectory": {
      "id": "jakarta-id",
      "name": "Jakarta",
      "type": "CITY"
    },
    "parent": {
      "id": "west-java-id",
      "name": "West Java",
      "type": "PROVINCE"
    },
    "children": [
      {
        "id": "central-jakarta-id",
        "name": "Central Jakarta",
        "type": "DISTRICT"
      },
      {
        "id": "north-jakarta-id", 
        "name": "North Jakarta",
        "type": "DISTRICT"
      }
    ]
  }
}
```

### Advanced Operations

#### Search Geodirectories
```bash
# Search by name, code, or postal code
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/search?q=jakarta&limit=10"

# Search with pagination
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/search?q=central&limit=20&offset=0"
```

#### Move a Geodirectory
```bash
# Move a district to a different city
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "new_parent_id": "new-city-id"
     }' \
     "http://localhost:8080/api/v1/geodirectories/district-id/move"
```

#### Rebuild Nested Set Structure
```bash
# Rebuild the entire nested set structure (admin operation)
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/geodirectories/rebuild"
```

## 🏦 Banks API

### Basic Bank Operations

#### Get All Banks
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/banks?limit=50&offset=0"
```

#### Create a New Bank
```bash
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Bank Central Asia",
       "alias": "BCA",
       "company": "PT Bank Central Asia Tbk",
       "code": "014"
     }' \
     "http://localhost:8080/api/v1/banks"
```

#### Get Bank by ID
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/banks/bank-id"
```

#### Get Bank by Code
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/banks/code/014"
```

#### Search Banks
```bash
# Search by name, alias, company, or code
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/banks/search?q=central&limit=10"
```

#### Update Bank
```bash
curl -X PUT \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Bank Central Asia Updated",
       "alias": "BCA",
       "company": "PT Bank Central Asia Tbk",
       "code": "014"
     }' \
     "http://localhost:8080/api/v1/banks/bank-id"
```

## 💰 Currencies API

### Currency Management

#### Get All Currencies
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/currencies"
```

#### Create a Currency
```bash
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Indonesian Rupiah",
       "code": "IDR",
       "symbol": "Rp",
       "decimal_places": 0,
       "is_active": true
     }' \
     "http://localhost:8080/api/v1/currencies"
```

#### Get Currency by Code
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/currencies/code/IDR"
```

#### Activate/Deactivate Currency
```bash
# Activate a currency
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/currencies/currency-id/activate"

# Deactivate a currency
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/currencies/currency-id/deactivate"
```

## 🗣️ Languages API

### Language Management

#### Get All Languages
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/languages"
```

#### Create a Language
```bash
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Indonesian",
       "code": "id",
       "is_active": true
     }' \
     "http://localhost:8080/api/v1/languages"
```

#### Get Language by Code
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/languages/code/id"
```

## 🔑 API Keys Management

### Managing API Keys

#### Get All API Keys
```bash
curl -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/api-keys"
```

#### Create New API Key
```bash
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Mobile App Key",
       "description": "API key for mobile application",
       "expires_at": "2024-12-31T23:59:59Z"
     }' \
     "http://localhost:8080/api/v1/api-keys"
```

#### Deactivate API Key
```bash
curl -X POST \
     -H "Authorization: Bearer $API_KEY" \
     "http://localhost:8080/api/v1/api-keys/api-key-id/deactivate"
```

## 🎯 Common Usage Patterns

### 1. Building a Location Picker

```bash
# Step 1: Get all countries
countries=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories/type/COUNTRY")

# Step 2: User selects Indonesia, get provinces  
provinces=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories/indonesia-id/children?type=PROVINCE")

# Step 3: User selects West Java, get cities
cities=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories/west-java-id/children?type=CITY")

# Step 4: User selects Jakarta, get districts
districts=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories/jakarta-id/children?type=DISTRICT")
```

### 2. Bank Dropdown with Search

```bash
# Get popular banks (first 20)
popular_banks=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/banks?limit=20")

# Search banks as user types
search_results=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/banks/search?q=mandiri&limit=5")
```

### 3. Multi-Currency Price Display

```bash
# Get all active currencies
active_currencies=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/currencies" | \
  jq '.data.currencies[] | select(.is_active == true)')

# Get specific currency details
currency_detail=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/currencies/code/USD")
```

### 4. Geographic Analytics

```bash
# Get all locations in a province for analysis
all_locations=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories/west-java-id/descendants?limit=1000")

# Get geographic breadcrumb trail
breadcrumb=$(curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories/village-id/ancestors")
```

## 🔧 Error Handling

### Common Error Responses

#### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized: Invalid or missing API key",
  "error": "UNAUTHORIZED"
}
```

#### 404 Not Found
```json
{
  "success": false,
  "message": "Geodirectory not found",
  "error": "NOT_FOUND"
}
```

#### 400 Bad Request
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "VALIDATION_ERROR",
  "details": {
    "name": ["Name is required"],
    "type": ["Type must be one of: CONTINENT, COUNTRY, PROVINCE, CITY, DISTRICT, VILLAGE"]
  }
}
```

### Error Handling in Scripts

```bash
#!/bin/bash

# Function to make API calls with error handling
api_call() {
  local url="$1"
  local method="${2:-GET}"
  local data="$3"
  
  local response
  local http_code
  
  if [ "$method" = "POST" ] && [ -n "$data" ]; then
    response=$(curl -s -w "%{http_code}" \
      -X POST \
      -H "Authorization: Bearer $API_KEY" \
      -H "Content-Type: application/json" \
      -d "$data" \
      "$url")
  else
    response=$(curl -s -w "%{http_code}" \
      -H "Authorization: Bearer $API_KEY" \
      "$url")
  fi
  
  http_code="${response: -3}"
  body="${response%???}"
  
  case $http_code in
    200|201)
      echo "$body" | jq .
      return 0
      ;;
    401)
      echo "Error: Unauthorized. Check your API key." >&2
      return 1
      ;;
    404)
      echo "Error: Resource not found." >&2
      return 1
      ;;
    *)
      echo "Error: HTTP $http_code" >&2
      echo "$body" | jq . >&2
      return 1
      ;;
  esac
}

# Example usage
if api_call "http://localhost:8080/api/v1/geodirectories"; then
  echo "Success!"
else
  echo "Failed to get geodirectories"
  exit 1
fi
```

## 📊 Batch Operations

### Bulk Data Import

```bash
#!/bin/bash

# Import multiple banks from CSV
while IFS=',' read -r name alias company code; do
  if [ "$name" != "Name" ]; then  # Skip header
    echo "Creating bank: $name"
    curl -X POST \
      -H "Authorization: Bearer $API_KEY" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"$name\",
        \"alias\": \"$alias\", 
        \"company\": \"$company\",
        \"code\": \"$code\"
      }" \
      "http://localhost:8080/api/v1/banks"
    echo
  fi
done < banks.csv
```

### Data Export

```bash
#!/bin/bash

# Export all geodirectories to JSON
curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories?limit=10000" | \
  jq '.data.geodirectories' > geodirectories_export.json

# Export with hierarchy information
curl -s -H "Authorization: Bearer $API_KEY" \
  "http://localhost:8080/api/v1/geodirectories" | \
  jq '.data.geodirectories[] | select(.type == "COUNTRY")' | \
  while read -r country; do
    country_id=$(echo "$country" | jq -r '.id')
    echo "Exporting $country_id..."
    curl -s -H "Authorization: Bearer $API_KEY" \
      "http://localhost:8080/api/v1/geodirectories/$country_id/descendants" \
      > "export_${country_id}.json"
  done
```

## 🧪 Testing with Different Tools

### Using HTTPie

```bash
# Install HTTPie: pip install httpie

# GET request
http GET localhost:8080/api/v1/banks Authorization:"Bearer $API_KEY"

# POST request  
http POST localhost:8080/api/v1/banks \
  Authorization:"Bearer $API_KEY" \
  name="Test Bank" alias="TEST" company="Test Corp" code="999"
```

### Using Python requests

```python
import requests
import json

API_KEY = "your-api-key-here"
BASE_URL = "http://localhost:8080/api/v1"

headers = {
    "Authorization": f"Bearer {API_KEY}",
    "Content-Type": "application/json"
}

# Get all countries
response = requests.get(f"{BASE_URL}/geodirectories/type/COUNTRY", headers=headers)
countries = response.json()

# Create a new bank
bank_data = {
    "name": "Test Bank",
    "alias": "TEST", 
    "company": "Test Corporation",
    "code": "999"
}
response = requests.post(f"{BASE_URL}/banks", headers=headers, json=bank_data)
print(response.json())
```

### Using JavaScript/Node.js

```javascript
const axios = require('axios');

const API_KEY = 'your-api-key-here';
const BASE_URL = 'http://localhost:8080/api/v1';

const headers = {
  'Authorization': `Bearer ${API_KEY}`,
  'Content-Type': 'application/json'
};

// Get all geodirectories
async function getGeodirectories() {
  try {
    const response = await axios.get(`${BASE_URL}/geodirectories`, { headers });
    console.log(response.data);
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
}

// Create a new currency
async function createCurrency() {
  const currencyData = {
    name: 'US Dollar',
    code: 'USD',
    symbol: '$',
    decimal_places: 2,
    is_active: true
  };
  
  try {
    const response = await axios.post(`${BASE_URL}/currencies`, currencyData, { headers });
    console.log('Currency created:', response.data);
  } catch (error) {
    console.error('Error:', error.response?.data || error.message);
  }
}

getGeodirectories();
createCurrency();
```

---

<div align="center">
  <strong>🚀 Start building amazing applications with the Master Data API!</strong>
  <br>
  <sub>Explore the interactive documentation at <a href="http://localhost:8080/swagger/">localhost:8080/swagger/</a></sub>
</div>
