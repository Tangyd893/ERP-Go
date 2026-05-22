package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/domain"
	sharedErrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/response"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	appService   *app.FileAppService
	fallbackMode bool
}

func NewFileHandler(appService *app.FileAppService) *FileHandler {
	return &FileHandler{appService: appService, fallbackMode: appService == nil}
}

func (h *FileHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/upload", h.upload)
	router.GET("/download/:id", h.download)
}

func (h *FileHandler) upload(c *gin.Context) {
	if h.fallbackMode { c.JSON(http.StatusOK, gin.H{"code": 0, "message": "接口已联通，对象存储待接入"}); return }

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, sharedErrors.CodeInvalidParameter, "未接收到文件")
		return
	}
	defer file.Close()

	f := &domain.File{
		ID: fmt.Sprintf("FL%d", time.Now().UnixNano()), TenantID: c.GetString("tenant_id"),
		FileName: header.Filename, FileSize: header.Size, MimeType: header.Header.Get("Content-Type"),
		CreatedBy: c.GetString("username"), CreatedAt: time.Now(),
	}
	if err := h.appService.Upload(c.Request.Context(), f); err != nil {
		response.Error(c, http.StatusInternalServerError, sharedErrors.CodeInternalError, "上传失败: "+err.Error())
		return
	}
	response.Success(c, f)
}

func (h *FileHandler) download(c *gin.Context) {
	if h.fallbackMode {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeNotFound, "文件不存在或对象存储未接入")
		return
	}
	f, err := h.appService.Download(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, sharedErrors.CodeNotFound, "文件不存在")
		return
	}
	response.Success(c, f)
}
