HOST="localhost"
PORT="8080"
API_KEY="test-1111"
BASE_URL="http://$HOST:$PORT"

curl -i "$BASE_URL/health"
curl -i -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rate?fromCurrency=USD&toCurrency=EUR"
curl -i -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rates"
curl -i -H "X-API-Key: $API_KEY" "$BASE_URL/api/v1/rates?baseCurrency=EUR"
curl -i -X POST -H "X-API-Key: $API_KEY" -H "Content-Type: application/json" \
  -d '{"fromCurrency":"USD","toCurrency":"EUR","amount":100}' \
  "$BASE_URL/api/v1/convert"