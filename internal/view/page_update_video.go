package view

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/charmbracelet/lipgloss"
	"github.com/northmule/pinetree/internal/controller"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg struct {
	all     int
	current int
}

func tickCmd(all int, current int) tea.Cmd {
	return tea.Tick(time.Millisecond*1, func(_ time.Time) tea.Msg {
		return tickMsg{
			all:     all,
			current: current,
		}
	})
}

// Ввод/редактирование текстовых данных
type pageUpdateVideo struct {
	Choice          int
	Chosen          bool
	mainPage        *pageIndex
	controller      *controller.Video
	responseMessage string

	videoDescription textarea.Model

	allVideos     int
	currentVideos int
	allPrepare    int
	videoList     *controller.ResponseGetAll
}

func newPageUpdateVideo(mainPage *pageIndex, controller *controller.Video) *pageUpdateVideo {

	text := textarea.New()
	text.Placeholder = "Новое описание для видео"
	text.CharLimit = 50000
	text.MaxHeight = 100
	text.MaxWidth = 200
	text.SetHeight(15)
	text.SetWidth(50)

	m := &pageUpdateVideo{}
	m.mainPage = mainPage
	m.controller = controller
	m.videoDescription = text

	return m
}

func (m *pageUpdateVideo) Init() tea.Cmd {
	return textinput.Blink
}

func (m *pageUpdateVideo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case clearErrorMsg:
		m.responseMessage = ""
		return m, nil
	case tea.WindowSizeMsg:
		m.videoDescription.SetWidth(msg.Width - 25)
		m.videoDescription.SetHeight(msg.Height - 25)

		return m, nil

	case tickMsg:
		m.allVideos = msg.all
		m.currentVideos = msg.current

		if len(m.videoList.Response.Items) == m.currentVideos {
			// Данные отправлены
			m.responseMessage = "Обновления завершены"
			m.mainPage.logger.Info("Все видео обработанны")
			return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
		}

		video := m.videoList.Response.Items[m.currentVideos]

		requestUpdateDescription := controller.RequestUpdateDescription{}
		requestUpdateDescription.Desc = m.videoDescription.Value()
		requestUpdateDescription.OwnerId = m.mainPage.cfg.Value().VK.GroupID
		requestUpdateDescription.AccessToken = m.mainPage.cfg.Value().VK.AccessToken
		requestUpdateDescription.Version = m.mainPage.cfg.Value().VK.ApiVersion

		m.mainPage.logger.Info(fmt.Sprintf("Подготовка обновления для видео: %s\n", video.Title))
		requestUpdateDescription.VideoId = video.ID
		updateResponse, err := m.controller.UpdateDescription(requestUpdateDescription)
		if err != nil {
			m.responseMessage = err.Error()
			m.mainPage.logger.Warn("Ошибка обновления видео", err.Error())
			return m, nil
		}
		if updateResponse.Response.Success == 1 {
			m.mainPage.logger.Info("Видео обновлено:\n")
		} else {
			m.mainPage.logger.Warn("Видео не обновлено", updateResponse)
		}
		m.currentVideos++
		return m, tea.Sequence(tickCmd(msg.all, m.currentVideos))
	}

	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "tab" {
			m.Choice++
			if m.Choice > 2 {
				m.Choice = 0
			}
		}
		if k == "enter" {
			if m.Choice == 1 {

				var newDescription string

				newDescription = m.videoDescription.Value()
				m.mainPage.logger.Info("Текст для отправки", newDescription)
				if newDescription == "" {
					m.responseMessage = "Необходимо заполнить текст"
					m.mainPage.logger.Warn("Не выбран текст для отправки")
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}

				requestGetAll := controller.RequestGetAll{
					AccessToken: m.mainPage.cfg.Value().VK.AccessToken,
					OwnerId:     m.mainPage.cfg.Value().VK.GroupID,
					AlbumId:     m.mainPage.cfg.Value().VK.AlbumID,
					Offset:      0,
					Version:     m.mainPage.cfg.Value().VK.ApiVersion,
					Count:       100,
				}

				videoList, err := m.controller.GetAll(requestGetAll)
				if err != nil {
					m.responseMessage = err.Error()
					m.mainPage.logger.Warn("Не удалось выполнить запрос для получения списка видео", err.Error())
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				if videoList.Error != nil {
					m.responseMessage = "Ошибка получения списка видео"
					m.mainPage.logger.Warn("Ошибка получения списка видео", videoList.Error)
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				if videoList.Response == nil {
					m.responseMessage = "Не известная ошибка, ответ пустой"
					m.mainPage.logger.Warn("Пустой response", videoList)
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				if len(videoList.Response.Items) == 0 {
					m.responseMessage = "Пустой список видео"
					m.mainPage.logger.Warn("Пустой список видео", videoList)
					return m, tea.Batch(cmd, clearErrorAfter(3*time.Second))
				}
				m.mainPage.logger.Info(fmt.Sprintf("Нашёл видео в количестве %d\n", videoList.Response.Count))

				m.videoList = videoList
				m.allPrepare = len(videoList.Response.Items)
				return m, tea.Sequence(tickCmd(videoList.Response.Count, 0))
			}

			if m.Choice == 2 {
				return m.mainPage, nil
			}
		}
	}

	if m.Choice == 0 {
		m.videoDescription, cmd = m.videoDescription.Update(msg)
		m.videoDescription.Focus()
		return m, cmd
	}

	m.videoDescription, cmd = m.videoDescription.Update(msg)

	return m, nil
}

func (m *pageUpdateVideo) View() string {

	c := m.Choice

	title := renderTitle("Обновление описания для видео ВК")

	tpl := "%s\n\n"
	tpl += subtleStyle.Render("tab: для переключения меню") + dotStyle +
		subtleStyle.Render("ctrl+v: для вставки текста") + dotStyle +
		subtleStyle.Render("enter: выбрать пункт меню") + dotStyle +
		responseTextStyle.Render("\n"+m.responseMessage) + dotStyle +
		responseTextStyle.Render(fmt.Sprintf("\nВсего видео: %d. \nПодготовлено для обновления: %d. Обновлено: %d", m.allVideos, m.allPrepare, m.currentVideos)) + dotStyle

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		renderCheckbox(m.videoDescription.View(), c == 0),
		renderCheckbox("Обновить", c == 1),
		renderCheckbox("Вернуться", c == 2),
	)

	s := fmt.Sprintf(tpl, choices)
	return mainStyle.Render(title + "\n" + s + "\n\n")

}
