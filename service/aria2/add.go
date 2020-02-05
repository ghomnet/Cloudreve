package aria2

import (
	model "github.com/HFO4/cloudreve/models"
	"github.com/HFO4/cloudreve/pkg/aria2"
	"github.com/HFO4/cloudreve/pkg/filesystem"
	"github.com/HFO4/cloudreve/pkg/serializer"
	"github.com/gin-gonic/gin"
)

// AddURLService 添加URL离线下载服务
type AddURLService struct {
	URL string `json:"url" binding:"required"`
	Dst string `json:"dst" binding:"required,min=1"`
}

// Add 创建新的链接离线下载任务
func (service *AddURLService) Add(c *gin.Context, taskType int) serializer.Response {
	// 创建文件系统
	fs, err := filesystem.NewFileSystemFromContext(c)
	if err != nil {
		return serializer.Err(serializer.CodePolicyNotAllowed, err.Error(), err)
	}
	defer fs.Recycle()

	// 检查用户组权限
	if !fs.User.Group.OptionsSerialized.Aria2 {
		return serializer.Err(serializer.CodeGroupNotAllowed, "当前用户组无法进行此操作", nil)
	}

	// 存放目录是否存在
	if exist, _ := fs.IsPathExist(service.Dst); !exist {
		return serializer.Err(serializer.CodeNotFound, "存放路径不存在", nil)
	}

	// 创建任务
	task := &model.Download{
		Status: aria2.Ready,
		Type:   taskType,
		Dst:    service.Dst,
		UserID: fs.User.ID,
		Source: service.URL,
	}
	if err := aria2.Instance.CreateTask(task); err != nil {
		return serializer.Err(serializer.CodeNotSet, "任务创建失败", err)
	}

	return serializer.Response{}
}
