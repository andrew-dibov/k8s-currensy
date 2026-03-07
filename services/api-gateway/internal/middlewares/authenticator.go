package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

/*
	Middleware : аутентификация запросов
		1. LoggerMiddleware(next, logger) : принять следующий обработчик + логгер
		2. Вернуть обработчик : логгирование запросов
*/

func AuthenticatorMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {

	/*
		HTTP : новый обработчик
			1. http.HandlerFunc(func) : создать обработчик из функции
			2. http.ResponseWriter : интерфейс отправки ответа клиенту
			3. http.Request : структура для хранения информации об HTTP запросе
	*/

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		/*
			Получил запрос
				1. Пропуск проверки для /health : вызов следующего в цепи обработчика
				2. Get("X-API-Key") : поиск ключа в заголовке
				3. Get("api_key") : поиск ключа в параметре
				4. Структура данных : пары "ключ-значение" : !!! ИСПОЛЬЗОВАТЬ СТОРОННИЙ СЕРВИС !!!
				5. Проверка на отсутствие ключа
				6. Проверка существования ключа в map :
					- Если да : пропустить запрос
					- Если нет : залогировать неверный ключ и вернуть ошибку
				7. ServeHTTP(res, req) : вызов следующего в цепи обработчика
			Отправил ответ
		*/

		if req.URL.Path == "/health" {
			next.ServeHTTP(res, req)
			return
		}

		apiKey := req.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = req.URL.Query().Get("api_key")
		}

		validKeys := map[string]bool{
			"secret-key-1234": true,
			"test-key-1234":   true,
			"dev-key-1234":    true,
		}

		if apiKey == "" {
			logger.Warn("Absent API key")
			http.Error(res, "Provide API key or go fuck yourself", http.StatusUnauthorized)
			return
		}

		if !validKeys[apiKey] {
			logger.WithField("key", apiKey).Warn("Wrong API key")
			http.Error(res, "Wrong API key", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}
