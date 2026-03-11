package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jbl1108/goRunner/usecases/datamodel"
	"github.com/jbl1108/goRunner/usecases/ports/input"
)

type TrainingRestService struct {
	trainingHandlingUsecase input.TrainingInputPort
	address                 string
}

func NewTrainingRestService(address string, trainingHandlingUsecase input.TrainingInputPort) *TrainingRestService {
	return &TrainingRestService{trainingHandlingUsecase: trainingHandlingUsecase, address: address}
}

func (s *TrainingRestService) Start() error {
	mux := http.NewServeMux()
	log.Printf("Starting Training REST Service on %s", s.address)
	s.RegisterRoutes(mux)
	return http.ListenAndServe(s.address, mux)
}

func (s *TrainingRestService) handleGetTraining(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	training, err := s.trainingHandlingUsecase.GetTraining(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(training)

}

func (s *TrainingRestService) handleGetAllTrainings(w http.ResponseWriter, r *http.Request) {
	trainings, err := s.trainingHandlingUsecase.GetAllTrainings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trainings)
}

func (s *TrainingRestService) handlePostTraining(w http.ResponseWriter, r *http.Request) {
	var training datamodel.Training
	err := json.NewDecoder(r.Body).Decode(&training)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	addedTraining, err := s.trainingHandlingUsecase.AddTraining(training)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(addedTraining)
}

func (s *TrainingRestService) handlePutTraining(w http.ResponseWriter, r *http.Request) {
	var training datamodel.Training
	err := json.NewDecoder(r.Body).Decode(&training)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	training.Uid = r.PathValue("uid")
	updatedTraining, err := s.trainingHandlingUsecase.UpdateTraining(training)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTraining)
}

func (s *TrainingRestService) handleDeleteTraining(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	err := s.trainingHandlingUsecase.DeleteTraining(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *TrainingRestService) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Runner REST Service"))
		w.Write([]byte("\nUse one of the following endpoints"))
		w.Write([]byte("\n/health/ - Health check endpoint"))
		w.Write([]byte("\nGET /training/{uid} - Retrieve a training by UID"))
		w.Write([]byte("\nPOST /training - Create a new training"))
		w.Write([]byte("\nPUT /training/{uid} - Update an existing training"))
		w.Write([]byte("\nDELETE /training/{uid} - Delete a training"))
		w.Write([]byte("\nGET /training - Retrieve all trainings"))

	})
	mux.HandleFunc("GET /training/{uid}", s.handleGetTraining)
	mux.HandleFunc("POST /training", s.handlePostTraining)
	mux.HandleFunc("PUT /training/{uid}", s.handlePutTraining)
	mux.HandleFunc("DELETE /training/{uid}", s.handleDeleteTraining)
	mux.HandleFunc("GET /training", s.handleGetAllTrainings)
	mux.HandleFunc("/health/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

}
