package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

var (
	S3Client     *s3.Client
	S3Uploader   *manager.Uploader
	S3BucketName string
	S3Region     string
)

// UploadResult contains the result of an S3 upload
type UploadResult struct {
	S3Key          string // Opaque S3 object key (UUID-based)
	OriginalFilename string // Original filename from upload
}

// InitializeS3 initializes the S3 client and uploader with credentials
// This function forces the use of static credentials from .env and prevents
// fallback to IAM role credentials (which would use temporary ASIA keys)
func InitializeS3() error {
	// Get credentials from environment variables
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	// Validate required environment variables
	if accessKeyID == "" {
		return fmt.Errorf("AWS_ACCESS_KEY_ID environment variable is required")
	}
	if secretAccessKey == "" {
		return fmt.Errorf("AWS_SECRET_ACCESS_KEY environment variable is required")
	}
	if bucketName == "" {
		return fmt.Errorf("AWS_S3_BUCKET_NAME environment variable is required")
	}
	if region == "" {
		return fmt.Errorf("AWS_REGION environment variable is required")
	}

	// CRITICAL: Unset temporary credential environment variables
	// These are set when IAM roles are used and would cause the SDK to use
	// temporary credentials (ASIA keys) instead of our static credentials (AKIA keys)
	tempCredVars := []string{
		"AWS_SESSION_TOKEN",
		"AWS_SECURITY_TOKEN",
		"AWS_ROLE_ARN",
		"AWS_WEB_IDENTITY_TOKEN_FILE",
	}

	for _, envVar := range tempCredVars {
		if val := os.Getenv(envVar); val != "" {
			os.Unsetenv(envVar)
			log.Printf("S3 Init: Unset %s to prevent IAM role credential fallback", envVar)
		}
	}

	// Create static credentials provider - explicitly force .env credentials
	credsProvider := credentials.NewStaticCredentialsProvider(
		accessKeyID,
		secretAccessKey,
		"", // Explicitly empty session token - ensures permanent credentials
	)

	// Create AWS config with static credentials provider
	// WithCredentialsProvider should prioritize our static credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credsProvider),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// CRITICAL: Verify which credentials are actually being used
	// This ensures we catch any credential chain fallback issues
	actualCreds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	// Mask credentials for logging (show first 8 characters only)
	maskKey := func(key string) string {
		if len(key) > 8 {
			return key[:8] + "***"
		}
		return key + "***"
	}

	expectedMasked := maskKey(accessKeyID)
	actualMasked := maskKey(actualCreds.AccessKeyID)

	// Log credential verification for debugging
	log.Printf("S3 Credentials Verification - Expected: %s, Actual: %s, Source: %s",
		expectedMasked, actualMasked, actualCreds.Source)

	// Verify access key matches exactly - this is the critical check
	// Allow both AKIA (permanent) and ASIA (temporary) credentials if explicitly set in environment
	if actualCreds.AccessKeyID != accessKeyID {
		log.Printf("ERROR: Access Key mismatch detected!")
		log.Printf("Expected: %s, Got: %s", expectedMasked, actualMasked)
		return fmt.Errorf("credentials mismatch: SDK is using %s instead of %s from .env", actualMasked, expectedMasked)
	}

	// Warn if using temporary credentials (ASIA) - but allow them if explicitly set in environment
	if !strings.HasPrefix(actualCreds.AccessKeyID, "AKIA") {
		log.Printf("WARNING: Using temporary credentials (ASIA prefix) instead of permanent (AKIA prefix)")
		log.Printf("WARNING: Temporary credentials will expire and may cause authentication failures")
		log.Printf("WARNING: Consider using permanent credentials (AKIA prefix) for production")
		// Don't return error - allow temporary credentials if explicitly set in environment
	}

	// Credentials verified - create S3 client
	S3Client = s3.NewFromConfig(cfg)
	S3Uploader = manager.NewUploader(S3Client)
	S3BucketName = bucketName
	S3Region = region

	log.Printf("S3 initialized successfully - Bucket: %s, Region: %s, Credentials: %s (verified)",
		bucketName, region, expectedMasked)

	// Verify bucket access and permissions
	if err := VerifyS3Connection(context.TODO()); err != nil {
		return fmt.Errorf("S3 bucket verification failed: %w", err)
	}

	log.Printf("✓ S3 bucket verification passed - bucket is accessible and has correct permissions")

	return nil
}

// VerifyS3Connection verifies that S3 bucket is accessible and has correct permissions
func VerifyS3Connection(ctx context.Context) error {
	if S3Client == nil {
		return fmt.Errorf("S3 client is not initialized")
	}

	// Test 1: Check if bucket exists and is accessible (HeadBucket)
	log.Printf("Verifying S3 bucket access: %s", S3BucketName)
	_, err := S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(S3BucketName),
	})
	if err != nil {
		return fmt.Errorf("cannot access bucket %s: %w. Check bucket name, region, and IAM permissions (s3:ListBucket)", S3BucketName, err)
	}
	log.Printf("✓ Bucket exists and is accessible")

	// Test 2: Verify we can list objects (tests s3:ListBucket permission)
	_, err = S3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(S3BucketName),
		MaxKeys: aws.Int32(1), // Only list 1 object to test permission
	})
	if err != nil {
		return fmt.Errorf("cannot list objects in bucket %s: %w. Check IAM permissions (s3:ListBucket)", S3BucketName, err)
	}
	log.Printf("✓ List objects permission verified")

	// Test 3: Verify we can generate presigned URLs (tests s3:GetObject permission)
	// Use a test key that might not exist - we're just testing permission, not object existence
	testKey := "test-permission-check-" + fmt.Sprintf("%d", time.Now().Unix())
	presignClient := s3.NewPresignClient(S3Client)
	_, err = presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(testKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 1 * time.Minute
	})
	if err != nil {
		return fmt.Errorf("cannot generate presigned URLs: %w. Check IAM permissions (s3:GetObject)", err)
	}
	log.Printf("✓ Presigned URL generation permission verified")

	// Test 4: Verify we can upload (tests s3:PutObject permission)
	// Create a minimal test upload to verify write permissions
	testData := []byte("test")
	testUploadKey := "test-upload-permission-" + fmt.Sprintf("%d", time.Now().Unix()) + ".txt"
	_, err = S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(S3BucketName),
		Key:         aws.String(testUploadKey),
		Body:        bytes.NewReader(testData),
		ContentType: aws.String("text/plain"),
	})
	if err != nil {
		return fmt.Errorf("cannot upload to bucket %s: %w. Check IAM permissions (s3:PutObject)", S3BucketName, err)
	}
	log.Printf("✓ Upload permission verified")

	// Clean up test file
	_, err = S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(testUploadKey),
	})
	if err != nil {
		log.Printf("Warning: Failed to delete test file %s: %v", testUploadKey, err)
		// Don't fail verification if cleanup fails
	} else {
		log.Printf("✓ Delete permission verified (test file cleaned up)")
	}

	return nil
}

// UploadFile uploads a file to S3 and returns the S3 key and original filename
// S3 keys are opaque UUID-based to decouple from original filenames
func UploadFile(ctx context.Context, fileData []byte, fileName string, contentType string, folder string) (*UploadResult, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return nil, fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	// Generate opaque, collision-safe S3 key using UUID
	// Format: {folder}/{uuid}.{ext}
	ext := filepath.Ext(fileName)
	s3Key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload file to S3 with Standard storage class for immediate access
	storageClass := types.StorageClassStandard
	putInput := &s3.PutObjectInput{
		Bucket:       aws.String(S3BucketName),
		Key:          aws.String(s3Key),
		Body:         bytes.NewReader(fileData),
		ContentType:  aws.String(contentType),
		StorageClass: storageClass,
		Metadata: map[string]string{
			"original-filename": fileName,
			"upload-date":       time.Now().Format(time.RFC3339),
		},
	}

	// Note: ACL is not set because the bucket has ACLs disabled
	// Public access should be configured via bucket policy instead
	// All access should use presigned URLs for security

	_, err := S3Uploader.Upload(ctx, putInput)
	if err != nil {
		// Return detailed error for debugging
		return nil, fmt.Errorf("S3 upload failed (bucket: %s, key: %s): %w", S3BucketName, s3Key, err)
	}

	return &UploadResult{
		S3Key:           s3Key,
		OriginalFilename: fileName,
	}, nil
}

// UploadFileLegacy uploads a file to S3 and returns the S3 URL (legacy compatibility)
// Deprecated: Use UploadFile() instead which returns S3 key separately
func UploadFileLegacy(ctx context.Context, fileData []byte, fileName string, contentType string, folder string) (string, error) {
	result, err := UploadFile(ctx, fileData, fileName, contentType, folder)
	if err != nil {
		return "", err
	}
	// Return legacy URL format for backward compatibility
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", S3BucketName, S3Region, result.S3Key)
	return url, nil
}

// GetPresignedURL generates a presigned URL for downloading a file
func GetPresignedURL(ctx context.Context, s3Key string, expiration time.Duration) (string, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return "", fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	// Validate S3 key
	if s3Key == "" {
		return "", fmt.Errorf("S3 key cannot be empty")
	}

	// Verify object exists (optional check - can be removed if it causes performance issues)
	// This helps identify permission issues early
	// Note: We don't fail - presigned URL might still work even if HeadObject fails
	_, err := S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		// Presigned URL generation might still succeed even if HeadObject fails
		// Continue without logging to avoid noise
	}

	presignClient := s3.NewPresignClient(S3Client)
	
	// Generate presigned URL with response headers for CORS support
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
		// Add response headers for CORS support
		ResponseCacheControl:       aws.String("public, max-age=3600"),
		ResponseContentDisposition: nil, // Let browser handle disposition
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL (bucket: %s, key: %s): %w. Check AWS IAM permissions for s3:GetObject", S3BucketName, s3Key, err)
	}

	return request.URL, nil
}

// DeleteFile deletes a file from S3
func DeleteFile(ctx context.Context, s3Key string) error {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	_, err := S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetS3KeyFromURL extracts the S3 key from a full S3 URL
func GetS3KeyFromURL(s3URL string) string {
	// Handle presigned URLs - extract key before query parameters
	// Format: https://bucket.s3.region.amazonaws.com/key?X-Amz-Algorithm=...
	if strings.Contains(s3URL, "?") {
		s3URL = strings.Split(s3URL, "?")[0]
	}
	
	// Extract key from URL like: https://bucket.s3.region.amazonaws.com/key
	parts := strings.Split(s3URL, ".amazonaws.com/")
	if len(parts) > 1 {
		key := parts[1]
		// URL decode the key in case it was encoded
		decodedKey, err := url.QueryUnescape(key)
		if err == nil {
			return decodedKey
		}
		return key
	}
	
	// Try alternative format: https://s3.region.amazonaws.com/bucket/key
	if strings.Contains(s3URL, "/"+S3BucketName+"/") {
		parts := strings.Split(s3URL, "/"+S3BucketName+"/")
		if len(parts) > 1 {
			key := parts[1]
			decodedKey, err := url.QueryUnescape(key)
			if err == nil {
				return decodedKey
			}
			return key
		}
	}
	
	return ""
}

// GetObjectMetadata retrieves metadata for an S3 object
func GetObjectMetadata(ctx context.Context, s3Key string) (map[string]string, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return nil, fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	result, err := S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	metadata := make(map[string]string)
	if result.Metadata != nil {
		for key, value := range result.Metadata {
			metadata[key] = value
		}
	}

	return metadata, nil
}

// GetOriginalFilename retrieves the original filename from S3 object metadata
func GetOriginalFilename(ctx context.Context, s3Key string) string {
	metadata, err := GetObjectMetadata(ctx, s3Key)
	if err != nil {
		return ""
	}
	return metadata["original-filename"]
}

// GetFileTypeFromContentType determines file type category from content type
func GetFileTypeFromContentType(contentType string) string {
	// Normalize content type
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	contentType = strings.TrimSpace(contentType)

	if strings.HasPrefix(contentType, "image/") {
		return "image"
	} else if strings.HasPrefix(contentType, "video/") {
		return "video"
	} else if strings.HasPrefix(contentType, "audio/") {
		return "audio"
	} else if strings.Contains(contentType, "pdf") ||
		strings.Contains(contentType, "word") ||
		strings.Contains(contentType, "excel") ||
		strings.Contains(contentType, "powerpoint") ||
		strings.Contains(contentType, "spreadsheet") ||
		strings.Contains(contentType, "presentation") {
		return "file"
	}
	return "file"
}

// GetFolderFromFileType returns the S3 folder based on file type
func GetFolderFromFileType(fileType string) string {
	switch fileType {
	case "image":
		return "images"
	case "video":
		return "videos"
	case "audio":
		return "audio"
	case "file":
		return "files"
	default:
		return "files"
	}
}

// ValidateFileType checks if the file type is allowed
func ValidateFileType(contentType string) bool {
	allowedTypes := []string{
		// Images
		"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp", "image/svg+xml",
		// Videos
		"video/mp4", "video/mpeg", "video/quicktime", "video/x-msvideo", "video/x-ms-wmv",
		"video/webm", "video/ogg", "video/x-matroska",
		// Audio
		"audio/mpeg", "audio/mp3", "audio/wav", "audio/ogg", "audio/webm", "audio/aac",
		"audio/x-m4a", "audio/flac", "audio/x-wav",
		// Documents
		"application/pdf",
		// Office documents (optional)
		"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	}

	// Normalize content type (remove charset, etc.)
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	contentType = strings.TrimSpace(contentType)

	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}

	return false
}

// ValidateFileSize checks if the file size is within allowed limits
func ValidateFileSize(size int64, fileType string) error {
	var maxSize int64

	switch fileType {
	case "image":
		maxSize = 10 * 1024 * 1024 // 10 MB for images
	case "video":
		maxSize = 500 * 1024 * 1024 // 500 MB for videos
	case "audio":
		maxSize = 50 * 1024 * 1024 // 50 MB for audio
	case "file":
		maxSize = 100 * 1024 * 1024 // 100 MB for PDFs and other files
	default:
		maxSize = 100 * 1024 * 1024 // 100 MB default
	}

	if size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d MB", maxSize/(1024*1024))
	}

	return nil
}
