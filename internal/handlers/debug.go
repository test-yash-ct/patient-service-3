package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/models"
)

var lastPatient sync.Map

func RecordLastPatient(p models.Patient) {
	lastPatient.Store("last", p)
}

func DebugLastPatient(c *gin.Context) {
	v, ok := lastPatient.Load("last")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "no_cached_patient"})
		return
	}
	c.JSON(http.StatusOK, v)
}
