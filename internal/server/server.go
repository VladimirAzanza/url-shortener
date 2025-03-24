package server

import (
	"context"
	"net/http"

	"github.com/VladimirAzanza/url-shortener/internal/controller"
	"go.uber.org/fx"
)

func NewServer(urlController *controller.URLController) *http.Server {
	return &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				urlController.HandlePost(w, r)
			case http.MethodGet:
				urlController.HandleGet(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusBadRequest)
			}
		}),
	}
}

func StartServer(lc fx.Lifecycle, server *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

}
