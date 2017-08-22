package action

//TODO
/*
import (
	"errors"
	"io/ioutil"
	"testing"

	model "geeny/api/model"
	cli "geeny/cli"
	mqtt "geeny/mqtt"
	util "geeny/util"
	tu "testing/util"
)

func TestTestThingNoTokenError(t *testing.T) {
	netrc := &tu.MockNetrc{
		Error: errors.New("read error"),
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: netrcReader,
		MQTT:        new(tu.MockMQTT),
	}
	f := util.CreateTempFile()
	_, err := action.TestThing(testTestThingValidContext(f.Name()))
	util.RemoveFile(f)
	if err != netrc.Error {
		t.Fatal("expected error:", netrc.Error, "got:", err)
	}
}

func TestTestThingNetworkError(t *testing.T) {
	apiManager := &tu.MockAPIManager{
		Error: errors.New("network error"),
	}
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		MQTT:        new(tu.MockMQTT),
	}
	f := util.CreateTempFile()
	_, err := action.TestThing(testTestThingValidContext(f.Name()))
	util.RemoveFile(f)
	if err != apiManager.Error {
		t.Fatal("expected error:", apiManager.Error, "got:", err)
	}
}

func TestTestFileError(t *testing.T) {
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: netrcReader,
	}
	_, err := action.TestThing(testTestThingValidContext("badFileName"))
	_, expectedErr := ioutil.ReadFile("badFileName")
	if err == nil || err.Error() != expectedErr.Error() {
		t.Fatal("expected error:", expectedErr, "got:", err)
	}
}

func TestTestMQTTError(t *testing.T) {
	apiManager := testThingValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	mqtt := tu.MockMQTT{
		Error: errors.New("mqtt error"),
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		MQTT:        mqtt,
	}
	f := util.CreateTempFile()
	_, err := action.TestThing(testTestThingValidContext(f.Name()))
	util.RemoveFile(f)
	if err != mqtt.Error {
		t.Fatal("expected error:", mqtt.Error, "got:", err)
	}
}

func TestTestThingArgsError(t *testing.T) {
	tid := "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	cid := "f81d4fae-7dec-11d0-a765-00a0c91e6bf7"
	badtid := "81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	badcid := "81d4fae-7dec-11d0-a765-00a0c91e6bf7"
	validTopic := "data/" + cid + "/" + tid

	f := util.CreateTempFile()
	testTestThingArgs(t, "", cid, 500, 1, "testPayload", f.Name(), validTopic, true, "testEndpoint", "thingID is missing")
	testTestThingArgs(t, tid, "", 500, 1, "testPayload", f.Name(), validTopic, true, "testEndpoint", "contentTypeID is missing")
	testTestThingArgs(t, tid, cid, 0, 1, "testPayload", f.Name(), validTopic, true, "testEndpoint", "interval must > 100")
	testTestThingArgs(t, tid, cid, 500, 0, "testPayload", f.Name(), validTopic, true, "testEndpoint", "number must be > 0 && < 1000, || == -1 ")
	testTestThingArgs(t, tid, cid, 500, 1001, "testPayload", f.Name(), validTopic, true, "testEndpoint", "number must be > 0 && < 1000, || == -1 ")
	testTestThingArgs(t, tid, cid, 500, -2, "testPayload", f.Name(), validTopic, true, "testEndpoint", "number must be > 0 && < 1000, || == -1 ")
	testTestThingArgs(t, tid, cid, 500, 1, "", "", validTopic, true, "testEndpoint", "please provide a file, or a text payload")
	_, ferr := ioutil.ReadFile("badFile")
	testTestThingArgs(t, tid, cid, 500, 1, "", "badFile", validTopic, true, "testEndpoint", ferr.Error())
	testTestThingArgs(t, tid, cid, 500, 1, "testPayload", f.Name(), "badTopic", true, "testEndpoint", "please make sure your topic is in this format: data/:contentTypeID/:thingID")
	testTestThingArgs(t, tid, cid, 500, 1, "testPayload", f.Name(), "a/bad/topic", true, "testEndpoint", "bad topic. the first part of your topic should be data/")
	testTestThingArgs(t, tid, cid, 500, 1, "testPayload", f.Name(), "data/differentid/topic", true, "testEndpoint", "bad content type id in topic. differentid doesn't match "+cid)
	testTestThingArgs(t, tid, cid, 500, 1, "testPayload", f.Name(), "data/"+cid+"/differentid", true, "testEndpoint", "bad thing id in topic. differentid doesn't match "+tid)
	testTestThingArgs(t, tid, badcid, 500, 1, "testPayload", f.Name(), "data/"+badcid+"/"+tid, true, "testEndpoint", badcid+" is not a valid UUID")
	testTestThingArgs(t, badtid, cid, 500, 1, "testPayload", f.Name(), "data/"+cid+"/"+badtid, true, "testEndpoint", badtid+" is not a valid UUID")
	testTestThingArgs(t, tid, cid, 500, 1, "testPayload", f.Name(), validTopic, true, "", "endpoint is missing")
	util.RemoveFile(f)
}

func TestTestThingArgCountError(t *testing.T) {
	action := &Action{
		APIManager:  new(tu.MockAPIManager),
		NetrcReader: new(tu.MockNetrcReader),
	}
	context := &cli.Context{
		Args: []*cli.Option{},
	}
	_, err := action.TestThing(context)
	expected := "expected 8 arguments"
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func TestTestThingSuccess(t *testing.T) {
	apiManager := testThingValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		MQTT:        new(tu.MockMQTT),
	}
	f := util.CreateTempFile()
	_, err := action.TestThing(testTestThingValidContext(f.Name()))
	util.RemoveFile(f)
	if err != nil {
		t.Fatal("got error:", err)
	}
}

// - private

func testTestThingArgs(t *testing.T, tid string, cid string, interval int, number int, payload string, file string, topic string, subscribe bool, endpoint string, expected string) {
	apiManager := testThingValidAPIManager()
	netrc := &tu.MockNetrc{
		Pwd: "test",
	}
	netrcReader := tu.MockNetrcReader{
		Netrc: netrc,
	}
	action := &Action{
		APIManager:  apiManager,
		NetrcReader: netrcReader,
		MQTT:        new(mqtt.MQTT),
	}
	context := &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Name: "thingID", Value: &tid},
			&cli.Option{Name: "contentTypeID", Value: &cid},
			&cli.Option{Name: "interval", Value: &interval},
			&cli.Option{Name: "number", Value: &number},
			&cli.Option{Name: "payload", Value: &payload},
			&cli.Option{Name: "file", Value: &file},
			&cli.Option{Name: "topic", Value: &topic},
			&cli.Option{Name: "subscribe", Value: &subscribe},
			&cli.Option{Name: "endpoint", Value: &endpoint},
		},
	}
	_, err := action.TestThing(context)
	if err == nil || err.Error() != expected {
		t.Fatal("expected error:", expected, "got:", err)
	}
}

func testTestThingValidContext(fileName string) *cli.Context {
	tid := "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"
	cid := "f81d4fae-7dec-11d0-a765-00a0c91e6bf7"
	interval := 500
	number := 5
	payload := "testPayload"
	file := fileName
	topic := "data/" + tid + "/" + cid
	subscribe := true
	endpoint := "default"
	return &cli.Context{
		Args: []*cli.Option{
			&cli.Option{Value: &tid},
			&cli.Option{Value: &cid},
			&cli.Option{Value: &interval},
			&cli.Option{Value: &number},
			&cli.Option{Value: &payload},
			&cli.Option{Value: &file},
			&cli.Option{Value: &topic},
			&cli.Option{Value: &subscribe},
			&cli.Option{Value: &endpoint},
		},
	}
}

func testThingValidAPIManager() *tu.MockAPIManager {
	return &tu.MockAPIManager{
		Payloads: []interface{}{
			&model.Thing{
				Attributes: model.Attributes{
					IOT: model.IOT{
						Certificate: model.Certificate{
							KeyPair: model.KeyPair{
								PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAkUZ9rzaJTOwdy6DDLEiQRMIIC+i4EWic84yMjqNc+1WCr43G\nLm3F1hvB+aJPriGMPyncNgW8dPMSEak75xlOz+TF8Yayo9ztGIHZT3ELyj4c05iZ\nkCP90omYaWYdrGsAX2FLGUrx8oO6I1no3CV6npEdO0y+eNRxfWmv88hpwtljFcen\nQrL+PM0mjIQ5GenIoB1MWNMuDruFjRO0ebFzh6ja/I4hXLS6vix/XrXYaSMrrQUX\n3QoyuwUkdhAiQ+as3KTf7S3iECAHzaOB3VDLrgotGuS2bzXp89tS6DEFHokhsdZj\nsxKN7MQD+e3qTA2pYdCmRJ+7WxWQeXJh+NyWMQIDAQABAoIBAEIrpHXZVmKmLdAh\nUqTqDGR1KKsb/jNXvh2lCCS0PlbwlQ0Xe4yMTY9/pzPofXioKTRgAcDNjhCK4tEO\nj41s1pwU7SBdJSNELu55vpzTtfYRopmyqUehTSzHpZdkfuGY/1tyen1zyI6Y94DA\npDMaEycEnBb1ltB8m5DvQci9oYaxQ9NyQ2WPvU2Kkia5xQGmNuaE8aUgZmrOzDkV\nA1TSv7uwuARiDETs/Z2IZ+IO6VovCX+NaFL6Tk7TgcQjvPqxzC3m8wCPf4XDd6u2\nBzLE4I8pTkasyyrnArAQ4jyf34o4k6ujybZ1gBHveFAGwtbFNQRuV0XBzAtENjrY\nrPB619ECgYEAzNX7RNgMWiGMAlUBp3Kp7cwZNZzaU3W0oEDE0z90oZPgZpix2lCQ\nDvc4UPcSGCkcQtNhobWt5wbbTmUlIMRo7Zd9bauOAF8u9zlpgl95/TqXi6V32akc\nFshvyFmOXQAda/vrTmppIVZb3pZHbGaEkOgvw0hTrN/t4gTKi7z4bQ0CgYEAtY/5\nVQxZUCVOeJn2XQlSZCgjeayT8eNg3ZBJl+4xs3Gl5WuARPa1s9NR9JlfiwkhD2/P\n+L6ibVz3CTCyMHbzlNsuJRpZ4CdiRZgHPOWdoL/gqiekF7xmZO8ZxpfGwFqFFA1y\n59ixU56IPCvWbH3Tkh6BlBhBB5XiTEcpP9qObLUCgYEAuGeLpha6SAobeZ39tznL\nMYGk0Fc9VhWJXxvwFh3yaeQvCS4+L0SH7HE0Ce/kIkwgXSjfpC1jObE6jgEQw8cO\nj5bqHGltlXgbWAYfrnc143t0Iwv8Mb/Ewd4AhGXbfG42Dpjk5dt2ecS9QE4aTsSc\nY7gkB7J9YgzWo7ucfODK5aECgYAmiTZsfMCAKErtghAuUwovw+0zqBOGpbIrlLJq\nEt5trdN+TEDfYlXSoymj9uG2iut/cUX9D8k92Rt90d4gNz8f+x1iNqOY1gMfrlLq\n8/lu4wr2Uo+bkhtRlQYtho1iylOwm+Iln3KTwvQ7gTpzsIk1XFA4dKVozMFJW1k/\n+k18eQKBgF+TgY+DI/Wbs1hmmynmvYu4gkpA5+geisJm4fO8RAC0zpXwI9ilEKY1\n0We1Fftjdg6VoJLQbDrX/XJ16x2CiXuYmoxlZoK1MP2Neu4zu0hJm2A+Vb/DukbT\nTRq3ku6nFu+f5Y7qDqQE4ksjcoBgw+bZi0IZKOdiM6yTPuYVK7ly\n-----END RSA PRIVATE KEY-----\n",
							},
							CertificatePem: "-----BEGIN CERTIFICATE-----\nMIIDWTCCAkGgAwIBAgIULDHgE6WNhARdou8Vrl80wqZoPBMwDQYJKoZIhvcNAQEL\nBQAwTTFLMEkGA1UECwxCQW1hem9uIFdlYiBTZXJ2aWNlcyBPPUFtYXpvbi5jb20g\nSW5jLiBMPVNlYXR0bGUgU1Q9V2FzaGluZ3RvbiBDPVVTMB4XDTE2MTIxNDE1MjYw\nOVoXDTQ5MTIzMTIzNTk1OVowHjEcMBoGA1UEAwwTQVdTIElvVCBDZXJ0aWZpY2F0\nZTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJFGfa82iUzsHcugwyxI\nkETCCAvouBFonPOMjI6jXPtVgq+Nxi5txdYbwfmiT64hjD8p3DYFvHTzEhGpO+cZ\nTs/kxfGGsqPc7RiB2U9xC8o+HNOYmZAj/dKJmGlmHaxrAF9hSxlK8fKDuiNZ6Nwl\nep6RHTtMvnjUcX1pr/PIacLZYxXHp0Ky/jzNJoyEORnpyKAdTFjTLg67hY0TtHmx\nc4eo2vyOIVy0ur4sf1612GkjK60FF90KMrsFJHYQIkPmrNyk3+0t4hAgB82jgd1Q\ny64KLRrktm816fPbUugxBR6JIbHWY7MSjezEA/nt6kwNqWHQpkSfu1sVkHlyYfjc\nljECAwEAAaNgMF4wHwYDVR0jBBgwFoAUGodIPwqfFDS2ch8RI/vktTUk7HIwHQYD\nVR0OBBYEFGozplNXqluC1lyIzU2droJonO/sMAwGA1UdEwEB/wQCMAAwDgYDVR0P\nAQH/BAQDAgeAMA0GCSqGSIb3DQEBCwUAA4IBAQCqr5d/7pXTz6cxESoC8ZLdXTEi\nIc1z4kzF9SgEVRIVNGCN3b1xHoEVpADCUkLy1P9f95xX1cAz1Whr4LJQGMSUjGta\n5nOGW+l+J8JaBHU3/UxUFYVuJGTnAjqs8DGbf810PDT7tJu+KpV8MgJfoJjQ0E/L\nenGEkj9teor9mcq5HyCo3CYCx9ANMO7joSBtatgHk90RgrtqJc90Sq4Crx5R47Ka\nGL5aWJkhfNHRj3yhOWGMs3f0wzw+twKB24hc4DSVzOSw5r5EBRwwp5G/Fsq0MzQw\nAxfa/eEsIOPNaI3fbr7zA5ZZLXwL98K9hkPjXkNXivueJlaSaHR/IL8q7BCD\n-----END CERTIFICATE-----\n",
						},
					},
				},
			},
		},
	}
}
*/
