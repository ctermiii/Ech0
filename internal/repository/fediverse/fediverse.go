package repository

import (
	"context"
	"errors"

	model "github.com/lin-snow/ech0/internal/model/fediverse"
	"github.com/lin-snow/ech0/internal/transaction"
	"gorm.io/gorm"
)

type FediverseRepository struct {
	db func() *gorm.DB
}

func NewFediverseRepository(dbProvider func() *gorm.DB) FediverseRepositoryInterface {
	return &FediverseRepository{
		db: dbProvider,
	}
}

func (r *FediverseRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB); ok {
		return tx
	}
	return r.db()
}

func (r *FediverseRepository) GetFollowers(userID uint) ([]model.Follower, error) {
	var followers []model.Follower
	if err := r.db().Where("user_id = ?", userID).Find(&followers).Error; err != nil {
		return nil, err
	}
	return followers, nil
}

func (r *FediverseRepository) GetFollowing(userID uint) ([]model.Follow, error) {
	var following []model.Follow
	if err := r.db().Where("user_id = ?", userID).Find(&following).Error; err != nil {
		return nil, err
	}
	return following, nil
}

func (r *FediverseRepository) SaveFollower(ctx context.Context, follower *model.Follower) error {
	return r.getDB(ctx).Create(follower).Error
}

func (r *FediverseRepository) FollowerExists(ctx context.Context, userID uint, actor string) (bool, error) {
	var count int64
	if err := r.getDB(ctx).Model(&model.Follower{}).Where("user_id = ? AND actor_id = ?", userID, actor).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FediverseRepository) SaveOrUpdateFollow(ctx context.Context, follow *model.Follow) error {
	if follow == nil {
		return errors.New("follow is nil")
	}

	db := r.getDB(ctx)
	var existing model.Follow
	err := db.Where("user_id = ? AND object_id = ?", follow.UserID, follow.ObjectID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(follow).Error
	}
	if err != nil {
		return err
	}

	existing.ActorID = follow.ActorID
	existing.ActivityID = follow.ActivityID
	existing.Status = follow.Status
	return db.Save(&existing).Error
}

func (r *FediverseRepository) GetFollowByUserAndObject(ctx context.Context, userID uint, objectID string) (*model.Follow, error) {
	var follow model.Follow
	err := r.getDB(ctx).Where("user_id = ? AND object_id = ?", userID, objectID).First(&follow).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &follow, nil
}

func (r *FediverseRepository) DeleteFollow(ctx context.Context, followID uint) error {
	return r.getDB(ctx).Delete(&model.Follow{}, followID).Error
}
