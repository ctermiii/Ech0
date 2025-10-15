package fediverse

import (
	echoRepository "github.com/lin-snow/ech0/internal/repository/echo"
	repository "github.com/lin-snow/ech0/internal/repository/fediverse"
	keyvalueRepository "github.com/lin-snow/ech0/internal/repository/keyvalue"
	userRepository "github.com/lin-snow/ech0/internal/repository/user"
)

// FediverseCore 联邦宇宙核心
type FediverseCore struct {
	repo           repository.FediverseRepositoryInterface
	userRepository userRepository.UserRepositoryInterface
	echoRepository echoRepository.EchoRepositoryInterface
	keyvalueRepo   keyvalueRepository.KeyValueRepositoryInterface
}

func NewFediverseCore(
	repo repository.FediverseRepositoryInterface,
	keyvalueRepo keyvalueRepository.KeyValueRepositoryInterface,
	userRepository userRepository.UserRepositoryInterface,
	echoRepository echoRepository.EchoRepositoryInterface,
) *FediverseCore {
	return &FediverseCore{
		repo:           repo,
		keyvalueRepo:   keyvalueRepo,
		userRepository: userRepository,
		echoRepository: echoRepository,
	}
}
