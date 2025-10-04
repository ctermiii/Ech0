package service

import (
	"errors"

	commonModel "github.com/lin-snow/ech0/internal/model/common"
	model "github.com/lin-snow/ech0/internal/model/fediverse"
)

// GetActorByUsername 通过用户名获取 Actor 信息
func (fediverseService *FediverseService) GetActorByUsername(username string) (model.Actor, error) {
	// 查询用户
	user, err := fediverseService.userRepository.GetUserByUsername(username)
	if err != nil {
		return model.Actor{}, errors.New(commonModel.USER_NOTFOUND)
	}

	// 构建 Actor 对象
	actor, _, err := fediverseService.BuildActor(&user)
	if err != nil {
		return model.Actor{}, err
	}

	return actor, nil
}
