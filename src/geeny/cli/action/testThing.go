package action

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"time"

	model "geeny/api/model"
	cli "geeny/cli"
	log "geeny/log"
	mqtt "geeny/mqtt"
	output "geeny/output"
)

// TestThing cli action
func (a *Action) TestThing(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	if c.Count() < 8 {
		return nil, errors.New("expected 8 arguments")
	}
	thingID, err := c.GetStringForFlag("tid")
	if err != nil {
		return nil, err
	}
	contentTypeID, err := c.GetStringForFlag("cid")
	if err != nil {
		return nil, err
	}
	interval, err := c.GetIntForFlag("i")
	if err != nil {
		return nil, err
	}
	number, err := c.GetIntForFlag("n")
	if err != nil {
		return nil, err
	}

	payload, _ := c.GetStringForFlag("p")
	fileName, _ := c.GetStringForFlag("f")
	if payload == nil && fileName == nil {
		return nil, errors.New("please provide a file, or a text payload")
	}
	if fileName != nil { // filename overrides text payload
		file, err := ioutil.ReadFile(*fileName)
		if err != nil {
			return nil, err
		}
		*payload = string(file)
	}

	topic, _ := c.GetStringForFlag("t")
	subscribe, err := c.GetBoolForFlag("s")
	if err != nil {
		return nil, err
	}
	endpoint, err := c.GetStringForFlag("e")
	if err != nil {
		return nil, err
	}

	pemFile, _ := c.GetStringForFlag("c")
	keyFile, _ := c.GetStringForFlag("k")

	return a.testThing(*thingID, *contentTypeID, *interval, *number, *payload, topic, *subscribe, *endpoint, pemFile, keyFile)
}

// - private

func (a *Action) testThing(thingID string, contentTypeID string, interval int, number int, payload string, topic *string, subscribe bool, endpoint string, pemFile *string, keyFile *string) (*cli.Meta, error) {
	// set defaults
	if topic == nil {
		t := "data/" + contentTypeID + "/" + thingID
		topic = &t
	}

	// validate
	if interval <= 100 {
		return nil, errors.New("interval must > 100")
	}
	if (number <= 0 && number != -1) || number == 0 || number > 1000 {
		return nil, errors.New("number must be > 0 && < 1000, || == -1 ")
	}
	err := a.MQTT.ValidateTopic(*topic, contentTypeID, thingID)
	if err != nil {
		return nil, err
	}

	// get thing
	cmd, err := a.Tree.CommandForPath([]string{"geeny", "things", "get"}, 0)
	if err != nil {
		return nil, err
	}
	err = cmd.SetValueForOptionWithFlag(&thingID, "id")
	if err != nil {
		return nil, err
	}
	var meta *cli.Meta
	output.DisableForAction(func() {
		meta, err = cmd.Exec()
	})
	if err != nil {
		return nil, err
	}
	thing := &model.Thing{}
	err = meta.UnmarshalRawJSON(thing)
	if err != nil {
		return nil, err
	}

	// generate certificate
	key, err := a.getKeyFile(keyFile, thing)
	if err != nil {
		return nil, err
	}
	pem, err := a.getPemFile(pemFile, thing)
	if err != nil {
		return nil, err
	}
	cert, err := tls.X509KeyPair(pem, key)
	if err != nil {
		return nil, err
	}

	// configure endpoint
	if strings.Compare(endpoint, "default") == 0 {
		endpoint = thing.Attributes.IOT.IOTEndpoint
	}

	// configure mqtt
	a.MQTT.SetEndpoint(endpoint)
	a.MQTT.SetClientID(thing.ID)
	a.MQTT.SetMessageHandler(func(client mqtt.MQTT, topic string, payload string) {
		output.Println("TOPIC:", topic)
		output.Println("MSG:", payload)
	})
	a.MQTT.SetConnectHandler(func(client mqtt.MQTT) {
		output.Println("mqtt connected")
	})
	a.MQTT.SetConnectionLostHandler(func(client mqtt.MQTT, err error) {
		panic("mqtt connection lost: " + err.Error())
	})
	a.MQTT.SetKeepAlive(time.Duration(interval))
	a.MQTT.SetCertificate(cert)

	// connect to mqtt
	err = a.MQTT.Connect()
	if err != nil {
		return nil, err
	}

	// subscribe to mqtt, send data, unsubscribe
	if subscribe == true {
		err = a.MQTT.Subscribe(*topic)
		if err != nil {
			return nil, err
		}
	}

	if number == -1 {
		a.MQTT.SendDataForever(payload, *topic, time.Duration(interval))
	} else {
		a.MQTT.SendData(payload, *topic, number, time.Duration(interval))
	}

	if subscribe == true {
		err = a.MQTT.Unsubscribe(*topic)
		if err != nil {
			return nil, err
		}
	}
	output.Println("mqtt finished and disconnected")
	return nil, nil
}

// returns pem file
// 1. tries to read customPath (if provided)
// 2. tries to read from ~/.geeny/<thingid>/cert.pem
// 3. tries to read from thing.Attributes.IOT.Certificate.CertificatePem
// 4. fails if none of the above work
func (a *Action) getPemFile(customPath *string, thing *model.Thing) ([]byte, error) {
	var pem []byte
	var err error
	if customPath != nil {
		pem, err = ioutil.ReadFile(*customPath)
		if err != nil {
			return nil, err
		}
		log.Trace("read cert.pem from: ", *customPath)
	} else {
		usr, err := user.Current()
		if err != nil {
			return nil, errors.New("Unable to determine your home directory")
		}
		path := (usr.HomeDir + string(os.PathSeparator) +
			".geeny" + string(os.PathSeparator) +
			"things" + string(os.PathSeparator) +
			thing.ID + string(os.PathSeparator) + "cert.pem")
		pem, err = ioutil.ReadFile(path)
		if err != nil {
			if len(thing.Attributes.IOT.Certificate.CertificatePem) == 0 {
				return nil, errors.New("cert.pem not found, or provided")
			}
			pem = []byte(thing.Attributes.IOT.Certificate.CertificatePem)
			log.Trace("read cert.pem from thing")
		} else {
			log.Trace("read cert.pem from: ", path)
		}
	}
	return pem, nil
}

// returns key file
// 1. tries to read customPath (if provided)
// 2. tries to read from ~/.geeny/<thingid>/cert.key
// 3. tries to read from thing.Attributes.IOT.Certificate.KeyPair.PrivateKey
// 4. fails if none of the above work
func (a *Action) getKeyFile(customPath *string, thing *model.Thing) ([]byte, error) {
	var key []byte
	var err error
	if customPath != nil {
		key, err = ioutil.ReadFile(*customPath)
		if err != nil {
			return nil, err
		}
		log.Trace("read cert.key from: ", *customPath)
	} else {
		usr, err := user.Current()
		if err != nil {
			return nil, errors.New("Unable to determine your home directory")
		}
		path := (usr.HomeDir + string(os.PathSeparator) +
			".geeny" + string(os.PathSeparator) +
			"things" + string(os.PathSeparator) +
			thing.ID + string(os.PathSeparator) + "cert.key")
		key, err = ioutil.ReadFile(path)
		if err != nil {
			if len(thing.Attributes.IOT.Certificate.KeyPair.PrivateKey) == 0 {
				return nil, errors.New("cert.key not found, or provided")
			}
			key = []byte(thing.Attributes.IOT.Certificate.KeyPair.PrivateKey)
			log.Trace("read cert.key from thing")
		} else {
			log.Trace("read cert.key from: ", path)
		}
	}
	return key, nil
}
