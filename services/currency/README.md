Currency Service

- Хранит и отдает актуальные курсы валют
- Обновляет актуальность курсов валют

- Есть курс доллара к другим валютам :
  - USD > RUB = 90.00 : за 1 доллар дают 90 рублей
  - USD > EUR = 0.85 : за 1 доллар дают 0.85 евро

```
base_currency | currency_code | rate
USD           | RUB           | 90.00
USD           | EUR           | 0.85
```

- USD -> RUB : amount * rate
  - 100 USD -> 100 * 90.00 = 9000 RUB
- RUB -> USD : (1 / rate) * amount
  - 100 RUB -> (1 / 90.00) * 100 = 1.11 USD
- RUB -> USD -> EUR : ((1 / rub_rate) * amount) * eur_rate :
  - 100 RUB -> ((1 / 90.00) * 100) * 0.85 = 0.94
