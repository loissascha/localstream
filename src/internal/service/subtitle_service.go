package service

import (
	"context"

	"github.com/loissascha/localstream/internal/provider"
)

type SubtitleService struct {
	subtitleProvider provider.SubtitleProvider
}

func NewSubtitleService(
	subtitleProvider provider.SubtitleProvider,
) *SubtitleService {
	return &SubtitleService{
		subtitleProvider: subtitleProvider,
	}
}

func (s *SubtitleService) SupportedLanguages(ctx context.Context) ([]provider.SubtitleSupportedLanguage, error) {
	return s.subtitleProvider.SupportedLanguages(ctx)
}
