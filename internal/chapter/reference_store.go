package chapter

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/sigpanic/goink/internal/character"
	"github.com/sigpanic/goink/internal/reader"
	"github.com/sigpanic/goink/internal/storyarc"
	"github.com/sigpanic/goink/internal/timeline"
)

// ChapterReference 是阻止删除章节的一条跨领域引用。
type ChapterReference struct {
	Kind  string `json:"kind"`
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// ReferenceStore 查询指向章节的跨领域引用。
type ReferenceStore struct {
	db *gorm.DB
}

// NewReferenceStore 创建章节引用查询存储。
func NewReferenceStore(db *gorm.DB) *ReferenceStore {
	return &ReferenceStore{db: db}
}

// ReferencesByChapter 返回当前小说中指向指定章节的全部阻塞引用。
func (s *ReferenceStore) ReferencesByChapter(ctx context.Context, novelID, chapterID int64) ([]ChapterReference, error) {
	var references []ChapterReference
	var entries []timeline.TimelineEntry
	if err := s.db.WithContext(ctx).Where("novel_id = ? AND (source_chapter_id = ? OR resolved_chapter_id = ?)", novelID, chapterID, chapterID).Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("query timeline references: %w", err)
	}
	for _, entry := range entries {
		if entry.SourceChapterID != nil && *entry.SourceChapterID == chapterID {
			references = append(references, ChapterReference{Kind: "timeline_source", ID: entry.ID, Label: entry.Title})
		}
		if entry.ResolvedChapterID != nil && *entry.ResolvedChapterID == chapterID {
			references = append(references, ChapterReference{Kind: "timeline_resolved", ID: entry.ID, Label: entry.Title})
		}
	}

	var nodes []storyarc.ArcNode
	if err := s.db.WithContext(ctx).Where("novel_id = ? AND actual_chapter_id = ?", novelID, chapterID).Find(&nodes).Error; err != nil {
		return nil, fmt.Errorf("query story arc references: %w", err)
	}
	for _, node := range nodes {
		references = append(references, ChapterReference{Kind: "story_arc", ID: node.ID, Label: node.Title})
	}

	var perspectives []reader.ReaderPerspective
	if err := s.db.WithContext(ctx).Where("novel_id = ? AND (planted_chapter_id = ? OR revealed_chapter_id = ?)", novelID, chapterID, chapterID).Find(&perspectives).Error; err != nil {
		return nil, fmt.Errorf("query reader references: %w", err)
	}
	for _, item := range perspectives {
		if item.PlantedChapterID != nil && *item.PlantedChapterID == chapterID {
			references = append(references, ChapterReference{Kind: "reader_planted", ID: item.ID, Label: item.Content})
		}
		if item.RevealedChapterID != nil && *item.RevealedChapterID == chapterID {
			references = append(references, ChapterReference{Kind: "reader_revealed", ID: item.ID, Label: item.Content})
		}
	}

	var relations []character.CharacterRelation
	if err := s.db.WithContext(ctx).Where("novel_id = ? AND chapter_id = ?", novelID, chapterID).Find(&relations).Error; err != nil {
		return nil, fmt.Errorf("query character relation references: %w", err)
	}
	for _, relation := range relations {
		references = append(references, ChapterReference{Kind: "character_relation", ID: relation.ID, Label: relation.RelationDescribe})
	}
	return references, nil
}
