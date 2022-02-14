package repository

import (
	"context"
	"github.com/STRRL/poc-reference-to-grafana/pkg/entity"
)

type GrafanaPanelRepository interface {
	ListAll(ctx context.Context) ([]entity.GrafanaPanel, error)
	GetByID(ctx context.Context, id uint) (*entity.GrafanaPanel, error)
}

type GrafanaPanelVariableRepository interface {
	ListByGrafanaPanelID(ctx context.Context, panelID uint) ([]entity.GrafanaPanelVariable, error)
	GetByID(ctx context.Context, id uint) (*entity.GrafanaPanelVariable, error)
}

type GrafanaPanelBindingRepository interface {
	ListAll(ctx context.Context) ([]entity.GrafanaPanelBinding, error)
	ListByGrafanaPanelID(ctx context.Context, panelID uint) ([]entity.GrafanaPanelBinding, error)
	ListByChaosExperimentID(ctx context.Context, chaosExperimentID uint) ([]entity.GrafanaPanelBinding, error)
	ListByScheduleID(ctx context.Context, scheduleID uint) ([]entity.GrafanaPanelBinding, error)
	ListByWorkflowID(ctx context.Context, workflowID uint) ([]entity.GrafanaPanelBinding, error)
}
