## Strategy

```go
package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"os"
)

const (
	TEXT_STRATEGY  = "text"
	IMAGE_STRATEGY = "image"
)

type PrintStrategy interface {
	Print() error
	SetLog(io.Writer)
	SetWriter(io.Writer)
}

type PrintOutput struct {
	Writer    io.Writer
	LogWriter io.Writer
}

func (d *PrintOutput) SetLog(w io.Writer) {
	d.LogWriter = w
}
func (d *PrintOutput) SetWriter(w io.Writer) {
	d.Writer = w
}

type TextSquare struct {
	PrintOutput
}

func (c *TextSquare) Print() error {
	c.Writer.Write([]byte("Square"))
	return nil
}

type ImageSquare struct {
	PrintOutput
	DestinationFilePath string
}

func (t *ImageSquare) Print() error {
	width := 800
	height := 600
	origin := image.Point{0, 0}
	bgImage := image.NewRGBA(image.Rectangle{
		Min: origin,
		Max: image.Point{X: width, Y: height},
	})
	bgColor := image.Uniform{color.RGBA{R: 70, G: 70, B: 70, A: 0}}
	quality := &jpeg.Options{Quality: 75}
	draw.Draw(bgImage, bgImage.Bounds(), &bgColor, origin, draw.Src)

	if t.Writer == nil {
		return errors.New("no writer stored on ImageSquare")
	}

	if err := jpeg.Encode(t.Writer, bgImage, quality); err != nil {
		return errors.New("error writing image to disk")
	}

	if t.LogWriter != nil {
		r := bytes.NewReader([]byte("image writen in provided writer"))
		io.Copy(t.LogWriter, r)
	}
	return nil
}

func NewPrinter(s string) (PrintStrategy, error) {
	switch s {
	case TEXT_STRATEGY:
		return &TextSquare{
			PrintOutput: PrintOutput{
				LogWriter: os.Stdout,
			},
		}, nil
	case IMAGE_STRATEGY:
		return &ImageSquare{
			PrintOutput: PrintOutput{
				LogWriter: os.Stdout,
			},
		}, nil
	default:
		return nil, fmt.Errorf("strategy '%s' not found", s)
	}
}

func main() {
	output := "text"
	activeStrategy, err := NewPrinter(output)
	if err != nil {
		log.Fatal(err)
	}
	switch output {
	case TEXT_STRATEGY:
		activeStrategy.SetWriter(os.Stdout)
	case IMAGE_STRATEGY:
		w, err := os.Create("/tmp/image.jpg")
		if err != nil {
			log.Fatal("error opening image")
		}
		defer w.Close()
		activeStrategy.SetWriter(w)
	}
	err = activeStrategy.Print()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Chain of Responsibility

```go
```
