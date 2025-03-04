package types

type ImageBatchInfo struct {
	Batchrsp Batchrsp `json:"batchrsp"`
}

type Batchrsp struct {
	Ver   string `json:"ver"`
	Items []Item `json:"items"`
}

type Item struct {
	Item string `json:"item"`
}

type ImageInfo struct {
	AD AD `json:"ad"`
}

type AD struct {
	EntityID       string         `json:"entityId"`
	LandscapeImage LandscapeImage `json:"landscapeImage"`
	PortraitImage  PortraitImage  `json:"portraitImage"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Copyright      string         `json:"copyright"`
}

type LandscapeImage struct {
	Asset string `json:"asset"`
}

type PortraitImage struct {
	Asset string `json:"asset"`
}
