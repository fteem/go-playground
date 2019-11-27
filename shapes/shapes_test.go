package shapes

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := rectangle.Perimeter()
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {
	checkArea := func(t *testing.T, shape Shape, want float64) {
		t.Helper()

		got := shape.Area()
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	}

	areaTests := []struct {
		name  string
		shape Shape
		want  float64
	}{
		{
			name:  "Rectangles",
			shape: Rectangle{10.0, 10.0},
			want:  100.0,
		},
		{
			name:  "Circles",
			shape: Circle{10.0},
			want:  314.1592653589793,
		},
		{
			name:  "Triangles",
			shape: Triangle{12, 6},
			want:  36.0,
		},
	}

	for _, test := range areaTests {
		checkArea(t, test.shape, test.want)
	}
}
