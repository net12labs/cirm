package mission

type MissionUnit struct {
	Id     string
	Name   string
	Status string
}

type Mission struct {
	*MissionUnit
}

type MissionGroup struct {
	*MissionUnit
}

type MissionSubGroup struct {
	*MissionUnit
	ParentId string
}

type SubMission struct {
	*MissionUnit
	ParentId string
}

func NewMission() *Mission {
	return &Mission{
		MissionUnit: &MissionUnit{},
	}
}

func NewMissionGroup() *MissionGroup {
	return &MissionGroup{
		MissionUnit: &MissionUnit{},
	}
}

func NewSubMission() *SubMission {
	return &SubMission{
		MissionUnit: &MissionUnit{},
	}
}

func NewMissionSubGroup() *MissionSubGroup {
	return &MissionSubGroup{
		MissionUnit: &MissionUnit{},
	}
}
