package drive_api

func (api *driveApiImpl) GetCredentialEmail() (string, error) {
	return api.DriveApi.GetCredentialEmail()
}
func (api *driveApiImpl) SetDomainAdminAccess(useDomainAdminAccess bool) {
	api.DriveApi.SetDomainAdminAccess(useDomainAdminAccess)
}

func (api *driveApiImpl) GetDomainAdminAccess() bool {
	return api.DriveApi.GetDomainAdminAccess()
}

func (api *driveApiImpl) SetImpersonate(impersonate string) {
	api.DriveApi.SetImpersonate(impersonate)
}

func (api *driveApiImpl) GetImpersonate() string {
	return api.DriveApi.GetImpersonate()
}
