package image_filter

import (
	"image"
	"image/color"
	"sort"
	"sync"
)

type MedianFilter struct {
}

// medianRGB calculates the median values for the red, green, blue, and alpha channels
// from the colors slice, which contains pixel colors within the median filter window.
func medianRGB(colors []color.Color) (r, g, b, a uint8) {
	// Separate slices for red, green, blue, and alpha channel values
	var redValues, greenValues, blueValues, alphaValues []uint8

	// Extract channel values from colors and store them in separate slices
	for _, c := range colors {
		r, g, b, a := c.RGBA()
		redValues = append(redValues, uint8(r>>8))
		greenValues = append(greenValues, uint8(g>>8))
		blueValues = append(blueValues, uint8(b>>8))
		alphaValues = append(alphaValues, uint8(a>>8))
	}

	// Sort channel values
	sortChannels := func(channelValues []uint8) {
		sort.Slice(channelValues, func(i, j int) bool {
			return channelValues[i] < channelValues[j]
		})
	}

	sortChannels(redValues)
	sortChannels(greenValues)
	sortChannels(blueValues)
	sortChannels(alphaValues)

	// Calculate median values and return them as color.RGBA
	middle := len(colors) / 2
	return redValues[middle], greenValues[middle], blueValues[middle], alphaValues[middle]
}

// MedianFilterParallel applies the median filter in parallel using multiple goroutines
// to the input image with a specified window size and number of threads.
func (*MedianFilter) MedianFilterParallel(input image.Image, windowSize, numThreads int) image.Image {
	bounds := input.Bounds()
	output := image.NewRGBA(bounds)
	radius := (windowSize - 1) / 2

	var wg sync.WaitGroup
	tasks := make(chan struct {
		x, y   int
		region image.Rectangle
	}, numThreads)

	// Worker function for parallel processing
	worker := func() {
		defer wg.Done()
		for task := range tasks {
			for x := task.region.Min.X; x < task.region.Max.X; x++ {
				for y := task.region.Min.Y; y < task.region.Max.Y; y++ {
					windowColors := make([]color.Color, 0)

					for i := -radius; i <= radius; i++ {
						for j := -radius; j <= radius; j++ {
							nx, ny := x+i, y+j
							nx = clamp(nx, bounds.Min.X, bounds.Max.X-1)
							ny = clamp(ny, bounds.Min.Y, bounds.Max.Y-1)
							windowColors = append(windowColors, input.At(nx, ny))
						}
					}

					r, g, b, a := medianRGB(windowColors)
					output.Set(x, y, color.RGBA{r, g, b, a})
				}
			}
		}
	}

	// Start worker goroutines
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go worker()
	}

	// Divide the image into regions and distribute them to goroutines
	regionWidth := bounds.Dx() / numThreads
	for i := 0; i < numThreads; i++ {
		minX := i * regionWidth
		maxX := (i + 1) * regionWidth
		tasks <- struct {
			x, y   int
			region image.Rectangle
		}{minX, 0, image.Rect(minX, 0, maxX, bounds.Dy())}
	}

	// Close the task channel and wait for all worker goroutines to finish
	close(tasks)
	wg.Wait()

	return output
}

// Helper function to clamp a value within a range
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
