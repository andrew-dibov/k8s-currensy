package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"

	proto "api-gateway/internal/protos/rates"
)

/*
	RatesClient : структура
		- gRPC : клиент : генерируется proto
		- conn : управление соединением : библиотека gRPC
*/

type RatesClient struct {
	gRPC proto.RatesServiceClient
	conn *grpc.ClientConn
}

/* NewRatesClient(url) : создание gRPC */

func NewRatesClient(url string) (*RatesClient, error) {

	/*
		Создание канала :
			- WithDefaultCallOptions : свойства вызовов
				- 10МБ лимит на получение
				- 10МБ лимит на отправку
			- WithTransportCredentials : шифрование вызовов
			- WithDefaultServiceConfig : конфигурация
	*/

	conn, err := grpc.NewClient(url,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10), grpc.MaxCallSendMsgSize(1024*1024*10)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		/*
			- loadBalancingPolicy :
				- DNS вернет несколько адресов подов
				- RoundRobin распределит запросы между ними
			- methodConfig : конфигурация по методу
				- name : название метода
				- retryPolicy : повторные запросы метода :
					- maxAttempts : до 3 попыток
					- initialBackoff : начальная задержка
					- maxBackoff : максимальная задержка
					- backoffMultiplier : прирост задержки
					- retryableStatusCodes : статусы для ретраев

			Попытка 1 : UNAVAILABLE
				ждем 0.1s
			Попытка 2 : неудача
				ждем 0.2s : initialBackoff * 2
			Попытка 3 : неудача
				ждем 0.4s : 0.2s * 2
			Попытка 4 : СТОП : maxAttempts = 3
		*/

		grpc.WithDefaultServiceConfig(`{
			"loadBalancingPolicy": "round_robin",
      "methodConfig": [{
        "name": [{"service": "rates.RatesService"}],
        "retryPolicy": {
          "maxAttempts": 3,
          "initialBackoff": "0.1s",
          "maxBackoff": "1s",
          "backoffMultiplier": 2,
          "retryableStatusCodes": ["UNAVAILABLE"]
        }
      }]
		}`),
	)

	/* Проверка создания канала */
	if err != nil {
		return nil, fmt.Errorf("Failed to create channel : %w", err)
	}

	/*
		Context : механизм управления временем жизни операций
			- WithTimeout(ctx, timeout) : новый контекст с таймаутом
			- defer cancel() : закрыть контекст перед завершением функции
	*/

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	/*
		Connect : подключение
			1. DNS запрос
			2. TCP рукопожатие : SYN, SYN-ACK, ACK
			3. TLS/SSL рукопожатие : этот шаг пропускается при insecure.NewCredentials()
			4. HTTP/2 рукопожатие : gRPC работает поверх HTTP/2 : обмен настройками
	*/

	conn.Connect()

	/* Проверка состояния канала */
	if err := waitForReady(ctx, conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to wait for ready : %w", err)
	}

	/* Инициализация и возврат структуры */
	return &RatesClient{
		gRPC: proto.NewRatesServiceClient(conn),
		conn: conn,
	}, nil // Если ошибка : nil
}

func waitForReady(ctx context.Context, conn *grpc.ClientConn) error {

	/*
		Отслеживать состояние :
			1. Бесконечный цикл
			2. Запрос состояние
			3. Обработка состояния
	*/

	for {
		state := conn.GetState()

		/* Соединение установлено : выход без ошибки */
		if state == connectivity.Ready {
			return nil
		}

		/* Канал закрылся : выход с ошибкой */
		if state == connectivity.Shutdown {
			return fmt.Errorf("Connection shutdown")
		}

		/* Блокировка : ожидание изменения состояния подключения */
		if !conn.WaitForStateChange(ctx, state) {

			/* Состояние изменилось : таймаут или контекст отменен */
			if ctx.Err() != nil {
				return fmt.Errorf("Connection timeout : %w", ctx.Err())
			}
		}

		/* Если дошли сюда : состояние изменилось : проверяем повторно */
	}
}

/* !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!11 */

func (client *RatesClient) GetRates(ctx context.Context, baseCurrency string) (*proto.GetRatesResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return client.gRPC.GetRates(callCtx, &proto.GetRatesRequest{
		BaseCurrency: baseCurrency,
	})
}
