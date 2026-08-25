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
	claims, err := auth.ParseBearer(c.GetHeader("Authorization"), a.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tenantID := claims.Tenant
	if headerTenant := c.GetHeader("X-Tenant-ID"); headerTenant != "" {
		if tenantID != "" && headerTenant != tenantID {
			c.JSON(http.StatusForbidden, gin.H{"error": "tenant_mismatch"})
			return
		}
		if tenantID == "" {
			tenantID = headerTenant
		}
	}
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_tenant"})
		return
	}

	id := c.Param("id")
	p, err := a.Store.GetByID(c.Request.Context(), tenantID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	log.Printf("patient_access patient_id=%s tenant=%s subject=%s", id, tenantID, claims.Sub)
	c.JSON(http.StatusOK, p)
}
