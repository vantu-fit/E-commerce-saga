package media

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/vantu-fit/saga-pattern/cmd/media/config"
)

type File struct {
	Name      string
	Data      io.Reader
	Bucket    string
	ProductID uuid.UUID
}

type Media interface {
	UploadObject(ctx context.Context, file *File) error
	GetOjectInfo(ctx context.Context, file *File) (*minio.ObjectInfo, error)
	GetUrl(ctx context.Context, file *File) (string, error)
	DeleteObject(ctx context.Context, file *File) error
	GetConfig() *config.Config
}

type media struct {
	config *config.Config
	minio  *minio.Client
}

func New(cfg *config.Config, minio *minio.Client) Media {
	return &media{
		config: cfg,
		minio:  minio,
	}
}

func (m *media) UploadObject(ctx context.Context, file *File) error {
	metadata := map[string]string{
		"product_id": file.ProductID.String(),
	}
	_, err := m.minio.PutObject(ctx, file.Bucket, file.Name, file.Data, -1, minio.PutObjectOptions{
		UserMetadata: metadata,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *media) GetOjectInfo(ctx context.Context, file *File) (*minio.ObjectInfo, error) {
	info , err := m.minio.StatObject(ctx, file.Bucket, file.Name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &info , nil
}

func (m *media) GetUrl(ctx context.Context, file *File) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	reqInfo, err := m.minio.PresignedGetObject(ctx, file.Bucket, file.Name, time.Hour*24*7, reqParams)
	if err != nil {
		return "", err
	}
	return reqInfo.String(), nil
}

func (m *media) DeleteObject(ctx context.Context, file *File) error {
	err := m.minio.RemoveObject(ctx, file.Bucket, file.Name , minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m *media) GetConfig() *config.Config {
	return m.config
}