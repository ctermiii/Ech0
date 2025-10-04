package service

import model "github.com/lin-snow/ech0/internal/model/fediverse"

// GetObjectByID 通过 ID 获取内容对象
func (fediverseService *FediverseService) GetObjectByID(id uint) (model.Object, error) {
	// 获取 Echo
	echo, err := fediverseService.echoRepository.GetEchosById(id)
	if err != nil || echo.Private {
		return model.Object{}, err
	}

	// 获取 Actor 和 setting
	user, err := fediverseService.userRepository.GetUserByUsername(echo.Username)
	if err != nil {
		return model.Object{}, err
	}
	actor, setting, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.Object{}, err
	}
	serverURL, err := normalizeServerURL(setting.ServerURL)
	if err != nil {
		return model.Object{}, err
	}

	// 转 Object
	return fediverseService.ConvertEchoToObject(echo, &actor, serverURL), nil
}
