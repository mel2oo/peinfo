package peinfo

var ResourceTypeName = map[uint32]string{
	1:  "RT_CURSOR",
	2:  "RT_BITMAP",
	3:  "RT_ICON",
	4:  "RT_MENU",
	5:  "RT_DIALOG",
	6:  "RT_STRING",
	7:  "RT_FONTDIR",
	8:  "RT_FONT",
	9:  "RT_ACCELERATOR",
	10: "RT_RCDATA",
	11: "RT_MESSAGETABLE",
	12: "RT_GROUP_CURSOR",
	14: "RT_GROUP_ICON",
	16: "RT_VERSION",
	17: "RT_DLGINCLUDE",
	19: "RT_PLUGPLAY",
	20: "RT_VXD",
	21: "RT_ANICURSOR",
	22: "RT_ANIICON",
	23: "RT_HTML",
	24: "RT_MANIFEST",
}

var LanguageTypeName = map[string]uint32{
	"LANG_NEUTRAL":        0x00,
	"LANG_INVARIANT":      0x7F,
	"LANG_AFRIKAANS":      0x36,
	"LANG_ALBANIAN":       0x1C,
	"LANG_ARABIC":         0x01,
	"LANG_ARMENIAN":       0x2B,
	"LANG_ASSAMESE":       0x4D,
	"LANG_AZERI":          0x2C,
	"LANG_BASQUE":         0x2D,
	"LANG_BELARUSIAN":     0x23,
	"LANG_BENGALI":        0x45,
	"LANG_BULGARIAN":      0x02,
	"LANG_CATALAN":        0x03,
	"LANG_CHINESE":        0x04,
	"LANG_CROATIAN":       0x1A,
	"LANG_CZECH":          0x05,
	"LANG_DANISH":         0x06,
	"LANG_DIVEHI":         0x65,
	"LANG_DUTCH":          0x13,
	"LANG_ENGLISH":        0x09,
	"LANG_ESTONIAN":       0x25,
	"LANG_FAEROESE":       0x38,
	"LANG_FARSI":          0x29,
	"LANG_FINNISH":        0x0B,
	"LANG_FRENCH":         0x0C,
	"LANG_GALICIAN":       0x56,
	"LANG_GEORGIAN":       0x37,
	"LANG_GERMAN":         0x07,
	"LANG_GREEK":          0x08,
	"LANG_GUJARATI":       0x47,
	"LANG_HEBREW":         0x0D,
	"LANG_HINDI":          0x39,
	"LANG_HUNGARIAN":      0x0E,
	"LANG_ICELANDIC":      0x0F,
	"LANG_INDONESIAN":     0x21,
	"LANG_ITALIAN":        0x10,
	"LANG_JAPANESE":       0x11,
	"LANG_KANNADA":        0x4B,
	"LANG_KASHMIRI":       0x60,
	"LANG_KAZAK":          0x3F,
	"LANG_KONKANI":        0x57,
	"LANG_KOREAN":         0x12,
	"LANG_KYRGYZ":         0x40,
	"LANG_LATVIAN":        0x26,
	"LANG_LITHUANIAN":     0x27,
	"LANG_MACEDONIAN":     0x2F,
	"LANG_MALAY":          0x3E,
	"LANG_MALAYALAM":      0x4C,
	"LANG_MANIPURI":       0x58,
	"LANG_MARATHI":        0x4E,
	"LANG_MONGOLIAN":      0x50,
	"LANG_NEPALI":         0x61,
	"LANG_NORWEGIAN":      0x14,
	"LANG_ORIYA":          0x48,
	"LANG_POLISH":         0x15,
	"LANG_PORTUGUESE":     0x16,
	"LANG_PUNJABI":        0x46,
	"LANG_ROMANIAN":       0x18,
	"LANG_RUSSIAN":        0x19,
	"LANG_SANSKRIT":       0x4F,
	"LANG_SERBIAN":        0x1A,
	"LANG_SINDHI":         0x59,
	"LANG_SLOVAK":         0x1B,
	"LANG_SLOVENIAN":      0x24,
	"LANG_SPANISH":        0x0A,
	"LANG_SWAHILI":        0x41,
	"LANG_SWEDISH":        0x1D,
	"LANG_SYRIAC":         0x5A,
	"LANG_TAMIL":          0x49,
	"LANG_TATAR":          0x44,
	"LANG_TELUGU":         0x4A,
	"LANG_THAI":           0x1E,
	"LANG_TURKISH":        0x1F,
	"LANG_UKRAINIAN":      0x22,
	"LANG_URDU":           0x20,
	"LANG_UZBEK":          0x43,
	"LANG_VIETNAMESE":     0x2A,
	"LANG_GAELIC":         0x3C,
	"LANG_MALTESE":        0x3A,
	"LANG_MAORI":          0x28,
	"LANG_RHAETO_ROMANCE": 0x17,
	"LANG_SAAMI":          0x3B,
	"LANG_SORBIAN":        0x2E,
	"LANG_SUTU":           0x30,
	"LANG_TSONGA":         0x31,
	"LANG_TSWANA":         0x32,
	"LANG_VENDA":          0x33,
	"LANG_XHOSA":          0x34,
	"LANG_ZULU":           0x35,
	"LANG_ESPERANTO":      0x8F,
	"LANG_WALON":          0x90,
	"LANG_CORNISH":        0x91,
	"LANG_WELSH":          0x92,
	"LANG_BRETON":         0x93,
}

var SubLanguageTypeName = map[string]uint32{
	"SUBLANG_NEUTRAL":                    0x00,
	"SUBLANG_DEFAULT":                    0x01,
	"SUBLANG_SYS_DEFAULT":                0x02,
	"SUBLANG_ARABIC_SAUDI_ARABIA":        0x01,
	"SUBLANG_ARABIC_IRAQ":                0x02,
	"SUBLANG_ARABIC_EGYPT":               0x03,
	"SUBLANG_ARABIC_LIBYA":               0x04,
	"SUBLANG_ARABIC_ALGERIA":             0x05,
	"SUBLANG_ARABIC_MOROCCO":             0x06,
	"SUBLANG_ARABIC_TUNISIA":             0x07,
	"SUBLANG_ARABIC_OMAN":                0x08,
	"SUBLANG_ARABIC_YEMEN":               0x09,
	"SUBLANG_ARABIC_SYRIA":               0x0A,
	"SUBLANG_ARABIC_JORDAN":              0x0B,
	"SUBLANG_ARABIC_LEBANON":             0x0C,
	"SUBLANG_ARABIC_KUWAIT":              0x0D,
	"SUBLANG_ARABIC_UAE":                 0x0E,
	"SUBLANG_ARABIC_BAHRAIN":             0x0F,
	"SUBLANG_ARABIC_QATAR":               0x10,
	"SUBLANG_AZERI_LATIN":                0x01,
	"SUBLANG_AZERI_CYRILLIC":             0x02,
	"SUBLANG_CHINESE_TRADITIONAL":        0x01,
	"SUBLANG_CHINESE_SIMPLIFIED":         0x02,
	"SUBLANG_CHINESE_HONGKONG":           0x03,
	"SUBLANG_CHINESE_SINGAPORE":          0x04,
	"SUBLANG_CHINESE_MACAU":              0x05,
	"SUBLANG_DUTCH":                      0x01,
	"SUBLANG_DUTCH_BELGIAN":              0x02,
	"SUBLANG_ENGLISH_US":                 0x01,
	"SUBLANG_ENGLISH_UK":                 0x02,
	"SUBLANG_ENGLISH_AUS":                0x03,
	"SUBLANG_ENGLISH_CAN":                0x04,
	"SUBLANG_ENGLISH_NZ":                 0x05,
	"SUBLANG_ENGLISH_EIRE":               0x06,
	"SUBLANG_ENGLISH_SOUTH_AFRICA":       0x07,
	"SUBLANG_ENGLISH_JAMAICA":            0x08,
	"SUBLANG_ENGLISH_CARIBBEAN":          0x09,
	"SUBLANG_ENGLISH_BELIZE":             0x0A,
	"SUBLANG_ENGLISH_TRINIDAD":           0x0B,
	"SUBLANG_ENGLISH_ZIMBABWE":           0x0C,
	"SUBLANG_ENGLISH_PHILIPPINES":        0x0D,
	"SUBLANG_FRENCH":                     0x01,
	"SUBLANG_FRENCH_BELGIAN":             0x02,
	"SUBLANG_FRENCH_CANADIAN":            0x03,
	"SUBLANG_FRENCH_SWISS":               0x04,
	"SUBLANG_FRENCH_LUXEMBOURG":          0x05,
	"SUBLANG_FRENCH_MONACO":              0x06,
	"SUBLANG_GERMAN":                     0x01,
	"SUBLANG_GERMAN_SWISS":               0x02,
	"SUBLANG_GERMAN_AUSTRIAN":            0x03,
	"SUBLANG_GERMAN_LUXEMBOURG":          0x04,
	"SUBLANG_GERMAN_LIECHTENSTEIN":       0x05,
	"SUBLANG_ITALIAN":                    0x01,
	"SUBLANG_ITALIAN_SWISS":              0x02,
	"SUBLANG_KASHMIRI_SASIA":             0x02,
	"SUBLANG_KASHMIRI_INDIA":             0x02,
	"SUBLANG_KOREAN":                     0x01,
	"SUBLANG_LITHUANIAN":                 0x01,
	"SUBLANG_MALAY_MALAYSIA":             0x01,
	"SUBLANG_MALAY_BRUNEI_DARUSSALAM":    0x02,
	"SUBLANG_NEPALI_INDIA":               0x02,
	"SUBLANG_NORWEGIAN_BOKMAL":           0x01,
	"SUBLANG_NORWEGIAN_NYNORSK":          0x02,
	"SUBLANG_PORTUGUESE":                 0x02,
	"SUBLANG_PORTUGUESE_BRAZILIAN":       0x01,
	"SUBLANG_SERBIAN_LATIN":              0x02,
	"SUBLANG_SERBIAN_CYRILLIC":           0x03,
	"SUBLANG_SPANISH":                    0x01,
	"SUBLANG_SPANISH_MEXICAN":            0x02,
	"SUBLANG_SPANISH_MODERN":             0x03,
	"SUBLANG_SPANISH_GUATEMALA":          0x04,
	"SUBLANG_SPANISH_COSTA_RICA":         0x05,
	"SUBLANG_SPANISH_PANAMA":             0x06,
	"SUBLANG_SPANISH_DOMINICAN_REPUBLIC": 0x07,
	"SUBLANG_SPANISH_VENEZUELA":          0x08,
	"SUBLANG_SPANISH_COLOMBIA":           0x09,
	"SUBLANG_SPANISH_PERU":               0x0A,
	"SUBLANG_SPANISH_ARGENTINA":          0x0B,
	"SUBLANG_SPANISH_ECUADOR":            0x0C,
	"SUBLANG_SPANISH_CHILE":              0x0D,
	"SUBLANG_SPANISH_URUGUAY":            0x0E,
	"SUBLANG_SPANISH_PARAGUAY":           0x0F,
	"SUBLANG_SPANISH_BOLIVIA":            0x10,
	"SUBLANG_SPANISH_EL_SALVADOR":        0x11,
	"SUBLANG_SPANISH_HONDURAS":           0x12,
	"SUBLANG_SPANISH_NICARAGUA":          0x13,
	"SUBLANG_SPANISH_PUERTO_RICO":        0x14,
	"SUBLANG_SWEDISH":                    0x01,
	"SUBLANG_SWEDISH_FINLAND":            0x02,
	"SUBLANG_URDU_PAKISTAN":              0x01,
	"SUBLANG_URDU_INDIA":                 0x02,
	"SUBLANG_UZBEK_LATIN":                0x01,
	"SUBLANG_UZBEK_CYRILLIC":             0x02,
	"SUBLANG_DUTCH_SURINAM":              0x03,
	"SUBLANG_ROMANIAN":                   0x01,
	"SUBLANG_ROMANIAN_MOLDAVIA":          0x02,
	"SUBLANG_RUSSIAN":                    0x01,
	"SUBLANG_RUSSIAN_MOLDAVIA":           0x02,
	"SUBLANG_CROATIAN":                   0x01,
	"SUBLANG_LITHUANIAN_CLASSIC":         0x02,
	"SUBLANG_GAELIC":                     0x01,
	"SUBLANG_GAELIC_SCOTTISH":            0x02,
	"SUBLANG_GAELIC_MANX":                0x03,
}

const (
	// https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#section-flags
	IMAGE_SCN_MEM_EXECUTE = 0x20000000
	IMAGE_SCN_MEM_READ    = 0x40000000
	IMAGE_SCN_MEM_WRITE   = 0x80000000
)

var (
	// https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#machine-types
	IMAGE_FILE_MACHINE_UNKNOWN     uint16 = 0x0
	IMAGE_FILE_MACHINE_AM33        uint16 = 0x1d3
	IMAGE_FILE_MACHINE_AMD64       uint16 = 0x8664
	IMAGE_FILE_MACHINE_ARM         uint16 = 0x1c0
	IMAGE_FILE_MACHINE_ARM64       uint16 = 0xaa64
	IMAGE_FILE_MACHINE_ARMNT       uint16 = 0x1c4
	IMAGE_FILE_MACHINE_EBC         uint16 = 0xebc
	IMAGE_FILE_MACHINE_I386        uint16 = 0x14c
	IMAGE_FILE_MACHINE_IA64        uint16 = 0x200
	IMAGE_FILE_MACHINE_LOONGARCH32 uint16 = 0x6232
	IMAGE_FILE_MACHINE_LOONGARCH64 uint16 = 0x6264
	IMAGE_FILE_MACHINE_M32R        uint16 = 0x9041
	IMAGE_FILE_MACHINE_MIPS16      uint16 = 0x266
	IMAGE_FILE_MACHINE_MIPSFPU     uint16 = 0x366
	IMAGE_FILE_MACHINE_MIPSFPU16   uint16 = 0x466
	IMAGE_FILE_MACHINE_POWERPC     uint16 = 0x1f0
	IMAGE_FILE_MACHINE_POWERPCFP   uint16 = 0x1f1
	IMAGE_FILE_MACHINE_R4000       uint16 = 0x166
	IMAGE_FILE_MACHINE_RISCV32     uint16 = 0x5032
	IMAGE_FILE_MACHINE_RISCV64     uint16 = 0x5064
	IMAGE_FILE_MACHINE_RISCV128    uint16 = 0x5128
	IMAGE_FILE_MACHINE_SH3         uint16 = 0x1a2
	IMAGE_FILE_MACHINE_SH3DSP      uint16 = 0x1a3
	IMAGE_FILE_MACHINE_SH4         uint16 = 0x1a6
	IMAGE_FILE_MACHINE_SH5         uint16 = 0x1a8
	IMAGE_FILE_MACHINE_THUMB       uint16 = 0x1c2
	IMAGE_FILE_MACHINE_WCEMIPSV2   uint16 = 0x169

	MachineTypeDesc = map[uint16]string{
		IMAGE_FILE_MACHINE_UNKNOWN:     "The content of this field is assumed to be applicable to any machine type",
		IMAGE_FILE_MACHINE_AM33:        "Matsushita AM33",
		IMAGE_FILE_MACHINE_AMD64:       "x64",
		IMAGE_FILE_MACHINE_ARM:         "ARM little endian",
		IMAGE_FILE_MACHINE_ARM64:       "ARM64 little endian",
		IMAGE_FILE_MACHINE_ARMNT:       "ARM Thumb-2 little endian",
		IMAGE_FILE_MACHINE_EBC:         "EFI byte code",
		IMAGE_FILE_MACHINE_I386:        "Intel 386 or later processors and compatible processors",
		IMAGE_FILE_MACHINE_IA64:        "Intel Itanium processor family",
		IMAGE_FILE_MACHINE_LOONGARCH32: "LoongArch 32-bit processor family",
		IMAGE_FILE_MACHINE_LOONGARCH64: "LoongArch 64-bit processor family",
		IMAGE_FILE_MACHINE_M32R:        "Mitsubishi M32R little endian",
		IMAGE_FILE_MACHINE_MIPS16:      "MIPS16",
		IMAGE_FILE_MACHINE_MIPSFPU:     "MIPS with FPU",
		IMAGE_FILE_MACHINE_MIPSFPU16:   "MIPS16 with FPU",
		IMAGE_FILE_MACHINE_POWERPC:     "Power PC little endian",
		IMAGE_FILE_MACHINE_POWERPCFP:   "Power PC with floating point support",
		IMAGE_FILE_MACHINE_R4000:       "MIPS little endian",
		IMAGE_FILE_MACHINE_RISCV32:     "RISC-V 32-bit address space",
		IMAGE_FILE_MACHINE_RISCV64:     "RISC-V 64-bit address space",
		IMAGE_FILE_MACHINE_RISCV128:    "RISC-V 128-bit address space",
		IMAGE_FILE_MACHINE_SH3:         "Hitachi SH3",
		IMAGE_FILE_MACHINE_SH3DSP:      "Hitachi SH3 DSP",
		IMAGE_FILE_MACHINE_SH4:         "Hitachi SH4",
		IMAGE_FILE_MACHINE_SH5:         "Hitachi SH5",
		IMAGE_FILE_MACHINE_THUMB:       "Thumb",
		IMAGE_FILE_MACHINE_WCEMIPSV2:   "MIPS little-endian WCE v2",
	}
)

var (
	// https://learn.microsoft.com/en-us/windows/win32/debug/pe-format#windows-subsystem
	IMAGE_SUBSYSTEM_UNKNOWN                  uint16 = 0
	IMAGE_SUBSYSTEM_NATIVE                   uint16 = 1
	IMAGE_SUBSYSTEM_WINDOWS_GUI              uint16 = 2
	IMAGE_SUBSYSTEM_WINDOWS_CUI              uint16 = 3
	IMAGE_SUBSYSTEM_OS2_CUI                  uint16 = 5
	IMAGE_SUBSYSTEM_POSIX_CUI                uint16 = 7
	IMAGE_SUBSYSTEM_NATIVE_WINDOWS           uint16 = 8
	IMAGE_SUBSYSTEM_WINDOWS_CE_GUI           uint16 = 9
	IMAGE_SUBSYSTEM_EFI_APPLICATION          uint16 = 10
	IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER  uint16 = 11
	IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER       uint16 = 12
	IMAGE_SUBSYSTEM_EFI_ROM                  uint16 = 13
	IMAGE_SUBSYSTEM_XBOX                     uint16 = 14
	IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION uint16 = 16

	SubsystemTypeDesc = map[uint16]string{
		IMAGE_SUBSYSTEM_UNKNOWN:                  "An unknown subsystem",
		IMAGE_SUBSYSTEM_NATIVE:                   "Device drivers and native Windows processes",
		IMAGE_SUBSYSTEM_WINDOWS_GUI:              "The Windows graphical user interface (GUI) subsystem",
		IMAGE_SUBSYSTEM_WINDOWS_CUI:              "The Windows character subsystem",
		IMAGE_SUBSYSTEM_OS2_CUI:                  "The OS/2 character subsystem",
		IMAGE_SUBSYSTEM_POSIX_CUI:                "The Posix character subsystem",
		IMAGE_SUBSYSTEM_NATIVE_WINDOWS:           "Native Win9x driver",
		IMAGE_SUBSYSTEM_WINDOWS_CE_GUI:           "Windows CE",
		IMAGE_SUBSYSTEM_EFI_APPLICATION:          "An Extensible Firmware Interface (EFI) application",
		IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER:  "An EFI driver with boot services",
		IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER:       "An EFI driver with run-time services",
		IMAGE_SUBSYSTEM_EFI_ROM:                  "An EFI ROM image",
		IMAGE_SUBSYSTEM_XBOX:                     "XBOX",
		IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION: "Windows boot application.",
	}
)
