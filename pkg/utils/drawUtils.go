package utils

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	pixelSize    = 25
	topMargin    = 60
	bottomMargin = 60
	leftMargin   = 60
	rightMargin  = 30
	inBetween    = 5
	blockSize    = pixelSize + inBetween

	totalX = 53
	totalY = 7

	canvasSizeWidth  = totalX*blockSize - inBetween + leftMargin + rightMargin
	canvasSizeHeight = totalY*blockSize - inBetween + topMargin + bottomMargin

	monthTextStartY    = topMargin / 2
	monthTextStartX    = int(leftMargin + 1.5*(pixelSize+inBetween))
	monthTextInBetween = int(4.4 * (pixelSize + inBetween))

	dayTextStartX    = leftMargin / 3
	dayTextStartY    = int(topMargin + 1.5*(pixelSize+inBetween))
	dayTextInBetween = 2 * (pixelSize + inBetween)

	legendTextStartX = canvasSizeWidth - 8*(pixelSize+inBetween)
	legendTextStartY = canvasSizeHeight - pixelSize - 3*inBetween
	legendTextAdjust = 2
)

const (
	background = iota
	text
)

type contributionLevel int

const (
	level0 contributionLevel = iota + 2
	level1
	level2
	level3
	level4
)

type weekday int

// Mapping weekday to index
const (
	Sunday weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var (
	months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	days   = []string{"Mon", "Wed", "Fri"}
)

var (
	// Theme Name -> color{ [ background text level0 level1 level2 level3 level4] }
	themes = make(map[string][]color.RGBA)
)

// ConstructMap function constructs and saves the contributions image
func ConstructMap(contributionList []int) {

	intensities := calculateContributionIntensity(contributionList)

	log.Println(len(intensities))

	themes["classic"] = []color.RGBA{
		color.RGBA{255, 255, 255, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{235, 237, 240, 255},
		color.RGBA{155, 233, 168, 255},
		color.RGBA{64, 196, 99, 255},
		color.RGBA{48, 161, 78, 255},
		color.RGBA{33, 110, 57, 255},
	}

	newPngFile := "./output.png"

	myImage := image.NewRGBA(image.Rect(0, 0, canvasSizeWidth, canvasSizeHeight))

	indexColor := 0
	locX := leftMargin

	// Painting the whole board
	draw.Draw(myImage, image.Rect(0, 0, canvasSizeWidth, canvasSizeHeight),
		&image.Uniform{themes["classic"][background]}, image.ZP, draw.Src)

	// Add month text
	x := monthTextStartX
	y := monthTextStartY
	for _, month := range months {
		addLabel(myImage, x, y, month)
		x += monthTextInBetween
	}

	// Add days text
	x = dayTextStartX
	y = dayTextStartY
	for _, day := range days {
		addLabel(myImage, x, y, day)
		y += dayTextInBetween
	}

	// Add legend
	x = legendTextStartX
	y = legendTextStartY
	addLabel(myImage, x-17*legendTextAdjust, y+8*legendTextAdjust, "Less")
	for color := 2; color < 7; color++ {
		draw.Draw(myImage, image.Rect(x, y, x+pixelSize, y+pixelSize),
			&image.Uniform{themes["classic"][color]}, image.ZP, draw.Src)
		x += inBetween + pixelSize
	}
	addLabel(myImage, x+2*legendTextAdjust, y+8*legendTextAdjust, "More")

	// Get starting day of the year
	date := "01-01-2020"
	t, err := time.Parse("01-02-2006", date)
	if err != nil {
		panic(err)
	}
	log.Println(t.Weekday())
	log.Println(int(t.Weekday()))

	intensitiesIndex := 0
	stop := false

	for currX := 0; currX < totalX; currX++ {

		locY := topMargin
		for currY := 0; currY < totalY; currY++ {
			if currY < int(t.Weekday()) && currX == 0 {
				log.Println(currX)
				log.Println(currY)
				log.Println(int(t.Weekday()))
				locY += blockSize
				continue
			}
			indexColor = intensities[intensitiesIndex]

			draw.Draw(myImage, image.Rect(locX, locY, locX+pixelSize, locY+pixelSize),
				&image.Uniform{themes["classic"][indexColor]}, image.ZP, draw.Src)

			locY += blockSize
			if intensitiesIndex == len(intensities)-1 {
				stop = true
				break
			}
			intensitiesIndex++
		}
		locX += blockSize
		if stop {
			break
		}
	}

	myfile, err := os.Create(newPngFile)
	if err != nil {
		panic(err.Error())
	}
	defer myfile.Close()
	png.Encode(myfile, myImage)
}

func calculateContributionIntensity(contributionList []int) []int {
	max := findMax(contributionList)
	breakPoint := float32(max / 4)
	intensities := make([]int, 0)

	log.Println(len(contributionList))
	log.Println(max)
	log.Println(breakPoint)

	for _, contribution := range contributionList {
		if contribution == 0 {
			intensities = append(intensities, 2)
			continue
		}
		contributionRange := float32(float32(contribution) / breakPoint)
		log.Println(contribution, "-->", contributionRange)
		if contributionRange <= 1 {
			intensities = append(intensities, 3)
		} else if contributionRange <= 2 {
			intensities = append(intensities, 4)
		} else if contributionRange <= 3 {
			intensities = append(intensities, 5)
		} else if contributionRange <= 4 {
			intensities = append(intensities, 6)
		}
	}
	return intensities
}

func findMax(intArray []int) float32 {
	max := intArray[0]
	for i := 1; i < len(intArray); i++ {
		if intArray[i] > max {
			max = intArray[i]
		}
	}
	return float32(max)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	color := themes["classic"][text]
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
