package netrc

// @see https://www.ibm.com/support/knowledgecenter/en/ssw_aix_71/com.ibm.aix.files/netrc.htm
type Machine struct {
	HostName        string
	DefaultHostName string
	UserName        string
	Password        string
	AccountPassword string
	MacroName       string
}
