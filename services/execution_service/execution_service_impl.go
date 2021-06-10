package execution_service

type commandExecutionServiceImpl struct {
}

func (service *commandExecutionServiceImpl) Executor(command string) CommandExecutor {
	return newExecutor(command)
}

func New() CommandExecutionService {
	return &commandExecutionServiceImpl{}
}
