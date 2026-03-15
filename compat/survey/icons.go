package survey

// Icon represents a survey icon with text and formatting.
type Icon struct {
	Text   string
	Format string
}

// IconSet configures the icons used by survey prompts.
type IconSet struct {
	HelpInput   Icon
	Error       Icon
	Help        Icon
	Question    Icon
	MarkedOption   Icon
	UnmarkedOption Icon
	SelectFocus    Icon
}

// DefaultIcons returns the default icon set matching survey/v2.
func DefaultIcons() *IconSet {
	return &IconSet{
		HelpInput:      Icon{Text: "?", Format: "default+hb"},
		Error:          Icon{Text: "X", Format: "red"},
		Help:           Icon{Text: "i", Format: "cyan"},
		Question:       Icon{Text: "?", Format: "green+hb"},
		MarkedOption:   Icon{Text: "[x]", Format: "green"},
		UnmarkedOption: Icon{Text: "[ ]", Format: "default"},
		SelectFocus:    Icon{Text: ">", Format: "cyan"},
	}
}
