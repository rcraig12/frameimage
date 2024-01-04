package frameimage

import (
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
)

func Render(printImage, frameImage string) {
	// Load the source image
	srcImage, err := imaging.Open(printImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// Load the frame image
	frameImg, err := imaging.Open(frameImage)
	if err != nil {
		log.Fatalf("failed to open frame image: %v", err)
	}

	frameWidth := frameImg.Bounds().Dx()
	frameHeight := frameImg.Bounds().Dy()

	// Create a new image to hold the result
	dstImage := imaging.New(srcImage.Bounds().Dx()+2*frameWidth, srcImage.Bounds().Dy()+2*frameHeight, color.NRGBA{0, 0, 0, 0})
	dstImage = imaging.Paste(dstImage, srcImage, image.Pt(frameWidth, frameHeight))

	// Function to repeat the frame along a side
	repeatFrame := func(img *image.NRGBA, frame image.Image, start image.Point, horizontal bool) *image.NRGBA {
		var i int
		for i = 0; i < img.Bounds().Dx(); i += frame.Bounds().Dx() {
			if horizontal {
				img = imaging.Overlay(img, frame, image.Pt(start.X+i, start.Y), 1.0)
			} else {
				img = imaging.Overlay(img, frame, image.Pt(start.X, start.Y+i), 1.0)
			}
		}
		return img
	}

	// Top side
	dstImage = repeatFrame(dstImage, frameImg, image.Pt(frameWidth, 0), true)

	// Right side (rotate frame 90 degrees)
	frameImg = imaging.Rotate270(frameImg)
	dstImage = repeatFrame(dstImage, frameImg, image.Pt(dstImage.Bounds().Dx()-frameWidth, frameHeight), false)

	// Bottom side (rotate frame 90 degrees)
	frameImg = imaging.Rotate270(frameImg)
	dstImage = repeatFrame(dstImage, frameImg, image.Pt(frameWidth, dstImage.Bounds().Dy()-frameHeight), true)

	// Left side (rotate frame 90 degrees)
	frameImg = imaging.Rotate90(frameImg)
	dstImage = repeatFrame(dstImage, frameImg, image.Pt(0, frameHeight), false)

	// Save the resulting image
	err = imaging.Save(dstImage, "./framed_image.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
