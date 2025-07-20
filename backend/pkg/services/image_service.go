package services

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/doggyclub/backend/config"
	"github.com/doggyclub/backend/pkg/utils"
	"github.com/nfnt/resize"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ImageService struct {
	db           *gorm.DB
	redis        *redis.Client
	cfg          config.Config
	cacheService *CacheService
}

func NewImageService(db *gorm.DB, redis *redis.Client, cfg config.Config) *ImageService {
	return &ImageService{
		db:           db,
		redis:        redis,
		cfg:          cfg,
		cacheService: NewCacheService(redis, cfg),
	}
}

// Image size configurations
type ImageSize struct {
	Width   uint
	Height  uint
	Quality int
	Suffix  string
}

var ImageSizes = map[string]ImageSize{
	"thumbnail": {Width: 150, Height: 150, Quality: 80, Suffix: "_thumb"},
	"small":     {Width: 300, Height: 300, Quality: 85, Suffix: "_small"},
	"medium":    {Width: 600, Height: 600, Quality: 90, Suffix: "_medium"},
	"large":     {Width: 1200, Height: 1200, Quality: 95, Suffix: "_large"},
}

// ProcessImageRequest represents image processing request
type ProcessImageRequest struct {
	ImageData   []byte            `json:"image_data"`
	Filename    string            `json:"filename"`
	ContentType string            `json:"content_type"`
	Sizes       []string          `json:"sizes"` // Which sizes to generate
	Options     ProcessingOptions `json:"options"`
}

type ProcessingOptions struct {
	MaxWidth      uint `json:"max_width"`
	MaxHeight     uint `json:"max_height"`
	Quality       int  `json:"quality"`
	AutoOrient    bool `json:"auto_orient"`
	StripMetadata bool `json:"strip_metadata"`
}

type ProcessedImage struct {
	Size        string `json:"size"`
	URL         string `json:"url"`
	Width       uint   `json:"width"`
	Height      uint   `json:"height"`
	FileSize    int64  `json:"file_size"`
	ContentType string `json:"content_type"`
}

type ImageProcessingResult struct {
	OriginalURL string           `json:"original_url"`
	Images      []ProcessedImage `json:"images"`
	Metadata    ImageMetadata    `json:"metadata"`
}

type ImageMetadata struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Format      string `json:"format"`
	Size        int64  `json:"size"`
	AspectRatio float64 `json:"aspect_ratio"`
}

// ProcessImage processes an image and generates multiple sizes
func (s *ImageService) ProcessImage(req ProcessImageRequest) (*ImageProcessingResult, error) {
	// Validate image
	if err := s.validateImage(req.ImageData, req.ContentType); err != nil {
		return nil, err
	}

	// Decode original image
	img, format, err := image.Decode(bytes.NewReader(req.ImageData))
	if err != nil {
		return nil, utils.WrapError(err, "failed to decode image")
	}

	// Get original dimensions
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// Create metadata
	metadata := ImageMetadata{
		Width:       originalWidth,
		Height:      originalHeight,
		Format:      format,
		Size:        int64(len(req.ImageData)),
		AspectRatio: float64(originalWidth) / float64(originalHeight),
	}

	var processedImages []ProcessedImage

	// Process requested sizes
	for _, sizeKey := range req.Sizes {
		sizeConfig, exists := ImageSizes[sizeKey]
		if !exists {
			continue
		}

		// Resize image
		resizedImg := s.resizeImage(img, sizeConfig, req.Options)

		// Encode resized image
		var buf bytes.Buffer
		quality := sizeConfig.Quality
		if req.Options.Quality > 0 {
			quality = req.Options.Quality
		}

		if err := s.encodeImage(&buf, resizedImg, format, quality); err != nil {
			return nil, utils.WrapError(err, "failed to encode resized image")
		}

		// Generate filename for resized image
		resizedFilename := s.generateResizedFilename(req.Filename, sizeConfig.Suffix)

		// Upload to storage (would integrate with Cloudflare R2 or similar)
		url, err := s.uploadImage(buf.Bytes(), resizedFilename, req.ContentType)
		if err != nil {
			return nil, err
		}

		// Get dimensions of resized image
		resizedBounds := resizedImg.Bounds()
		processedImage := ProcessedImage{
			Size:        sizeKey,
			URL:         url,
			Width:       uint(resizedBounds.Dx()),
			Height:      uint(resizedBounds.Dy()),
			FileSize:    int64(buf.Len()),
			ContentType: req.ContentType,
		}

		processedImages = append(processedImages, processedImage)
	}

	// Upload original image
	originalURL, err := s.uploadImage(req.ImageData, req.Filename, req.ContentType)
	if err != nil {
		return nil, err
	}

	result := &ImageProcessingResult{
		OriginalURL: originalURL,
		Images:      processedImages,
		Metadata:    metadata,
	}

	return result, nil
}

// OptimizeImage optimizes a single image without generating multiple sizes
func (s *ImageService) OptimizeImage(imageData []byte, filename string, contentType string, options ProcessingOptions) (*ProcessedImage, error) {
	// Validate image
	if err := s.validateImage(imageData, contentType); err != nil {
		return nil, err
	}

	// Decode image
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, utils.WrapError(err, "failed to decode image")
	}

	// Apply optimizations
	optimizedImg := s.optimizeImageQuality(img, options)

	// Encode optimized image
	var buf bytes.Buffer
	quality := 85 // Default quality
	if options.Quality > 0 {
		quality = options.Quality
	}

	if err := s.encodeImage(&buf, optimizedImg, format, quality); err != nil {
		return nil, utils.WrapError(err, "failed to encode optimized image")
	}

	// Upload optimized image
	url, err := s.uploadImage(buf.Bytes(), filename, contentType)
	if err != nil {
		return nil, err
	}

	// Get dimensions
	bounds := optimizedImg.Bounds()
	result := &ProcessedImage{
		Size:        "optimized",
		URL:         url,
		Width:       uint(bounds.Dx()),
		Height:      uint(bounds.Dy()),
		FileSize:    int64(buf.Len()),
		ContentType: contentType,
	}

	return result, nil
}

// GetImageInfo returns metadata about an image
func (s *ImageService) GetImageInfo(imageData []byte) (*ImageMetadata, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("image_info:%x", imageData[:min(len(imageData), 32)])
	var metadata ImageMetadata
	if err := s.cacheService.Get(cacheKey, &metadata); err == nil {
		return &metadata, nil
	}

	// Decode image to get metadata
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, utils.WrapError(err, "failed to decode image")
	}

	bounds := img.Bounds()
	metadata = ImageMetadata{
		Width:       bounds.Dx(),
		Height:      bounds.Dy(),
		Format:      format,
		Size:        int64(len(imageData)),
		AspectRatio: float64(bounds.Dx()) / float64(bounds.Dy()),
	}

	// Cache metadata
	s.cacheService.Set(cacheKey, metadata, VeryLongCacheDuration)

	return &metadata, nil
}

// ValidateImageUpload validates an image upload
func (s *ImageService) ValidateImageUpload(imageData []byte, contentType string, maxSize int64) error {
	// Check file size
	if int64(len(imageData)) > maxSize {
		return utils.NewValidationError([]utils.ValidationError{
			{Field: "file_size", Message: fmt.Sprintf("Image size exceeds maximum allowed size of %d bytes", maxSize)},
		})
	}

	// Validate content type
	allowedTypes := []string{"image/jpeg", "image/png", "image/webp"}
	if !contains(allowedTypes, contentType) {
		return utils.NewValidationError([]utils.ValidationError{
			{Field: "content_type", Message: "Unsupported image format. Allowed formats: JPEG, PNG, WebP"},
		})
	}

	// Try to decode image to ensure it's valid
	_, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return utils.NewValidationError([]utils.ValidationError{
			{Field: "image", Message: "Invalid image file"},
		})
	}

	return nil
}

// Helper methods

func (s *ImageService) validateImage(imageData []byte, contentType string) error {
	// Maximum file size: 10MB
	maxSize := int64(10 * 1024 * 1024)
	return s.ValidateImageUpload(imageData, contentType, maxSize)
}

func (s *ImageService) resizeImage(img image.Image, size ImageSize, options ProcessingOptions) image.Image {
	// Calculate target dimensions maintaining aspect ratio
	bounds := img.Bounds()
	originalWidth := uint(bounds.Dx())
	originalHeight := uint(bounds.Dy())

	targetWidth := size.Width
	targetHeight := size.Height

	// Apply max width/height constraints if specified
	if options.MaxWidth > 0 && targetWidth > options.MaxWidth {
		targetWidth = options.MaxWidth
	}
	if options.MaxHeight > 0 && targetHeight > options.MaxHeight {
		targetHeight = options.MaxHeight
	}

	// Don't upscale images
	if targetWidth > originalWidth && targetHeight > originalHeight {
		targetWidth = originalWidth
		targetHeight = originalHeight
	}

	// Resize maintaining aspect ratio
	return resize.Thumbnail(targetWidth, targetHeight, img, resize.Lanczos3)
}

func (s *ImageService) optimizeImageQuality(img image.Image, options ProcessingOptions) image.Image {
	// Apply quality optimizations
	bounds := img.Bounds()
	
	// If max dimensions are specified, resize
	if options.MaxWidth > 0 || options.MaxHeight > 0 {
		width := options.MaxWidth
		height := options.MaxHeight
		
		if width == 0 {
			width = uint(bounds.Dx())
		}
		if height == 0 {
			height = uint(bounds.Dy())
		}
		
		return resize.Thumbnail(width, height, img, resize.Lanczos3)
	}

	return img
}

func (s *ImageService) encodeImage(w io.Writer, img image.Image, format string, quality int) error {
	switch format {
	case "jpeg":
		options := &jpeg.Options{Quality: quality}
		return jpeg.Encode(w, img, options)
	case "png":
		return png.Encode(w, img)
	default:
		return errors.New("unsupported image format")
	}
}

func (s *ImageService) generateResizedFilename(originalFilename, suffix string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	return fmt.Sprintf("%s%s%s", name, suffix, ext)
}

func (s *ImageService) uploadImage(imageData []byte, filename, contentType string) (string, error) {
	// This would integrate with your storage service (Cloudflare R2, AWS S3, etc.)
	// For now, return a placeholder URL
	
	// In a real implementation, you would:
	// 1. Upload to cloud storage
	// 2. Return the public URL
	// 3. Handle errors appropriately
	
	baseURL := s.cfg.External.CloudflareR2Endpoint
	if baseURL == "" {
		baseURL = "https://storage.example.com"
	}
	
	// Generate a unique filename to avoid conflicts
	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), filename)
	url := fmt.Sprintf("%s/images/%s", baseURL, uniqueFilename)
	
	return url, nil
}

// Image caching helpers
func (s *ImageService) CacheProcessedImage(key string, result *ImageProcessingResult) error {
	cacheKey := fmt.Sprintf("processed_image:%s", key)
	return s.cacheService.Set(cacheKey, result, VeryLongCacheDuration)
}

func (s *ImageService) GetCachedProcessedImage(key string) (*ImageProcessingResult, error) {
	cacheKey := fmt.Sprintf("processed_image:%s", key)
	var result ImageProcessingResult
	err := s.cacheService.Get(cacheKey, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Utility functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}