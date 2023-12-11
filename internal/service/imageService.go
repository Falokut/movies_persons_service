package service

import (
	"context"
	"net/url"

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

	u, err := url.Parse(s.cfg.BasePhotoUrl)
	if err != nil {
		s.logger.Errorf("can't parse url. error: %s", err.Error())
		return ""
	}

	q := u.Query()
	q.Add("image_id", pictureID)
	q.Add("category", s.cfg.PicturesCategory)
	u.RawQuery = q.Encode()
	return u.String()
}
