package peinfo

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/mel2oo/peinfo/resource"
	"github.com/saferwall/pe"
	"github.com/vimeo/go-magic/magic"
)

type PEInfo struct {
	*pe.File
	Version   resource.VersionResources
	Imports   []Import
	Resources []Resource
	Sections  []Section
	Exports   []Function
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

type Import struct {
	Name      string
	Desc      string
	Count     int
	Functions []Function
}

type Function struct {
	Name    string `json:"name"`
	Address string `json:"address"`
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
				Name: fuc.Name,
				// todo
				Address: fmt.Sprintf("0x%08x", fuc.ThunkRVA),
			})
		}

		p.Imports = append(p.Imports, imp)
	}
}

func (p *PEInfo) ParseExport() {
	for _, fuc := range p.File.Export.Functions {
		p.Exports = append(p.Exports, Function{
			Name: fuc.Name,
			// todo
			Address: fmt.Sprintf("0x%08x", fuc.FunctionRVA),
			Index:   fuc.Ordinal,
		})
	}
}

type Section struct {
	Name string
	VA   string
	VS   string
	PA   string
	PS   string
	Hash string
}

func (p *PEInfo) ParseSection() {
	for _, s := range p.File.Sections {
		data, err := p.File.ReadBytesAtOffset(s.Header.PointerToRawData, s.Header.SizeOfRawData)
		if err != nil {
			return
		}

		h := md5.New()
		h.Write(data)

		p.Sections = append(p.Sections, Section{
			Name: s.NameString(),
			VA:   fmt.Sprintf("0x%08x", s.Header.VirtualAddress),
			VS:   fmt.Sprintf("0x%08x", s.Header.VirtualSize),
			PA:   fmt.Sprintf("0x%08x", s.Header.PointerToRawData),
			PS:   fmt.Sprintf("0x%08x", s.Header.SizeOfRawData),
			Hash: hex.EncodeToString(h.Sum(nil)),
		})
	}
}

type Resource struct {
	Name        string
	Filetype    string
	Size        string
	Offset      string
	Language    string
	SubLanguage string
	data        []byte
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
				res.Filetype = magic.MimeFromBytes(data)
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
