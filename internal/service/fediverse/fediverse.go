package service

import (
	echoRepository "github.com/lin-snow/ech0/internal/repository/echo"
	repository "github.com/lin-snow/ech0/internal/repository/fediverse"
	userRepository "github.com/lin-snow/ech0/internal/repository/user"
	commonService "github.com/lin-snow/ech0/internal/service/common"
	settingService "github.com/lin-snow/ech0/internal/service/setting"
	"github.com/lin-snow/ech0/internal/transaction"
)

type FediverseService struct {
	txManager           transaction.TransactionManager
	fediverseRepository repository.FediverseRepositoryInterface
	userRepository      userRepository.UserRepositoryInterface
	settingService      settingService.SettingServiceInterface
	echoRepository      echoRepository.EchoRepositoryInterface
	commonService       commonService.CommonServiceInterface
}

func NewFediverseService(
	txManager transaction.TransactionManager,
	fediverseRepository repository.FediverseRepositoryInterface,
	userRepository userRepository.UserRepositoryInterface,
	settingService settingService.SettingServiceInterface,
	echoRepository echoRepository.EchoRepositoryInterface,
	commonService commonService.CommonServiceInterface,
) FediverseServiceInterface {
	return &FediverseService{
		txManager:           txManager,
		fediverseRepository: fediverseRepository,
		userRepository:      userRepository,
		settingService:      settingService,
		echoRepository:      echoRepository,
		commonService:       commonService,
	}
}
