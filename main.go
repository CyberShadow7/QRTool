package main

import (
	"bufio"
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

var (
	c                = flag.Bool("c", false, "Create a simple QR Code")
	i                = flag.Bool("i", false, "Create a QR Code with custom image")
	a                = flag.Bool("a", false, "Create a QR Code with control over all aspects available")
	r                = flag.Bool("r", false, "Read a QR Code from a file")
	p                = flag.String("path", "", "The path to the QR Code (only with -r)")
	w                = flag.Bool("wc", false, "Read a QR Code from the webcam")
	transparent      = flag.Bool("t", false, "Use transparent background (only with -i)")
	data             = flag.String("data", "", "Data to be encoded in the QR Code")
	img              = flag.String("img", "", "Image to use in the QR Code (only with -i)")
	availableAspects = flag.Bool("av", false, "Show all available aspects")
	help             = flag.Bool("h", false, "Show this help and exit")
)

// Recon for QR codes in image file
func reck(img image.Image) {
	qrcodes, e := goqr.Recognize(img)
	if e != nil {
		log.Fatal("Recognition of QR codes failed : ", e)
		return
	}
	dcd(qrcodes)
}

// Scan the QR code
func scn(path string) {
	// Read image data from file
	dt, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file : %v\n", err)
		return
	}

	// Decode image (dt)
	img, _, e := image.Decode(bytes.NewReader(dt))
	if e != nil {
		log.Fatalf("Error with image.Decode() : %v\n", e)
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

// Insert data for created QR Code
func cre() {
	// init scanner to accept input
	sc := bufio.NewScanner(os.Stdin)
	// Ask user for data to be encoded
	fmt.Println("Please type the data you wish to encode : ")
	sc.Scan()
	txdt := sc.Text()
	if sc.Err() != nil {
		log.Fatalf("\nI/O Error reading from Scanner : %v\n", sc.Err())
		os.Exit(3)
	}
	fmt.Println("Please give a name to the output file including the extension : ")
	sc.Scan()
	fn := sc.Text()
	if sc.Err() != nil {
		log.Fatalf("\nI/O Error reading from Scanner : %v\n", sc.Err())
		os.Exit(3)
	}
	imp(txdt, fn)
}

// Encode QR Code data and Write-To-File
func imp(txd string, fn string) {
	dt := txd
	qrcd, err := qrcode.NewWith(dt,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest))
	if err != nil {
		panic(err)
	}
	file, e := os.Create(fn)
	if e != nil {
		log.Fatalf("\nError creating file : %v\n", e)
		os.Exit(4)
	}
	defer file.Close()

	//file.Write(qrcd)
	w, e := standard.New(fn, standard.WithQRWidth(40))
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

	file := "./custom-qr.jpg"

	if *transparent {
		options = append(
			options,
			standard.WithBuiltinImageEncoder(standard.JPEG_FORMAT),
			standard.WithBgTransparent(),
		)

		file = "./custom-qr-transparent.jpg"
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

func main() {

	flag.Parse()

	if *help {
		flag.CommandLine.SetOutput(os.Stdout)
		flag.PrintDefaults()
		return
	}
	if *c {
		if *data != "" {
			cre()
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
		return
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
		//fmt.Println("hi")
		return
	}
	if !*c && !*i && !*r && !*a && !*w && !*help && !*availableAspects {
		flag.PrintDefaults()
		return
	}

}
