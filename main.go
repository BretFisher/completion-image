package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var (
	csvFilePath   = flag.String("names", "", "Path to the CSV file containing names, one name per row")
	imagePath     = flag.String("image", "", "Path to the input image file (png only)")
	outputPath    = flag.String("output", ".", "Path to the output directory")
	fontPath      = flag.String("font", "", "Path to the TrueType (ttf) font file")
	fontSize      = flag.Float64("size", 75, "Font size in points")
	srcColorName  = flag.String("color", "black", "Font color name")
	widthPercent  = flag.Float64("width", 0.5, "Percentage of image width to start printing the name")
	heightPercent = flag.Float64("height", 0.5, "Percentage of image height to start printing the name")
	centerText    = flag.Bool("center", false, "Center the text instead of aligning to the left")
)

func main() {
	flag.Parse()

	// Check if the required inputs are provided
	if *csvFilePath == "" || *imagePath == "" || (*fontPath == "") {
		fmt.Println("Please provide the required inputs: names, image, and font")
		// provide an example command to run
		fmt.Println("Example: gencerts -names names.csv -image image.png -font font.ttf")
		fmt.Println("Run \"gencerts -h\" for more information")
		return
	}

	// Check if the provided inputs are valid
	if _, err := os.Stat(*csvFilePath); os.IsNotExist(err) {
		log.Fatalf("The names CSV file \"%s\" does not exist", *csvFilePath)
	}
	if _, err := os.Stat(*imagePath); os.IsNotExist(err) {
		log.Fatalf("The image file \"%s\" does not exist", *imagePath)
	}
	if _, err := os.Stat(*fontPath); os.IsNotExist(err) {
		log.Fatalf("The font file \"%s\" does not exist", *fontPath)
	}

	if _, err := os.Stat(*outputPath); os.IsNotExist(err) {
		log.Fatalf("Error: output directory path does not exist")
	}

	// Validate the source color name
	*srcColorName = strings.ToLower(*srcColorName)
	srcColor, ok := colornames.Map[*srcColorName]
	fmt.Println(srcColor)
	if !ok {
		log.Fatalf("Invalid color name: \"%s\". Check the list of valid color names here: https://godoc.org/golang.org/x/image/colornames", *srcColorName)
	}

	// Read the names from the CSV file
	names, err := readNamesFromCSV(*csvFilePath)
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	// Load the original image
	file, err := os.Open(*imagePath)
	if err != nil {
		log.Fatalf("Error loading image file: %v", err)
	}
	defer file.Close()

	originalImg, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding image file: %v", err)
	}

	// Get the original image width and height
	imgWidth := originalImg.Bounds().Dx()
	imgHeight := originalImg.Bounds().Dy()

	// Load the font file
	var fontBytes []byte
	fontBytes, err = os.ReadFile(*fontPath)
	if err != nil {
		log.Fatalf("Error reading font file: %v", err)
	}

	// Determine the font type
	fontType, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalf("Error reading font type: %v", err)
	}

	// Print out the input summary
	fmt.Println("Input summary:")
	fmt.Printf("  Found %d names in the CSV file\n", len(names))
	fmt.Printf("  Original image size: %d x %d\n", imgWidth, imgHeight)
	fmt.Printf("  Font size: %v\n", *fontSize)
	fmt.Printf("  Font color: %s\n", *srcColorName)
	fmt.Printf("  Printing names %v from left, and %v up from bottom\n", *widthPercent*100, *heightPercent*100)
	fmt.Printf("  Center text is %v\n", *centerText)

	// Loop through each name and print it on a new image
	for _, name := range names {
		// Create a copy of the original image
		imgCopy := image.NewRGBA(originalImg.Bounds())
		draw.Draw(imgCopy, imgCopy.Bounds(), originalImg, image.Point{}, draw.Over)

		c := freetype.NewContext()
		c.SetFont(fontType)
		c.SetFontSize(*fontSize)
		c.SetDPI(72)
		//c.SetSrc(image.White)
		// Set the source color to RGB red
		//red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		// Create a color by its friendly name
		//srcColor := colornames.Red
		c.SetSrc(image.NewUniform(srcColor))
		c.SetClip(imgCopy.Bounds())
		c.SetDst(imgCopy)
		// Set this to None to avoid errors. Set to Full for better quality
		c.SetHinting(font.HintingNone)

		// Calculate the starting point to center or align the text
		textWidth := getTextWidth(name, fontType, *fontSize)
		startY := imgHeight - int(float64(imgHeight)**heightPercent)
		var startX = int(float64(imgWidth) * *widthPercent)
		// if center text is enabled, calculate the starting point from the center of the name text
		if *centerText {
			startX = int(float64(imgWidth)**widthPercent) - textWidth/2
		}
		pt := freetype.Pt(startX, startY)
		_, err = c.DrawString(name, pt)

		if err != nil {
			log.Fatalf("Error drawing name on image: %v", err)
		}

		// Save to a unique file
		filename := filepath.Join(*outputPath, sanitizeFilename(name)+".png")
		outputFile, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer outputFile.Close()
		err = png.Encode(outputFile, imgCopy)
		if err != nil {
			log.Fatalf("Error encoding PNG file: %v", err)
		}
		fmt.Printf("Image saved successfully as %s!\n", filename)
	}
}

// sanitize the output filename so it follows the OS rules
func sanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	reg := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	return reg.ReplaceAllString(filename, "")
}

// creates a slice of names from the CSV file
func readNamesFromCSV(filepath string) (names []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		names = append(names, record...)
	}

	return names, nil
}

// calculates the width of the text in pixels for each name
func getTextWidth(text string, font *truetype.Font, size float64) int {
	width := 0
	for _, ch := range text {
		aw := font.HMetric(fixed.Int26_6(size), font.Index(ch)).AdvanceWidth
		width += int(aw)
	}
	return width
}
