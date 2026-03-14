curl -H "X-API-Key: test-a" "http://localhost:8080/api/v1/rate?fromCurrency=USD&toCurrency=RUB"
curl -H "X-API-Key: test-b" -H "Content-Type: application/json" -d '{"from_currency": "USD", "to_currency": "RUB", "amount": 100}' "http://localhost:8080/api/v1/convert"

# curl -H "X-API-Key: test-wrong-a" "http://localhost:8080/api/v1/rate?fromCurrency=USD&toCurrency=RUB"
# curl -H "X-API-Key: test-wrong-b" -H "Content-Type: application/json" -d '{"from_currency": "USD", "to_currency": "RUB", "amount": 100}' "http://localhost:8080/api/v1/convert"