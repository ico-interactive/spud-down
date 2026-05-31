package types

type Profile struct {
	AccountID                 int64   `json:"accountId"`
	Name                      string  `json:"name"`
	AvatarURL                 string  `json:"avatarUrl"`
	PerformanceRankMessage    *string `json:"performanceRankMessage"`
	LastUpdated               string  `json:"lastUpdated"`
	CalibrationMatches        int     `json:"calibrationMatches"`
	TwitchUsername            *string `json:"twitchUsername"`
	YoutubeChannelURL         *string `json:"youtubeChannelUrl"`
	Region                    string  `json:"region"`
	CalibrationMatchID        int64   `json:"calibrationMatchId"`
	CalibrationResetMatchID   *int64  `json:"calibrationResetMatchId"`
	FontID                    string  `json:"fontId"`
	GlowStyleID               string  `json:"glowStyleId"`
	GlowColorID               string  `json:"glowColorId"`
	OutlineID                 string  `json:"outlineId"`
	AnimationID               string  `json:"animationId"`
	AvatarShapeID             string  `json:"avatarShapeId"`
	AvatarEffectID            string  `json:"avatarEffectId"`
	AvatarEffectColorID       string  `json:"avatarEffectColorId"`
	TitleID                   string  `json:"titleId"`
	UpdatedWithinLast1Minutes bool    `json:"updatedWithinLast1Minutes"`
	LocalisedLastUpdated      string  `json:"localisedLastUpdated"`
	PPScore                   int     `json:"ppScore"`
	SelectedWidgets           *string `json:"selectedWidgets"`
	EstimatedRankNumber       int     `json:"estimatedRankNumber"`
}
