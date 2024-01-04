package frameimage

import (
	"image"
	"log"

	"github.com/disintegration/imaging"
)

// Render framed image
func Render(printImage, frameImage string) {
	    // Load the source image
			srcImage, err := imaging.Open(printImage)
			if err != nil {
					log.Fatalf("failed to open image: %v", err)
			}
	
			// Load the frame image (this is a simplified example)
			frameImg, err := imaging.Open(frameImage)
			if err != nil {
					log.Fatalf("failed to open frame image: %v", err)
			}
	
			// Resize frame to desired dimensions (this is a simplified example)
			frameResized := imaging.Resize(frameImg, srcImage.Bounds().Dx(), 50, imaging.Lanczos)
	
			// Apply the frame to the top of the source image
			// For a full implementation, you'd need to apply it to all sides and handle corners
			dstImage := imaging.Overlay(srcImage, frameResized, image.Pt(0, 0), 1.0)
	
			// Save the resulting image
			err = imaging.Save(dstImage, "./framed_image.jpg")
			if err != nil {
					log.Fatalf("failed to save image: %v", err)
			}
			
}