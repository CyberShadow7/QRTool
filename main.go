package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	// QR reading library goqr
	"github.com/liyue201/goqr"
	// QR generation library go-qrcode (should be changed with customizations-able library github.com/yeqown/go-qrcode)
	qrcode "github.com/skip2/go-qrcode" // Change with github.com/yeqown/go-qrcode
)

func reck(img image.Image) {
	qrcodes, e := goqr.Recognize(img)
	if e != nil {
		log.Fatal("Recognition of QR codes failed : ", e)
		return
	}
	dcd(qrcodes)
}

// Scan the qr code
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

func cre() {
	//a
	sc := bufio.NewScanner(os.Stdin)
	// Ask user for data to be encoded
	fmt.Println("Please type the data you wish to encode : ")
	sc.Scan()
	txdt := sc.Text()
	if sc.Err() != nil {
		log.Fatalf("\nError reading from Scanner : %v\n", sc.Err())
		os.Exit(3)
	}
	fmt.Println("Please give a name to the output file including the extension : ")
	sc.Scan()
	fn := sc.Text()
	if sc.Err() != nil {
		log.Fatalf("\nError reading from Scanner : %v\n", sc.Err())
		os.Exit(3)
	}
	imp(txdt, fn)
}

func imp(txd string, fn string) {
	dt := txd
	qrcd, _ := qrcode.Encode(dt, qrcode.Highest, 256)
	file, e := os.Create(fn)
	if e != nil {
		log.Fatalf("\nError creating file : %v\n", e)
		os.Exit(4)
	}
	defer file.Close()

	file.Write(qrcd)
}

func main() {
	rd := bufio.NewReader(os.Stdin)
	s := bufio.NewScanner(os.Stdin)
	acts := []string{"[c] Create QR Code", "[r] Read QR Code from file", "[w] Read QR Code from webcam", "[q] Quit"}
	for {
		fmt.Printf("Choose an action to perform : \n%s, %s, %s\n", acts[0], acts[1], acts[3])
		c, e := rd.ReadString('\n')
		if e != nil {
			log.Fatalf("Error reading input from STDIN : %v\n", e)
			os.Exit(2)
		}
		c = strings.TrimSpace(c)
		switch c {
		case "c":
			cre()
		case "r":
			//read from file
			fmt.Printf("Please type the /path/to/file : ")
			s.Scan()
			path := s.Text()
			if s.Err() != nil {
				log.Fatalf("\nError reading from Scanner : %v\n", s.Err())
				os.Exit(3)
			}
			scn(path)
			os.Exit(0)
		case "w":
			//w()
			// PENDING, IF IMPLEMENTED AT ALL
		case "q":
			// Quit the program
			os.Exit(0)
		default:
			fmt.Printf("Input not recognized.\nPlease choose from available inputs.\n")
		}
	}
}
