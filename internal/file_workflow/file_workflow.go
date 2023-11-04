package file_workflow

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
)

type FileWorkflow struct{}

func (*FileWorkflow) SaveImage(outputFilepath string, outputImage image.Image) {
	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, outputImage, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Output image saved in", outputFilepath)
}

func (*FileWorkflow) ReadImage(inputFilepath string) image.Image {
	file, err := os.Open(inputFilepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	return img
}
