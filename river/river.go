package river

type RiverOutput struct {
	NextChangeID string `json:"next_change_id"`
	Stashes      []Stash
}

type Stash struct {
	Id string

	AccountName       *string
	LastCharacterName *string
	League            *string
	Public            bool
	Stash             *string
	StashType         string

	Items []Item
}

type Item struct {
	AbyssJewel           bool
	AdditionalProperties []Property
	ArtFilename          string
	CisRaceReward        bool
	Colour               *string
	Corrupted            bool
	CosmeticMods         []string
	CraftedMods          []string
	Delve                bool
	DescrText            string
	Duplicated           bool
	Elder                bool
	EnchantMods          []string
	ExplicitMods         []string
	Extended             struct {
		BaseType      string
		Category      string
		Prefixes      int
		Subcategories []string
		Suffixes      int
	}
	FlavourText   []string
	Fractured     bool
	FracturedMods []string
	FrameType     int
	H             int
	Hybrid        struct {
		BaseTypeName string
		ExplicitMods []string
		IsVaalGem    bool
		Properties   []Property
		SecDescrText string
	}
	Icon          string
	Identified    bool
	Id            string
	Ilvl          int
	ImplicitMods  []string
	IncubatedItem struct {
		Level    int
		Name     string
		Progress int
		Total    int
	}
	Influences struct {
		Crusader bool
		Elder    bool
		Hunter   bool
		Redeemer bool
		Shaper   bool
		Warlord  bool
	}
	InventoryId           string
	IsRelic               bool
	ItemLevel             int
	League                string
	MaxStackSize          int
	Name                  string
	NextLevelRequirements []Property
	Note                  string
	Properties            []Property
	ProphecyText          string
	Requirements          []Property
	SeaRaceReward         string
	SecDescrText          string
	Shaper                bool
	Socket                int
	SocketedItems         []Item
	Sockets               []struct {
		Attr    string
		Group   int
		SColour string
	}
	StackSize    int
	Support      bool
	Synthesised  bool
	TalismanTier int
	ThRaceReward bool
	TypeLine     string
	UtilityMods  []string
	Veiled       bool
	VeiledMods   []string
	Verified     bool
	W            int
	X            int
	Y            int
}

type Property struct {
	Name        string
	Values      [][]interface{} // array of [string, int] pairs
	DisplayMode int
	Type        int
	Progress    float64 // used for experience on skill gems
}
