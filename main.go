package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/headblockhead/landmine/backend/airtable"
	"github.com/headblockhead/landmine/handlers/createrecords"
	"github.com/headblockhead/landmine/handlers/deleterecords"
	"github.com/headblockhead/landmine/handlers/listrecords"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Landmine struct {
	EtcdClient      *clientv3.Client
	CacheExpiryTime int
	Port            int
	Mux             *http.ServeMux
}

func NewLandmine(log *slog.Logger, client *clientv3.Client) *Landmine {
	mux := http.NewServeMux()
	airtableClient := airtable.New(log, http.DefaultClient, os.Getenv("AIRTABLE_API_KEY"))
	lrecs := listrecords.New(log, airtableClient)
	mux.Handle("GET /{baseID}/{tableIDOrName}", lrecs)
	crecs := createrecords.New(log, airtableClient)
	mux.Handle("POST /{baseID}/{tableIDOrName}", crecs)
	drecs := deleterecords.New(log, airtableClient)
	mux.Handle("DELETE /{baseID}/{tableIDOrName}", drecs)
	return &Landmine{
		Mux:        mux,
		EtcdClient: client,
	}
}

func (l *Landmine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.Mux.ServeHTTP(w, r)
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // TODO: use Kong for this
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err) // TODO: don't panic
	}
	l := NewLandmine(log, client)
	s := &http.Server{
		Addr:    ":8080", // TODO: take in listen addr
		Handler: l,
	}
	s.ListenAndServe()
}
