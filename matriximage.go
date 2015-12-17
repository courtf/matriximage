package matriximage

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

type Image struct {
	image.Image
}

func (m Image) DFT() FourierImage {
	return FourierImage{Matrix: m.fftn()}
}

// Work with gray for now
// Returns a matrix with values scaled from 0.0 - 1.0
func (m Image) toGrayMatrix() *dsputils.Matrix {
	// Generate 0-based dimensions
	min, max := m.Bounds().Min, m.Bounds().Max
	lenY, lenX := max.Y-min.Y, max.X-min.X

	matrix := dsputils.MakeEmptyMatrix([]int{lenY, lenX})

	scale := 1.0

	for i := 0; i < lenX; i++ {
		for j := 0; j < lenY; j++ {

			v := scale * float64(m.Image.(*image.Gray).GrayAt(i+min.X, j+min.Y).Y)

			matrix.SetValue(complex(v, 0), []int{j, i})
		}
	}

	return matrix
}

func (m Image) fftn() *dsputils.Matrix {
	matrix := m.toGrayMatrix()
	return fft.FFTN(matrix)
}

func FromFile(filename string) (*Image, error) {
	infile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}

	grayImage := imageToGray(src)

	return &Image{Image: grayImage}, nil
}

func imageToGray(m image.Image) *image.Gray {
	b := m.Bounds()
	gray := image.NewGray(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			gray.SetGray(x, y, color.GrayModel.Convert(m.At(x, y)).(color.Gray))
		}
	}
	return gray
}

func (m Image) ToFile(named string) error {
	outfile, err := os.Create(named)
	if err != nil {
		return err
	}
	defer outfile.Close()

	return png.Encode(outfile, m.Image)
}
