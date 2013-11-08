package rayTracer

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"
)

type Material struct {
	Diffuse    []float64
	Diffuse2   []float64
	Color      ColorV
	Reflection float64
	Specular   []float64
	Power      float64
	Bump       float64
	Type       string
}
type Sphere struct {
	Center   []float64
	Pos      Point
	Size     int
	Material int
}
type Blob struct {
	Centers       []Point
	invSizeSquare float64
	Size          float64
	Material      int
}
type Light struct {
	Position  []float64
	Pos       Point
	Intensity []float64
	Color     ColorV
}
type ImageS struct {
	Width  int
	Height int
}
type Scene struct {
	Name        string
	CubeMapFile string
	Materials   []Material
	Spheres     []Sphere
	Blobs       []Blob
	Lights      []Light
	Exposure    float64
	CM          CubeMap
	Image       ImageS
}

func contents(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result []byte
	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...)

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result, nil
}
func jsonCheck(f interface{}) {

	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is an object")
			jsonCheck(vv)
		}
	}
	return
}

func NewScene(file string) (s *Scene, err error) {
	cfgJson, err := contents(file)
	// println(string(cfgJson))

	if err != nil {
		return
	}
	s = new(Scene)
	err = json.Unmarshal(cfgJson, &s)
	if err != nil {
		return
	}

	cmF, err := os.Open("../img/" + s.CubeMapFile)
	if err != nil {
		return
	}
	cmImg, _, err := image.Decode(cmF)
	if err != nil {
		return
	}
	s.CM = CubeMapInit(cmImg)
	s.Image = ImageS{int(s.Image.Width), int(s.Image.Height)}
	InitBlobZones()
	println("Unmarshaled scene")
	// s.Image = ImageS{img["Width"].(float64) ,img["Height"].(float64)}
	// fmt.Printf("%+v\n", s)

	// println("Unmarshaled scene")
	// jsonCheck(f)

	return s, nil

}

func (s *Scene) Strings() string {
	st := "Scene: \n"
	st += "Image : \n"
	st += "\tWidth " + fmt.Sprint(s.Image.Width) + "\n"
	st += "\tHeight " + fmt.Sprint(s.Image.Height) + "\n"

	stt := "Spheres\n"
	for i := range s.Spheres {
		stt += "\t" + fmt.Sprint(s.Spheres[i].Center) + "\n"
		stt += "\t" + fmt.Sprint(s.Spheres[i].Size) + "\n"
		stt += "\t" + fmt.Sprint(s.Spheres[i].Material) + "\n"

	}
	st += stt

	stt = "Lights\n"
	for i := range s.Lights {
		stt += "\t" + fmt.Sprint(s.Lights[i].Position) + "\n"
		stt += "\t" + fmt.Sprint(s.Lights[i].Intensity) + "\n"

	}
	stt = "Materials\n"
	for i := range s.Materials {
		stt += "\t" + fmt.Sprint(s.Materials[i].Diffuse) + "\n"
		stt += "\t" + fmt.Sprint(s.Materials[i].Reflection) + "\n"

	}
	st += stt

	return st
}
