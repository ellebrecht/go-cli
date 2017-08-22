package netrc

import (
	"io/ioutil"
	"strings"
	"testing"

	"geeny/util"
)

func TestWrite(t *testing.T) {
	testCases := []struct {
		n *Netrc
	}{
		{
			&Netrc{
				machines: []*Machine{
					&Machine{
						HostName: "random",
						UserName: "test1",
						Password: "test2",
					},
					&Machine{
						HostName: "api.geeny.io",
						UserName: "lee@geeny.io",
						Password: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc",
					},
					&Machine{
						HostName: "connect.geeny.io",
						UserName: "lee@geeny.io",
						Password: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc",
					},
				},
			},
		},
	}

	for ti, tc := range testCases {
		file := util.CreateTempFile()
		getPath = func() (string, error) {
			return file.Name(), nil
		}
		w := &Writer{}
		err := w.Write(tc.n)
		if err != nil {
			t.Fatal("got:", err)
		}
		expected := `machine random
  login test1
  password test2
machine api.geeny.io
  login lee@geeny.io
  password eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc
machine connect.geeny.io
  login lee@geeny.io
  password eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc
`
		fileWritten, err := ioutil.ReadFile(file.Name())
		if err != nil {
			t.Fatal("got:", err)
		}

		if strings.Compare(expected, string(fileWritten)) != 0 {
			t.Fatal(ti, "::", "expected:\n", expected, "\ngot:\n", string(fileWritten))
		}
	}
}
