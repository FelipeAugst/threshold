package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"os"
	"slices"
	"sync"
)

var white = color.RGBA{255, 255, 255, 255}
var black = color.RGBA{0, 0, 0, 255}

func threshold(original image.Image) image.Image {
	img, lums := toGray(original)
	median := (slices.Max(lums) + slices.Min(lums)) / 2

	X := img.Bounds().Max.X
	Y := img.Bounds().Max.Y
	for y := range Y {
		for x := range X {
			r, g, b, _ := img.At(x, y).RGBA()
			if getLum(r, g, b) > median {
				img.Set(x, y, white)
			} else {
				img.Set(x, y, black)

			}

		}

	}
	return img

}

func toGray(img image.Image) (*image.RGBA, []float64) {

	X := img.Bounds().Max.X
	Y := img.Bounds().Max.Y
	filtered := image.NewRGBA(img.Bounds())
	lums := make([]float64, X*Y)
	pixel := 0

	for y := range Y {

		for x := range X {
			r, g, b, _ := img.At(x, y).RGBA()
			lum := getLum(r, g, b)
			lums[pixel] = lum
			pixel++

			var grayPixel = color.Gray{uint8(lum / 256)}

			filtered.Set(x, y, grayPixel)

		}
	}

	return filtered, lums

}

func getLum(r, g, b uint32) float64 {
	return float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114
}

func main() {

	var wg sync.WaitGroup
	if len(os.Args) <= 1 {
		fmt.Println("No images passed as arguments")
		return
	}

	for idx, path := range os.Args[1:] {
		wg.Add(1)
		go func() {

			defer wg.Done()
			file, err := os.Open(path)
			if err != nil {
				fmt.Printf("failed to open file %s:\n %s\n", path, err.Error())
				return
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {

				fmt.Printf("failed to decode image from %s:\n %s\n", path, err.Error())
				return

			}
			img = threshold(img)
			newfile := fmt.Sprintf("results/result%d.jpg", idx)
			save, err := os.Create(newfile)
			if err != nil {
				fmt.Printf("failed to create file %s: \n %s", path, err.Error())
				return

			}
			defer save.Close()
			if err := jpeg.Encode(save, img, nil); err != nil {

				fmt.Printf("failed to save file %s: \n %s \n", newfile, err.Error())
				return

			}

		}()
	}
	wg.Wait()
}
