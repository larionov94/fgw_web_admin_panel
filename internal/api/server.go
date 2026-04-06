package api

import (
	"context"
	"errors"
	"fgw_web_aforms_panel/pkg/logg"
	"fgw_web_aforms_panel/pkg/msg"
	"net/http"
	"os"
	"time"
)

const (
	defaultReadTimeOut  = 15 * time.Second
	defaultWriteTimeOut = 15 * time.Second
	defaultIdlerTimeOut = 180 * time.Second
)

type Server struct {
	httpServer *http.Server
	logger     *logg.Logger
}

// NewServer создаёт и инициализирует новый экземпляр сервера.
//
// Параметры:
//   - addr: адрес для прослушивания (:8080);
//   - handler: HTTP-обработчик (роутер/mux);
//   - logger: логирование записи событий сервера.
func NewServer(addr string, handler http.Handler, logger *logg.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  defaultReadTimeOut,
			WriteTimeout: defaultWriteTimeOut,
			IdleTimeout:  defaultIdlerTimeOut,
		},
		logger: logger,
	}
}

// StartServer запускает HTTP-сервер и блокирует выполнение до его остановки.
//
// Параметры:
//   - ctx: контекст для управления жизненным циклом сервера.
//
// Описание работы метода:
//   - Логируем запуск сервера;
//   - Создаем канал для получения ошибок из горутины;
//   - Запускаем сервер в отдельной горутине, чтобы не блокировать основной поток;
//   - Блокируем выполнение до остановки сервера или ошибки;
//   - Закрываем канал;
//   - Ожидаем либо отмены контекста, либо ошибки от сервера.
func (s *Server) StartServer(ctx context.Context) error {
	s.logger.LogIf(logg.SkipNofS, "%s: %s", msg.IS2000, os.Getenv("PORT"))

	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.LogEf(logg.SkipNofS, err, "%s: %s", msg.ES5000, os.Getenv("PORT"))

			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// StopServer выполняет корректное завершение работы сервера.
//
// Параметры:
//   - ctx: контекст для управления жизненным циклом сервера.
//
// Описание работы метода:
//   - Ожидает завершение всех активных запросов в течении заданного таймаута;
//   - Потом останавливает сервер, избегая обрыва соединений и потери данных.
func (s *Server) StopServer(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.LogEf(logg.SkipNofS, err, "%s: %s", msg.ES5001, os.Getenv("PORT"))

		return err
	}

	s.logger.LogIf(logg.SkipNofS, "%s: %s", msg.IS2002, os.Getenv("PORT"))

	return nil
}
