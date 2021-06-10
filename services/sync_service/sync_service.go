package sync_service

type SyncService interface {
	PushUpload(folderName string, uploadFolderId string, speedUpload bool, deleteExtraFiles bool)
	PullDownload(folderName string, downloadFolderId string)
}
