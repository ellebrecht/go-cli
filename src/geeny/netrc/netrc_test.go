package netrc

import (
	"strings"
	"testing"
)

func TestRemove(t *testing.T) {
	n := &Netrc{
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
	}
	n.Remove("random")
	if strings.Compare(n.machines[0].HostName, "random") == 0 {
		t.Fatal("couldnt remove machine")
	}
}

func TestAdd(t *testing.T) {
	n := &Netrc{
		machines: []*Machine{
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
	}
	n.Add("test1", "test2", "test3")
	if len(n.machines) != 3 {
		t.Fatal("couldnt add machine")
	}
	if strings.Compare(n.machines[2].HostName, "test1") != 0 {
		t.Fatal("couldnt add HostName")
	}
	if strings.Compare(n.machines[2].UserName, "test2") != 0 {
		t.Fatal("couldnt add UserName")
	}
	if strings.Compare(n.machines[2].Password, "test3") != 0 {
		t.Fatal("couldnt add Password")
	}
}

func TestEdit(t *testing.T) {
	n := &Netrc{
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
	}
	n.Add("random", "test1Changed", "test2Changed")
	if strings.Compare(n.machines[0].UserName, "test1Changed") != 0 {
		t.Fatal("couldnt edit UserName")
	}
	if strings.Compare(n.machines[0].Password, "test2Changed") != 0 {
		t.Fatal("couldnt edit Password")
	}
}
