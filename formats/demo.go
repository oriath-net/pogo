package formats

// I haven't finished converting formats from PyPoE, so this is a demo of what works

type rowID uint64
type shortID uint32

type ActiveSkills struct {
	Id                                string
	DisplayedName                     string
	Description                       string
	Index3                            string
	Icon_DDSFile                      string
	ActiveSkillTargetTypes            []shortID
	ActiveSkillTypes                  []shortID
	WeaponRestriction_ItemClassesKeys []rowID
	WebsiteDescription                string
	WebsiteImage                      string
	Flag0                             bool
	Unknown0                          string
	Flag1                             bool
	SkillTotemId                      int32
	IsManuallyCasted                  bool
	Input_StatKeys                    []rowID
	Output_StatKeys                   []rowID
	MinionActiveSkillTypes            []shortID
	Flag2                             bool
	Flag3                             bool
	Keys0                             []rowID
	Unknown1                          int32
	Key0                              rowID
	Flag4                             bool
}
