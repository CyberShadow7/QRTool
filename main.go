package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	// QR reading library goqr
	"github.com/liyue201/goqr"
	// QR generation library go-qrcode (should be changed with customizations-able library github.com/yeqown/go-qrcode)
	//qrcode "github.com/skip2/go-qrcode" // Change with github.com/yeqown/go-qrcode
	// OpenCV
	"github.com/yeqown/go-qrcode/v2"              // Main library
	"github.com/yeqown/go-qrcode/writer/standard" // Writer for QR files
)

// Flags to be used in command line
var (
	c                = flag.Bool("c", false, "Create a simple QR Code")
	i                = flag.Bool("i", false, "Create a QR Code with custom image")
	a                = flag.Bool("a", false, "Create a QR Code with control over all aspects available")
	l                = flag.Bool("logo", false, "Create a QR Code with a logo (JPEG file only)")
	r                = flag.Bool("r", false, "Read a QR Code from a file")
	p                = flag.String("path", "", "The path to the QR Code (only with -r)")
	w                = flag.Bool("wc", false, "Read a QR Code from the webcam")
	transparent      = flag.Bool("t", false, "Use transparent background (only with -i)")
	data             = flag.String("data", "", "Data to be encoded in the QR Code")
	img              = flag.String("img", "", "Image to use in the QR Code (only with -i)")
	format           = flag.String("format", "", "Use this format for the output image")
	out              = flag.String("o", "", "Output file name")
	availableAspects = flag.Bool("av", false, "Show all available aspects")
	help             = flag.Bool("h", false, "Show this help and exit")
)

// Recon for QR codes in image file
func reck(img image.Image) {
	qrcodes, e := goqr.Recognize(img)
	if e != nil {
		log.Fatalf("\nRecognition of QR codes failed : %v\n", e)
		return
	}
	dcd(qrcodes)
}

// Scan the QR code
func scn(path string) {
	// Read image data from file
	dt, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("\nError reading file : %v\n", err)
		return
	}

	// Decode image (dt)
	img, _, e := image.Decode(bytes.NewReader(dt))
	if e != nil {
		log.Fatalf("\nError with image.Decode() : %v\n", e)
		return
	}

	// Recon QR Codes
	reck(img)

}

// Take path as argument and start decoding
func dcd(qrc []*goqr.QRData) {
	// Print data from QR Code(s) in a list
	for _, qrc := range qrc {
		fmt.Printf("\nQR Code data : %s\n\n", qrc.Payload)
	}
}

// Encode QR Code data and Write-To-File
func createSimple() {
	qrcd, err := qrcode.NewWith(
		*data,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)

	if err != nil {
		panic(err)
	}

	var file string

	if *out != "" {
		file = *out
	} else {
		file = "qr.jpg"
	}

	w, e := standard.New(file, standard.WithQRWidth(40))
	if e != nil {
		log.Fatalf("\nError creating file : %v\n", e)
		os.Exit(4)
	}
	qrcd.Save(w)
}

// create a QR Code with a custom image
func customImg() {
	qrc, e := qrcode.New(*data)
	if e != nil {
		log.Fatalf("\nError parsing data : %v\n", e)
		os.Exit(1)
	}

	options := []standard.ImageOption{
		standard.WithHalftone(*img),
		standard.WithQRWidth(21),
	}

	file := *out

	if *transparent {
		options = append(
			options,
			standard.WithBuiltinImageEncoder(standard.JPEG_FORMAT),
			standard.WithBgTransparent(),
		)

		file = *out
	}

	w0, e := standard.New(file, options...)
	if e != nil {
		log.Fatalf("\nError creating file : %v\n", e)
		os.Exit(2)
	}

	e = qrc.Save(w0)
	if e != nil {
		log.Fatalf("\nError saving data : %v\n", e)
	}

	e = w0.Close()
	if e != nil {
		log.Fatalf("Failed to close file : %v\n", e)
	}
}

// Add a logo to the QR Code
func logo() {
	qrc, e := qrcode.NewWith(*data,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)

	if e != nil {
		log.Printf("Error reading data : %v\n", e)
	}
	file := *out
	w, e := standard.New(
		file,
		standard.WithLogoImageFileJPEG(*img),
		standard.WithLogoSizeMultiplier(2),
	)

	if e != nil {
		log.Printf("An error occured : %v\n", e)
	}

	if e = qrc.Save(w); e != nil {
		panic(e)
	}
}

func customAll() {
	qrc, e := qrcode.New(*data)
	if e != nil {
		log.Panicf("\n[!] Error reading data : %v\n", e)
		os.Exit(1)
	}

	options := []standard.ImageOption{}

	if *i {
		options = append(
			options,
			standard.WithHalftone(*img),
		)
	}
	if *transparent {
		options = append(
			options,
			standard.WithBgTransparent(),
		)
	}

	file := "./qr-code.jpg"

	if *format != "" {
		// Translate JPEG and PNG into standard.JPEG_FORMAT or standard.PNG_FORMAT accordingly.
		// Default is JPEG, by the library, but we are giving the end-user as much control as
		// the library allows.
		//
		// This could become its own function in a future update.
		if *format == "JPEG" {
			options = append(
				options,
				standard.WithBuiltinImageEncoder(standard.JPEG_FORMAT),
			)
			file = *out
		}
		if *format == "PNG" {
			options = append(
				options,
				standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
			)
			file = *out
		}
	}

	w1, e := standard.New(file, options...)
	if e != nil {
		log.Fatalf("\n[!] Error creating file : %v\n", e)
	}

	e = qrc.Save(w1)
	if e != nil {
		log.Fatalf("\n[!] Error saving data : %v\n", e)
	}

	e = w1.Close()
	if e != nil {
		log.Fatalf("\n[!] Failed to close file : %v\n", e)
	}
}

func main() {

	flag.Parse()

	if *help {
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		return
	}
	if *c {
		if *data != "" {
			createSimple()
		} else {
			fmt.Printf("\nMissing argument for -data\n")
			return
		}
	}
	if *i {
		if *data != "" || *img != "" {
			customImg()
		} else {
			fmt.Printf("\nMissing argument for -data or -img\n")
			return
		}
	}
	if *a {
		//if *as1 {
		//
		//}
		// check if all needed aspects are present.
		fmt.Println("This is a work-in-progress")
		if *data != "" || *out != "" || *format != "" {
			customAll()
		} else {
			fmt.Printf("\nMissing a needed argument!\n")
		}
	}
	if *r {
		if *p != "" {
			scn(*p)
		} else {
			fmt.Printf("\nMissing argument for -path\n")
			return
		}
	}
	if *w {
		fmt.Println("This feature is a work-in-progress")
		return
	}
	if *availableAspects {
		fmt.Println("This is a work-in-progress")
		showAllAspects()
		return
	}
	if *l {
		fmt.Println("This is a work-in-progress")
		if *img != "" || *data != "" || *out != "" {
			logo()
		} else {
			fmt.Println("Missing needed arguments")
		}
		return
	}
	if !*c && !*i && !*r && !*a && !*w && !*l && !*help && !*availableAspects {
		flag.PrintDefaults()
		return
	}

}

func showAllAspects() {
	fmt.Printf(`Image Ecnoder : JPEG, PNG
Background : Transparent
QR Width : Integer (if used with -i, must be 21)
Encoding Mode : Byte/Binary, Kenji, Numeric, Alphanumeric
Error Correction Level : L (7%%  error recovery), M (15%% error recovery - default), Q (25%% error recovery), H (30%% error recovery)
LogoImage : Add an icon at the center of the QR Code
LogoImageFilePNG / LogoImageFileJPEG : Same functionality as above
Background Color : You can specify a custom background color
Foreground Color : You can specify a custom foreground color
Shape : Rectangle (default), Circle, Custom
`)
}
