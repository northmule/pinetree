package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/northmule/pinetree/internal/config"
	"github.com/northmule/pinetree/internal/controller"
	"github.com/northmule/pinetree/internal/service"
	"golang.org/x/net/context"
)

// View вьюха клиента
type View struct {
	log service.LogPusher
	cfg *config.Config
}

// NewView конструктор
func NewView(log service.LogPusher, cfg *config.Config) *View {
	instance := &View{
		log: log,
		cfg: cfg,
	}

	return instance
}

// InitMain подготовка консольных форм
func (v *View) InitMain(ctx context.Context) error {
	var err error

	client := service.NewClientWithLimitForOneSecond(20)
	videoUpdater := controller.NewVideo(client, v.log)
	pi := newPageIndex(videoUpdater, v.cfg, v.log)
	p := tea.NewProgram(pi, tea.WithContext(ctx))
	pi.pg = p
	if _, err = p.Run(); err != nil {
		return err
	}

	return nil
}
