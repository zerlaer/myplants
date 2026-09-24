package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"myplants/internal/config"
)

// r2 Cloudflare R2 兼容 S3 API,这里手写 SigV4 签名,避免引入整个 AWS SDK
type r2 struct {
	cfg    config.R2Config
	host   string
	client *http.Client
}

func NewR2(cfg config.R2Config) *r2 {
	return &r2{
		cfg:    cfg,
		host:   fmt.Sprintf("%s.r2.cloudflarestorage.com", cfg.AccountID),
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *r2) Driver() string { return "r2" }

func (s *r2) Put(key string, data []byte, contentType string) error {
	if contentType == "" {
		contentType = ContentType(key)
	}
	return s.do("PUT", key, data, contentType)
}

func (s *r2) Delete(key string) error {
	return s.do("DELETE", key, nil, "")
}

// List 列举 bucket 内所有 object key (ListObjectsV2,自动翻页)。prefix 为空则全部。
func (s *r2) List(prefix string) ([]string, error) {
	var keys []string
	token := ""
	for {
		query := map[string]string{"list-type": "2", "max-keys": "1000"}
		if prefix != "" {
			query["prefix"] = prefix
		}
		if token != "" {
			query["continuation-token"] = token
		}
		body, err := s.getWithQuery("/", query)
		if err != nil {
			return keys, err
		}
		var result struct {
			IsTruncated           bool   `xml:"IsTruncated"`
			NextContinuationToken string `xml:"NextContinuationToken"`
			Contents              []struct {
				Key string `xml:"Key"`
			} `xml:"Contents"`
		}
		if err := xml.Unmarshal(body, &result); err != nil {
			return keys, err
		}
		for _, c := range result.Contents {
			keys = append(keys, c.Key)
		}
		if !result.IsTruncated || result.NextContinuationToken == "" {
			break
		}
		token = result.NextContinuationToken
	}
	return keys, nil
}

// getWithQuery 对 bucket 根路径发起带查询参数的签名 GET,返回响应体
func (s *r2) getWithQuery(path string, query map[string]string) ([]byte, error) {
	// 规范查询串:按 key 排序,逐段转义
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, escapeSeg(k)+"="+escapeSeg(query[k]))
	}
	canonicalQuery := strings.Join(pairs, "&")

	fullPath := "/" + s.cfg.Bucket + path
	payloadHash := hashHex(nil)
	amzDate := time.Now().UTC().Format("20060102T150405Z")
	dateStamp := amzDate[:8]
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", s.host, payloadHash, amzDate)
	canonicalRequest := strings.Join([]string{"GET", fullPath, canonicalQuery, canonicalHeaders, signedHeaders, payloadHash}, "\n")
	scope := dateStamp + "/auto/s3/aws4_request"
	stringToSign := strings.Join([]string{"AWS4-HMAC-SHA256", amzDate, scope, hashHex([]byte(canonicalRequest))}, "\n")
	signature := hex.EncodeToString(hmacBytes(s.signingKey(dateStamp), []byte(stringToSign)))

	url := "https://" + s.host + fullPath
	if canonicalQuery != "" {
		url += "?" + canonicalQuery
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	req.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s,SignedHeaders=%s,Signature=%s",
		s.cfg.AccessKeyID, scope, signedHeaders, signature))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("R2 GET %s: %s %s", url, resp.Status, strings.TrimSpace(string(body)))
	}
	return body, err
}

func (s *r2) do(method, key string, body []byte, contentType string) error {
	path := "/" + s.cfg.Bucket + "/" + encodeKeyPath(key)
	payloadHash := hashHex(body)
	amzDate := time.Now().UTC().Format("20060102T150405Z")
	dateStamp := amzDate[:8]

	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n", s.host, payloadHash, amzDate)
	if contentType != "" {
		signedHeaders = "content-type;" + signedHeaders
		canonicalHeaders = "content-type:" + contentType + "\n" + canonicalHeaders
	}

	canonicalRequest := strings.Join([]string{method, path, "", canonicalHeaders, signedHeaders, payloadHash}, "\n")
	scope := dateStamp + "/auto/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256", amzDate, scope, hashHex([]byte(canonicalRequest)),
	}, "\n")

	signature := hex.EncodeToString(hmacBytes(s.signingKey(dateStamp), []byte(stringToSign)))

	req, err := http.NewRequest(method, "https://"+s.host+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("x-amz-date", amzDate)
	req.Header.Set("x-amz-content-sha256", payloadHash)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Authorization", fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s,SignedHeaders=%s,Signature=%s",
		s.cfg.AccessKeyID, scope, signedHeaders, signature))

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		if resp.StatusCode == http.StatusNotFound && method == "DELETE" {
			return nil
		}
		return fmt.Errorf("R2 %s %s: %s %s", method, key, resp.Status, strings.TrimSpace(string(msg)))
	}
	return nil
}

func (s *r2) signingKey(dateStamp string) []byte {
	k := hmacBytes([]byte("AWS4"+s.cfg.SecretAccessKey), []byte(dateStamp))
	k = hmacBytes(k, []byte("auto"))
	k = hmacBytes(k, []byte("s3"))
	return hmacBytes(k, []byte("aws4_request"))
}

func hmacBytes(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func hashHex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// encodeKeyPath 按 S3 规范逐段转义 object key,保留 '/'
func encodeKeyPath(key string) string {
	segs := strings.Split(strings.TrimPrefix(key, "/"), "/")
	for i, s := range segs {
		segs[i] = escapeSeg(s)
	}
	return strings.Join(segs, "/")
}

func escapeSeg(s string) string {
	var b strings.Builder
	for _, c := range []byte(s) {
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' ||
			c == '-' || c == '.' || c == '_' || c == '~' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}
