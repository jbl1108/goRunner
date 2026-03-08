package delivery

import (
	"io"
	"log"
	"net/http"

	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/datamodel"
	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/ports/input"
)

type KeyValueRestService struct {
	keyValueHandlingUsecase input.KeyValInputPort
	address                 string
}

func NewKeyValueRestService(address string, keyValueHandlingUsecase input.KeyValInputPort) *KeyValueRestService {
	return &KeyValueRestService{keyValueHandlingUsecase: keyValueHandlingUsecase, address: address}
}

func (s *KeyValueRestService) Start() error {
	mux := http.NewServeMux()
	log.Printf("Starting Key Value REST Service on %s", s.address)
	s.RegisterRoutes(mux)
	return http.ListenAndServe(s.address, mux)
}

func (s *KeyValueRestService) handleGetKey(w http.ResponseWriter, r *http.Request) {
	result, err := s.keyValueHandlingUsecase.GetKey(r.PathValue("topic") + ":" + r.PathValue("key"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(result)
}

func (s *KeyValueRestService) handleSetKey(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Printf("Url %v", r.URL)
	if r.PathValue("topic") == "" || r.PathValue("key") == "" {
		http.Error(w, "Missing topic or key in URL", http.StatusBadRequest)
		return
	}

	message := datamodel.Message{
		Topic: r.PathValue("topic"),
		Data: datamodel.KeyValue{
			Key:   r.PathValue("key"),
			Value: string(body),
		},
	}
	err = s.keyValueHandlingUsecase.SetKey(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *KeyValueRestService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the KeyValue REST Service"))
		w.Write([]byte("Use one of the following endpoints"))
		w.Write([]byte("\nGET /key/{topic}/{key} - Retrieve a key value by key"))
		w.Write([]byte("\nPOST /key/{topic}/{key} - Set a key value by key"))
		w.Write([]byte("\n/health/ - Health check endpoint"))
	})
	mux.HandleFunc("GET /key/{topic}/{key}", s.handleGetKey)
	mux.HandleFunc("POST /key/{topic}/{key}", s.handleSetKey)
	mux.HandleFunc("/health/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

}
