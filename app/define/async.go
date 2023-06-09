package define

import "asset-management/app/model"

type CreateAsyncTaskReq struct {
	DownloadLink string           `json:"download_link"`
	Type         uint             `json:"type" binding:"oneof=0 1 2"`
	LogType      uint             `json:"log_type" binding:"oneof=0 1 2"` // 0-all 1-success 2-failed
	ObjectKey    string           `json:"object_key"`
	DepartmentID uint             `json:"department_id"`
	EntityID     uint             `json:"entity_id"`
	FromTime     *model.ModelTime `json:"from_time" binding:"omitempty"`
}

type ModifyAsyncTaskStateReq struct {
	State uint `json:"state" binding:"oneof=0 4"`
}

type AsyncTaskInfo struct {
	ID           uint   `json:"async_id"`
	Type         uint   `json:"type"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	State        uint   `json:"state"`
	DownloadLink string `json:"download_link"`
	Message      string `json:"message"`
	LogType      uint   `json:"log_type"`
}

type AsyncTaskListResponse struct {
	AsyncList []AsyncTaskInfo `json:"async_list"`
}
