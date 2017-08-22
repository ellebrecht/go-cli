package output

import (
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	config "geeny/config"
	log "geeny/log"
)

type Spinner struct {
	dest   *os.File
	text   string
	pad    int
	style  []string
	color  []string
	delay  time.Duration
	active bool
	spin   bool
	lock   *sync.RWMutex
}

func NewSpinner() *Spinner {
	if !config.CurrentInt.SpinnerOutput {
		return SilentSpinner()
	}
	if config.CurrentExt.Spinner {
		return &Spinner{
			dest: os.Stdout,
			text: " ",
			pad:  0,
			style: []string{
				"█     ", " █    ", "  █   ", "   █  ", "    █ ", "     █", "    █ ", "   █  ", "  █   ", " █    ",
				"▓     ", " ▓    ", "  ▓   ", "   ▓  ", "    ▓ ", "     ▓", "    ▓ ", "   ▓  ", "  ▓   ", " ▓    ",
				"▒     ", " ▒    ", "  ▒   ", "   ▒  ", "    ▒ ", "     ▒", "    ▒ ", "   ▒  ", "  ▒   ", " ▒    ",
				"░     ", " ░    ", "  ░   ", "   ░  ", "    ░ ", "     ░", "    ░ ", "   ░  ", "  ░   ", " ░    ",
				"▒     ", " ▒    ", "  ▒   ", "   ▒  ", "    ▒ ", "     ▒", "    ▒ ", "   ▒  ", "  ▒   ", " ▒    ",
				"▓     ", " ▓    ", "  ▓   ", "   ▓  ", "    ▓ ", "     ▓", "    ▓ ", "   ▓  ", "  ▓   ", " ▓    ",
			},
			color:  []string{"\x1b[34;1m", "\x1b[36;1m", "\x1b[37;1m", "\x1b[33;1m", "\x1b[32;1m", "\x1b[35;1m", "\x1b[31;1m"},
			delay:  70 * time.Millisecond,
			active: false,
			spin:   true,
			lock:   &sync.RWMutex{},
		}
	}
	return NoSpinner()
}

func NoSpinner() *Spinner {
	return &Spinner{
		dest:   os.Stdout,
		text:   " ",
		pad:    0,
		style:  nil,
		color:  nil,
		delay:  0,
		active: false,
		spin:   false,
		lock:   &sync.RWMutex{},
	}
}

func SilentSpinner() *Spinner {
	return &Spinner{
		dest:   nil,
		text:   "",
		pad:    0,
		style:  nil,
		color:  nil,
		delay:  0,
		active: false,
		spin:   false,
		lock:   &sync.RWMutex{},
	}
}

func (spinner *Spinner) Start() {
	if spinner.active {
		return
	}
	if spinner.spin {
		go spinner.run()
	}
}

func (spinner *Spinner) Stop(erase bool) {
	spinner.lock.Lock()
	spinner.active = false
	spinner.Clear(erase)
	spinner.lock.Unlock()
}

func (spinner *Spinner) Clear(erase bool) {
	suffix := "\n"
	if erase {
		suffix = "\r" + strings.Repeat(" ", utf8.RuneCountInString(spinner.text)) + "\r"
	}
	if spinner.spin {
		Fprintf(spinner.dest, "\r%s%s%s", spinner.text, strings.Repeat(" ", spinner.pad+utf8.RuneCountInString(spinner.style[0])+1), suffix)
	} else {
		Fprintf(spinner.dest, "%s", suffix)
	}
}

func (spinner *Spinner) Text(newline bool, text ...string) {
	if spinner.dest == nil {
		return
	}
	newText := strings.Join(text, " ")
	log.Info(newText)
	spinner.lock.Lock()
	if newline {
		spinner.Clear(false)
	} else {
		p := utf8.RuneCountInString(spinner.text) - utf8.RuneCountInString(newText)
		if p > 0 {
			spinner.pad = p
		}
	}
	spinner.text = newText
	if !spinner.spin {
		Fprintf(spinner.dest, "\r%s%s ", spinner.text, spinner.padding())
	}
	spinner.lock.Unlock()
}

func (spinner *Spinner) run() {
	spinner.active = true
	for spinner.active {
		for _, c := range spinner.color {
			for _, s := range spinner.style {
				if spinner.active {
					spinner.lock.Lock()
					Fprintf(spinner.dest, "\r%s%s\x1b[0m %s%s ", c, s, spinner.text, spinner.padding())
					spinner.lock.Unlock()
				}
				if spinner.active {
					time.Sleep(spinner.delay)
				}
			}
		}
	}
}

func (spinner *Spinner) padding() string {
	if spinner.pad > 0 {
		padding := strings.Repeat(" ", spinner.pad)
		spinner.pad = 0
		return padding
	}
	return ""
}
