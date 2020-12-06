package utils

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Measurements for drawing the contribution chart
const (
	pixelSize    = 25
	topMargin    = 100
	bottomMargin = 60
	leftMargin   = 60
	rightMargin  = 30
	inBetween    = 5
	blockSize    = pixelSize + inBetween

	textAdjust = 20

	totalX = 53
	totalY = 7

	canvasSizeWidth  = totalX*blockSize - inBetween + leftMargin + rightMargin
	canvasSizeHeight = totalY*blockSize - inBetween + topMargin + bottomMargin

	monthTextStartY    = 4 * textAdjust
	monthTextStartX    = int(leftMargin + 1.5*(pixelSize+inBetween))
	monthTextInBetween = int(4.4 * (pixelSize + inBetween))

	dayTextStartX    = leftMargin / 3
	dayTextStartY    = int(topMargin + 1.5*(pixelSize+inBetween))
	dayTextInBetween = 2 * (pixelSize + inBetween)

	totalTextStartX = leftMargin/2 + 2*textAdjust
	totalTextStartY = topMargin/2 - textAdjust

	legendTextStartX = canvasSizeWidth - 8*(pixelSize+inBetween)
	legendTextStartY = canvasSizeHeight - pixelSize - 3*inBetween
	legendTextAdjust = 2
)

const (
	newPngFile = "./output.png" // Output file location
	date       = "01-01-2020"   // first date of the year
)

// Mapping of description to the index
// Used in themes
const (
	background = iota
	text
)

// intensity type to describe intensity of contribution
type intensity int

// Mapping of intensity names to the intensity levels
const (
	level0 intensity = iota + 2
	level1
	level2
	level3
	level4
)

// weekday type to describe the day of the week
type weekday int

// Mapping of weekday to indexes
const (
	Sunday weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// List of "months" and "days"
var (
	months = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	days   = []string{"Mon", "Wed", "Fri"}
)

// Mapping of theme name to the color palette values
var (
	themes = make(map[string][]color.RGBA)
)

func init() {
	themes["classic"] = []color.RGBA{
		color.RGBA{255, 255, 255, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{235, 237, 240, 255},
		color.RGBA{155, 233, 168, 255},
		color.RGBA{64, 196, 99, 255},
		color.RGBA{48, 161, 78, 255},
		color.RGBA{33, 110, 57, 255},
	}
}

// ConstructMap function constructs and saves the contributions image
func ConstructMap(contributionList []int) error {

	intensities := findIntensities(contributionList)

	myImage := image.NewRGBA(image.Rect(0, 0, canvasSizeWidth, canvasSizeHeight))

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

	// Add "total contributions" text
	x = totalTextStartX
	y = totalTextStartY
	addLabel(myImage, x, y, "x contributions this year")

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
	t, err := time.Parse("01-02-2006", date)
	if err != nil {
		return err
	}

	indexColor := level0 // Initialize intensity color to default "level0"
	intensitiesIndex := 0
	stop := false
	locX := leftMargin

	for currX := 0; currX < totalX; currX++ {

		locY := topMargin
		for currY := 0; currY < totalY; currY++ {

			// Skip weekdays until the starting weekday of the year
			if currY < int(t.Weekday()) && currX == 0 {
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

	myFile, err := os.Create(newPngFile)
	if err != nil {
		panic(err.Error())
	}
	defer myFile.Close()
	png.Encode(myFile, myImage)

	return nil
}

func findIntensities(contributions []int) []intensity {
	max := findMax(contributions)
	breakPoint := float32(max / 4)
	intensities := make([]intensity, 0)

	for _, contribution := range contributions {
		if contribution == 0 {
			intensities = append(intensities, level0)
			continue
		}
		contributionRange := float32(float32(contribution) / breakPoint)
		if contributionRange <= 1 {
			intensities = append(intensities, level1)
		} else if contributionRange <= 2 {
			intensities = append(intensities, level2)
		} else if contributionRange <= 3 {
			intensities = append(intensities, level3)
		} else if contributionRange <= 4 {
			intensities = append(intensities, level4)
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
