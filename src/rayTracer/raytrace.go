package rayTracer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
)

type Ray struct {
	Start Point
	Dir   Vector
}

func min16(f float64) uint16 {
	t := f * 65535
	if t > 65535 {
		return 65535
	}

	return uint16(t)
}
func min8(f float64) uint8 {
	t := f * 255
	if t > 255 {
		return 255
	}

	return uint8(t)
}
func (s *Sphere) hitSphere(r Ray, t float64) (bool, float64) {
	// Intersection of a ray and a sphere
	// Check the articles for the rationale
	// NB : this is probably a naive solution
	// that could cause precision problems
	// but that will do it for now. 
	dist := s.Pos.Sub(Vector{r.Start.X, r.Start.Y, r.Start.Z})
	B := r.Dir.MulV(dist)
	D := B*B - dist.MulV(dist) + float64(s.Size*s.Size)
	if D < 0.0 {

		return false, t
	}

	t0 := B - math.Sqrt(D)
	t1 := B + math.Sqrt(D)
	retvalue := false
	if t0 > 0.1 && t0 < t {
		t = t0
		retvalue = true
	}

	if t1 > 0.1 && t1 < t {
		t = t1
		retvalue = true
	}

	return retvalue, t
}

func addRay(viewRay Ray, s *Scene) ColorV {
	output := ColorV{0.0, 0.0, 0.0}
	coef := 1.0
	level := 0

	for coef > 0.0 && level < 10 {
		// Looking for the closest intersection
		t := 2000.0
		currentSphere := -1
		currentBlob := -1

		for i, b := range s.Blobs {
			hit := false

			hit, t = b.isBlobIntersected(viewRay, t)
			if hit {

				currentBlob = i
			}
		}

		for i := range s.Spheres {
			hit := false
			hit, t = s.Spheres[i].hitSphere(viewRay, t)

			if hit {

				currentSphere = i
			}
		}

		ptHitPoint := viewRay.Start.Add(viewRay.Dir.MulC(t))
		// ptHitVector := viewRay.Start.AddP(viewRay.Dir.MulC(t))

		var normal Vector
		var currentMat Material
		if currentBlob != -1 {
			// ptHitPoint = 
			print("Blob")
			normal = s.Blobs[currentBlob].BlobInterpolation(Point{ptHitPoint.X, ptHitPoint.Y, ptHitPoint.Z})
			temp := normal.MulV(normal)
			if temp == 0.0 {
				break
			}
			temp = math.Pow(temp, -0.5)
			normal = normal.MulC(temp)
			currentMat = s.Materials[s.Blobs[currentBlob].Material]
		} else if currentSphere != -1 {

			// ptHitPoint = viewRay.Start.Add(viewRay.Dir.MulC(t))
			tmP := s.Spheres[currentSphere].Pos

			normal = ptHitPoint.Sub(Vector{tmP.X, tmP.Y, tmP.Z})
			temp := normal.MulV(normal)
			if temp == 0.0 {
				break
			}
			temp = math.Pow(temp, -0.5)
			normal = normal.MulC(temp)
			currentMat = s.Materials[s.Spheres[currentSphere].Material]
		} else {
			// print("FU")
			break

		}
		// ptHitPoint := Point{viewRay.Start.Add(viewRay.Dir.MulC(t)).X, viewRay.Start.Add(viewRay.Dir.MulC(t)).Y, viewRay.Start.Add(viewRay.Dir.MulC(t)).Z}

		// What is the normal vector at the point of intersection ?
		// It's pretty simple because we're dealing with spheres
		// normal := ptHitPoint.Sub(Vector{s.Spheres[currentSphere].Pos.X, s.Spheres[currentSphere].Pos.Y, s.Spheres[currentSphere].Pos.Z})
		// temp := normal.MulV(normal)
		// if temp == 0.0 {
		// break
		// }

		// temp = 1.0 / math.Sqrt(temp)
		// normal = normal.MulC(temp)

		// currentMat := s.Materials[s.Spheres[currentSphere].Material]

		if currentMat.Bump > 0.0 {
			noiseCoefx := Noise(0.1*ptHitPoint.X, 0.1*ptHitPoint.Y, 0.1*ptHitPoint.Z)
			noiseCoefy := Noise(0.1*ptHitPoint.Y, 0.1*ptHitPoint.Z, 0.1*ptHitPoint.X)
			noiseCoefz := Noise(0.1*ptHitPoint.Z, 0.1*ptHitPoint.X, 0.1*ptHitPoint.Y)

			normal.X = (1.0-currentMat.Bump)*normal.X + currentMat.Bump*noiseCoefx
			normal.Y = (1.0-currentMat.Bump)*normal.Y + currentMat.Bump*noiseCoefy
			normal.Z = (1.0-currentMat.Bump)*normal.Z + currentMat.Bump*noiseCoefz

			temp := normal.MulV(normal)
			if temp == 0.0 {
				break
			}
			temp = math.Pow(temp, -0.5)
			normal = normal.MulC(temp)
		}

		var lightRay Ray
		lightRay.Start = Point{ptHitPoint.X, ptHitPoint.Y, ptHitPoint.Z}

		for j := range s.Lights {
			currentLight := s.Lights[j]
			lightRay.Dir = currentLight.Pos.Sub(Vector{ptHitPoint.X, ptHitPoint.Y, ptHitPoint.Z})
			fLightProjection := lightRay.Dir.MulV(normal)

			if fLightProjection <= 0.0 {
				continue
			}
			// t := dist.MulV(dist);
			// if  t <= 0.0 {
			//     continue
			// }
			lightDist := lightRay.Dir.MulV(lightRay.Dir)
			temp := lightDist
			temp = math.Pow(temp, -0.5)
			lightRay.Dir = lightRay.Dir.MulC(temp)
			fLightProjection = temp * fLightProjection

			// computation of the shadows
			inShadow := false
			t = lightDist
			for i := range s.Spheres {
				hit := false
				hit, t = s.Spheres[i].hitSphere(lightRay, t)
				if hit {
					inShadow = true
					break
				}
			}
			if !inShadow {
				// lambert
				lambert := lightRay.Dir.MulV(normal) * coef
				noiseCoef := 0.0
				switch currentMat.Type {
				case "turbulence":

					for level := 1.0; level < 10.0; level++ {
						noiseCoef += (1.0 / level) * math.Abs(Noise(float64(level)*0.05*ptHitPoint.X, float64(level)*0.05*ptHitPoint.Y, float64(level)*0.05*ptHitPoint.Z))
					}
					output.Red += coef * (lambert * currentLight.Intensity[0]) * (noiseCoef*currentMat.Diffuse[0] + (1.0-noiseCoef)*currentMat.Diffuse2[0])
					output.Green += coef * (lambert * currentLight.Intensity[1]) * (noiseCoef*currentMat.Diffuse[1] + (1.0-noiseCoef)*currentMat.Diffuse2[1])
					output.Blue += coef * (lambert * currentLight.Intensity[2]) * (noiseCoef*currentMat.Diffuse[2] + (1.0-noiseCoef)*currentMat.Diffuse2[2])

				case "marble":

					for level := 1.0; level < 10.0; level++ {
						noiseCoef += (1.0 / level) * math.Abs(Noise(float64(level)*0.05*ptHitPoint.X, float64(level)*0.05*ptHitPoint.Y, float64(level)*0.05*ptHitPoint.Z))
					}
					noiseCoef = 0.5*math.Sin((ptHitPoint.X+ptHitPoint.Y)*0.05+noiseCoef) + 0.5
					output.Red += coef * (lambert * currentLight.Intensity[0]) * (noiseCoef*currentMat.Diffuse[0] + (1.0-noiseCoef)*currentMat.Diffuse2[0])
					output.Green += coef * (lambert * currentLight.Intensity[1]) * (noiseCoef*currentMat.Diffuse[1] + (1.0-noiseCoef)*currentMat.Diffuse2[1])
					output.Blue += coef * (lambert * currentLight.Intensity[2]) * (noiseCoef*currentMat.Diffuse[2] + (1.0-noiseCoef)*currentMat.Diffuse2[2])

				default:

					output.Red += lambert * currentLight.Color.Red * currentMat.Color.Red
					output.Green += lambert * currentLight.Color.Green * currentMat.Color.Green
					output.Blue += lambert * currentLight.Color.Blue * currentMat.Color.Blue

				}
				// output.Red += lambert * currentLight.Color.Red * currentMat.Color.Red
				// output.Green += lambert * currentLight.Color.Green * currentMat.Color.Green
				// output.Blue += lambert * currentLight.Color.Blue * currentMat.Color.Blue

				fViewProjection := viewRay.Dir.MulV(normal)
				blinnDir := lightRay.Dir.Sub(viewRay.Dir)
				temp = blinnDir.MulV(blinnDir)
				if temp != 0.0 {
					blinn := math.Pow(temp, -0.5) * math.Max(fLightProjection-fViewProjection, 0.0)
					blinn = coef * math.Pow(blinn, currentMat.Power)
					output.Red += blinn * currentMat.Specular[0] * currentLight.Intensity[0]
					output.Green += blinn * currentMat.Specular[1] * currentLight.Intensity[1]
					output.Blue += blinn * currentMat.Specular[2] * currentLight.Intensity[2]
				}
			}
		}

		// We iterate on the next reflection
		coef *= currentMat.Reflection
		reflect := 2.0 * viewRay.Dir.MulV(normal)
		viewRay.Start = Point{ptHitPoint.X, ptHitPoint.Y, ptHitPoint.Z}
		viewRay.Dir = viewRay.Dir.Sub(normal.MulC(reflect))

		level++
	}
	if coef > 0.0 {
		output.AddEq(s.CM.ReadRay(viewRay).MulC(coef))
		// println("coef")
	}
	return output
}

func AutoExposure(s *Scene) float64 {
	// #define ACCUMULATION_SIZE 16
	ACCUMULATION_SIZE := 16.0
	exposure := -1.0
	accufacteur := math.Max(float64(s.Image.Width), float64(s.Image.Height))

	accufacteur = accufacteur / ACCUMULATION_SIZE

	mediumPoint := 0.0
	mediumPointWeight := 1.0 / (ACCUMULATION_SIZE * ACCUMULATION_SIZE)
	for y := 0; float64(y) < ACCUMULATION_SIZE; y++ {
		for x := 0; float64(x) < ACCUMULATION_SIZE; x++ {
			viewRay := Ray{Point{float64(x) * accufacteur, float64(y) * accufacteur, -1000.0}, Vector{0.0, 0.0, 1.0}}
			currentColor := addRay(viewRay, s)
			luminance := 0.2126*currentColor.Red + 0.715160*currentColor.Green + 0.072169*currentColor.Blue
			mediumPoint = mediumPoint + mediumPointWeight*(luminance*luminance)
		}
	}

	mediumLuminance := math.Sqrt(mediumPoint)

	if mediumLuminance > 0.0 {
		exposure = math.Log(0.5) / mediumLuminance
	}

	return exposure
}

func DrawWorker(s *Scene, img *image.RGBA64, rect chan image.Rectangle, quit chan int) {

	for r := range rect {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			for x := r.Min.X; x < r.Max.X; x++ {

				// red, green, blue := 0.0,0.0,0.0;

				output := ColorV{0.0, 0.0, 0.0}
				for fragmentx := float64(x); fragmentx < float64(x)+1.0; fragmentx += 0.5 {
					for fragmenty := float64(y); fragmenty < float64(y)+1.0; fragmenty += 0.5 {
						sampleRatio := 0.25
						viewRay := Ray{Point{fragmentx, fragmenty, -2000.0}, Vector{0.0, 0.0, 1.0}}
						temp := addRay(viewRay, s)
						// pseudo photo exposure

						temp.Blue = (1.0 - math.Exp(temp.Blue*s.Exposure))
						temp.Red = (1.0 - math.Exp(temp.Red*s.Exposure))
						temp.Green = (1.0 - math.Exp(temp.Green*s.Exposure))

						output.AddEq(temp.MulC(sampleRatio))
					}
				}
				c := color.RGBA{min8(output.Red), min8(output.Green), min8(output.Blue), uint8(255)}
				// c := color.RGBA64{min16(output.Red), min16(output.Green), min16(output.Blue), uint16(65535)}

				img.Set(x, y, c)
				// img.SetRGBA64(x, y, c)

			}
		}
	}
	quit <- 1
}

func Draw(s *Scene) *image.RGBA64 {
	rect := image.Rect(0, 0, s.Image.Width, s.Image.Height)
	img := image.NewRGBA64(rect)
	N := runtime.NumCPU()
	runtime.GOMAXPROCS(N)
	// N := 4
	rChan := make(chan image.Rectangle)
	qChan := make(chan int, N)
	// Scanning 
	s.CM.Exposure = 3.0
	s.Exposure = AutoExposure(s)
	// s.CM.Exposure = s.Exposure

	fmt.Println("Running on ", N, " cores!")
	for i := 0; i < N; i++ {
		go DrawWorker(s, img, rChan, qChan)
	}

	// distribute image rectangles to workers
	// (line by line)
	x := img.Bounds().Dx()
	for y := 0; y < img.Bounds().Dy(); y++ {
		rChan <- image.Rect(0, y, x, y+1)
	}
	close(rChan)

	for i := 0; i < N; i++ {
		<-qChan
	}
	return img

}
