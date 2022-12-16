package peinfo

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/lazybeaver/entropy"
	"github.com/mel2oo/peinfo/resource"
	"github.com/saferwall/pe"
)

type PEInfo struct {
	*pe.File   `json:"-"`
	imageBase  uint   `json:"-"`
	entryPoint uint32 `json:"-"`

	Header struct {
		Machine             string `json:"machine,omitempty"`
		Subsystem           string `json:"subsystem,omitempty"`
		TimeDateStamp       string `json:"time_date_stamp,omitempty"`
		EntryPoint          string `json:"entry_point,omitempty"`
		EntryPointInSection string `json:"entry_point_in_section,omitempty"`
		ImageBase           string `json:"image_base,omitempty"`
		NumberOfSections    int    `json:"number_of_sections,omitempty"`
		LinkerVersion       string `json:"linker_version,omitempty"`
	} `json:"header,omitempty"`
	Signed    SignedInfo                `json:"signed,omitempty"`
	Version   resource.VersionResources `json:"version,omitempty"`
	Debugs    []Debug                   `json:"debugs,omitempty"`
	Imports   []Import                  `json:"imports,omitempty"`
	Exports   []Function                `json:"exports,omitempty"`
	Sections  []Section                 `json:"sections,omitempty"`
	Resources []Resource                `json:"resources,omitempty"`
}

func New(path string) (*PEInfo, error) {
	pedata, err := pe.New(path, &pe.Options{})
	if err != nil {
		return nil, err
	}

	if err := pedata.Parse(); err != nil {
		return nil, err
	}

	peinfo := PEInfo{File: pedata}
	peinfo.ParseHeader()
	peinfo.ParseSigned()
	peinfo.ParseDebugs()
	peinfo.ParseImport()
	peinfo.ParseExport()
	peinfo.ParseSection()
	peinfo.ParseResource()

	rc := resource.NewReader()
	rcinfo, err := rc.Read(path)
	if err == nil {
		peinfo.Version = rcinfo
	}

	return &peinfo, nil
}

func (p *PEInfo) ParseHeader() {
	switch h := p.File.NtHeader.OptionalHeader.(type) {
	case pe.ImageOptionalHeader32:
		p.imageBase = uint(h.ImageBase)
		p.entryPoint = h.AddressOfEntryPoint

		p.Header.Subsystem = SubsystemTypeDesc[h.Subsystem]
		p.Header.EntryPoint = fmt.Sprintf("0x%08x", h.AddressOfEntryPoint)
		p.Header.LinkerVersion = fmt.Sprintf("%d.%d", h.MajorLinkerVersion, h.MinorLinkerVersion)
	case pe.ImageOptionalHeader64:
		p.imageBase = uint(h.ImageBase)
		p.entryPoint = h.AddressOfEntryPoint

		p.Header.Subsystem = SubsystemTypeDesc[h.Subsystem]
		p.Header.EntryPoint = fmt.Sprintf("0x%08x", h.AddressOfEntryPoint)
		p.Header.LinkerVersion = fmt.Sprintf("%d.%d", h.MajorLinkerVersion, h.MinorLinkerVersion)
	}

	p.Header.Machine = MachineTypeDesc[p.File.NtHeader.FileHeader.Machine]
	tm := time.Unix(int64(p.File.NtHeader.FileHeader.TimeDateStamp), 0)
	p.Header.TimeDateStamp = tm.Format("2006-01-02 15:04:05")
	p.Header.ImageBase = fmt.Sprintf("0x%08x", p.imageBase)
	p.Header.NumberOfSections = int(p.File.NtHeader.FileHeader.NumberOfSections)
}

type SignedInfo struct {
	Subject string `json:"subject,omitempty"`
	Issuer  string `json:"issuer,omitempty"`
}

func (p *PEInfo) ParseSigned() {
	if p.File.Certificates.Info.Subject != "" {
		p.Signed.Subject = p.File.Certificates.Info.Subject
		p.Signed.Issuer = p.File.Certificates.Info.Issuer
	}
}

type Debug struct {
	PDB  string `json:"pdb,omitempty"`
	GUID string `json:"guid,omitempty"`
}

func (p *PEInfo) ParseDebugs() {
	for _, d := range p.File.Debugs {
		info, ok := d.Info.(pe.CvInfoPDB70)
		if ok {
			p.Debugs = append(p.Debugs, Debug{
				PDB:  info.PDBFileName,
				GUID: parseGUID(info.Signature),
			})
		}
	}
}

func parseGUID(guid pe.GUID) string {
	return fmt.Sprintf("%x-%x-%x-%s", guid.Data1, guid.Data2, guid.Data3, hex.EncodeToString(guid.Data4[:]))
}

type Import struct {
	Name      string     `json:"name,omitempty"`
	Desc      string     `json:"desc,omitempty"`
	Count     int        `json:"count,omitempty"`
	Functions []Function `json:"functions,omitempty"`
}

type Function struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Index   uint32 `json:"index,omitempty"`
}

func (p *PEInfo) ParseImport() {
	for _, m := range p.File.Imports {
		imp := Import{
			Name:      m.Name,
			Count:     len(m.Functions),
			Functions: make([]Function, 0),
		}

		for _, fuc := range m.Functions {
			imp.Functions = append(imp.Functions, Function{
				Name:    fuc.Name,
				Address: fmt.Sprintf("0x%08x", p.imageBase+uint(fuc.ThunkRVA)),
			})
		}

		p.Imports = append(p.Imports, imp)
	}
}

func (p *PEInfo) ParseExport() {
	for _, fuc := range p.File.Export.Functions {
		p.Exports = append(p.Exports, Function{
			Name:    fuc.Name,
			Address: fmt.Sprintf("0x%08x", p.imageBase+uint(fuc.FunctionRVA)),
			Index:   fuc.Ordinal,
		})
	}
}

type Section struct {
	Name    string  `json:"name,omitempty"`
	VA      string  `json:"va,omitempty"`
	VS      string  `json:"vs,omitempty"`
	PA      string  `json:"pa,omitempty"`
	PS      string  `json:"ps,omitempty"`
	Perm    string  `json:"perm,omitempty"`
	Entropy float64 `json:"entropy,omitempty"`
	Hash    string  `json:"hash,omitempty"`
}

func (p *PEInfo) ParseSection() {
	for _, s := range p.File.Sections {
		data, err := p.File.ReadBytesAtOffset(s.Header.PointerToRawData, s.Header.SizeOfRawData)
		if err != nil {
			return
		}

		h := md5.New()
		h.Write(data)

		if p.entryPoint >= s.Header.VirtualAddress && p.entryPoint <= s.Header.VirtualAddress+s.Header.VirtualSize {
			p.Header.EntryPointInSection = s.NameString()
		}

		p.Sections = append(p.Sections, Section{
			Name:    s.NameString(),
			VA:      fmt.Sprintf("0x%08x", s.Header.VirtualAddress),
			VS:      fmt.Sprintf("0x%08x", s.Header.VirtualSize),
			PA:      fmt.Sprintf("0x%08x", s.Header.PointerToRawData),
			PS:      fmt.Sprintf("0x%08x", s.Header.SizeOfRawData),
			Perm:    parseSectionFlags(s.Header.Characteristics),
			Entropy: shannon(data),
			Hash:    hex.EncodeToString(h.Sum(nil)),
		})
	}
}

func parseSectionFlags(flags uint32) string {
	var s strings.Builder

	for _, v := range []uint{IMAGE_SCN_MEM_READ, IMAGE_SCN_MEM_WRITE, IMAGE_SCN_MEM_EXECUTE} {
		if uint(flags)&v != v {
			s.WriteRune('-')
			continue
		}
		switch v {
		case IMAGE_SCN_MEM_READ:
			s.WriteRune('R')
		case IMAGE_SCN_MEM_WRITE:
			s.WriteRune('W')
		case IMAGE_SCN_MEM_EXECUTE:
			s.WriteRune('E')
		}
	}
	return s.String()
}

func shannon(b []byte) float64 {
	estimator := entropy.NewShannonEstimator()
	if _, err := estimator.Write(b); err != nil {
		return 0.0
	}
	return estimator.Value()
}

type Resource struct {
	Name        string `json:"name,omitempty"`
	Filetype    string `json:"filetype,omitempty"`
	Size        string `json:"size,omitempty"`
	Offset      string `json:"offset,omitempty"`
	Language    string `json:"language,omitempty"`
	SubLanguage string `json:"sub_language,omitempty"`
	data        []byte `json:"-"`
}

func (p *PEInfo) ParseResource() {
	for _, resource_type := range p.File.Resources.Entries {
		var res Resource
		var ok bool

		if len(resource_type.Name) == 0 {
			res.Name, ok = ResourceTypeName[resource_type.ID]
			if !ok {
				res.Name = "None"
			}
		} else {
			res.Name = resource_type.Name
		}

		for _, resource_id := range resource_type.Directory.Entries {
			for _, resource_lang := range resource_id.Directory.Entries {
				rva := p.File.GetOffsetFromRva(resource_lang.Data.Struct.OffsetToData)
				data, err := p.File.ReadBytesAtOffset(rva, resource_lang.Data.Struct.Size)
				if err != nil {
					continue
				}

				res.data = data
				res.Filetype = mimetype.Detect(data).String()
				res.Size = fmt.Sprintf("0x%08x", resource_lang.Data.Struct.Size)
				res.Offset = fmt.Sprintf("0x%08x", resource_lang.Data.Struct.OffsetToData)

			loop1:
				for str, i := range LanguageTypeName {
					if i == resource_lang.Data.Lang {
						res.Language = str
						break loop1
					}
				}

			loop2:
				for str, i := range SubLanguageTypeName {
					if i == resource_lang.Data.Sublang {
						if strings.Contains(str, res.Language) {
							res.SubLanguage = str
							break loop2
						}
					}
				}

				p.Resources = append(p.Resources, res)
			}
		}
	}
}
