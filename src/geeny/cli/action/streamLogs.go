package action

import (
	"errors"
	"time"

	model "geeny/api/model"
	cli "geeny/cli"
	json "geeny/json"
	output "geeny/output"
)

// StreamLogs cli action
func (a *Action) StreamLogs(c *cli.Context) (*cli.Meta, error) {
	if c == nil {
		panic("context is nil")
	}
	if c.Count() < 1 {
		return nil, errors.New("expected 1 argument")
	}
	thingIDs, err := c.GetStringForFlag("tids")
	if err != nil {
		return nil, err
	}
	serviceIDs, _ := c.GetStringForFlag("sids")
	if serviceIDs == nil {
		emptyString := "*"
		serviceIDs = &emptyString
	}

	return a.streamLogs(*thingIDs, *serviceIDs)
}

// - private

func (a *Action) streamLogs(thingIDs string, serviceIDs string) (*cli.Meta, error) {
	output.Printf("logs will appear when they are received\n\n")

	cmd, err := a.Tree.CommandForPath([]string{"geeny", "logs", "search"}, 0)
	if err != nil {
		return nil, err
	}

	printedLogs := []model.Log{}
	for true {
		// search logs
		var meta *cli.Meta
		offset := 0
		number := -1
		err = cmd.SetValueForOptionWithFlag(&offset, "o")
		if err != nil {
			return nil, err
		}
		err = cmd.SetValueForOptionWithFlag(&serviceIDs, "sids")
		if err != nil {
			return nil, err
		}
		err = cmd.SetValueForOptionWithFlag(&number, "n")
		if err != nil {
			return nil, err
		}
		err = cmd.SetValueForOptionWithFlag(&thingIDs, "tids")
		if err != nil {
			return nil, err
		}
		output.DisableForAction(func() {
			meta, err = cmd.Exec()
		})
		if err != nil {
			return nil, err
		}
		logs := []model.Log{}
		err = meta.UnmarshalRawJSON(&logs)
		if err != nil {
			return nil, err
		}

		// if there are logs
		if len(logs) > 0 {
			// find the new logs
			newLogs := findNewLogs(printedLogs, logs)
			if len(newLogs) > 0 {
				str, err := json.PrettyJSON(newLogs)
				if err != nil {
					return nil, err
				}
				// append all logs that are printed to screen
				printedLogs = append(printedLogs, newLogs...)
				output.Println(*str)
			}
		}
		time.Sleep(time.Second * 5)
	}
	return nil, nil
}

func findNewLogs(a []model.Log, b []model.Log) []model.Log {
	newLogs := []model.Log{}
	for _, bLog := range b {
		found := false
		for _, aLog := range a {
			if bLog.RequestID == aLog.RequestID {
				found = true
				break
			}
		}
		if found == false {
			newLogs = append(newLogs, bLog)
		}
	}
	return newLogs
}
