package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/healthops/patient-service/internal/auth"
	"github.com/healthops/patient-service/internal/obs"
	"github.com/healthops/patient-service/internal/service"
)

type PatientAPI struct {
	Patients *service.Patients
	Secret   string
}

func (a *PatientAPI) Register(r *gin.RouterGroup) {
	r.GET("/patients/:id", a.Get)
}

func (a *PatientAPI) Get(c *gin.Context) {
	claims, err := auth.ParseBearer(c.Request.Context(), c.GetHeader("Authorization"), a.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	p, err := a.Patients.Get(c.Request.Context(), claims.Tenant, c.GetHeader(obs.HeaderTenantID), c.Param("id"), claims.Sub)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "tenant_mismatch"})
	case errors.Is(err, service.ErrMissingTenant):
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_tenant"})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
	case errors.Is(err, service.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
	}
}
