package draw

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/shreyas-sriram/gh-contributions-aggregator/pkg/data"
)

// Measurements for drawing the contribution chart
const (
	pixelSize    = 25
	topMargin    = 100
	bottomMargin = 80
	leftMargin   = 80
	rightMargin  = 40
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
	monthTextFontSize  = 14.0

	dayTextStartX    = leftMargin / 2
	dayTextStartY    = int(topMargin + 1.5*(pixelSize+inBetween))
	dayTextInBetween = 2 * (pixelSize + inBetween)
	dayTextFontSize  = 14.0

	totalTextStartX   = leftMargin/2 + 2*textAdjust
	totalTextStartY   = topMargin/2 - textAdjust/2
	totalTextFontSize = 16.0

	legendTextStartX   = canvasSizeWidth - 8*(pixelSize+inBetween)
	legendTextStartY   = canvasSizeHeight - pixelSize - 5*inBetween
	legendTextFontSize = 14.0
	legendTextAdjust   = 2

	dpi          = 72
	fontFileProd = "./data/Raleway-Regular.ttf"
	fontFileTest = "../../data/Raleway-Regular.ttf"

	firstDate  = "01-01-"     // first date of the year minus the year
	dateFormat = "01-02-2006" // date format of variable "date"
)

// Mapping of color-usage to the index
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

	font *truetype.Font
)

func init() {
	themes["light"] = []color.RGBA{
		color.RGBA{255, 255, 255, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{235, 237, 240, 255},
		color.RGBA{155, 233, 168, 255},
		color.RGBA{64, 196, 99, 255},
		color.RGBA{48, 161, 78, 255},
		color.RGBA{33, 110, 57, 255},
	}

	themes["dark"] = []color.RGBA{
		color.RGBA{16, 16, 16, 255},
		color.RGBA{255, 255, 255, 255},
		color.RGBA{33, 33, 33, 255},
		color.RGBA{33, 110, 57, 255},
		color.RGBA{48, 161, 78, 255},
		color.RGBA{64, 196, 99, 255},
		color.RGBA{155, 233, 168, 255},
	}

	// Initialize data required for drawing
	// Read the font data
	var fontFile string
	if env := os.Getenv("ENV"); env == "test" {
		fontFile = fontFileTest
	} else {
		fontFile = fontFileProd
	}

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// findIntensities function calculates the intensity values for the contributions
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

// findMax function finds the maximum element in an integer-array
func findMax(intArray []int) float32 {
	max := intArray[0]
	for i := 1; i < len(intArray); i++ {
		if intArray[i] > max {
			max = intArray[i]
		}
	}
	return float32(max)
}

// addLabel function writes a given text on the image
func addLabel(img *image.RGBA, x, y int, label string, fontSize float64, theme string) {
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())

	if theme == "light" {
		c.SetSrc(image.Black)
	} else {
		c.SetSrc(image.White)
	}
	c.SetDst(img)
	point := freetype.Pt(x, y)

	c.DrawString(label, point)
}

// ConstructMap function constructs and saves the contributions image
func ConstructMap(request data.Request) (string, error) {

	total, contributions := data.Aggregate(request.ContributionList)

	intensities := findIntensities(contributions)

	// Create the base image
	myImage := image.NewRGBA(image.Rect(0, 0, canvasSizeWidth, canvasSizeHeight))

	draw.Draw(myImage, myImage.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Painting the background and border
	// Draw a bigger rectangle and smaller rectangle to get a border
	draw.Draw(myImage, image.Rect(0, 0, canvasSizeWidth, canvasSizeHeight),
		&image.Uniform{themes[request.Theme][text]}, image.Point{}, draw.Src)
	draw.Draw(myImage, image.Rect(0+legendTextAdjust, 0+legendTextAdjust, canvasSizeWidth-legendTextAdjust, canvasSizeHeight-legendTextAdjust),
		&image.Uniform{themes[request.Theme][background]}, image.Point{}, draw.Src)

	// Add month text
	x := monthTextStartX
	y := monthTextStartY
	for _, month := range months {
		addLabel(myImage, x, y, month, monthTextFontSize, request.Theme)
		x += monthTextInBetween
	}

	// Add days text
	x = dayTextStartX
	y = dayTextStartY
	for _, day := range days {
		addLabel(myImage, x, y, day, dayTextFontSize, request.Theme)
		y += dayTextInBetween
	}

	// Add "total contributions" text
	x = totalTextStartX
	y = totalTextStartY
	totalContributionsLabel := strconv.Itoa(total) + " contributions in " + request.Year
	addLabel(myImage, x, y, totalContributionsLabel, totalTextFontSize, request.Theme)

	// Add legend
	x = legendTextStartX
	y = legendTextStartY
	addLabel(myImage, x-18*legendTextAdjust, y+8*legendTextAdjust, "Less", legendTextFontSize, request.Theme) // Add "Less"
	for color := 2; color < 7; color++ {
		draw.Draw(myImage, image.Rect(x, y, x+pixelSize, y+pixelSize),
			&image.Uniform{themes[request.Theme][color]}, image.Point{}, draw.Src)
		x += inBetween + pixelSize
	}
	addLabel(myImage, x+legendTextAdjust, y+8*legendTextAdjust, "More", legendTextFontSize, request.Theme) // Add "More"

	// Get starting day of the year
	date := firstDate + request.Year
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", err
	}

	indexColor := level0  // Initialize intensity color to default "level0"
	intensitiesIndex := 0 // Index variable to iterate over "intensities"
	stop := false
	locationX := leftMargin

	// Paint the contributions
	for currentX := 0; currentX < totalX; currentX++ {

		locationY := topMargin
		for currentY := 0; currentY < totalY; currentY++ {

			// Skip weekdays until the starting weekday of the year
			if currentY < int(t.Weekday()) && currentX == 0 {
				locationY += blockSize
				continue
			}
			indexColor = intensities[intensitiesIndex]

			draw.Draw(myImage, image.Rect(locationX, locationY, locationX+pixelSize, locationY+pixelSize),
				&image.Uniform{themes[request.Theme][indexColor]}, image.Point{}, draw.Src)

			if intensitiesIndex == len(intensities)-1 {
				stop = true
				break
			}
			intensitiesIndex++
			locationY += blockSize
		}
		locationX += blockSize
		if stop {
			break
		}
	}

	// Encode image to buffer
	buff := new(bytes.Buffer)
	_ = png.Encode(buff, myImage)
	buf := buff.Bytes()

	// Convert image bytes to Base64
	imgBase64String := base64.StdEncoding.EncodeToString(buf)

	return imgBase64String, nil
}
