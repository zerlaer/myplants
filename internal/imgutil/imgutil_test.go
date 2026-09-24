package imgutil

import (
	"image"
	"image/color"
	"testing"
)

func makeJPEG(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 128, 255})
		}
	}
	data, err := EncodeJPEG(img, 82)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestThumbnailsOnlyShrink(t *testing.T) {
	// 窄图不应生成任何档位(避免与原件重复的对象)
	if got := Thumbnails(makeJPEG(t, 100, 80)); len(got) != 0 {
		t.Fatalf("100px 图应无缩略图档位, got %v", got)
	}
	// 400px 宽:只有 128 档需要缩小,其余档位与原图等宽或更大,应跳过
	if got := Thumbnails(makeJPEG(t, 400, 300)); len(got) != 1 || got[128] == nil {
		t.Fatalf("400px 图应只生成 w128, got %d 档", len(got))
	}
	// 宽图生成全部四档
	big := Thumbnails(makeJPEG(t, 2000, 1500))
	for _, w := range TierWidths {
		if big[w] == nil {
			t.Fatalf("2000px 图缺少 w%d 档", w)
		}
	}
}
