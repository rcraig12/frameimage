package frameimage

import (
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
)

const (
	TL = 0
	TR = 1
	BR = 2
	BL = 3
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

	placeMitre := func( img *image.NRGBA, frame image.Image, position int) *image.NRGBA {

		newFrame := image.NewRGBA(image.Rect(0, 0, frameWidth, frameHeight))
		
		maxWidth := frameWidth
		startPos := 0

		for y := 0; y < frameHeight; y++ {
			for x:= startPos; x < maxWidth; x++ {
				pixel := frame.At(x,y)
				newFrame.Set(x,y,pixel)
			}
			startPos++
		}

		if position == TL {
			img = imaging.Overlay(img, newFrame, image.Pt(0, 0), 1.0)
		} else if position == TR {
			frameImg = imaging.Rotate270(newFrame)
			img = imaging.Overlay(img, frameImg, image.Pt(dstImage.Bounds().Dx()-frameWidth,0), 1.0 )
		} else if position == BL {
			for y := 0; y < frameHeight; y++ {
				for x:= 0; x < frameWidth; x++ {
					pixel := img.At(x,y)
					newFrame.Set(x,y,pixel)
				}
			}
			frameImg = imaging.Rotate90(newFrame)
			img = imaging.Overlay(img, frameImg, image.Pt(0, dstImage.Bounds().Dy()-frameHeight), 1.0 )
		} else if position == BR {
			for y := 0; y < frameHeight; y++ {
				for x:= 0; x < frameWidth; x++ {
					pixel := img.At(x,y)
					newFrame.Set(x,y,pixel)
				}
			}
			frameImg = imaging.Rotate180(newFrame)
			img = imaging.Overlay(img, frameImg, image.Pt(dstImage.Bounds().Dx()-frameWidth, dstImage.Bounds().Dy()-frameHeight), 1.0 )
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
	frameImg = imaging.Rotate270(frameImg)
	dstImage = repeatFrame(dstImage, frameImg, image.Pt(0, frameHeight - frameWidth), false)

	// Set mitres Top Left, Top Right, Bottom Left, Bottom Right
	frameImg = imaging.Rotate270(frameImg)
	dstImage = placeMitre( dstImage, frameImg, TL )
	dstImage = placeMitre( dstImage, frameImg, TR )
	dstImage = placeMitre( dstImage, frameImg, BL )
	dstImage = placeMitre( dstImage, frameImg, BR )


	// Save the resulting image
	err = imaging.Save(dstImage, "./framed_image.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
