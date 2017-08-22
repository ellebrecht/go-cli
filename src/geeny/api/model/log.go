package model

// Log represents the data type of a log
type Log struct {
	RequestID      string      `requestId:"id"`
	CreatedAt      string      `createdAt:"identityId"`
	ServiceID      string      `serviceId:"name"`
	Level          string      `json:"level"`
	EventName      string      `json:"eventName"`
	Data           interface{} `json:"data"`
	Runtime        int         `json:"runtime"`
	ParentID       string      `json:"parentId"`
	ThingID        string      `json:"thingId"`
	ThingTypeID    string      `json:"thingTypeId"`
	ApplicationID  string      `json:"applicationId"`
	SubscriptionID string      `json:"subscriptionId"`
	AddonID        string      `json:"addonId"`
	IdentityID     string      `json:"identityId"`
}

// - ValidationInterface

// IsValid validates the data structure
func (l *Log) IsValid() bool {
	return (len(l.RequestID) > 0 &&
		len(l.CreatedAt) > 0 &&
		len(l.ServiceID) > 0 &&
		len(l.Level) > 0 &&
		len(l.EventName) > 0 &&
		l.Data != nil &&
		len(l.ParentID) > 0 &&
		len(l.ThingID) > 0 &&
		len(l.ThingTypeID) > 0 &&
		len(l.ApplicationID) > 0 &&
		len(l.SubscriptionID) > 0 &&
		len(l.AddonID) > 0 &&
		len(l.IdentityID) > 0)
}
