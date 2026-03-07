package middlewares

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

/*
	Middleware : логирование запросов
		1. LoggerMiddleware(next, logger) : принять следующий обработчик + логгер
		2. Вернуть обработчик : логгирование запросов
*/

func LoggerMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {

	/*
		HTTP : новый обработчик
			1. http.HandlerFunc(func) : создать обработчик из функции
			2. http.ResponseWriter : интерфейс отправки ответа клиенту
			3. http.Request : структура для хранения информации об HTTP запросе
	*/

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		/*
			Получил запрос
				1. received : время начала обработки запроса
				2. Логирование запроса
				3. ServeHTTP(res, req) : вызов следующего в цепи обработчика
				4. Завершение обработчиков
				5. duration : интервал от received до настоящего времени
				6. Логирование длительности обработки запроса
			Отправил ответ
		*/

		received := time.Now()

		logger.WithFields(logrus.Fields{
			"method": req.Method,
			"path":   req.URL.Path,
			"ip":     req.RemoteAddr,
		}).Info("Request received")

		next.ServeHTTP(res, req)

		duration := time.Since(received)

		logger.WithFields(logrus.Fields{
			"duration": duration,
			"path":     req.URL.Path,
		}).Info("Request processed")
	})
}
