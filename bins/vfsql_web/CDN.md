# VFSQL CDN Feature

## Overview

The CDN endpoint allows you to serve files from VFSQL with optimized caching headers, making it perfect for hosting static assets that can be linked from external websites.

## Quick Start

### 1. Upload a file
```bash
# Via API
curl -X POST http://localhost:8080/api/upload \
  -F "path=/assets" \
  -F "files=@logo.png"

# Or use the web interface (drag & drop)
```

### 2. Get CDN URL
```bash
curl http://localhost:8080/api/cdn-url/assets/logo.png
```

Response:
```json
{
  "url": "/cdn/a1b2c3d4/assets/logo.png",
  "full_url": "http://localhost:8080/cdn/a1b2c3d4/assets/logo.png",
  "path": "/assets/logo.png",
  "hash": "a1b2c3d4",
  "size": 12345,
  "modified": 1234567890
}
```

### 3. Use the URL
```html
<img src="http://localhost:8080/cdn/a1b2c3d4/assets/logo.png" alt="Logo">
```

## Features

### Cache Headers
```http
Cache-Control: public, max-age=31536000, immutable
Access-Control-Allow-Origin: *
ETag: "timestamp-size"
X-Content-Type-Options: nosniff
```

### Content Type Detection
Automatically sets correct `Content-Type` for:
- Images (JPG, PNG, GIF, SVG, WEBP, ICO)
- Documents (PDF, TXT, JSON, XML, HTML, MD, CSV)
- Styles (CSS)
- Scripts (JS)
- Fonts (WOFF, WOFF2, TTF, OTF)
- Media (MP4, MP3)

### Performance Optimizations
- **1-year caching** - Reduces server load
- **ETag validation** - Returns 304 Not Modified when unchanged
- **Immutable directive** - Browser can cache without revalidation
- **HEAD support** - Check file without downloading

## URL Format

```
/cdn/{hash}/{path}
```

- `hash` - Security identifier (8-16 chars)
- `path` - File path in VFS (without leading slash)

Example:
```
/cdn/7a8b9c0d/images/photo.jpg
```

## Use Cases

### 1. Blog Images
```html
<img src="http://cdn.example.com/cdn/abc123/blog/header.jpg">
```

### 2. Web Assets
```html
<link rel="stylesheet" href="http://cdn.example.com/cdn/abc123/css/style.css">
<script src="http://cdn.example.com/cdn/abc123/js/app.js"></script>
```

### 3. Fonts
```css
@font-face {
  font-family: 'MyFont';
  src: url('http://cdn.example.com/cdn/abc123/fonts/font.woff2');
}
```

### 4. API Data
```javascript
fetch('http://cdn.example.com/cdn/abc123/data/config.json')
  .then(r => r.json())
  .then(data => console.log(data));
```

### 5. Email Images
```html
<!-- In HTML emails -->
<img src="http://cdn.example.com/cdn/abc123/email/banner.png">
```

### 6. Markdown
```markdown
![Diagram](http://cdn.example.com/cdn/abc123/docs/diagram.png)
```

## Security Notes

### Current Implementation
- Hash is generated from file path (basic security)
- No authentication required for public access
- CORS enabled for all origins

### Production Recommendations
1. Use HMAC-SHA256 with secret key for hash generation
2. Include volume/filesystem ID in hash
3. Implement rate limiting
4. Add optional token-based authentication
5. Log access for monitoring

### Example HMAC Hash
```go
import "crypto/hmac"
import "crypto/sha256"
import "encoding/base64"

func generateHash(volumeID, fsID int64, path, secret string) string {
    data := fmt.Sprintf("%d:%d:%s", volumeID, fsID, path)
    h := hmac.New(sha256.New, []byte(secret))
    h.Write([]byte(data))
    hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
    return hash[:16] // Use first 16 chars
}
```

## Web Interface

### Get CDN URL Button
1. Select a file in the file browser
2. Click "🔗 Get CDN URL" button
3. Modal shows:
   - Full URL (with domain)
   - Relative URL (for same-domain use)
   - Copy to clipboard button
   - Feature list

### Features Dialog
- Shows both full and relative URLs
- Click to select entire URL
- One-click copy to clipboard
- Lists CDN features

## API Endpoints

### Generate CDN URL
```http
GET /api/cdn-url/{path}
```

Response:
```json
{
  "url": "/cdn/hash/path",
  "full_url": "http://host/cdn/hash/path",
  "path": "/original/path",
  "hash": "hash-value",
  "size": 12345,
  "modified": 1234567890
}
```

### Serve File
```http
GET /cdn/{hash}/{path}
HEAD /cdn/{hash}/{path}
```

Response Headers:
```http
Content-Type: image/jpeg
Content-Length: 12345
Cache-Control: public, max-age=31536000, immutable
Access-Control-Allow-Origin: *
ETag: "1234567890-12345"
```

## Testing

### Test Cache Headers
```bash
curl -I http://localhost:8080/cdn/abc123/test.jpg
```

### Test ETag
```bash
# First request
curl -I http://localhost:8080/cdn/abc123/test.jpg

# Second request with ETag
curl -I -H 'If-None-Match: "1234567890-12345"' \
  http://localhost:8080/cdn/abc123/test.jpg
# Should return 304 Not Modified
```

### Test CORS
```javascript
// From different origin (e.g., CodePen)
fetch('http://localhost:8080/cdn/abc123/data.json')
  .then(r => r.json())
  .then(data => console.log('CORS works!', data));
```

## Performance Tips

1. **Use for static files** - Files that don't change frequently
2. **Leverage caching** - Let browsers cache aggressively
3. **Optimize images** - Compress before uploading
4. **Use variants** - Store multiple sizes (thumbnail, full)
5. **Monitor usage** - Track popular files
6. **CDN proxy** - Put Cloudflare/Fastly in front for global distribution

## Limitations

### Current
- Simple hash (not cryptographically secure)
- No rate limiting
- No authentication option
- No bandwidth tracking
- No usage analytics

### Future Enhancements
- HMAC-based hash validation
- Optional authentication tokens
- Rate limiting per IP/hash
- Bandwidth usage tracking
- Access logs and analytics
- Automatic image optimization
- WebP conversion
- Thumbnail generation via variants
