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

// LightTheme represents the default Code Garden light palette.
var LightTheme = ThemeTokens{
	Primary:          "#10b981",
	PrimaryHover:     "#059669",
	PrimaryContent:   "#ffffff",
	PrimarySubtle:    "rgba(16, 185, 129, 0.12)",
	Secondary:        "#0d9488",
	SecondaryHover:   "#0f766e",
	SecondaryContent: "#ffffff",
	SecondarySubtle:  "rgba(13, 148, 136, 0.12)",
	Accent:           "#f59e0b",
	AccentHover:      "#d97706",
	AccentContent:    "#ffffff",
	AccentSubtle:     "rgba(245, 158, 11, 0.14)",
	Success:          "#10b981",
	Info:             "#0284c7",
	Warning:          "#f59e0b",
	Error:            "#ef4444",

	BG:           "#fafaf9",
	BG2:          "#f5f5f4",
	BG3:          "#e7e5e4",
	FG:           "#1c1917",
	FGMuted:      "#78716c",
	FGSubtle:     "#a8a29e",
	Border:       "#e7e5e4",
	BorderSubtle: "#f5f5f4",
	BorderFocus:  "#10b981",
	Surface:      "#ffffff",
	SurfaceMuted: "#f5f5f4",

	ShadowSM: "0 1px 2px 0 rgba(0, 0, 0, 0.05)",
	Shadow:   "0 1px 3px 0 rgba(0, 0, 0, 0.08), 0 1px 2px -1px rgba(0, 0, 0, 0.05)",
	ShadowMD: "0 4px 6px -1px rgba(0, 0, 0, 0.08), 0 2px 4px -2px rgba(0, 0, 0, 0.04)",
	ShadowLG: "0 10px 15px -3px rgba(0, 0, 0, 0.08), 0 4px 6px -4px rgba(0, 0, 0, 0.03)",
	ShadowXL: "0 20px 25px -5px rgba(0, 0, 0, 0.08), 0 8px 10px -6px rgba(0, 0, 0, 0.03)",

	RadiusXS: "4px",
	RadiusSM: "6px",
	Radius:   "8px",
	RadiusMD: "10px",
	RadiusLG: "12px",
	RadiusXL: "16px",
}

// DarkTheme represents the Code Garden dark/moss palette.
var DarkTheme = ThemeTokens{
	Primary:          "#34d399",
	PrimaryHover:     "#6ee7b7",
	PrimaryContent:   "#0a0f0d",
	PrimarySubtle:    "rgba(52, 211, 153, 0.16)",
	Secondary:        "#2dd4bf",
	SecondaryHover:   "#5eead4",
	SecondaryContent: "#0a0f0d",
	SecondarySubtle:  "rgba(45, 212, 191, 0.16)",
	Accent:           "#fbbf24",
	AccentHover:      "#fcd34d",
	AccentContent:    "#0a0f0d",
	AccentSubtle:     "rgba(251, 191, 36, 0.16)",
	Success:          "#34d399",
	Info:             "#38bdf8",
	Warning:          "#fbbf24",
	Error:            "#f87171",

	BG:           "#0a0f0d",
	BG2:          "#121a16",
	BG3:          "#1a2620",
	FG:           "#f5f5f4",
	FGMuted:      "#a8a29e",
	FGSubtle:     "#78716c",
	Border:       "#24362e",
	BorderSubtle: "#1a2620",
	BorderFocus:  "#34d399",
	Surface:      "#111814",
	SurfaceMuted: "#18221d",

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
