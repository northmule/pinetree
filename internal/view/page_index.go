package view

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/pinetree/internal/config"
	"github.com/northmule/pinetree/internal/controller"
	"github.com/northmule/pinetree/internal/service"
)

// Экран выбора действий (доступные действия по добавлению данных)
type pageIndex struct {
	Choice       int
	Chosen       bool
	Quitting     bool
	cfg          *config.Config
	videoUpdater *controller.Video
	logger       service.LogPusher

	pg *tea.Program
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

type clearFieldMsg struct {
}

func clearFieldAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearFieldMsg{}
	})
}

func newPageIndex(videoUpdater *controller.Video, cfg *config.Config, logger service.LogPusher) *pageIndex {
	return &pageIndex{
		videoUpdater: videoUpdater,
		cfg:          cfg,
		logger:       logger,
	}
}

// Init инициализация модели
func (m *pageIndex) Init() tea.Cmd {
	return nil
}

// Update изменение модели
func (m *pageIndex) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "tab" {
			m.Choice++
			if m.Choice > 1 {
				m.Choice = 0
			}
		}
		if k == "enter" {
			if m.Choice == 0 {
				m.logger.Info("Переход к пункту 0")
				return newPageUpdateVideo(m, m.videoUpdater), nil
			}
			// выход
			if m.Choice == 1 {
				m.Quitting = true
				return m, tea.Quit
			}

		}
	}

	return m, nil
}

// View вид модели( в том числе при старте)
func (m *pageIndex) View() string {
	c := m.Choice

	title := renderTitle("Доступные действия")
	tpl := "%s\n\n"
	tpl += subtleStyle.Render("tab: для переключения") + dotStyle +
		subtleStyle.Render("enter: выбрать")

	choices := fmt.Sprintf(
		"%s\n%s\n\n",
		renderCheckbox("Изменить описание для всех видео в ВК", c == 0),
		renderCheckbox("Выйти", c == 1),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")
}
