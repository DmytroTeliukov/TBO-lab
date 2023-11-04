package main

import (
	"TBO-lab/internal/file_workflow"
	"TBO-lab/internal/image_filter"
	"fmt"
	"runtime"
	"time"
)

const (
	numCPU         = 1
	windowsSize    = 3
	inputFileName  = "media_resources/komaru_cat.jpg"
	outputFileName = "media_resources/komaru_cat_output.jpg"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	fileWorkflow := &file_workflow.FileWorkflow{}
	medianFilter := &image_filter.MedianFilter{}
	inputImage := fileWorkflow.ReadImage(inputFileName)

	startTime := time.Now()
	outputImage := medianFilter.MedianFilterParallel(inputImage, windowsSize, numCPU)
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)

	fmt.Printf("Execution time: %f\n", elapsedTime.Seconds())
	fileWorkflow.SaveImage(outputFileName, outputImage)
}
