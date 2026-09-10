package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/healthops/patient-service/internal/events"
	"github.com/healthops/patient-service/internal/models"
	"github.com/healthops/patient-service/internal/obs"
	"github.com/jackc/pgx/v5"
)

type PatientRepository interface {
	GetByID(ctx context.Context, tenantID, id string) (models.Patient, error)
	ListByTenant(ctx context.Context, tenantID string) ([]models.Patient, error)
	Update(ctx context.Context, tenantID, id, fullName, email, phone, notes string) error
}

type Patients struct {
	repo   PatientRepository
	outbox events.Outbox
}

func NewPatients(repo PatientRepository) *Patients {
	return NewPatientsWithOutbox(repo, events.Nop{})
}

func NewPatientsWithOutbox(repo PatientRepository, outbox events.Outbox) *Patients {
	if outbox == nil {
		outbox = events.Nop{}
	}
	return &Patients{repo: repo, outbox: outbox}
}

func ResolveTenant(claimTenant, headerTenant string) (string, error) {
	tenantID := claimTenant
	if headerTenant != "" {
		if tenantID != "" && headerTenant != tenantID {
			return "", ErrForbidden
		}
		if tenantID == "" {
			tenantID = headerTenant
		}
	}
	if tenantID == "" {
		return "", ErrMissingTenant
	}
	return tenantID, nil
}

func (s *Patients) Get(ctx context.Context, claimTenant, headerTenant, id, subject string) (models.Patient, error) {
	tenantID, err := ResolveTenant(claimTenant, headerTenant)
	if err != nil {
		return models.Patient{}, err
	}
	ctx = obs.WithTenant(ctx, tenantID)
	if id == "" {
		return models.Patient{}, ErrInvalidInput
	}
	p, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Patient{}, ErrNotFound
		}
		return models.Patient{}, err
	}
	logPatientEvent(obs.RequestIDFromContext(ctx), tenantID, id, subject, "patient_access")
	return p, nil
}

func (s *Patients) List(ctx context.Context, claimTenant, headerTenant string) ([]models.Patient, error) {
	tenantID, err := ResolveTenant(claimTenant, headerTenant)
	if err != nil {
		return nil, err
	}
	ctx = obs.WithTenant(ctx, tenantID)
	rows, err := s.repo.ListByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	logPatientEvent(obs.RequestIDFromContext(ctx), tenantID, "", "", "patient_list")
	return rows, nil
}

type PatientUpdate struct {
	FullName     string
	Email        string
	Phone        string
	MedicalNotes string
}

func (s *Patients) Update(ctx context.Context, claimTenant, headerTenant, id, subject string, upd PatientUpdate) error {
	tenantID, err := ResolveTenant(claimTenant, headerTenant)
	if err != nil {
		return err
	}
	ctx = obs.WithTenant(ctx, tenantID)
	if id == "" {
		return ErrInvalidInput
	}
	if err := s.repo.Update(ctx, tenantID, id, upd.FullName, upd.Email, upd.Phone, upd.MedicalNotes); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	requestID := obs.RequestIDFromContext(ctx)
	logPatientEvent(requestID, tenantID, id, subject, "patient_update")
	ev := events.New(events.TypePatientUpdated, tenantID, requestID, map[string]any{
		"patient_id": id,
	})
	if err := s.outbox.Append(ctx, ev); err != nil {
		log.Printf("outbox append failed request_id=%s event_type=%s", requestID, ev.EventType)
	}
	return nil
}

func logPatientEvent(requestID, tenant, patientID, subject, event string) {
	entry := map[string]string{
		"event":      event,
		"request_id": requestID,
		"tenant":     tenant,
		"patient_id": patientID,
		"subject":    subject,
	}
	b, _ := json.Marshal(entry)
	log.New(os.Stdout, "", 0).Println(string(b))
}
