package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2)
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Background(lipgloss.Color("#25A065")).Padding(0, 1)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type item struct {
	title       string
	description string
	id          string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	selected map[string]struct{}
	exit     bool
}

func newModel(plugins []ToolPlugin, title string) model {
	items := make([]list.Item, len(plugins))
	for i, p := range plugins {
		items[i] = item{title: p.Name, description: p.Description, id: p.ID}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.Styles.Title = titleStyle

	// Use a custom delegate to render the list items with checkboxes
	delegate := &customDelegate{selected: make(map[string]struct{})}
	l.SetDelegate(delegate)

	return model{
		list:     l,
		selected: delegate.selected,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.exit = true
			return m, tea.Quit

		case "enter":
			return m, tea.Quit

		case " ":
			if i, ok := m.list.SelectedItem().(item); ok {
				if _, ok := m.selected[i.id]; ok {
					delete(m.selected, i.id)
				} else {
					m.selected[i.id] = struct{}{}
				}
			}
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.exit {
		return ""
	}
	return appStyle.Render(m.list.View())
}

type customDelegate struct {
	selected map[string]struct{}
}

func (d *customDelegate) Height() int                               { return 2 }
func (d *customDelegate) Spacing() int                              { return 1 }
func (d *customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d *customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	// Determine if the item is selected
	checked := " "
	if _, ok := d.selected[i.id]; ok {
		checked = "x"
	}

	// Render the item
	title := fmt.Sprintf("[%s] %s", checked, i.title)
	desc := i.description

	if index == m.Index() {
		// Selected item style
		fmt.Fprintf(w, selectedItemStyle.Render("> "+title+"\n  "+desc))
	} else {
		// Normal item style
		fmt.Fprintf(w, itemStyle.Render("  "+title+"\n  "+desc))
	}
}

func runTUI(plugins []ToolPlugin, title string) ([]ToolPlugin, error) {
	m := newModel(plugins, title)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(model); ok {
		if m.exit {
			return nil, nil // User quit
		}

		var selectedPlugins []ToolPlugin
		for _, p := range plugins {
			if _, ok := m.selected[p.ID]; ok {
				selectedPlugins = append(selectedPlugins, p)
			}
		}
		return selectedPlugins, nil
	}
	return nil, fmt.Errorf("could not cast final model")
}
