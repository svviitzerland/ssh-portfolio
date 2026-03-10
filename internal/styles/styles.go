package styles

import (
	"charm.land/lipgloss/v2"
)

// ── Color Palette ──────────────────────────────────────────────────────
var (
	Primary    = lipgloss.Color("#7C3AED") // Purple
	Secondary  = lipgloss.Color("#06B6D4") // Cyan
	Accent     = lipgloss.Color("#F59E0B") // Amber
	Success    = lipgloss.Color("#10B981") // Emerald
	Danger     = lipgloss.Color("#EF4444") // Red
	Muted      = lipgloss.Color("#6B7280") // Gray
	TextLight  = lipgloss.Color("#F9FAFB") // Almost white
	TextDim    = lipgloss.Color("#9CA3AF") // Gray-400
	BgDark     = lipgloss.Color("#0F172A") // Slate-900
	BgCard     = lipgloss.Color("#1E293B") // Slate-800
	Border     = lipgloss.Color("#334155") // Slate-700
	Highlight  = lipgloss.Color("#A78BFA") // Purple-400
	Pink       = lipgloss.Color("#EC4899") // Pink-500
	Blue       = lipgloss.Color("#3B82F6") // Blue-500
	Green      = lipgloss.Color("#22C55E") // Green-500
)

// ── Base Styles ────────────────────────────────────────────────────────
var (
	// App container
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Header / title bar
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextLight).
			Background(Primary).
			Padding(0, 2).
			Align(lipgloss.Center)

	// Logo / Name big title
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary)

	// Subtitle
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	// Section title
	SectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Accent).
				Padding(0, 0, 1, 0)

	// Normal text
	TextStyle = lipgloss.NewStyle().
			Foreground(TextLight)

	// Dimmed text
	DimTextStyle = lipgloss.NewStyle().
			Foreground(TextDim)

	// Bold text
	BoldTextStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(TextLight)

	// Link style
	LinkStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Underline(true)

	// Card style with border
	CardStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	// Highlighted card
	HighlightCardStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Primary).
				Padding(1, 2)

	// Tag / badge style
	TagStyle = lipgloss.NewStyle().
			Foreground(BgDark).
			Background(Secondary).
			Padding(0, 1)

	// Skill tag
	SkillTagStyle = lipgloss.NewStyle().
			Foreground(TextLight).
			Background(lipgloss.Color("#374151")).
			Padding(0, 1)

	// Active tab
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(BgDark).
			Background(Primary).
			Padding(0, 2)

	// Inactive tab
	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(TextDim).
				Background(lipgloss.Color("#1E293B")).
				Padding(0, 2)

	// Status bar at bottom
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(TextDim).
			Padding(1, 0, 0, 0)

	// Bullet point
	BulletStyle = lipgloss.NewStyle().
			Foreground(Primary)

	// Date range style
	DateStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	// Divider
	DividerStyle = lipgloss.NewStyle().
			Foreground(Border)

	// Scrollbar indicator
	ScrollStyle = lipgloss.NewStyle().
			Foreground(Primary)

	// Achievement highlight
	AchievementStyle = lipgloss.NewStyle().
				Foreground(Accent).
				Bold(true)

	// Company name
	CompanyStyle = lipgloss.NewStyle().
			Foreground(Highlight).
			Bold(true)

	// Position title
	PositionStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	// Project name
	ProjectNameStyle = lipgloss.NewStyle().
				Foreground(Pink).
				Bold(true)

	// Stack/tech
	StackStyle = lipgloss.NewStyle().
			Foreground(Green).
			Italic(true)
)

// ── Helper Functions ───────────────────────────────────────────────────

func Divider(width int) string {
	line := ""
	for i := 0; i < width; i++ {
		line += "─"
	}
	return DividerStyle.Render(line)
}

func DoubleDivider(width int) string {
	line := ""
	for i := 0; i < width; i++ {
		line += "═"
	}
	return DividerStyle.Render(line)
}

func DottedDivider(width int) string {
	line := ""
	for i := 0; i < width; i++ {
		line += "·"
	}
	return DividerStyle.Render(line)
}

func RenderTab(label string, active bool) string {
	if active {
		return ActiveTabStyle.Render(label)
	}
	return InactiveTabStyle.Render(label)
}
