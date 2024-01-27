package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type ImageServiceConfig struct {
	BasePhotoUrl     string
	PicturesCategory string
}

type imageService struct {
	cfg    ImageServiceConfig
	logger *logrus.Logger
}

func NewImageService(cfg ImageServiceConfig, logger *logrus.Logger) *imageService {
	return &imageService{
		cfg:    cfg,
		logger: logger,
	}
}

// Returns picture url for GET request
func (s *imageService) GetPictureURL(ctx context.Context, pictureID string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "imageService.GetPictureURL")
	defer span.Finish()

	if pictureID == "" {
		return ""
	}

	return s.cfg.BasePhotoUrl + "/" + s.cfg.PicturesCategory + "/" + pictureID
}
