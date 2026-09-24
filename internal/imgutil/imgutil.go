package imgutil

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png" // 注册 PNG 解码器
	"io"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // 注册 webp 解码器
)

// TierWidths R2 模式下上传时预生成的缩略图宽度档位
var TierWidths = []int{128, 400, 800, 1600}

// Thumbnails 为原图生成缩略图,只包含真正需要缩小的档位(原图窄于某档则跳过该档,避免重复对象);
// 无法解码(gif/webp 异常等)时返回空,前端会回退加载原图
func Thumbnails(src []byte) map[int][]byte {
	out := make(map[int][]byte, len(TierWidths))
	img, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return out
	}
	srcW := img.Bounds().Dx()
	for _, w := range TierWidths {
		if srcW <= w {
			continue
		}
		h := img.Bounds().Dy() * w / srcW
		data, err := EncodeJPEG(Resample(img, w, h), 82)
		if err != nil {
			continue
		}
		out[w] = data
	}
	return out
}

// Resample 等比缩放(CatmullRom 高质量滤波),透明像素铺白底
func Resample(src image.Image, w, h int) image.Image {
	bounds := src.Bounds()
	flat := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(flat, flat.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(flat, flat.Bounds(), src, bounds.Min, draw.Over)

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), flat, flat.Bounds(), draw.Src, nil)
	return dst
}

func EncodeJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeConfigOr 读取图片宽高;失败时返回 0,0
func Width(r io.Reader) int {
	cfg, _, err := image.DecodeConfig(r)
	if err != nil {
		return 0
	}
	return cfg.Width
}

// IsAnimated 粗略判断是否 gif(本系统只存首帧缩略图,原图仍动图直出)
func HasExt(path string, exts ...string) bool {
	lower := strings.ToLower(path)
	for _, e := range exts {
		if strings.HasSuffix(lower, e) {
			return true
		}
	}
	return false
}
