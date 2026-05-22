package dto

type TaskStatsResponse struct {
	AsAssignee int `json:"as_assignee"`
	AsWatcher  int `json:"as_watcher"`
}
