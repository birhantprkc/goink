package app

import (
	"fmt"

	"github.com/sigpanic/goink/internal/volume"
)

// CreateVolume 为指定小说新建卷，追加到现有卷末尾。
func (a *App) CreateVolume(novelID int64, name string) (*volume.Volume, error) {
	v, err := a.volume.Create(a.ctx, nil, novelID, name)
	if err != nil {
		return nil, fmt.Errorf("create volume: %w", err)
	}
	return v, nil
}

// UpdateVolume 重命名指定小说中的卷。
func (a *App) UpdateVolume(novelID, volumeID int64, name string) error {
	if err := a.volume.Update(a.ctx, nil, novelID, volumeID, name); err != nil {
		return fmt.Errorf("update volume: %w", err)
	}
	return nil
}

// DeleteVolume 删除指定小说中的空卷。
func (a *App) DeleteVolume(novelID, volumeID int64) error {
	if err := a.volume.Delete(a.ctx, nil, novelID, volumeID); err != nil {
		return fmt.Errorf("delete volume: %w", err)
	}
	return nil
}

// GetVolumes 返回指定小说的全部卷，按阅读顺序排列。
func (a *App) GetVolumes(novelID int64) ([]volume.Volume, error) {
	volumes, err := a.volume.ListByNovel(a.ctx, nil, novelID)
	if err != nil {
		return nil, fmt.Errorf("list volumes: %w", err)
	}
	return volumes, nil
}

// ReorderVolumes 按传入的完整卷 ID 顺序重排指定小说的卷。
func (a *App) ReorderVolumes(novelID int64, volumeIDs []int64) error {
	if err := a.volume.Reorder(a.ctx, nil, novelID, volumeIDs); err != nil {
		return fmt.Errorf("reorder volumes: %w", err)
	}
	return nil
}
