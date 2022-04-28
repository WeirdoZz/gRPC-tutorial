package sample

import (
	"gRPC-tutori/pb"
	"github.com/google/uuid"
	"math/rand"
)

// randomKeyboardLayout generate random keyboard layout type
func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY

	}
}

// randomBool return a random bool value
func randomBool() bool {
	return rand.Intn(2) == 1
}

// randomCPUBrand generate a random CPU brand
func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

// randomStringFromSet return a random value in the parameter a
func randomStringFromSet(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

// randomCPUName return a random CPU Name
func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon E-2286M",
			"Core i9-9980HK",
			"Core i7-9750H",
			"Core i5-9400F",
			"Core i3-1005G1",
		)
	}

	return randomStringFromSet(
		"Ryzen 7 PRO 2700U",
		"Ryzen 5 PRO 3500U",
		"Ryzen 3 PRO 3200GE",
	)
}

// randomInt return a random int between min and max
func randomInt(min, max int) int {
	return min + rand.Int()%(max-min+1)
}

// randomFloat64 return a random float64 between min and max
func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// randomGPUBrand return a random GPU brand
func randomGPUBrand() string {
	return randomStringFromSet("Nvidia", "AMD")
}

// randomGPUName return a random GPU name
func randomGPUName(brand string) string {
	if brand == "Nvidia" {
		return randomStringFromSet(
			"RTX 2060",
			"RTX 2070",
			"GTX 1660-Ti",
			"GTX 1070",
		)
	}

	return randomStringFromSet(
		"RX 590",
		"RX 580",
		"RX 5700-XT",
		"RX Vega-56",
	)
}

// randomFloat32 return a random float32
func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

// randomScreenResolution return a random resolution of screen
func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	resolution := &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}

	return resolution
}

// randomScreenPanel return a random screen panel
func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

// randomId return a uuid
func randomId() string {
	return uuid.New().String()
}

// randomLaptopBrand return a random laptop brand
func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}

// randomLaptopName return a random laptop name
func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1", "Thinkpad P53")
	}
}
