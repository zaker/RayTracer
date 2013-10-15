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
	Size          int
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
	Name      string
	Materials []Material
	Spheres   []Sphere
	Blobs     []Blob
	Lights    []Light
	Exposure  float64
	CM        CubeMap
	Image     ImageS
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
	var f interface{}
	err = json.Unmarshal(cfgJson, &f)

	if err != nil {
		return
	}
	println("Unmarshaled scene")
	m := f.(map[string]interface{})
	s = new(Scene)
	sc := m["Scene"].(map[string]interface{})
	img := sc["Image"].(map[string]interface{})
	s.Name = sc["Name"].(string)
	cmFn := sc["CubeMap"].(string)
	cmF, err := os.Open("../img/" + cmFn)
	if err != nil {
		return
	}
	// format := ""
	cmImg, format, err := image.Decode(cmF)
	fmt.Println("CubeMap Format is ", format)
	if err != nil {
		return
	}

	s.CM = CubeMapInit(cmImg)
	// fmt.Println("CubeMap Format is ",format)
	s.Image = ImageS{int(img["Width"].(float64)), int(img["Height"].(float64))}

	materials := sc["Materials"].([]interface{})
	for i := range materials {
		l := materials[i].(map[string]interface{})
		d := l["Diffuse"].([]interface{})
		d2 := l["Diffuse2"].([]interface{})
		spec := l["Specular"].([]interface{})
		r := l["Reflection"].(float64)
		p := l["Power"].(float64)
		bump := l["Bump"].(float64)
		typ := l["Type"].(string)
		df := make([]float64, 0)
		for j := range d {
			df = append(df, d[j].(float64))

		}
		df2 := make([]float64, 0)
		for j := range d2 {
			df2 = append(df2, d2[j].(float64))

		}
		specular := make([]float64, 0)
		for j := range spec {
			specular = append(specular, spec[j].(float64))

		}

		s.Materials = append(s.Materials, Material{df, df2, ColorV{df[0], df[1], df[2]}, r, specular, p, bump, typ})
	}

	lights := sc["Lights"].([]interface{})
	for i := range lights {
		l := lights[i].(map[string]interface{})
		p := l["Position"].([]interface{})
		in := l["Intensity"].([]interface{})
		pf := make([]float64, 0)
		for j := range p {
			pf = append(pf, p[j].(float64))
		}
		inf := make([]float64, 0)
		for j := range in {
			inf = append(inf, in[j].(float64))
		}
		s.Lights = append(s.Lights, Light{pf, Point{pf[0], pf[1], pf[2]}, inf, ColorV{inf[0], inf[1], inf[2]}})
	}

	spheres := sc["Spheres"].([]interface{})
	for i := range spheres {
		l := spheres[i].(map[string]interface{})
		c := l["Center"].([]interface{})

		cf := make([]float64, 0)
		for j := range c {
			cf = append(cf, c[j].(float64))
		}
		si := int(l["Size"].(float64))
		ma := int(l["Material"].(float64))
		println("Sphere", si, ma)
		s.Spheres = append(s.Spheres, Sphere{cf, Point{cf[0], cf[1], cf[2]}, si, ma})
	}

	InitBlobZones()
	blobs := sc["Blobs"].([]interface{})
	for i := range blobs {
		l := blobs[i].(map[string]interface{})
		cs := l["Centers"].([]interface{})

		pts := make([]Point, 0)

		si := int(l["Size"].(float64))
		ma := int(l["Material"].(float64))
		// println(cfs, cf, cs, si, ma)
		for _, c := range cs {
			floats := c.([]interface{})
			cf := make([]float64, 0)

			for j := range floats {
				cf = append(cf, floats[j].(float64))

			}
			pts = append(pts, Point{cf[0], cf[1], cf[2]})
		}

		s.Blobs = append(s.Blobs, Blob{pts, 1.0 / float64(si*si), si, ma})
	}

	println(s.Blobs[0].Centers[0].X)
	// s.Image = ImageS{img["Width"].(float64) ,img["Height"].(float64)}
	// println(s.Image.Width)

	// jsonCheck(f)

	return s, nil

}

func (s *Scene) String() string {
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
