package ignore_service

type IgnoreService interface {
	LoadRules(baseDir string) FileNameGuard
}

type FileNameGuard interface {
	ShouldIgnore(relativeDir string, fileName string) (*entry, bool)
}
