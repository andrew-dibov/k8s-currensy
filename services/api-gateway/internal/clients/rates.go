package clients

import (
	rates "api-gateway/internal/protos/rates"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type RatesClient struct {
	client rates.RatesServiceClient
	conn   *grpc.ClientConn
}

func NewRatesClient(serviceURL string) (*RatesClient, error) {

	/*
		Канал :
			1. Создание канала
			2. Передача конфигурационных файлов
	*/

	conn, err := grpc.NewClient(serviceURL,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10), grpc.MaxCallSendMsgSize(1024*1024*10)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC client : %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	/*
		Железо :
			1. Процессор формирует сетевые пакеты
			2. Операционная система выделяет порты и буферы для сокета
			3. Сетевая карта отправляет и принимает пакеты

		Подключение :
			1. DNS запрос
			2. TCP рукопожатие : SYN, SYN-ACK, ACK
			3. TLS/SSL рукопожатие : этот шаг пропускается при insecure.NewCredentials()
			4. HTTP/2 рукопожатие : gRPC работает поверх HTTP/2 : обмен настройками
	*/

	conn.Connect()

	if err := waitForReady(ctx, conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to wait for connection : %w", err)
	}

	return &RatesClient{
		client: rates.NewRatesServiceClient(conn),
		conn:   conn,
	}, nil
}

func waitForReady(ctx context.Context, conn *grpc.ClientConn) error {

	/*
		Мониторим состояние канала :
			1. Создаем бесконечный цикл
			2. Запрашиваем состояние
			3. Обрабатываем состояние
	*/

	for {
		state := conn.GetState()

		/* Соединение установлено : выход без ошибки */

		if state == connectivity.Ready {
			return nil // Соединение установлено : выходим без ошибки
		}

		/* Канал закрылся : выход с ошибкой */

		if state == connectivity.Shutdown {
			return fmt.Errorf("Connection shutdown")
		}

		/* Горутина блокируется и висит : ждем сигнал ОС или таймаута */

		if !conn.WaitForStateChange(ctx, state) {

			/* Состояние изменилось : сигнал ОС или таймаут пришел */

			if ctx.Err() != nil {
				return fmt.Errorf("Timeout waiting for connection : %w", ctx.Err())
			}
		}

		/*
			1. Если сигнал от ОС : смена состояния
			2. Если не ready : соединение не установлено
			3. Проверка состояний : новый цикл
		*/
	}
}
