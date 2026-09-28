package app

import (
	"errors"
	"testing"

	"github.com/sigpanic/goink/internal/volume"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVolumeManagementAPI(t *testing.T) {
	app := setupTestApp(t)
	novel := createTestNovel(t, app)

	first, err := app.CreateVolume(novel.ID, "第一卷")
	require.NoError(t, err)
	second, err := app.CreateVolume(novel.ID, "第二卷")
	require.NoError(t, err)

	volumes, err := app.GetVolumes(novel.ID)
	require.NoError(t, err)
	require.Len(t, volumes, 2)
	assert.Equal(t, []int64{first.ID, second.ID}, []int64{volumes[0].ID, volumes[1].ID})

	require.NoError(t, app.UpdateVolume(novel.ID, second.ID, "终卷"))
	require.NoError(t, app.ReorderVolumes(novel.ID, []int64{second.ID, first.ID}))

	volumes, err = app.GetVolumes(novel.ID)
	require.NoError(t, err)
	require.Len(t, volumes, 2)
	assert.Equal(t, "终卷", volumes[0].Name)
	assert.Equal(t, []int64{second.ID, first.ID}, []int64{volumes[0].ID, volumes[1].ID})

	require.NoError(t, app.DeleteVolume(novel.ID, first.ID))
	volumes, err = app.GetVolumes(novel.ID)
	require.NoError(t, err)
	require.Len(t, volumes, 1)
	assert.Equal(t, second.ID, volumes[0].ID)
}

func TestDeleteVolumeRejectsNonEmptyOrForeignVolume(t *testing.T) {
	app := setupTestApp(t)
	novel := createTestNovel(t, app)
	otherNovel := createTestNovel(t, app)

	v, err := app.CreateVolume(novel.ID, "第一卷")
	require.NoError(t, err)
	_, err = app.CreateChapter(CreateChapterInput{NovelID: novel.ID, Title: "卷内章节"})
	require.NoError(t, err)

	err = app.DeleteVolume(novel.ID, v.ID)
	assert.ErrorIs(t, err, volume.ErrHasChapters)

	err = app.UpdateVolume(otherNovel.ID, v.ID, "越权改名")
	assert.True(t, errors.Is(err, volume.ErrNotFound), "err = %v, want ErrNotFound", err)
}
