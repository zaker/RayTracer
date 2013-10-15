package rayTracer

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

type CubeMap struct {
	L, R, U, D, B, F *image.RGBA64
	Exposure         float64
}

func CubeMapInit(imgIn image.Image) (cm CubeMap) {
	b := imgIn.Bounds()
	x := b.Dx()
	y := b.Dy()

	img := image.NewRGBA64(b)
	draw.Draw(img, img.Bounds(), imgIn, b.Min, draw.Src)

	diffX := x / 4
	diffY := y / 3
	println("dff : ", diffX, diffY)
	r := image.Rect(0, diffY, diffX, 2*diffY)
	tmp := img.SubImage(r)
	cm.L = image.NewRGBA64(tmp.Bounds())
	cm.L.Rect = image.Rect(0, 0, diffX, diffY)
	draw.Draw(cm.L, cm.L.Bounds(), tmp, tmp.Bounds().Min, draw.Src)

	r = image.Rect(diffX, diffY, 2*diffX, 2*diffY)
	tmp = img.SubImage(r)
	cm.F = image.NewRGBA64(tmp.Bounds())
	cm.F.Rect = image.Rect(0, 0, diffX, diffY)
	draw.Draw(cm.F, cm.F.Bounds(), tmp, tmp.Bounds().Min, draw.Src)

	r = image.Rect(diffX, 0, 2*diffX, diffY)
	tmp = img.SubImage(r)
	cm.U = image.NewRGBA64(tmp.Bounds())
	cm.U.Rect = image.Rect(0, 0, diffX, diffY)
	draw.Draw(cm.U, cm.U.Bounds(), tmp, tmp.Bounds().Min, draw.Src)

	r = image.Rect(2*diffX, diffY, 3*diffX, 2*diffY)
	tmp = img.SubImage(r)
	cm.R = image.NewRGBA64(tmp.Bounds())
	cm.R.Rect = image.Rect(0, 0, diffX, diffY)
	draw.Draw(cm.R, cm.R.Bounds(), tmp, tmp.Bounds().Min, draw.Src)

	r = image.Rect(3*diffX, diffY, 4*diffX, 2*diffY)
	tmp = img.SubImage(r)
	cm.B = image.NewRGBA64(tmp.Bounds())
	cm.B.Rect = image.Rect(0, 0, diffX, diffY)
	draw.Draw(cm.B, cm.B.Bounds(), tmp, tmp.Bounds().Min, draw.Src)

	r = image.Rect(diffX, 2*diffY, 2*diffX, 3*diffY)
	tmp = img.SubImage(r)
	cm.D = image.NewRGBA64(tmp.Bounds())
	cm.D.Rect = image.Rect(0, 0, diffX, diffY)
	draw.Draw(cm.D, cm.D.Bounds(), tmp, tmp.Bounds().Min, draw.Src)

	// cm = CubeMap{L,R,U,D,B,F}
	return
}

func convertRGBA64(c color.Color) (out *ColorV) {
	r, g, b, _ := c.RGBA()

	// out = &ColorV{float64(r/65535),float64(g/65535),float64(b/65535)}
	out = &ColorV{float64(r / 255), float64(g / 255), float64(b / 255)}
	return

}

func inRange(n, min, max int) int {

	if n < min {
		return min
	}
	if n > max {
		return max
	}

	return n

}
func (cm *CubeMap) ReadRay(myRay Ray) (outputColor *ColorV) {
	outputColor = &ColorV{0.0, 0.0, 0.0}

	if math.Abs(myRay.Dir.X) >= math.Abs(myRay.Dir.Y) && math.Abs(myRay.Dir.X) >= math.Abs(myRay.Dir.Z) {
		dxy := myRay.Dir.Y / myRay.Dir.X
		if myRay.Dir.X > 0.0 {
			xf := ((myRay.Dir.Z/myRay.Dir.X + 1.0) * 0.5) * float64(cm.R.Bounds().Dx())
			xff := math.Floor(xf)
			x := int(xff)
			yf := ((dxy + 1.0) * 0.5) * float64(cm.R.Bounds().Dy())
			yff := math.Floor(yf)
			y := int(yff)
			y = inRange(y, 1, cm.L.Bounds().Dy()-1)
			// 
			// println(myRay.Dir.X,myRay.Dir.Y,myRay.Dir.Z)
			// println(x,y)
			// r,g,b,_ := cm.R.At(x,y).RGBA()
			// println(r,g,b)
			outputColor = convertRGBA64(cm.R.At(x, y))
			// outputColor = convertRGBA64(color.White)
		} else if myRay.Dir.X < 0.0 {
			xf := ((myRay.Dir.Z/myRay.Dir.X + 1.0) * 0.5) * float64(cm.L.Bounds().Dx())
			xff := math.Floor(xf)
			x := int(xff)
			yf := (1.0 - (dxy+1.0)*0.5) * float64(cm.L.Bounds().Dy())
			yff := math.Floor(yf)
			y := int(yff)
			y = inRange(y, 1, cm.L.Bounds().Dy()-1)

			outputColor = convertRGBA64(cm.L.At(x, y))

		}

	} else if math.Abs(myRay.Dir.Y) >= math.Abs(myRay.Dir.X) && math.Abs(myRay.Dir.Y) >= math.Abs(myRay.Dir.Z) {
		if myRay.Dir.Y < 0.0 {
			xf := (1.0 - (myRay.Dir.X/myRay.Dir.Y+1.0)*0.5) * float64(cm.U.Bounds().Dx())
			xff := math.Floor(xf)
			x := int(xff)
			yf := ((myRay.Dir.Z/myRay.Dir.Y + 1.0) * 0.5) * float64(cm.U.Bounds().Dy())
			yff := math.Floor(yf)
			y := int(yff)

			outputColor = convertRGBA64(cm.U.At(x, y))

		} else if myRay.Dir.Y > 0.0 {
			xf := ((myRay.Dir.X/myRay.Dir.Y + 1.0) * 0.5) * float64(cm.D.Bounds().Dx())
			xff := math.Floor(xf)

			x := int(xff)
			yf := ((myRay.Dir.Z/myRay.Dir.Y + 1.0) * 0.5) * float64(cm.D.Bounds().Dy())
			yff := math.Floor(yf)
			y := int(yff)

			outputColor = convertRGBA64(cm.D.At(x, y))

		}

	} else if math.Abs(myRay.Dir.Z) >= math.Abs(myRay.Dir.X) && math.Abs(myRay.Dir.Z) >= math.Abs(myRay.Dir.Y) {
		if myRay.Dir.Z < 0.0 {
			xf := (1.0 - (myRay.Dir.X/myRay.Dir.Z+1.0)*0.5) * float64(cm.F.Bounds().Dx())
			xff := math.Floor(xf)
			x := int(xff)
			yf := (1.0 - (myRay.Dir.Y/myRay.Dir.Z+1.0)*0.5) * float64(cm.F.Bounds().Dy())
			yff := math.Floor(yf)
			y := int(yff)
			outputColor = convertRGBA64(cm.F.At(x, y))
			// outputColor = convertRGBA64(color.White)
		} else if myRay.Dir.Z > 0.0 {

			// (myRay.dir.x / myRay.dir.z + 1.0f) * 0.5f,  
			// 1.0f - (myRay.dir.y /myRay.dir.z+1) * 0.5f

			xf := (1.0 - (myRay.Dir.X/myRay.Dir.Z+1.0)*0.5) * float64(cm.B.Bounds().Dx())
			xff := math.Floor(xf)
			x := int(xff)
			yf := ((myRay.Dir.Y/myRay.Dir.Z + 1.0) * 0.5) * float64(cm.B.Bounds().Dy())
			yff := math.Floor(yf)
			y := int(yff)

			outputColor = convertRGBA64(cm.B.At(x, y))
			// outputColor = convertRGBA64(color.White)
		}

	}

	// We make sure the data that was in sRGB storage mode is brought back to a 
	// linear format. We don't need the full accuracy of the sRGBEncode function
	// so a powf should be sufficient enough.
	// outputColor.Blue = math.Pow(outputColor.Blue, 2.2)
	// outputColor.Red = math.Pow(outputColor.Red, 2.2)
	// outputColor.Green = math.Pow(outputColor.Green, 2.2)

	//  // The LDR (low dynamic range) images were supposedly already
	//  // exposed, but we need to make the inverse transformation
	//  // so that we can expose them a second time.
	// outputColor.Blue  = -math.Log(1.001 - outputColor.Blue);
	// outputColor.Red   = -math.Log(1.001 - outputColor.Red);
	// outputColor.Green = -math.Log(1.001 - outputColor.Green)


	outputColor.Blue /= cm.Exposure
	outputColor.Red /= cm.Exposure
	outputColor.Green /= cm.Exposure

	return
}
