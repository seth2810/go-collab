package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seth2810/go-collab/internal/config"
	"github.com/seth2810/go-collab/internal/logger"
	"go.uber.org/zap"
	"gopkg.in/olahol/melody.v1"
)

func ServeContext(ctx context.Context, cfg *config.Config) error {
	log, err := logger.New(&cfg.Logger)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	addr := net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: createHandler(log),
	}

	errCh := make(chan error, 1)

	go func() {
		log.Info("server listen", zap.String("address", addr))

		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errCh:
	}

	if !errors.Is(err, context.Canceled) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFn()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}

func createHandler(log *zap.Logger) http.Handler {
	var usersCount uint64
	var docBodyMx sync.RWMutex
	var docBody json.RawMessage

	ws := melody.New()
	ws.Config.MaxMessageSize = 0

	ws.HandleConnect(func(s *melody.Session) {
		var document json.RawMessage

		atomic.AddUint64(&usersCount, 1)

		docBodyMx.RLock()
		document = docBody
		docBodyMx.RUnlock()

		loadDocumentMessageBody, _ := json.Marshal(&Message{
			Type:    MessageTypeLoadDocument,
			Payload: document,
		})

		_ = ws.BroadcastFilter(loadDocumentMessageBody, func(q *melody.Session) bool {
			return q == s
		})

		userCountJSON, _ := json.Marshal(&usersCount)

		updateUserCountMessageBody, _ := json.Marshal(&Message{
			Type:    MessageTypeUpdateUserCount,
			Payload: userCountJSON,
		})

		_ = ws.Broadcast(updateUserCountMessageBody)
	})

	ws.HandleDisconnect(func(s *melody.Session) {
		atomic.AddUint64(&usersCount, ^uint64(0))

		userCountJSON, _ := json.Marshal(&usersCount)

		updateUserCountMessageBody, _ := json.Marshal(&Message{
			Type:    MessageTypeUpdateUserCount,
			Payload: userCountJSON,
		})

		_ = ws.Broadcast(updateUserCountMessageBody)
	})

	ws.HandleMessage(func(s *melody.Session, messageJSON []byte) {
		var msg Message

		_ = json.Unmarshal(messageJSON, &msg)

		log.Debug("handle message", zap.String("type", string(msg.Type)), zap.ByteString("payload", msg.Payload))

		switch msg.Type {
		case MessageTypeUpdateDocument:
			_ = ws.BroadcastOthers(messageJSON, s)
		case MessageTypeSaveDocument:
			docBodyMx.Lock()
			docBody = msg.Payload
			docBodyMx.Unlock()
		}
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := ws.HandleRequest(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return mux
}
