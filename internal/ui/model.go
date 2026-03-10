package ui

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"ssh-portfolio/internal/data"
	"ssh-portfolio/internal/styles"
)

// ── Pages ──────────────────────────────────────────────────────────────

type Page int

const (
	PageHome Page = iota
	PageAbout
	PageExperience
	PageProjects
	PageSkills
	PageAchievements
	PageContact
)

var pageNames = []string{
	"  Home",
	"  About",
	"  Experience",
	"  Projects",
	"  Skills",
	"  Achievements",
	"  Contact",
}

var pageIcons = []string{
	"◆",
	"●",
	"▶",
	"◈",
	"⬟",
	"★",
	"◉",
}

// ── Model ──────────────────────────────────────────────────────────────

type Model struct {
	cv       *data.CV
	page     Page
	width    int
	height   int
	scrollY  int
	maxScroll int
	username string
	ready    bool
}

func NewModel(cv *data.CV, username string) Model {
	return Model{
		cv:       cv,
		page:     PageHome,
		width:    80,
		height:   24,
		username: username,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// ── Update ─────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.scrollY = 0

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// Tab navigation
		case "tab", "l", "right":
			m.page = (m.page + 1) % Page(len(pageNames))
			m.scrollY = 0
		case "shift+tab", "h", "left":
			m.page = (m.page - 1 + Page(len(pageNames))) % Page(len(pageNames))
			m.scrollY = 0

		// Direct page jumps
		case "1":
			m.page = PageHome
			m.scrollY = 0
		case "2":
			m.page = PageAbout
			m.scrollY = 0
		case "3":
			m.page = PageExperience
			m.scrollY = 0
		case "4":
			m.page = PageProjects
			m.scrollY = 0
		case "5":
			m.page = PageSkills
			m.scrollY = 0
		case "6":
			m.page = PageAchievements
			m.scrollY = 0
		case "7":
			m.page = PageContact
			m.scrollY = 0

		// Scrolling
		case "j", "down":
			m.scrollY++
		case "k", "up":
			if m.scrollY > 0 {
				m.scrollY--
			}
		case "g":
			m.scrollY = 0
		case "G":
			m.scrollY = 9999
		}
	}

	return m, nil
}

// ── View ───────────────────────────────────────────────────────────────

func (m Model) View() tea.View {
	if !m.ready {
		v := tea.NewView("\n  Loading...")
		v.AltScreen = true
		return v
	}

	w := m.width
	if w > 100 {
		w = 100
	}

	contentWidth := w - 4

	// Build layout
	var b strings.Builder

	// Top bar
	b.WriteString(m.renderTopBar(w))
	b.WriteString("\n")

	// Tab bar
	b.WriteString(m.renderTabBar(w))
	b.WriteString("\n")
	b.WriteString(styles.Divider(w))
	b.WriteString("\n")

	// Page content
	content := m.renderPage(contentWidth)

	// Apply scrolling
	lines := strings.Split(content, "\n")
	availableHeight := m.height - 7 // top bar + tab bar + divider + status bar + margins

	if m.scrollY > len(lines)-availableHeight {
		m.scrollY = max(0, len(lines)-availableHeight)
	}

	endLine := m.scrollY + availableHeight
	if endLine > len(lines) {
		endLine = len(lines)
	}

	if m.scrollY < len(lines) {
		visibleLines := lines[m.scrollY:endLine]
		b.WriteString(strings.Join(visibleLines, "\n"))
	}

	b.WriteString("\n")

	// Status bar
	b.WriteString(m.renderStatusBar(w, len(lines), availableHeight))

	// Center the whole thing
	fullContent := b.String()
	if m.width > 100 {
		fullContent = lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(fullContent)
	}

	v := tea.NewView(fullContent)
	v.AltScreen = true
	return v
}

// ── Top Bar ────────────────────────────────────────────────────────────

func (m Model) renderTopBar(width int) string {
	left := styles.LogoStyle.Render("  " + m.cv.Basics.Name)
	right := styles.DimTextStyle.Render(fmt.Sprintf("ssh visitor: %s ", m.username))

	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	return left + strings.Repeat(" ", gap) + right
}

// ── Tab Bar ────────────────────────────────────────────────────────────

func (m Model) renderTabBar(width int) string {
	var tabs []string
	for i, name := range pageNames {
		label := fmt.Sprintf(" %s%s ", pageIcons[i], name)
		if width < 80 {
			// Short labels for narrow terminals
			shortNames := []string{" ◆ ", " ● ", " ▶ ", " ◈ ", " ⬟ ", " ★ ", " ◉ "}
			label = shortNames[i]
		}
		tabs = append(tabs, styles.RenderTab(label, Page(i) == m.page))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// ── Status Bar ─────────────────────────────────────────────────────────

func (m Model) renderStatusBar(width, totalLines, visibleLines int) string {
	nav := styles.DimTextStyle.Render("  ←/→ navigate • ↑/↓ scroll • 1-7 jump • q quit")

	scrollInfo := ""
	if totalLines > visibleLines {
		pct := 0
		if totalLines-visibleLines > 0 {
			pct = m.scrollY * 100 / (totalLines - visibleLines)
		}
		if pct > 100 {
			pct = 100
		}
		scrollInfo = styles.ScrollStyle.Render(fmt.Sprintf(" %d%% ", pct))
	}

	gap := width - lipgloss.Width(nav) - lipgloss.Width(scrollInfo)
	if gap < 0 {
		gap = 0
	}

	return nav + strings.Repeat(" ", gap) + scrollInfo
}

// ── Page Router ────────────────────────────────────────────────────────

func (m Model) renderPage(width int) string {
	switch m.page {
	case PageHome:
		return m.renderHome(width)
	case PageAbout:
		return m.renderAbout(width)
	case PageExperience:
		return m.renderExperience(width)
	case PageProjects:
		return m.renderProjects(width)
	case PageSkills:
		return m.renderSkills(width)
	case PageAchievements:
		return m.renderAchievements(width)
	case PageContact:
		return m.renderContact(width)
	default:
		return ""
	}
}

// ── Page: Home ─────────────────────────────────────────────────────────

func (m Model) renderHome(width int) string {
	var b strings.Builder

	// ASCII Art Name
	ascii := `
   ███████╗ █████╗ ██████╗ ██╗  ██╗ █████╗ ███╗   ██╗
   ██╔════╝██╔══██╗██╔══██╗██║  ██║██╔══██╗████╗  ██║
   █████╗  ███████║██████╔╝███████║███████║██╔██╗ ██║
   ██╔══╝  ██╔══██║██╔══██╗██╔══██║██╔══██║██║╚██╗██║
   ██║     ██║  ██║██║  ██║██║  ██║██║  ██║██║ ╚████║
   ╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝`

	asciiStyled := lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Render(ascii)

	b.WriteString(asciiStyled)
	b.WriteString("\n\n")

	// Title
	title := lipgloss.NewStyle().
		Foreground(styles.Secondary).
		Bold(true).
		Render("   " + m.cv.Basics.Label)
	b.WriteString(title)
	b.WriteString("\n\n")

	// Divider
	b.WriteString("   ")
	b.WriteString(styles.DottedDivider(width - 6))
	b.WriteString("\n\n")

	// Summary in a card
	summaryCard := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.Border).
		Padding(1, 2).
		Width(width - 4).
		Render(
			styles.DimTextStyle.Render("   \"") +
				styles.TextStyle.Render(m.cv.Basics.Summary) +
				styles.DimTextStyle.Render("\""),
		)
	b.WriteString("  ")
	b.WriteString(summaryCard)
	b.WriteString("\n\n")

	// Quick stats
	b.WriteString("  ")
	b.WriteString(styles.SectionTitleStyle.Render("  Quick Overview"))
	b.WriteString("\n")

	stats := []struct {
		icon  string
		label string
		value string
	}{
		{"▸", "Projects", fmt.Sprintf("%d shipped products", len(m.cv.Projects))},
		{"▸", "Experience", fmt.Sprintf("%d roles across companies", len(m.cv.Work))},
		{"▸", "Skills", fmt.Sprintf("%d+ technologies", len(m.cv.Skills))},
		{"▸", "Achievements", fmt.Sprintf("%d awards & competitions", len(m.cv.Achievements))},
		{"▸", "Community", fmt.Sprintf("%d organizations", len(m.cv.ExperiencesInOrganization))},
	}

	for _, s := range stats {
		line := fmt.Sprintf("   %s %s  %s",
			styles.BulletStyle.Render(s.icon),
			styles.BoldTextStyle.Render(s.label),
			styles.DimTextStyle.Render(s.value),
		)
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Navigation hint
	hint := lipgloss.NewStyle().
		Foreground(styles.Muted).
		Italic(true).
		Render("   Use ← → or Tab to navigate sections, ↑ ↓ to scroll")
	b.WriteString(hint)
	b.WriteString("\n")

	return b.String()
}

// ── Page: About ────────────────────────────────────────────────────────

func (m Model) renderAbout(width int) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.SectionTitleStyle.Render("  ● About Me"))
	b.WriteString("\n\n")

	// Bio
	bioCard := styles.CardStyle.Width(width - 4).Render(
		styles.TextStyle.Render(m.cv.Basics.Summary),
	)
	b.WriteString("  ")
	b.WriteString(bioCard)
	b.WriteString("\n\n")

	// Education
	b.WriteString(styles.SectionTitleStyle.Render("  Education"))
	b.WriteString("\n")

	for _, edu := range m.cv.Education {
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			styles.BulletStyle.Render("▸"),
			styles.BoldTextStyle.Render(edu.Institution),
		))
		b.WriteString(fmt.Sprintf("     %s · %s\n",
			styles.PositionStyle.Render(edu.Area),
			styles.DateStyle.Render(edu.StartDate+" — "+edu.EndDate),
		))
		b.WriteString(fmt.Sprintf("     %s\n\n",
			styles.DimTextStyle.Render(edu.StudyType),
		))
	}

	// Certifications
	b.WriteString(styles.SectionTitleStyle.Render("  Certifications"))
	b.WriteString("\n")

	for _, cert := range m.cv.Certifications {
		b.WriteString(fmt.Sprintf("  %s  %s  %s  %s\n",
			styles.BulletStyle.Render("▸"),
			styles.BoldTextStyle.Render(cert.Name),
			styles.DateStyle.Render("("+cert.Date+")"),
			styles.StackStyle.Render("Score: "+cert.Score),
		))
	}
	b.WriteString("\n")

	// Talks
	if len(m.cv.Talks) > 0 {
		b.WriteString(styles.SectionTitleStyle.Render("  Speaking"))
		b.WriteString("\n")

		for _, talk := range m.cv.Talks {
			talkCard := styles.HighlightCardStyle.Width(width - 4).Render(
				styles.BoldTextStyle.Render(talk.Title) + "\n" +
					styles.PositionStyle.Render(talk.Event) + " · " +
					styles.DateStyle.Render(talk.Date) + "\n\n" +
					styles.DimTextStyle.Render(talk.Summary),
			)
			b.WriteString("  ")
			b.WriteString(talkCard)
			b.WriteString("\n")
		}
	}

	// Organizations
	b.WriteString("\n")
	b.WriteString(styles.SectionTitleStyle.Render("  Community & Organizations"))
	b.WriteString("\n")

	for _, org := range m.cv.ExperiencesInOrganization {
		b.WriteString(fmt.Sprintf("  %s  %s  %s\n",
			styles.BulletStyle.Render("▸"),
			styles.CompanyStyle.Render(org.Organization),
			styles.DateStyle.Render(org.StartDate+" — "+org.EndDate),
		))
		b.WriteString(fmt.Sprintf("     %s\n",
			styles.PositionStyle.Render(org.Position),
		))
		if org.Summary != "" {
			b.WriteString(fmt.Sprintf("     %s\n",
				styles.DimTextStyle.Render(org.Summary),
			))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// ── Page: Experience ───────────────────────────────────────────────────

func (m Model) renderExperience(width int) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.SectionTitleStyle.Render("  ▶ Work Experience"))
	b.WriteString("\n\n")

	for i, work := range m.cv.Work {
		// Company header
		header := fmt.Sprintf("  %s  %s",
			styles.CompanyStyle.Render(work.Company),
			styles.DateStyle.Render(work.StartDate+" — "+work.EndDate),
		)
		b.WriteString(header)
		b.WriteString("\n")

		// Position
		b.WriteString(fmt.Sprintf("     %s\n\n",
			styles.PositionStyle.Render(work.Position),
		))

		// Highlights
		for _, hl := range work.Highlights {
			wrapped := wordWrap(hl, width-10)
			lines := strings.Split(wrapped, "\n")
			b.WriteString(fmt.Sprintf("     %s %s\n",
				styles.BulletStyle.Render("▸"),
				styles.TextStyle.Render(lines[0]),
			))
			for _, line := range lines[1:] {
				b.WriteString(fmt.Sprintf("       %s\n",
					styles.TextStyle.Render(line),
				))
			}
		}

		if i < len(m.cv.Work)-1 {
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(styles.DottedDivider(width - 4))
			b.WriteString("\n\n")
		}
	}

	return b.String()
}

// ── Page: Projects ─────────────────────────────────────────────────────

func (m Model) renderProjects(width int) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.SectionTitleStyle.Render("  ◈ Projects"))
	b.WriteString("\n\n")

	for _, proj := range m.cv.Projects {
		// Project card
		var cardContent strings.Builder

		// Name and date
		nameAndDate := fmt.Sprintf("%s  %s",
			styles.ProjectNameStyle.Render(proj.Name),
			styles.DateStyle.Render(proj.Date),
		)
		if proj.EndDate != "" {
			nameAndDate = fmt.Sprintf("%s  %s",
				styles.ProjectNameStyle.Render(proj.Name),
				styles.DateStyle.Render(proj.Date+" — "+proj.EndDate),
			)
		}
		cardContent.WriteString(nameAndDate)
		cardContent.WriteString("\n")

		// URL
		if proj.URL != "" {
			cardContent.WriteString(styles.LinkStyle.Render("  "+proj.URL))
			cardContent.WriteString("\n")
		}
		cardContent.WriteString("\n")

		// Summary
		if proj.Summary != "" {
			wrapped := wordWrap(proj.Summary, width-12)
			cardContent.WriteString(styles.TextStyle.Render(wrapped))
			cardContent.WriteString("\n")
		}

		// Highlights
		if len(proj.Highlights) > 0 {
			cardContent.WriteString("\n")
			for _, hl := range proj.Highlights {
				wrapped := wordWrap(hl, width-14)
				lines := strings.Split(wrapped, "\n")
				cardContent.WriteString(fmt.Sprintf("%s %s\n",
					styles.BulletStyle.Render("▸"),
					styles.DimTextStyle.Render(lines[0]),
				))
				for _, line := range lines[1:] {
					cardContent.WriteString(fmt.Sprintf("  %s\n",
						styles.DimTextStyle.Render(line),
					))
				}
			}
		}

		// Stack
		if proj.Stack != "" {
			cardContent.WriteString("\n")
			cardContent.WriteString(styles.StackStyle.Render("⚙ " + proj.Stack))
		}

		card := styles.CardStyle.Width(width - 4).Render(cardContent.String())
		b.WriteString("  ")
		b.WriteString(card)
		b.WriteString("\n\n")
	}

	return b.String()
}

// ── Page: Skills ───────────────────────────────────────────────────────

func (m Model) renderSkills(width int) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.SectionTitleStyle.Render("  ⬟ Skills & Technologies"))
	b.WriteString("\n\n")

	// Group skills into categories
	categories := map[string][]string{
		"Languages":      {},
		"Frontend":       {},
		"Backend":        {},
		"Databases":      {},
		"Cloud & DevOps": {},
		"AI/ML":          {},
	}

	categoryOrder := []string{"Languages", "Frontend", "Backend", "Databases", "Cloud & DevOps", "AI/ML"}

	langKeywords := []string{"JavaScript", "TypeScript", "Python", "HTML", "CSS", "GraphQL"}
	frontendKeywords := []string{"Svelte", "SvelteKit", "Next.js", "React"}
	backendKeywords := []string{"FastAPI", "Node.js"}
	dbKeywords := []string{"MongoDB", "PostgreSQL", "Supabase", "Redis", "ClickHouse", "Kafka"}
	cloudKeywords := []string{"Docker", "AWS", "GCP", "Cloudflare", "DigitalOcean", "Terraform", "CI/CD", "Git", "Linux"}
	aiKeywords := []string{"LiteLLM", "LangChain", "Strands", "AI/ML"}

	for _, skill := range m.cv.Skills {
		placed := false
		for _, kw := range langKeywords {
			if strings.Contains(skill, kw) {
				categories["Languages"] = append(categories["Languages"], skill)
				placed = true
				break
			}
		}
		if placed {
			continue
		}
		for _, kw := range frontendKeywords {
			if strings.Contains(skill, kw) {
				categories["Frontend"] = append(categories["Frontend"], skill)
				placed = true
				break
			}
		}
		if placed {
			continue
		}
		for _, kw := range backendKeywords {
			if strings.Contains(skill, kw) {
				categories["Backend"] = append(categories["Backend"], skill)
				placed = true
				break
			}
		}
		if placed {
			continue
		}
		for _, kw := range dbKeywords {
			if strings.Contains(skill, kw) {
				categories["Databases"] = append(categories["Databases"], skill)
				placed = true
				break
			}
		}
		if placed {
			continue
		}
		for _, kw := range cloudKeywords {
			if strings.Contains(skill, kw) {
				categories["Cloud & DevOps"] = append(categories["Cloud & DevOps"], skill)
				placed = true
				break
			}
		}
		if placed {
			continue
		}
		for _, kw := range aiKeywords {
			if strings.Contains(skill, kw) {
				categories["AI/ML"] = append(categories["AI/ML"], skill)
				placed = true
				break
			}
		}
		if !placed {
			categories["Languages"] = append(categories["Languages"], skill)
		}
	}

	catColors := map[string]color.Color{
		"Languages":      styles.Primary,
		"Frontend":       styles.Pink,
		"Backend":        styles.Secondary,
		"Databases":      styles.Accent,
		"Cloud & DevOps": styles.Success,
		"AI/ML":          styles.Blue,
	}

	catIcons := map[string]string{
		"Languages":      "◆",
		"Frontend":       "◇",
		"Backend":        "▣",
		"Databases":      "◈",
		"Cloud & DevOps": "☁",
		"AI/ML":          "◎",
	}

	for _, cat := range categoryOrder {
		skills := categories[cat]
		if len(skills) == 0 {
			continue
		}

		color := catColors[cat]
		icon := catIcons[cat]

		catTitle := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Render(fmt.Sprintf("  %s %s", icon, cat))
		b.WriteString(catTitle)
		b.WriteString("\n")

		// Render skills as tags in a row
		line := "     "
		lineLen := 5
		for i, skill := range skills {
			tag := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F9FAFB")).
				Background(color).
				Padding(0, 1).
				Render(skill)
			tagLen := len(skill) + 2
			sep := "  "
			if i == len(skills)-1 {
				sep = ""
			}

			if lineLen+tagLen+2 > width-4 {
				b.WriteString(line)
				b.WriteString("\n")
				line = "     "
				lineLen = 5
			}
			line += tag + sep
			lineLen += tagLen + 2
		}
		b.WriteString(line)
		b.WriteString("\n\n")
	}

	// Total count
	totalBadge := lipgloss.NewStyle().
		Foreground(styles.TextDim).
		Italic(true).
		Render(fmt.Sprintf("  Total: %d+ technologies and growing", len(m.cv.Skills)))
	b.WriteString(totalBadge)
	b.WriteString("\n")

	return b.String()
}

// ── Page: Achievements ─────────────────────────────────────────────────

func (m Model) renderAchievements(width int) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.SectionTitleStyle.Render("  ★ Achievements & Awards"))
	b.WriteString("\n\n")

	for _, ach := range m.cv.Achievements {
		var cardContent strings.Builder

		// Trophy icon based on placement
		icon := "🏆"
		if strings.Contains(ach.Title, "1st") {
			icon = "🥇"
		} else if strings.Contains(ach.Title, "4th") || strings.Contains(ach.Title, "Finalist") {
			icon = "🏅"
		} else if strings.Contains(ach.Title, "Best") {
			icon = "⭐"
		}

		cardContent.WriteString(fmt.Sprintf("%s  %s\n",
			icon,
			styles.AchievementStyle.Render(ach.Title),
		))
		cardContent.WriteString(fmt.Sprintf("   %s\n",
			styles.DateStyle.Render(ach.Date),
		))

		if ach.Summary != "" {
			cardContent.WriteString(fmt.Sprintf("\n   %s\n",
				styles.TextStyle.Render(ach.Summary),
			))
		}

		if len(ach.Highlights) > 0 {
			cardContent.WriteString("\n")
			for _, hl := range ach.Highlights {
				wrapped := wordWrap(hl, width-14)
				lines := strings.Split(wrapped, "\n")
				cardContent.WriteString(fmt.Sprintf("   %s %s\n",
					styles.BulletStyle.Render("▸"),
					styles.DimTextStyle.Render(lines[0]),
				))
				for _, line := range lines[1:] {
					cardContent.WriteString(fmt.Sprintf("     %s\n",
						styles.DimTextStyle.Render(line),
					))
				}
			}
		}

		if ach.Stack != "" {
			cardContent.WriteString(fmt.Sprintf("\n   %s\n",
				styles.StackStyle.Render("⚙ "+ach.Stack),
			))
		}

		card := styles.HighlightCardStyle.Width(width - 4).Render(cardContent.String())
		b.WriteString("  ")
		b.WriteString(card)
		b.WriteString("\n\n")
	}

	return b.String()
}

// ── Page: Contact ──────────────────────────────────────────────────────

func (m Model) renderContact(width int) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(styles.SectionTitleStyle.Render("  ◉ Get In Touch"))
	b.WriteString("\n\n")

	// Contact card
	contactInfo := fmt.Sprintf(
		"%s  %s\n\n%s  %s\n\n%s  %s\n\n%s  %s",
		styles.BulletStyle.Render("✉"),
		styles.LinkStyle.Render(m.cv.Basics.Email),
		styles.BulletStyle.Render("☎"),
		styles.TextStyle.Render(m.cv.Basics.Phone),
		styles.BulletStyle.Render("◆"),
		styles.LinkStyle.Render(m.cv.Basics.Website),
		styles.BulletStyle.Render("⌂"),
		styles.TextStyle.Render("Indonesia"),
	)

	contactCard := styles.HighlightCardStyle.Width(width - 4).Render(contactInfo)
	b.WriteString("  ")
	b.WriteString(contactCard)
	b.WriteString("\n\n")

	// Social links
	b.WriteString(styles.SectionTitleStyle.Render("  Social Links"))
	b.WriteString("\n\n")

	socialIcons := map[string]string{
		"LinkedIn":  "in",
		"GitHub":    "gh",
		"Blog":     "bg",
		"Instagram": "ig",
		"Medium":    "md",
	}

	for _, social := range m.cv.Socials {
		icon := socialIcons[social.Network]
		if icon == "" {
			icon = "◆"
		}

		badge := lipgloss.NewStyle().
			Foreground(styles.BgDark).
			Background(styles.Secondary).
			Bold(true).
			Padding(0, 1).
			Render(icon)

		line := fmt.Sprintf("  %s  %s  %s",
			badge,
			styles.BoldTextStyle.Render(social.Network),
			styles.LinkStyle.Render(social.URL),
		)
		b.WriteString(line)
		b.WriteString("\n\n")
	}

	// Farewell message
	b.WriteString("\n")
	b.WriteString("  ")
	farewell := lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(styles.Accent).
		Padding(1, 3).
		Width(width - 4).
		Align(lipgloss.Center).
		Render(
			styles.AchievementStyle.Render("Thanks for visiting via SSH!") + "\n" +
				styles.DimTextStyle.Render("Let's build something awesome together.") + "\n\n" +
				styles.TextStyle.Render("— Farhan Aulianda"),
		)
	b.WriteString(farewell)
	b.WriteString("\n")

	return b.String()
}

// ── Helpers ────────────────────────────────────────────────────────────

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) > width {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine += " " + word
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
