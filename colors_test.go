package colorful

import (
	"image/color"
	"math"
	"math/rand"
	"strings"
	"testing"
)

var bench_result float64 // Dummy for benchmarks to avoid optimization

// Checks whether the relative error is below eps
func almosteq_eps(v1, v2, eps float64) bool {
	if math.Abs(v1) > delta {
		return math.Abs((v1-v2)/v1) < eps
	}
	return true
}

// Checks whether the relative error is below the 8bit RGB delta, which should be good enough.
const delta = 1.0 / 256.0

func almosteq(v1, v2 float64) bool {
	return almosteq_eps(v1, v2, delta)
}

// Note: the XYZ, L*a*b*, etc. are using D65 white and D50 white if postfixed by "50".
// See http://www.brucelindbloom.com/index.html?ColorCalcHelp.html
// For d50 white, no "adaptation" and the sRGB model are used in colorful
// HCL values form http://www.easyrgb.com/index.php?X=CALC and missing ones hand-computed from lab ones
var vals = []struct {
	c      Color
	hsl    [3]float64
	hsv    [3]float64
	hex    string
	xyz    [3]float64
	xyy    [3]float64
	lab    [3]float64
	lab50  [3]float64
	luv    [3]float64
	luv50  [3]float64
	hcl    [3]float64
	hcl50  [3]float64
	rgba   [4]uint32
	rgb255 [3]uint8
}{
	{c: Color{1.0, 1.0, 1.0}, hsl: [3]float64{0.0, 0.0, 1.00}, hsv: [3]float64{0.0, 0.0, 1.0}, hex: "#ffffff", xyz: [3]float64{0.950470, 1.000000, 1.088830}, xyy: [3]float64{0.312727, 0.329023, 1.000000}, lab: [3]float64{1.000000, -0.000004, -0.000067}, lab50: [3]float64{1.000000, -0.018661, -0.151341}, luv: [3]float64{1.00000, 0.00000, 0.00000}, luv50: [3]float64{1.00000, -0.14716, -0.25658}, hcl: [3]float64{0.000000, 0.000067, 1.000000}, hcl50: [3]float64{262.970824, 0.152487, 1.000000}, rgba: [4]uint32{65535, 65535, 65535, 65535}, rgb255: [3]uint8{255, 255, 255}},
	{c: Color{0.5, 1.0, 1.0}, hsl: [3]float64{180.0, 1.0, 0.75}, hsv: [3]float64{180.0, 0.5, 1.0}, hex: "#80ffff", xyz: [3]float64{0.626296, 0.832848, 1.073634}, xyy: [3]float64{0.247276, 0.328828, 0.832848}, lab: [3]float64{0.931395, -0.276009, -0.085174}, lab50: [3]float64{0.931395, -0.292244, -0.235741}, luv: [3]float64{0.93139, -0.53909, -0.11630}, luv50: [3]float64{0.93139, -0.67615, -0.35528}, hcl: [3]float64{197.149836, 0.288852, 0.931395}, hcl50: [3]float64{218.891684, 0.375474, 0.931395}, rgba: [4]uint32{32768, 65535, 65535, 65535}, rgb255: [3]uint8{128, 255, 255}},
	{c: Color{1.0, 0.5, 1.0}, hsl: [3]float64{300.0, 1.0, 0.75}, hsv: [3]float64{300.0, 0.5, 1.0}, hex: "#ff80ff", xyz: [3]float64{0.669430, 0.437920, 0.995150}, xyy: [3]float64{0.318397, 0.208285, 0.437920}, lab: [3]float64{0.720889, 0.509121, -0.329867}, lab50: [3]float64{0.720889, 0.492522, -0.476672}, luv: [3]float64{0.72089, 0.60047, -0.77626}, luv50: [3]float64{0.72089, 0.49438, -0.96123}, hcl: [3]float64{327.060255, 0.606644, 0.720889}, hcl50: [3]float64{315.936917, 0.685415, 0.720889}, rgba: [4]uint32{65535, 32768, 65535, 65535}, rgb255: [3]uint8{255, 128, 255}},
	{c: Color{1.0, 1.0, 0.5}, hsl: [3]float64{60.0, 1.0, 0.75}, hsv: [3]float64{60.0, 0.5, 1.0}, hex: "#ffff80", xyz: [3]float64{0.808654, 0.943273, 0.341930}, xyy: [3]float64{0.386203, 0.450496, 0.943273}, lab: [3]float64{0.977634, -0.129552, 0.470291}, lab50: [3]float64{0.977634, -0.147230, 0.367469}, luv: [3]float64{0.97764, 0.05759, 0.79816}, luv50: [3]float64{0.97764, -0.08628, 0.54731}, hcl: [3]float64{105.401396, 0.487808, 0.977634}, hcl50: [3]float64{111.834030, 0.395867, 0.977634}, rgba: [4]uint32{65535, 65535, 32768, 65535}, rgb255: [3]uint8{255, 255, 128}},
	{c: Color{0.5, 0.5, 1.0}, hsl: [3]float64{240.0, 1.0, 0.75}, hsv: [3]float64{240.0, 0.5, 1.0}, hex: "#8080ff", xyz: [3]float64{0.345256, 0.270768, 0.979954}, xyy: [3]float64{0.216329, 0.169656, 0.270768}, lab: [3]float64{0.590461, 0.260064, -0.497795}, lab50: [3]float64{0.590461, 0.246752, -0.643849}, luv: [3]float64{0.59045, -0.07568, -1.04877}, luv50: [3]float64{0.59045, -0.16257, -1.20027}, hcl: [3]float64{297.584033, 0.561634, 0.590461}, hcl50: [3]float64{290.969081, 0.689513, 0.590461}, rgba: [4]uint32{32768, 32768, 65535, 65535}, rgb255: [3]uint8{128, 128, 255}},
	{c: Color{1.0, 0.5, 0.5}, hsl: [3]float64{0.0, 1.0, 0.75}, hsv: [3]float64{0.0, 0.5, 1.0}, hex: "#ff8080", xyz: [3]float64{0.527613, 0.381193, 0.248250}, xyy: [3]float64{0.455996, 0.329451, 0.381193}, lab: [3]float64{0.681075, 0.378015, 0.178332}, lab50: [3]float64{0.681075, 0.362682, 0.085917}, luv: [3]float64{0.68108, 0.92148, 0.19879}, luv50: [3]float64{0.68106, 0.82106, 0.02393}, hcl: [3]float64{25.255982, 0.417968, 0.681075}, hcl50: [3]float64{13.327381, 0.372719, 0.681075}, rgba: [4]uint32{65535, 32768, 32768, 65535}, rgb255: [3]uint8{255, 128, 128}},
	{c: Color{0.5, 1.0, 0.5}, hsl: [3]float64{120.0, 1.0, 0.75}, hsv: [3]float64{120.0, 0.5, 1.0}, hex: "#80ff80", xyz: [3]float64{0.484480, 0.776121, 0.326734}, xyy: [3]float64{0.305216, 0.488946, 0.776121}, lab: [3]float64{0.906028, -0.469433, 0.389809}, lab50: [3]float64{0.906028, -0.484336, 0.288533}, luv: [3]float64{0.90603, -0.58869, 0.76102}, luv50: [3]float64{0.90603, -0.72202, 0.52855}, hcl: [3]float64{140.294351, 0.610179, 0.906028}, hcl50: [3]float64{149.216497, 0.563767, 0.906028}, rgba: [4]uint32{32768, 65535, 32768, 65535}, rgb255: [3]uint8{128, 255, 128}},
	{c: Color{0.5, 0.5, 0.5}, hsl: [3]float64{0.0, 0.0, 0.50}, hsv: [3]float64{0.0, 0.0, 0.5}, hex: "#808080", xyz: [3]float64{0.203440, 0.214041, 0.233054}, xyy: [3]float64{0.312727, 0.329023, 0.214041}, lab: [3]float64{0.533890, -0.000002, -0.000040}, lab50: [3]float64{0.533890, -0.011162, -0.090529}, luv: [3]float64{0.53389, 0.00000, 0.00000}, luv50: [3]float64{0.53389, -0.07857, -0.13699}, hcl: [3]float64{0.000000, 0.000040, 0.533890}, hcl50: [3]float64{262.970824, 0.091215, 0.533890}, rgba: [4]uint32{32768, 32768, 32768, 65535}, rgb255: [3]uint8{128, 128, 128}},
	{c: Color{0.0, 1.0, 1.0}, hsl: [3]float64{180.0, 1.0, 0.50}, hsv: [3]float64{180.0, 1.0, 1.0}, hex: "#00ffff", xyz: [3]float64{0.538014, 0.787327, 1.069496}, xyy: [3]float64{0.224656, 0.328760, 0.787327}, lab: [3]float64{0.911140, -0.375650, -0.110458}, lab50: [3]float64{0.911140, -0.391084, -0.260831}, luv: [3]float64{0.91113, -0.70477, -0.15204}, luv50: [3]float64{0.91113, -0.83886, -0.38582}, hcl: [3]float64{196.385720, 0.391553, 0.911140}, hcl50: [3]float64{213.701127, 0.470084, 0.911140}, rgba: [4]uint32{0, 65535, 65535, 65535}, rgb255: [3]uint8{0, 255, 255}},
	{c: Color{1.0, 0.0, 1.0}, hsl: [3]float64{300.0, 1.0, 0.50}, hsv: [3]float64{300.0, 1.0, 1.0}, hex: "#ff00ff", xyz: [3]float64{0.592894, 0.284848, 0.969638}, xyy: [3]float64{0.320938, 0.154190, 0.284848}, lab: [3]float64{0.603236, 0.767463, -0.475274}, lab50: [3]float64{0.603236, 0.751522, -0.620814}, luv: [3]float64{0.60324, 0.84071, -1.08683}, luv50: [3]float64{0.60324, 0.75194, -1.24161}, hcl: [3]float64{328.230959, 0.902710, 0.603236}, hcl50: [3]float64{320.440744, 0.974780, 0.603236}, rgba: [4]uint32{65535, 0, 65535, 65535}, rgb255: [3]uint8{255, 0, 255}},
	{c: Color{1.0, 1.0, 0.0}, hsl: [3]float64{60.0, 1.0, 0.50}, hsv: [3]float64{60.0, 1.0, 1.0}, hex: "#ffff00", xyz: [3]float64{0.770033, 0.927825, 0.138526}, xyy: [3]float64{0.419320, 0.505246, 0.927825}, lab: [3]float64{0.971388, -0.168420, 0.738104}, lab50: [3]float64{0.971388, -0.185812, 0.662024}, luv: [3]float64{0.97139, 0.07706, 1.06787}, luv50: [3]float64{0.97139, -0.06590, 0.81862}, hcl: [3]float64{102.853641, 0.757075, 0.971388}, hcl50: [3]float64{105.677996, 0.687606, 0.971388}, rgba: [4]uint32{65535, 65535, 0, 65535}, rgb255: [3]uint8{255, 255, 0}},
	{c: Color{0.0, 0.0, 1.0}, hsl: [3]float64{240.0, 1.0, 0.50}, hsv: [3]float64{240.0, 1.0, 1.0}, hex: "#0000ff", xyz: [3]float64{0.180437, 0.072175, 0.950304}, xyy: [3]float64{0.150000, 0.060000, 0.072175}, lab: [3]float64{0.322994, 0.618683, -0.842699}, lab50: [3]float64{0.322994, 0.607960, -0.987265}, luv: [3]float64{0.32297, -0.09405, -1.30342}, luv50: [3]float64{0.32297, -0.14158, -1.38629}, hcl: [3]float64{306.284932, 1.045423, 0.322994}, hcl50: [3]float64{301.624825, 1.159443, 0.322994}, rgba: [4]uint32{0, 0, 65535, 65535}, rgb255: [3]uint8{0, 0, 255}},
	{c: Color{0.0, 1.0, 0.0}, hsl: [3]float64{120.0, 1.0, 0.50}, hsv: [3]float64{120.0, 1.0, 1.0}, hex: "#00ff00", xyz: [3]float64{0.357576, 0.715152, 0.119192}, xyy: [3]float64{0.300000, 0.600000, 0.715152}, lab: [3]float64{0.877350, -0.673304, 0.649840}, lab50: [3]float64{0.877350, -0.686773, 0.577478}, luv: [3]float64{0.87735, -0.83078, 1.07398}, luv50: [3]float64{0.87735, -0.95989, 0.84887}, hcl: [3]float64{136.015956, 0.935751, 0.877350}, hcl50: [3]float64{139.940931, 0.897295, 0.877350}, rgba: [4]uint32{0, 65535, 0, 65535}, rgb255: [3]uint8{0, 255, 0}},
	{c: Color{1.0, 0.0, 0.0}, hsl: [3]float64{0.0, 1.0, 0.50}, hsv: [3]float64{0.0, 1.0, 1.0}, hex: "#ff0000", xyz: [3]float64{0.412456, 0.212673, 0.019334}, xyy: [3]float64{0.640000, 0.330000, 0.212673}, lab: [3]float64{0.532390, 0.625707, 0.525011}, lab50: [3]float64{0.532390, 0.611582, 0.485548}, luv: [3]float64{0.53241, 1.75015, 0.37756}, luv50: [3]float64{0.53241, 1.67180, 0.24096}, hcl: [3]float64{39.998956, 0.816790, 0.532390}, hcl50: [3]float64{38.446811, 0.780890, 0.532390}, rgba: [4]uint32{65535, 0, 0, 65535}, rgb255: [3]uint8{255, 0, 0}},
	{c: Color{0.0, 0.0, 0.0}, hsl: [3]float64{0.0, 0.0, 0.00}, hsv: [3]float64{0.0, 0.0, 0.0}, hex: "#000000", xyz: [3]float64{0.000000, 0.000000, 0.000000}, xyy: [3]float64{0.312727, 0.329023, 0.000000}, lab: [3]float64{0.000000, 0.000000, 0.000000}, lab50: [3]float64{0.000000, 0.000000, 0.000000}, luv: [3]float64{0.00000, 0.00000, 0.00000}, luv50: [3]float64{0.00000, 0.00000, 0.00000}, hcl: [3]float64{0.000000, 0.000000, 0.000000}, hcl50: [3]float64{0.000000, 0.000000, 0.000000}, rgba: [4]uint32{0, 0, 0, 65535}, rgb255: [3]uint8{0, 0, 0}},
}

// For testing short-hex values, since the above contains colors which don't
// have corresponding short hexes.
var shorthexvals = []struct {
	c   Color
	hex string
}{
	{Color{1.0, 1.0, 1.0}, "#fff"},
	{Color{0.6, 1.0, 1.0}, "#9ff"},
	{Color{1.0, 0.6, 1.0}, "#f9f"},
	{Color{1.0, 1.0, 0.6}, "#ff9"},
	{Color{0.6, 0.6, 1.0}, "#99f"},
	{Color{1.0, 0.6, 0.6}, "#f99"},
	{Color{0.6, 1.0, 0.6}, "#9f9"},
	{Color{0.6, 0.6, 0.6}, "#999"},
	{Color{0.0, 1.0, 1.0}, "#0ff"},
	{Color{1.0, 0.0, 1.0}, "#f0f"},
	{Color{1.0, 1.0, 0.0}, "#ff0"},
	{Color{0.0, 0.0, 1.0}, "#00f"},
	{Color{0.0, 1.0, 0.0}, "#0f0"},
	{Color{1.0, 0.0, 0.0}, "#f00"},
	{Color{0.0, 0.0, 0.0}, "#000"},
}

/// RGBA ///
////////////

func TestRGBAConversion(t *testing.T) {
	for i, tt := range vals {
		r, g, b, a := tt.c.RGBA()
		if r != tt.rgba[0] || g != tt.rgba[1] || b != tt.rgba[2] || a != tt.rgba[3] {
			t.Errorf("%v. %v.RGBA() => (%v), want %v (delta %v)", i, tt.c, []uint32{r, g, b, a}, tt.rgba, delta)
		}
	}
}

/// RGB255 ///
////////////

func TestRGB255Conversion(t *testing.T) {
	for i, tt := range vals {
		r, g, b := tt.c.RGB255()
		if r != tt.rgb255[0] || g != tt.rgb255[1] || b != tt.rgb255[2] {
			t.Errorf("%v. %v.RGB255() => (%v), want %v (delta %v)", i, tt.c, []uint8{r, g, b}, tt.rgb255, delta)
		}
	}
}

/// HSV ///
///////////

func TestHsvCreation(t *testing.T) {
	for i, tt := range vals {
		c := Hsv(tt.hsv[0], tt.hsv[1], tt.hsv[2])
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Hsv(%v) => (%v), want %v (delta %v)", i, tt.hsv, c, tt.c, delta)
		}
	}
}

func TestHsvConversion(t *testing.T) {
	for i, tt := range vals {
		h, s, v := tt.c.Hsv()
		if !almosteq(h, tt.hsv[0]) || !almosteq(s, tt.hsv[1]) || !almosteq(v, tt.hsv[2]) {
			t.Errorf("%v. %v.Hsv() => (%v), want %v (delta %v)", i, tt.c, []float64{h, s, v}, tt.hsv, delta)
		}
	}
}

/// HSL ///
///////////

func TestHslCreation(t *testing.T) {
	for i, tt := range vals {
		c := Hsl(tt.hsl[0], tt.hsl[1], tt.hsl[2])
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Hsl(%v) => (%v), want %v (delta %v)", i, tt.hsl, c, tt.c, delta)
		}
	}
}

func TestHslConversion(t *testing.T) {
	for i, tt := range vals {
		h, s, l := tt.c.Hsl()
		if !almosteq(h, tt.hsl[0]) || !almosteq(s, tt.hsl[1]) || !almosteq(l, tt.hsl[2]) {
			t.Errorf("%v. %v.Hsl() => (%v), want %v (delta %v)", i, tt.c, []float64{h, s, l}, tt.hsl, delta)
		}
	}
}

/// Hex ///
///////////

func TestHexCreation(t *testing.T) {
	for i, tt := range vals {
		c, err := Hex(tt.hex)
		if err != nil || !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Hex(%v) => (%v), want %v (delta %v)", i, tt.hex, c, tt.c, delta)
		}
	}
}

func TestHEXCreation(t *testing.T) {
	for i, tt := range vals {
		c, err := Hex(strings.ToUpper(tt.hex))
		if err != nil || !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. HEX(%v) => (%v), want %v (delta %v)", i, strings.ToUpper(tt.hex), c, tt.c, delta)
		}
	}
}

func TestShortHexCreation(t *testing.T) {
	for i, tt := range shorthexvals {
		c, err := Hex(tt.hex)
		if err != nil || !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Hex(%v) => (%v), want %v (delta %v)", i, tt.hex, c, tt.c, delta)
		}
	}
}

func TestShortHEXCreation(t *testing.T) {
	for i, tt := range shorthexvals {
		c, err := Hex(strings.ToUpper(tt.hex))
		if err != nil || !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Hex(%v) => (%v), want %v (delta %v)", i, strings.ToUpper(tt.hex), c, tt.c, delta)
		}
	}
}

func TestHexConversion(t *testing.T) {
	for i, tt := range vals {
		hex := tt.c.Hex()
		if hex != tt.hex {
			t.Errorf("%v. %v.Hex() => (%v), want %v (delta %v)", i, tt.c, hex, tt.hex, delta)
		}
	}
}

/// Linear ///
//////////////

// LinearRgb itself is implicitly tested by XYZ conversions below (they use it).
// So what we do here is just test that the FastLinearRgb approximation is "good enough"
func TestFastLinearRgb(t *testing.T) {
	const eps = 6.0 / 255.0 // We want that "within 6 RGB values total" is "good enough".

	for r := 0.0; r < 256.0; r++ {
		for g := 0.0; g < 256.0; g++ {
			for b := 0.0; b < 256.0; b++ {
				c := Color{r / 255.0, g / 255.0, b / 255.0}
				r_want, g_want, b_want := c.LinearRgb()
				r_appr, g_appr, b_appr := c.FastLinearRgb()
				dr, dg, db := math.Abs(r_want-r_appr), math.Abs(g_want-g_appr), math.Abs(b_want-b_appr)
				if dr+dg+db > eps {
					t.Errorf("FastLinearRgb not precise enough for %v: differences are (%v, %v, %v), allowed total difference is %v", c, dr, dg, db, eps)
					return
				}

				c_want := LinearRgb(r/255.0, g/255.0, b/255.0)
				c_appr := FastLinearRgb(r/255.0, g/255.0, b/255.0)
				dr, dg, db = math.Abs(c_want.R-c_appr.R), math.Abs(c_want.G-c_appr.G), math.Abs(c_want.B-c_appr.B)
				if dr+dg+db > eps {
					t.Errorf("FastLinearRgb not precise enough for (%v, %v, %v): differences are (%v, %v, %v), allowed total difference is %v", r, g, b, dr, dg, db, eps)
					return
				}
			}
		}
	}
}

// Also include some benchmarks to make sure the `Fast` versions are actually significantly faster!
// (Sounds silly, but the original ones weren't!)

func BenchmarkColorToLinear(bench *testing.B) {
	var r, g, b float64
	for n := 0; n < bench.N; n++ {
		r, g, b = Color{rand.Float64(), rand.Float64(), rand.Float64()}.LinearRgb()
	}
	bench_result = r + g + b
}

func BenchmarkFastColorToLinear(bench *testing.B) {
	var r, g, b float64
	for n := 0; n < bench.N; n++ {
		r, g, b = Color{rand.Float64(), rand.Float64(), rand.Float64()}.FastLinearRgb()
	}
	bench_result = r + g + b
}

func BenchmarkLinearToColor(bench *testing.B) {
	var c Color
	for n := 0; n < bench.N; n++ {
		c = LinearRgb(rand.Float64(), rand.Float64(), rand.Float64())
	}
	bench_result = c.R + c.G + c.B
}

func BenchmarkFastLinearToColor(bench *testing.B) {
	var c Color
	for n := 0; n < bench.N; n++ {
		c = FastLinearRgb(rand.Float64(), rand.Float64(), rand.Float64())
	}
	bench_result = c.R + c.G + c.B
}

/// XYZ ///
///////////
func TestXyzCreation(t *testing.T) {
	for i, tt := range vals {
		c := Xyz(tt.xyz[0], tt.xyz[1], tt.xyz[2])
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Xyz(%v) => (%v), want %v (delta %v)", i, tt.xyz, c, tt.c, delta)
		}
	}
}

func TestXyzConversion(t *testing.T) {
	for i, tt := range vals {
		x, y, z := tt.c.Xyz()
		if !almosteq(x, tt.xyz[0]) || !almosteq(y, tt.xyz[1]) || !almosteq(z, tt.xyz[2]) {
			t.Errorf("%v. %v.Xyz() => (%v), want %v (delta %v)", i, tt.c, [3]float64{x, y, z}, tt.xyz, delta)
		}
	}
}

/// xyY ///
///////////
func TestXyyCreation(t *testing.T) {
	for i, tt := range vals {
		c := Xyy(tt.xyy[0], tt.xyy[1], tt.xyy[2])
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Xyy(%v) => (%v), want %v (delta %v)", i, tt.xyy, c, tt.c, delta)
		}
	}
}

func TestXyyConversion(t *testing.T) {
	for i, tt := range vals {
		x, y, Y := tt.c.Xyy()
		if !almosteq(x, tt.xyy[0]) || !almosteq(y, tt.xyy[1]) || !almosteq(Y, tt.xyy[2]) {
			t.Errorf("%v. %v.Xyy() => (%v), want %v (delta %v)", i, tt.c, [3]float64{x, y, Y}, tt.xyy, delta)
		}
	}
}

/// L*a*b* ///
//////////////
func TestLabCreation(t *testing.T) {
	for i, tt := range vals {
		c := Lab(tt.lab[0], tt.lab[1], tt.lab[2])
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Lab(%v) => (%v), want %v (delta %v)", i, tt.lab, c, tt.c, delta)
		}
	}
}

func TestLabConversion(t *testing.T) {
	for i, tt := range vals {
		l, a, b := tt.c.Lab()
		if !almosteq(l, tt.lab[0]) || !almosteq(a, tt.lab[1]) || !almosteq(b, tt.lab[2]) {
			t.Errorf("%v. %v.Lab() => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, a, b}, tt.lab, delta)
		}
	}
}

func TestLabWhiteRefCreation(t *testing.T) {
	for i, tt := range vals {
		c := LabWhiteRef(tt.lab50[0], tt.lab50[1], tt.lab50[2], D50)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. LabWhiteRef(%v, D50) => (%v), want %v (delta %v)", i, tt.lab50, c, tt.c, delta)
		}
	}
}

func TestLabWhiteRefConversion(t *testing.T) {
	for i, tt := range vals {
		l, a, b := tt.c.LabWhiteRef(D50)
		if !almosteq(l, tt.lab50[0]) || !almosteq(a, tt.lab50[1]) || !almosteq(b, tt.lab50[2]) {
			t.Errorf("%v. %v.LabWhiteRef(D50) => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, a, b}, tt.lab50, delta)
		}
	}
}

/// L*u*v* ///
//////////////
func TestLuvCreation(t *testing.T) {
	for i, tt := range vals {
		c := Luv(tt.luv[0], tt.luv[1], tt.luv[2])
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Luv(%v) => (%v), want %v (delta %v)", i, tt.luv, c, tt.c, delta)
		}
	}
}

func TestLuvConversion(t *testing.T) {
	for i, tt := range vals {
		l, u, v := tt.c.Luv()
		if !almosteq(l, tt.luv[0]) || !almosteq(u, tt.luv[1]) || !almosteq(v, tt.luv[2]) {
			t.Errorf("%v. %v.Luv() => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, u, v}, tt.luv, delta)
		}
	}
}

func TestLuvWhiteRefCreation(t *testing.T) {
	for i, tt := range vals {
		c := LuvWhiteRef(tt.luv50[0], tt.luv50[1], tt.luv50[2], D50)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. LuvWhiteRef(%v, D50) => (%v), want %v (delta %v)", i, tt.luv50, c, tt.c, delta)
		}
	}
}

func TestLuvWhiteRefConversion(t *testing.T) {
	for i, tt := range vals {
		l, u, v := tt.c.LuvWhiteRef(D50)
		if !almosteq(l, tt.luv50[0]) || !almosteq(u, tt.luv50[1]) || !almosteq(v, tt.luv50[2]) {
			t.Errorf("%v. %v.LuvWhiteRef(D50) => (%v), want %v (delta %v)", i, tt.c, [3]float64{l, u, v}, tt.luv50, delta)
		}
	}
}

/// HCL ///
///////////
// CIE-L*a*b* in polar coordinates.
func TestHclCreation(t *testing.T) {
	for i, tt := range vals {
		c := Hcl(tt.hcl[0], tt.hcl[1], tt.hcl[2])
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. Hcl(%v) => (%v), want %v (delta %v)", i, tt.hcl, c, tt.c, delta)
		}
	}
}

func TestHclConversion(t *testing.T) {
	for i, tt := range vals {
		h, c, l := tt.c.Hcl()
		if !almosteq(h, tt.hcl[0]) || !almosteq(c, tt.hcl[1]) || !almosteq(l, tt.hcl[2]) {
			t.Errorf("%v. %v.Hcl() => (%v), want %v (delta %v)", i, tt.c, [3]float64{h, c, l}, tt.hcl, delta)
		}
	}
}

func TestHclWhiteRefCreation(t *testing.T) {
	for i, tt := range vals {
		c := HclWhiteRef(tt.hcl50[0], tt.hcl50[1], tt.hcl50[2], D50)
		if !c.AlmostEqualRgb(tt.c) {
			t.Errorf("%v. HclWhiteRef(%v, D50) => (%v), want %v (delta %v)", i, tt.hcl50, c, tt.c, delta)
		}
	}
}

func TestHclWhiteRefConversion(t *testing.T) {
	for i, tt := range vals {
		h, c, l := tt.c.HclWhiteRef(D50)
		if !almosteq(h, tt.hcl50[0]) || !almosteq(c, tt.hcl50[1]) || !almosteq(l, tt.hcl50[2]) {
			t.Errorf("%v. %v.HclWhiteRef(D50) => (%v), want %v (delta %v)", i, tt.c, [3]float64{h, c, l}, tt.hcl50, delta)
		}
	}
}

/// Test distances ///
//////////////////////

// Ground-truth from http://www.brucelindbloom.com/index.html?ColorDifferenceCalcHelp.html
var dists = []struct {
	c1  Color
	c2  Color
	d76 float64 // That's also dLab
	d94 float64
	d00 float64
}{
	{Color{1.0, 1.0, 1.0}, Color{1.0, 1.0, 1.0}, 0.0, 0.0, 0.0},
	{Color{0.0, 0.0, 0.0}, Color{0.0, 0.0, 0.0}, 0.0, 0.0, 0.0},

	// Just pairs of values of the table way above.
	{Lab(1.000000, 0.000000, 0.000000), Lab(0.931390, -0.353319, -0.108946), 47.82075103713838, 47.82075103713826, 24.95839826021627},
	{Lab(0.720892, 0.651673, -0.422133), Lab(0.977637, -0.165795, 0.602017), 169.6842805730313, 71.73241839963067, 85.19582299499467},
	{Lab(0.590453, 0.332846, -0.637099), Lab(0.681085, 0.483884, 0.228328), 112.81367663807802, 47.30269333036521, 41.38821174674354},
	{Lab(0.906026, -0.600870, 0.498993), Lab(0.533890, 0.000000, 0.000000), 106.6758044178174, 41.41738433437593, 41.690253231681176},
	{Lab(0.911132, -0.480875, -0.141312), Lab(0.603242, 0.982343, -0.608249), 198.99353255481768, 98.12426621634205, 63.89635218932751},
	{Lab(0.971393, -0.215537, 0.944780), Lab(0.322970, 0.791875, -1.078602), 296.49556584164037, 119.21658629889, 110.6721165119005},
	{Lab(0.877347, -0.861827, 0.831793), Lab(0.532408, 0.800925, 0.672032), 216.57695311092084, 73.21043676199285, 96.18944220977916},
}

func TestLabDistance(t *testing.T) {
	for i, tt := range dists {
		d := tt.c1.DistanceCIE76(tt.c2)
		if !almosteq(d, tt.d76) {
			t.Errorf("%v. %v.DistanceCIE76(%v) => (%v), want %v (delta %v)", i, tt.c1, tt.c2, d, tt.d76, delta)
		}
	}
}

func TestCIE94Distance(t *testing.T) {
	for i, tt := range dists {
		d := tt.c1.DistanceCIE94(tt.c2)
		if !almosteq(d, tt.d94) {
			t.Errorf("%v. %v.DistanceCIE94(%v) => (%v), want %v (delta %v)", i, tt.c1, tt.c2, d, tt.d94, delta)
		}
	}
}

func TestCIEDE2000Distance(t *testing.T) {
	for i, tt := range dists {
		d := tt.c1.DistanceCIEDE2000(tt.c2)
		if !almosteq(d, tt.d00) {
			t.Errorf("%v. %v.DistanceCIEDE2000(%v) => (%v), want %v (delta %v)", i, tt.c1, tt.c2, d, tt.d00, delta)
		}
	}
}

/// Test utilities ///
//////////////////////

func TestClamp(t *testing.T) {
	c_orig := Color{1.1, -0.1, 0.5}
	c_want := Color{1.0, 0.0, 0.5}
	if c_orig.Clamped() != c_want {
		t.Errorf("%v.Clamped() => %v, want %v", c_orig, c_orig.Clamped(), c_want)
	}
}

func TestMakeColor(t *testing.T) {
	c_orig_nrgba := color.NRGBA{123, 45, 67, 255}
	c_ours, ok := MakeColor(c_orig_nrgba)
	r, g, b := c_ours.RGB255()
	if r != 123 || g != 45 || b != 67 || !ok {
		t.Errorf("NRGBA->Colorful->RGB255 error: %v became (%v, %v, %v, %t)", c_orig_nrgba, r, g, b, ok)
	}

	c_orig_nrgba64 := color.NRGBA64{123 << 8, 45 << 8, 67 << 8, 0xffff}
	c_ours, ok = MakeColor(c_orig_nrgba64)
	r, g, b = c_ours.RGB255()
	if r != 123 || g != 45 || b != 67 || !ok {
		t.Errorf("NRGBA64->Colorful->RGB255 error: %v became (%v, %v, %v, %t)", c_orig_nrgba64, r, g, b, ok)
	}

	c_orig_gray := color.Gray{123}
	c_ours, ok = MakeColor(c_orig_gray)
	r, g, b = c_ours.RGB255()
	if r != 123 || g != 123 || b != 123 || !ok {
		t.Errorf("Gray->Colorful->RGB255 error: %v became (%v, %v, %v, %t)", c_orig_gray, r, g, b, ok)
	}

	c_orig_gray16 := color.Gray16{123 << 8}
	c_ours, ok = MakeColor(c_orig_gray16)
	r, g, b = c_ours.RGB255()
	if r != 123 || g != 123 || b != 123 || !ok {
		t.Errorf("Gray16->Colorful->RGB255 error: %v became (%v, %v, %v, %t)", c_orig_gray16, r, g, b, ok)
	}

	c_orig_rgba := color.RGBA{255, 255, 255, 0}
	c_ours, ok = MakeColor(c_orig_rgba)
	r, g, b = c_ours.RGB255()
	if r != 0 || g != 0 || b != 0 || ok {
		t.Errorf("RGBA->Colorful->RGB255 error: %v became (%v, %v, %v, %t)", c_orig_rgba, r, g, b, ok)
	}
}

/// Issues raised on github ///
///////////////////////////////

// https://github.com/lucasb-eyer/go-colorful/issues/11
func TestIssue11(t *testing.T) {
	c1hex := "#1a1a46"
	c2hex := "#666666"

	c1, _ := Hex(c1hex)
	c2, _ := Hex(c2hex)

	blend := c1.BlendHsv(c2, 0).Hex()
	if blend != c1hex {
		t.Errorf("Issue11: %v --Hsv-> %v = %v, want %v", c1hex, c2hex, blend, c1hex)
	}
	blend = c1.BlendHsv(c2, 1).Hex()
	if blend != c2hex {
		t.Errorf("Issue11: %v --Hsv-> %v = %v, want %v", c1hex, c2hex, blend, c2hex)
	}

	blend = c1.BlendLuv(c2, 0).Hex()
	if blend != c1hex {
		t.Errorf("Issue11: %v --Luv-> %v = %v, want %v", c1hex, c2hex, blend, c1hex)
	}
	blend = c1.BlendLuv(c2, 1).Hex()
	if blend != c2hex {
		t.Errorf("Issue11: %v --Luv-> %v = %v, want %v", c1hex, c2hex, blend, c2hex)
	}

	blend = c1.BlendRgb(c2, 0).Hex()
	if blend != c1hex {
		t.Errorf("Issue11: %v --Rgb-> %v = %v, want %v", c1hex, c2hex, blend, c1hex)
	}
	blend = c1.BlendRgb(c2, 1).Hex()
	if blend != c2hex {
		t.Errorf("Issue11: %v --Rgb-> %v = %v, want %v", c1hex, c2hex, blend, c2hex)
	}

	blend = c1.BlendLinearRgb(c2, 0).Hex()
	if blend != c1hex {
		t.Errorf("Issue11: %v --LinearRgb-> %v = %v, want %v", c1hex, c2hex, blend, c1hex)
	}
	blend = c1.BlendLinearRgb(c2, 1).Hex()
	if blend != c2hex {
		t.Errorf("Issue11: %v --LinearRgb-> %v = %v, want %v", c1hex, c2hex, blend, c2hex)
	}

	blend = c1.BlendLab(c2, 0).Hex()
	if blend != c1hex {
		t.Errorf("Issue11: %v --Lab-> %v = %v, want %v", c1hex, c2hex, blend, c1hex)
	}
	blend = c1.BlendLab(c2, 1).Hex()
	if blend != c2hex {
		t.Errorf("Issue11: %v --Lab-> %v = %v, want %v", c1hex, c2hex, blend, c2hex)
	}

	blend = c1.BlendHcl(c2, 0).Hex()
	if blend != c1hex {
		t.Errorf("Issue11: %v --Hcl-> %v = %v, want %v", c1hex, c2hex, blend, c1hex)
	}
	blend = c1.BlendHcl(c2, 1).Hex()
	if blend != c2hex {
		t.Errorf("Issue11: %v --Hcl-> %v = %v, want %v", c1hex, c2hex, blend, c2hex)
	}

	blend = c1.BlendLuvLCh(c2, 0).Hex()
	if blend != c1hex {
		t.Errorf("Issue11: %v --LuvLCh-> %v = %v, want %v", c1hex, c2hex, blend, c1hex)
	}
	blend = c1.BlendLuvLCh(c2, 1).Hex()
	if blend != c2hex {
		t.Errorf("Issue11: %v --LuvLCh-> %v = %v, want %v", c1hex, c2hex, blend, c2hex)
	}
}

// For testing angular interpolation internal function
// NOTE: They are being tested in both directions.
var anglevals = []struct {
	a0 float64
	a1 float64
	t  float64
	at float64
}{
	{0.0, 1.0, 0.0, 0.0},
	{0.0, 1.0, 0.25, 0.25},
	{0.0, 1.0, 0.5, 0.5},
	{0.0, 1.0, 1.0, 1.0},
	{0.0, 90.0, 0.0, 0.0},
	{0.0, 90.0, 0.25, 22.5},
	{0.0, 90.0, 0.5, 45.0},
	{0.0, 90.0, 1.0, 90.0},
	{0.0, 178.0, 0.0, 0.0}, // Exact 0-180 is ambiguous.
	{0.0, 178.0, 0.25, 44.5},
	{0.0, 178.0, 0.5, 89.0},
	{0.0, 178.0, 1.0, 178.0},
	{0.0, 182.0, 0.0, 0.0}, // Exact 0-180 is ambiguous.
	{0.0, 182.0, 0.25, 315.5},
	{0.0, 182.0, 0.5, 271.0},
	{0.0, 182.0, 1.0, 182.0},
	{0.0, 270.0, 0.0, 0.0},
	{0.0, 270.0, 0.25, 337.5},
	{0.0, 270.0, 0.5, 315.0},
	{0.0, 270.0, 1.0, 270.0},
	{0.0, 359.0, 0.0, 0.0},
	{0.0, 359.0, 0.25, 359.75},
	{0.0, 359.0, 0.5, 359.5},
	{0.0, 359.0, 1.0, 359.0},
}

func TestInterpolation(t *testing.T) {
	// Forward
	for i, tt := range anglevals {
		res := interp_angle(tt.a0, tt.a1, tt.t)
		if !almosteq_eps(res, tt.at, 1e-15) {
			t.Errorf("%v. interp_angle(%v, %v, %v) => (%v), want %v", i, tt.a0, tt.a1, tt.t, res, tt.at)
		}
	}
	// Backward
	for i, tt := range anglevals {
		res := interp_angle(tt.a1, tt.a0, 1.0-tt.t)
		if !almosteq_eps(res, tt.at, 1e-15) {
			t.Errorf("%v. interp_angle(%v, %v, %v) => (%v), want %v", i, tt.a1, tt.a0, 1.0-tt.t, res, tt.at)
		}
	}
}
