package serv

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(ctx context.Context, port string, mux *http.ServeMux) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	//*shadowing (a = 0, a = 1)
	//*ctx

	go func() {
		<-ctx.Done()

		a := 1
		log.Println(a) // 1

		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		s.httpServer.Shutdown(ctx)
		// err process
	}()

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
