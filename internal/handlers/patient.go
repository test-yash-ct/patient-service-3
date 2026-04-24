package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/auth"
	"github.com/healthops/patient-service/internal/store"
)

type PatientAPI struct {
	Store  *store.PatientStore
	Secret string
}

func (a *PatientAPI) Register(r *gin.RouterGroup) {
	r.GET("/patients/:id", a.Get)
}

func (a *PatientAPI) Get(c *gin.Context) {
	_, err := auth.ParseBearer(c.GetHeader("Authorization"), a.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	p, err := a.Store.GetByID(c.Request.Context(), c.GetHeader("X-Tenant-ID"), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	RecordLastPatient(p)
	log.Printf("patient_access patient_id=%s tenant_header=%s record_tenant=%s", id, c.GetHeader("X-Tenant-ID"), p.TenantID)
	c.JSON(http.StatusOK, p)
}
