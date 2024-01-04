package frameimage

import (
	"image"
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

	// Define frame width
	frameWidth := 50 // Adjust as needed

	// Resize frame for each side
	frameTopBottom := imaging.Resize(frameImg, srcImage.Bounds().Dx(), frameWidth, imaging.Lanczos)
	frameLeftRight := imaging.Resize(frameImg, frameWidth, srcImage.Bounds().Dy(), imaging.Lanczos)

	// Rotate frame for left and right sides
	frameLeft := imaging.Rotate90(frameLeftRight)
	frameRight := imaging.Rotate270(frameLeftRight)

	// Create corner pieces by cropping the frame image
	cornerPiece := imaging.Crop(frameImg, image.Rect(0, 0, frameWidth, frameWidth))

	// Apply frames to each side
	dstImage := imaging.Overlay(srcImage, frameTopBottom, image.Pt(0, 0), 1.0) // Top
	dstImage = imaging.Overlay(dstImage, frameTopBottom, image.Pt(0, srcImage.Bounds().Dy()-frameWidth), 1.0) // Bottom
	dstImage = imaging.Overlay(dstImage, frameLeft, image.Pt(0, frameWidth), 1.0) // Left
	dstImage = imaging.Overlay(dstImage, frameRight, image.Pt(srcImage.Bounds().Dx()-frameWidth, frameWidth), 1.0) // Right

	// Apply corner pieces
	dstImage = imaging.Overlay(dstImage, cornerPiece, image.Pt(0, 0), 1.0) // Top-left
	dstImage = imaging.Overlay(dstImage, cornerPiece, image.Pt(srcImage.Bounds().Dx()-frameWidth, 0), 1.0) // Top-right
	dstImage = imaging.Overlay(dstImage, cornerPiece, image.Pt(0, srcImage.Bounds().Dy()-frameWidth), 1.0) // Bottom-left
	dstImage = imaging.Overlay(dstImage, cornerPiece, image.Pt(srcImage.Bounds().Dx()-frameWidth, srcImage.Bounds().Dy()-frameWidth), 1.0) // Bottom-right

	// Save the resulting image
	err = imaging.Save(dstImage, "./framed_image.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
