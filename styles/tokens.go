package styles

// TokenDefinition holds the design token values for the Code Garden aesthetic.
type ThemeTokens struct {
	// Colors
	Primary          string
	PrimaryHover     string
	PrimaryContent   string
	PrimarySubtle    string
	Secondary        string
	SecondaryHover   string
	SecondaryContent string
	SecondarySubtle  string
	Accent           string
	AccentHover      string
	AccentContent    string
	AccentSubtle     string
	Success          string
	Info             string
	Warning          string
	Error            string

	// Neutrals
	BG           string
	BG2          string
	BG3          string
	FG           string
	FGMuted      string
	FGSubtle     string
	Border       string
	BorderSubtle string
	BorderFocus  string
	Surface      string
	SurfaceMuted string

	// Shadows
	ShadowSM string
	Shadow   string
	ShadowMD string
	ShadowLG string
	ShadowXL string

	// Radii
	RadiusXS string
	RadiusSM string
	Radius   string
	RadiusMD string
	RadiusLG string
	RadiusXL string
}

// LightTheme represents the Code Garden light palette with Krewire brand colors.
var LightTheme = ThemeTokens{
	Primary:          "#39D353",
	PrimaryHover:     "#2ecc50",
	PrimaryContent:   "#0B1F3B",
	PrimarySubtle:    "rgba(57, 211, 83, 0.14)",
	Secondary:        "#00D1C1",
	SecondaryHover:   "#00b8aa",
	SecondaryContent: "#0B1F3B",
	SecondarySubtle:  "rgba(0, 209, 193, 0.14)",
	Accent:           "#FF3B2E",
	AccentHover:      "#e0281b",
	AccentContent:    "#ffffff",
	AccentSubtle:     "rgba(255, 59, 46, 0.14)",
	Success:          "#39D353",
	Info:             "#0284c7",
	Warning:          "#f59e0b",
	Error:            "#FF3B2E",

	BG:           "#FAF8F4",
	BG2:          "#F2EDE4",
	BG3:          "#E8E2D5",
	FG:           "#0B1F3B",
	FGMuted:      "#57677d",
	FGSubtle:     "#8fa2bc",
	Border:       "#e2ded4",
	BorderSubtle: "#ede9e0",
	BorderFocus:  "#39D353",
	Surface:      "#ffffff",
	SurfaceMuted: "#F2EDE4",

	ShadowSM: "0 1px 2px 0 rgba(11, 31, 59, 0.05)",
	Shadow:   "0 1px 3px 0 rgba(11, 31, 59, 0.08), 0 1px 2px -1px rgba(11, 31, 59, 0.05)",
	ShadowMD: "0 4px 6px -1px rgba(11, 31, 59, 0.08), 0 2px 4px -2px rgba(11, 31, 59, 0.04)",
	ShadowLG: "0 10px 15px -3px rgba(11, 31, 59, 0.08), 0 4px 6px -4px rgba(11, 31, 59, 0.03)",
	ShadowXL: "0 20px 25px -5px rgba(11, 31, 59, 0.08), 0 8px 10px -6px rgba(11, 31, 59, 0.03)",

	RadiusXS: "4px",
	RadiusSM: "6px",
	Radius:   "8px",
	RadiusMD: "10px",
	RadiusLG: "12px",
	RadiusXL: "16px",
}

// DarkTheme represents the Code Garden dark palette with Krewire brand colors.
var DarkTheme = ThemeTokens{
	Primary:          "#39D353",
	PrimaryHover:     "#5eed75",
	PrimaryContent:   "#0B1F3B",
	PrimarySubtle:    "rgba(57, 211, 83, 0.18)",
	Secondary:        "#00D1C1",
	SecondaryHover:   "#33ddd0",
	SecondaryContent: "#0B1F3B",
	SecondarySubtle:  "rgba(0, 209, 193, 0.18)",
	Accent:           "#FF5C52",
	AccentHover:      "#ff7870",
	AccentContent:    "#ffffff",
	AccentSubtle:     "rgba(255, 92, 82, 0.18)",
	Success:          "#39D353",
	Info:             "#38bdf8",
	Warning:          "#fbbf24",
	Error:            "#f87171",

	BG:           "#0B1F3B",
	BG2:          "#071527",
	BG3:          "#0f284a",
	FG:           "#FAF8F4",
	FGMuted:      "#9cb0cb",
	FGSubtle:     "#657e9f",
	Border:       "#1c3a63",
	BorderSubtle: "#132c4d",
	BorderFocus:  "#39D353",
	Surface:      "#0f284a",
	SurfaceMuted: "#071527",

	ShadowSM: "0 1px 2px 0 rgba(0, 0, 0, 0.4)",
	Shadow:   "0 4px 6px -1px rgba(0, 0, 0, 0.5)",
	ShadowMD: "0 6px 10px -2px rgba(0, 0, 0, 0.55)",
	ShadowLG: "0 10px 15px -3px rgba(0, 0, 0, 0.6)",
	ShadowXL: "0 20px 25px -5px rgba(0, 0, 0, 0.7)",

	RadiusXS: "4px",
	RadiusSM: "6px",
	Radius:   "8px",
	RadiusMD: "10px",
	RadiusLG: "12px",
	RadiusXL: "16px",
}
