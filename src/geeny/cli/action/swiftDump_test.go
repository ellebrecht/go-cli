package action

//TODO
/*
func TestSwiftDumpSuccess(t *testing.T) {
	apiManager := testSwiftDumpValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
	}
	_, err := action.SwiftDump(nil)
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// - private

func testSwiftDumpValidAPIManager() *tu.MockAPIManager {
	return &tu.MockAPIManager{
		Payloads: []interface{}{
			&[]model.ContentType{},
			&[]model.ThingType{},
			&[]model.Thing{},
			&[]model.Project{},
			&[]model.App{},
			&[]model.Addon{},
		},
	}
}
*/
