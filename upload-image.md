# Upload Image

## Generating POST presigned-url with minio

```go
type PresignedPostPolicyRequest struct {
	Bucket   string
	Prefix   string
	Key      string
	Duration time.Duration
}

// PresignedPostPolicy creates a presigned POST URL with additional metadata to
// allow the client side to upload the image. It reduces load on the server,
// however we lose the ability to extract metadata such as width/height of the
// images with precision, as well as getting the versionID of the image upon
// upload.
func presignedPostPolicy(client *minio.Client, req PresignedPostPolicyRequest) error {
	if req.Duration == time.Duration(0) {
		req.Duration = presignedURLValidity
	}

	// Only allow certain image type.
	ext := filepath.Ext(req.Key)
	contentType := mime.TypeByExtension(ext)
	if !strings.HasPrefix(contentType, "image/") {
		return errors.New("Content-Type invalid")
	}

	policy := minio.NewPostPolicy()
	policy.SetBucket(req.Bucket)

	buildKey := func() string {
		// <bucket>/<prefix?>/<uuid><extension>
		return filepath.Join(req.Prefix, fmt.Sprintf("%s%s", uuid.New().String(), ext))
	}

	// Overrides all the filename to default.png
	policy.SetKey(buildKey())
	policy.SetExpires(time.Now().UTC().Add(req.Duration)) // Expires in 1 day.

	policy.SetContentType(contentType)

	// Only allow content size in range 1BK to 5MB.
	policy.SetContentLengthRange(minSize, maxSize)

	// Add a user metadata using the key "custom" and value "user".
	policy.SetUserMetadata("custom", "user")

	ctx := context.Background()
	// Get the POST form key/value object.
	url, formData, err := client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return err
	}

	fmt.Println(url)
	for k, v := range formData {
		fmt.Printf("-F %s=%s\n", k, v)
	}

	return nil
}
```
