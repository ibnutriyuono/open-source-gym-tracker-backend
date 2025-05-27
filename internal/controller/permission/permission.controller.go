package permission

import (
	"caloria-backend/internal/helper/response"
	"net/http"

	"gorm.io/gorm"
)

type PermissionController struct {
	DB *gorm.DB
}

func (uc *PermissionController) FindAll(w http.ResponseWriter, r *http.Request) {
	response.SendJSON(w, http.StatusOK, nil, "Initial Commit Find All")
}

func (uc *PermissionController) Create(w http.ResponseWriter, r *http.Request) {
	response.SendJSON(w, http.StatusOK, nil, "Initial Commit Create")
}

func (uc *PermissionController) Update(w http.ResponseWriter, r *http.Request) {
	response.SendJSON(w, http.StatusOK, nil, "Initial Commit Update")
}

func (uc *PermissionController) Delete(w http.ResponseWriter, r *http.Request) {
	response.SendJSON(w, http.StatusOK, nil, "Initial Commit Delete")
}

func (uc *PermissionController) FindById(w http.ResponseWriter, r *http.Request) {
	response.SendJSON(w, http.StatusOK, nil, "Initial Commit Find By Id")
}