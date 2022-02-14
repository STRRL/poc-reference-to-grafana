package entity

import "gorm.io/gorm"

type GrafanaPanel struct {
	gorm.Model
	Name        string
	RawURL      string
	Description string
}
type GrafanaPanelVariable struct {
	gorm.Model
	GrafanaPanelID int
	VariableName   string
	VariableType   string
}
type GrafanaPanelBinding struct {
	gorm.Model
	GrafanaPanelID    int
	BindingTargetType string
	ChaosExperimentID int
	ScheduleID        int
	WorkflowID        int
}

type GrafanaPanelBindingVariableWithValue struct {
	gorm.Model
	BindingID     int
	VariableID    int
	VariableValue string
}
