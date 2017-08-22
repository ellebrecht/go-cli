package netrc

import (
	"strings"

	"geeny/util"
	"testing"
)

func TestRead(t *testing.T) {
	file := util.CreateTempFile()
	getPath = func() (string, error) {
		return file.Name(), nil
	}
	r := &Reader{}
	err := r.Read()
	if err != nil {
		t.Fatal("got err:", err)
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		data string
	}{
		{`machine random
  login test1
  password test2
machine api.geeny.io
  login lee@geeny.io
  password eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc
machine connect.geeny.io
  login lee@geeny.io
  password eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc
`}, {`     
	
	
	machine      random
login test1
  password        test2
machine   api.geeny.io


login      lee@geeny.io


  password    eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc
   machine     connect.geeny.io
  login     lee@geeny.io

  
password    eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0OTEzMTg1MTUsImF1ZCI6ImNvbm5lY3QuZ2VlbnkuaW8iLCJpc3MiOiJjb25uZWN0LmdlZW55LmlvIiwic3ViIjoiYTQxMzEyYzctMWI3OC00Zjk3LWI3NTAtMjNlOGY1NTBhN2NiIn0.NSJg2GQEqEBrNB4uTTmKJZ2OB1ljoMeXhGTHNuKwuXc
`},
	}

	for ti, tc := range testCases {
		r := &Reader{
			data: []byte(tc.data),
		}
		machines := r.Parse()
		expected := []*Machine{
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
		}
		if len(expected) != len(machines) {
			t.Fatal(ti, "::", "expected:", expected, "got:", machines)
		}
		for i, m := range machines {
			e := expected[i]
			if strings.Compare(e.HostName, m.HostName) != 0 {
				t.Fatal(ti, "::", "expected:", e.HostName, "got:", m.HostName)
			}
			if strings.Compare(e.UserName, m.UserName) != 0 {
				t.Fatal(ti, "::", "expected:", e.HostName, "got:", m.HostName)
			}
			if strings.Compare(e.Password, m.Password) != 0 {
				t.Fatal(ti, "::", "expected:", e.HostName, "got:", m.HostName)
			}
		}
	}
}
