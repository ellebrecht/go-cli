package tree

import (
	"geeny/cli"
	"geeny/cli/action"
	"geeny/output"
)

// NewCommandTree creates a new geeny cli tree
func NewCommandTree(apiUrl string, connectUrl string) *cli.Command {
	a := action.NewAction(apiUrl, connectUrl)
	tree := &cli.Command{
		Name:        "geeny",
		Summary:     "A cross-platform CLI that interacts with the Geeny API.\nFor more info go to: https://developers.geeny.io/documentation/cli/",
		Action:      a.Version,
		NonCategory: true,
		Options: []*cli.Option{
			{
				Name:         "version",
				Description:  "displays version information",
				Flag:         "v",
				Aliases:      []string{"version"},
				DefaultValue: true,
				Type:         cli.OptionTypeBool,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "swift",
				Summary: "Fast setup of geeny components",
				Commands: []*cli.Command{
					{
						Name:        "dump",
						Summary:     "Dump all your geeny data",
						Action:      a.SwiftDump,
						NonCategory: true,
					},
					{
						Name:        "thing",
						Summary:     "setup a content-type, thing-type, thing and then pair it with an account",
						Action:      a.SwiftThing,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:         "name",
								Description:  "the name of your thing",
								Flag:         "n",
								Aliases:      []string{"name"},
								DefaultValue: "",
							},
						},
					},
					{
						Name:        "project",
						Summary:     "setup an app, client, project with mediation handler & pipeline, clone your project, register a mediation handler & pipeline, authorize your app & pipeline, create a webhook",
						Action:      a.SwiftProject,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:         "name",
								Description:  "the name of your project",
								Flag:         "n",
								Aliases:      []string{"name"},
								DefaultValue: "",
							},
						},
					},
				},
			},
			{
				Name:        "login",
				Summary:     "Login to your Geeny account",
				Interactive: true,
				Action:      a.Login,
				NonCategory: true,
				Options: []*cli.Option{
					{
						Name:        "email",
						Description: "your Geeny account email address",
						Flag:        "e",
						Aliases:     []string{"email"},
					},
					{
						Name:        "password",
						Description: "your Geeny account password",
						Flag:        "p",
						IsSecure:    true,
						Aliases:     []string{"password"},
					},
				},
				Commands: []*cli.Command{
					{
						Name:    "help",
						Summary: "Provides help for this command",
						Action: func(c *cli.Context) (*cli.Meta, error) {
							output.Println("To use the Geeny CLI you need to have an account on the Geeny platform and be logged in. Create an account at https://api.geeny.io/documentation. Then use the same credentials to login on the command line.")
							return nil, nil
						},
						NonCategory: true,
					},
				},
			},
			{
				Name:        "logout",
				Summary:     "Logout of your Geeny account",
				Action:      a.Logout,
				NonCategory: true,
				Commands: []*cli.Command{
					{
						Name:    "help",
						Summary: "Provides help for this command",
						Action: func(c *cli.Context) (*cli.Meta, error) {
							output.Println("You can logout from geeny and login again with a different username.")
							return nil, nil
						},
					},
				},
			},
			{
				Name:    "things",
				Summary: "Manage things registered with the platform",
				Commands: []*cli.Command{
					{
						Name:        "test",
						Summary:     "Test a thing",
						Action:      a.TestThing,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "thing id",
								Description: "your thing UUID",
								Flag:        "tid",
								Aliases:     []string{"thing-identifier"},
							},
							{
								Name:        "content type id",
								Description: "your content type UUID",
								Flag:        "cid",
								Aliases:     []string{"content-type-identifier"},
							},
							{
								Name:         "interval",
								Description:  "time interval for sending data (milliseconds). Defaults to 1000ms",
								Flag:         "i",
								Aliases:      []string{"interval"},
								DefaultValue: 1000,
								Type:         cli.OptionTypeInt,
							},
							{
								Name:         "number",
								Description:  "number of payloads. Defaults to -1 (unlimited)",
								Flag:         "n",
								Aliases:      []string{"number"},
								DefaultValue: -1,
								Type:         cli.OptionTypeInt,
							},
							{
								Name:         "payload",
								Description:  "text payload. Defaults to 'Hello, Geeny!'",
								Flag:         "p",
								Aliases:      []string{"payload"},
								DefaultValue: "Hello, Geeny!",
							},
							{
								Name:         "file",
								Description:  "file continaing payload. if a string payload and file are given, the file takes priority",
								Flag:         "f",
								Aliases:      []string{"file"},
								DefaultValue: "",
							},
							{
								Name:         "topic",
								Description:  "topic for payload. if not provided, defaults to data/:contentTypeID/:thingID",
								Flag:         "t",
								Aliases:      []string{"topic"},
								DefaultValue: "",
							},
							{
								Name:         "subscribe",
								Description:  "subscribe to topic. Defaults to false. (true: -s, -s=true) (false: -s=false)",
								Flag:         "s",
								Aliases:      []string{"subscribe"},
								DefaultValue: false,
								Type:         cli.OptionTypeBool,
							},
							{
								Name:         "endpoint",
								Description:  "custom mqtt endpoint. use 'default' for default endpoint. Defaults to 'default'",
								Flag:         "e",
								Aliases:      []string{"endpoint"},
								DefaultValue: "default",
							},
							{
								Name:         "certificate file (.pem)",
								Description:  "the certificate saved on creation of the thing. if not provided, will attempt to use the key in ~/.geeny/<thing id>/cert.pem",
								Flag:         "c",
								Aliases:      []string{"pem-file"},
								DefaultValue: "",
							},
							{
								Name:         "key file (.key)",
								Description:  "the key saved on creation of the thing. if not provided, will attempt to use the key in ~/.geeny/<thing id>/cert.key",
								Flag:         "k",
								Aliases:      []string{"key-file"},
								DefaultValue: "",
							},
						},
					},
					{
						Name:        "create",
						Summary:     "Creates a thing",
						Action:      a.CreateThing,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "thing type id",
								Description: "your thing type UUID",
								Flag:        "id",
								Aliases:     []string{"thing-type-identifier"},
							},
							{
								Name:         "store secrets",
								Description:  "true stores thing secrets in Geeny's database",
								Flag:         "s",
								Type:         cli.OptionTypeBool,
								DefaultValue: true,
								Aliases:      []string{"store-secrets"},
							},
						},
					},
					{
						Name:        "create",
						Summary:     "Creates a thing",
						Action:      a.CreateThing,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "thing type id",
								Description: "your thing type UUID",
								Flag:        "id",
								Aliases:     []string{"thing-type-identifier"},
							},
							{
								Name:        "store secrets",
								Description: "true stores thing secrets in Geeny's database",
								Flag:        "s",
								Type:        cli.OptionTypeBool,
								Aliases:     []string{"store-secrets"},
							},
						},
					},
				},
			},
			{
				Name:    "projects",
				Summary: "Manage your Geeny projects",
				Commands: []*cli.Command{
					{
						Name:        "deploy",
						Summary:     "Deploy your project changes to Geeny",
						Action:      a.DeployProject,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "message",
								Description: "a short message detailing changes",
								Flag:        "m",
								Aliases:     []string{"message"},
							},
						},
					},
				},
			},
			{
				Name:    "logs",
				Summary: "See activity within the platform",
				Commands: []*cli.Command{
					{
						Name:        "stream",
						Summary:     "Stream logs as they arrive",
						Action:      a.StreamLogs,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "thing ids",
								Description: "your thing UUIDs",
								Flag:        "tids",
								Aliases:     []string{"thing-identifiers"},
							},
							{
								Name:         "service ids",
								Description:  "your service ids (optional), e.g. 'io-mediation', 'iot-gateway'",
								Flag:         "sids",
								DefaultValue: "",
								Aliases:      []string{"service-identifiers"},
							},
						},
					},
					{
						Name:    "help",
						Summary: "Provides help for this command",
						Action: func(c *cli.Context) (*cli.Meta, error) {
							output.Println("As data from your things is processed by platform components they emit messages which can be seen in the logs. You can either ask for a number of most recent log messages or stream all messages continuously. You need to say which thing you are referring to. Also you can narrow the messages to a particular processing stage.")
							return nil, nil
						},
						NonCategory: true,
					},
				},
			},
			{
				Name:    "check",
				Summary: "Check CLI environment",
				Commands: []*cli.Command{
					{
						Name:        "update",
						Summary:     "Check for CLI update",
						Action:      a.CheckUpdate,
						NonCategory: true,
					},
					{
						Name:    "help",
						Summary: "Provides help for this command",
						Action: func(c *cli.Context) (*cli.Meta, error) {
							output.Println("You can ask the Geeny CLI to check various things by interrogating the Geeny API.")
							return nil, nil
						},
						NonCategory: true,
					},
				},
			},
			{
				Name:    "generate",
				Summary: "Generates extensions for an existing project",
				Commands: []*cli.Command{
					{
						Name:        "project",
						Summary:     "Generates a project template",
						Action:      a.GenerateProject,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "project id",
								Description: "The UUID of your project",
								Flag:        "id",
								Aliases:     []string{"project-identifier"},
							},
						},
					},
					{
						Name:        "handler",
						Summary:     "Generates a mediation handler (/mediation/<your handler name>)",
						Action:      a.GenerateHandler,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "name",
								Description: "The name of your mediation handler",
								Flag:        "n",
								Aliases:     []string{"name"},
							},
						},
					},
					{
						Name:        "actor",
						Summary:     "Generates a pipeline actor (/pipeline/<your actor name>)",
						Action:      a.GenerateActor,
						NonCategory: true,
						Options: []*cli.Option{
							{
								Name:        "name",
								Description: "The name of your pipeline actor",
								Flag:        "n",
								Aliases:     []string{"name"},
							},
						},
					},
					{
						Name:    "help",
						Summary: "Provides help for this command",
						Action: func(c *cli.Context) (*cli.Meta, error) {
							output.Println("The Geeny CLI can create skeleton code for projects, mediation handlers and actors. This code can be edited by you and then deployed to the platform.")
							return nil, nil
						},
						NonCategory: true,
					},
				},
			},
			// :)
			{
				Name:   "make",
				Hidden: true,
				Commands: []*cli.Command{
					{
						Name:   "wish",
						Hidden: true,
						Action: func(c *cli.Context) (*cli.Meta, error) {
							output.Println("You need to state your wish as a working Whitespace program.")
							return nil, nil
						},
						NonCategory: true,
					},
				},
			},
		},
	}
	a.Tree = tree
	return tree
}
