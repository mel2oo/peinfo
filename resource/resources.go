package resource

import (
	"encoding/binary"
	"errors"
	"expvar"
	"fmt"
	"io"
	"strings"
)

const (
	seekStart = 0
	// maxAllowedResourceEntries determines the maximum number of directory entries we are allowed to process
	maxAllowedResourceEntries = 4096
)

var (
	maxResourceEntriesExceeded = expvar.NewInt("pe.max.entries.exceeded")
	failedResourceEntryReads   = expvar.NewInt("pe.failed.entry.reads")

	errMalformedDir         = errors.New("malformed directory")
	errMaxAllowedDirEntries = func(nbEntries int) error {
		return fmt.Errorf("the directory contains %d entries. Max allowed entries %d", nbEntries, maxAllowedResourceEntries)
	}
)

// rdir represents the resource directory with its entries
type rdir struct {
	directory Directory
	entries   []rdirEntry
}

type rdirEntry struct {
	id    ID
	level uint16
	entry DirectoryEntry
	data  DataEntry
	dir   *rdir
}

// readResources a plenty of logic in this code is inspired by pefile tool. The native stdlib package doesn't offer
// any kind of introspection on PE resources, so this function takes care of parsing the resource directory and extracting
// the version resources from it.
func (r *reader) readResources(rva uint32) (map[string]string, error) {
	var vers map[string]string
	dir, err := r.readResourcesDirectory(rva, 0, 0, nil)
	if err != nil {
		return nil, err
	}

	for _, entry := range dir.entries {
		if entry.dir == nil {
			continue
		}
		if entry.level == 0 {
			switch entry.id {
			case Version:
				// read version information resources
				vers, err = r.readVersionInfo(entry.dir)
				if err != nil {
					continue
				}
			}
		}
	}

	return vers, nil
}

func (r *reader) readResourcesDirectory(rva uint32, baseRVA uint32, level uint16, dirs []uint32) (*rdir, error) {
	if dirs == nil {
		dirs = []uint32{rva}
	}
	if baseRVA == 0 {
		baseRVA = rva
	}
	sr := io.NewSectionReader(r.f, 0, 1<<63-1)
	offset, err := r.FindOffsetByRVA(rva)
	if err != nil {
		return nil, fmt.Errorf("couldn't read resources directory: %v", err)
	}
	// try to read the resource directory structure that is basically
	// a header of the table preceding the actual resource entries
	if _, err := sr.Seek(offset, seekStart); err != nil {
		return nil, err
	}
	var dir Directory
	if err := binary.Read(sr, binary.LittleEndian, &dir); err != nil {
		return nil, err
	}

	nbEntries := int(dir.NumberIDEntries + dir.NumberNamedEntries)
	dirents := make([]rdirEntry, nbEntries)

	// we have to protect us against reading a huge number of entries
	if nbEntries > maxAllowedResourceEntries {
		maxResourceEntriesExceeded.Add(1)
		return nil, errMaxAllowedDirEntries(nbEntries)
	}
	// advance the RVA to the position following the directory table
	// header that points the the first entry in the table of entries
	rva += uint32(dir.Size())

loop:
	for i := 0; i < nbEntries; i++ {
		res, err := r.readResourceEntry(sr, rva)
		if err != nil {
			failedResourceEntryReads.Add(1)
			continue
		}
		// the entry is a directory so we have to parse it recursively
		if res.IsDir() {
			// The following comment is from pefile.py
			//
			// OC Patch:
			//
			// One trick malware can do is to recursively reference
			// the next directory. This causes hilarity to ensue when
			// trying to parse everything correctly.
			// If the original RVA given to this function is equal to
			// the next one to parse, we assume that it's a trick.
			for _, dir := range dirs {
				if baseRVA+res.DirOffset() == dir {
					break loop
				}
			}
			dirs = append(dirs, baseRVA+res.DirOffset())
			dir, err := r.readResourcesDirectory(baseRVA+res.DirOffset(), baseRVA, level+1, dirs)
			if err != nil {
				break
			}
			if dir == nil {
				break
			}
			dirents[i] = rdirEntry{
				id:    res.ID(),
				entry: res,
				dir:   dir,
				level: level,
			}
		} else {
			// if we reached the actual directory data, let's read the structure that
			// contains the offset and size of the resource's data
			data, err := r.readResourceData(sr, baseRVA+res.DirOffset())
			if err != nil {
				break
			}
			dirents[i] = rdirEntry{
				id:    res.ID(),
				entry: res,
				data:  data,
				level: level,
			}
		}
		// increment the RVA to the next directory entry
		rva += uint32(res.Size())
	}

	return &rdir{directory: dir, entries: dirents}, nil
}

func (r *reader) readResourceEntry(sr *io.SectionReader, rva uint32) (DirectoryEntry, error) {
	offset, err := r.FindOffsetByRVA(rva)
	if err != nil {
		return DirectoryEntry{}, err
	}
	if _, err := sr.Seek(offset, seekStart); err != nil {
		return DirectoryEntry{}, err
	}
	var entry DirectoryEntry
	if err := binary.Read(sr, binary.LittleEndian, &entry); err != nil {
		return DirectoryEntry{}, fmt.Errorf("invalid directory entry at RVA 0x%x: %v", rva, err)
	}
	return entry, nil
}

func (r *reader) readResourceData(sr *io.SectionReader, rva uint32) (DataEntry, error) {
	offset, err := r.FindOffsetByRVA(rva)
	if err != nil {
		return DataEntry{}, err
	}
	if _, err := sr.Seek(offset, seekStart); err != nil {
		return DataEntry{}, err
	}
	var dataEntry DataEntry
	if err := binary.Read(sr, binary.LittleEndian, &dataEntry); err != nil {
		return DataEntry{}, fmt.Errorf("invalid resource data at RVA 0x%x: %v", rva, err)
	}
	return dataEntry, nil
}

func (r *reader) readVersionInfo(vsDir *rdir) (map[string]string, error) {
	if vsDir == nil || len(vsDir.entries) == 0 {
		return nil, errMalformedDir
	}
	dir := vsDir.entries[0].dir
	if dir == nil {
		return nil, errMalformedDir
	}

	errs := make([]error, 0)
	vers := make(map[string]string)
	vdents := dir.entries
	sr := io.NewSectionReader(r.f, 0, 1<<63-1)

	for _, ve := range vdents {
		offsetToData := ve.data.OffsetToData
		startOffset, err := r.FindOffsetByRVA(offsetToData)
		if err != nil {
			return nil, err
		}
		// read the version info structure and the VS_VERSION_INFO string
		versionInfo, versionString, err := r.parseVersionInfo(sr, startOffset, offsetToData)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// if we've able to correctly parse the VS_VERSION_INFO string, the next step
		// is to process the fixed version information by getting the offset of the struct
		fixedFileinfoOffset := dwordAlign(int64(versionInfo.Size()+(2*len(versionString))+1), int64(offsetToData))
		fixedFileinfo, err := r.parseFixedFileinfoStruct(sr, fixedFileinfoOffset)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// now the real work begins. To reach version keys/values, we first have to parse all
		// of the StringFileInfo and VarFileInfo structures until we get to string table whose
		// entries store the data we're after
		stringFileinfoOffset := dwordAlign(fixedFileinfoOffset+int64(fixedFileinfo.Size()), int64(offsetToData))
		for {
			// process StringFileInfo/VarFileInfo structures. The file info string determines whether we
			// should process StringFileInfo or VarFileInfo items
			stringFileinfo, fileInfoString, err := r.parseStringFileinfo(
				sr,
				startOffset+stringFileinfoOffset,
				uint32(int64(offsetToData)+stringFileinfoOffset)+uint32(versionInfo.Size()),
			)
			if err != nil {
				errs = append(errs, err)
				break
			}

			switch {
			case strings.HasPrefix(fileInfoString, "StringFileInfo"):
				if stringFileinfo.Skip() {
					continue
				}
				stringTableOffset := dwordAlign(stringFileinfoOffset+int64(stringFileinfo.Size()+2*(len(fileInfoString)+1)), int64(offsetToData))
				// now we can start processing all the StringTable entries that contain the k/v pairs
				for {
					stringTable, langID, err := r.parseStringTable(
						sr,
						startOffset+stringTableOffset,
						offsetToData+uint32(stringTableOffset),
					)
					if err != nil {
						errs = append(errs, err)
						break
					}
					// now we can process all the entries in the string table and populate the result map
					entryOffset := dwordAlign(stringTableOffset+int64(stringTable.Size()+(2*len(langID)+1)), int64(offsetToData))
					for entryOffset < stringTableOffset+int64(stringTable.Length) {
						if _, err := sr.Seek(startOffset+entryOffset, 0); err != nil {
							break
						}
						var str String
						if err := binary.Read(sr, binary.LittleEndian, &str); err != nil {
							break
						}

						key, err := r.readUTF16String(offsetToData + uint32(entryOffset) + uint32(str.Size()))
						if err != nil {
							break
						}
						valueOffset := dwordAlign(int64(2*(len(key)+1))+entryOffset+int64(str.Size()), int64(offsetToData))
						value, err := r.readUTF16String(offsetToData + uint32(valueOffset))
						if err != nil {
							// couldn't read the value but still index the key
							vers[key] = ""
							break
						}
						vers[key] = strings.ReplaceAll(value, "�", "")
						if str.Length == 0 {
							entryOffset = stringTableOffset + int64(stringTable.Length)
						} else {
							entryOffset = dwordAlign(int64(str.Length)+entryOffset, int64(offsetToData))
						}
					}

					// these checks breaks on the entries that could lead to infinite loops
					newStringtableOffset := dwordAlign(int64(stringTable.Length)+stringTableOffset, int64(offsetToData))
					if newStringtableOffset == stringTableOffset {
						break
					}

					stringTableOffset = newStringtableOffset
					if stringTableOffset >= int64(stringFileinfo.Length) {
						break
					}
				}
			case strings.HasPrefix(fileInfoString, "VarFileInfo"):
				// TODO
			default:
				errs = append(errs, fmt.Errorf("unknown StringFileInfo string: %s", fileInfoString))
				// TODO
			}
			// increment and align the string file info offset. Use the offset to check if we've
			// consumed all the StringFileInfo and VarFileinfo items so we can break the loops
			stringFileinfoOffset = dwordAlign(int64(stringFileinfo.Length)+stringFileinfoOffset, int64(offsetToData))
			if stringFileinfo.Length == 0 || stringFileinfoOffset >= int64(versionInfo.Length) {
				break
			}
		}
	}

	if len(vers) == 0 && len(errs) > 0 {
		return nil, Wrap(errs...)
	}

	return vers, nil
}

func (r *reader) parseVersionInfo(sr *io.SectionReader, startOffset int64, rva uint32) (VersionInfo, string, error) {
	if _, err := sr.Seek(startOffset, seekStart); err != nil {
		return VersionInfo{}, "", err
	}
	var versionInfo VersionInfo
	if err := binary.Read(sr, binary.LittleEndian, &versionInfo); err != nil {
		return VersionInfo{}, "", err
	}
	versionString, err := r.readUTF16String(rva + uint32(versionInfo.Size()))
	if err != nil || versionString != "VS_VERSION_INFO" {
		return VersionInfo{}, "", fmt.Errorf("invalid VS_VERSION_INFO block: %s", versionString)
	}
	return versionInfo, versionString, nil
}

func (r *reader) parseFixedFileinfoStruct(sr *io.SectionReader, offset int64) (FixedFileinfo, error) {
	if _, err := sr.Seek(offset, seekStart); err != nil {
		return FixedFileinfo{}, err
	}
	var fixedFileinfo FixedFileinfo
	if err := binary.Read(sr, binary.LittleEndian, &fixedFileinfo); err != nil {
		return FixedFileinfo{}, err
	}
	return fixedFileinfo, nil
}

func (r *reader) parseStringFileinfo(sr *io.SectionReader, offset int64, rva uint32) (StringFileInfo, string, error) {
	var stringFileinfo StringFileInfo
	if _, err := sr.Seek(offset, seekStart); err != nil {
		return StringFileInfo{}, "", err
	}
	if err := binary.Read(sr, binary.LittleEndian, &stringFileinfo); err != nil {
		return StringFileInfo{}, "", fmt.Errorf("couldn't parse StringFileInfo/VarFileInfo structure: %v", err)
	}
	str, err := r.readUTF16String(rva)
	if err != nil {
		return StringFileInfo{}, "", fmt.Errorf("couldn't read StringFileInfo unicode string at RVA 0x%x: %v", rva, err)
	}
	return stringFileinfo, str, nil
}

func (r *reader) parseStringTable(sr *io.SectionReader, offset int64, rva uint32) (StringTable, string, error) {
	if _, err := sr.Seek(offset, seekStart); err != nil {
		return StringTable{}, "", err
	}
	var stringTable StringTable
	if err := binary.Read(sr, binary.LittleEndian, &stringTable); err != nil {
		return StringTable{}, "", fmt.Errorf("couldn't parse StringTable structure: %v", err)
	}
	langID, err := r.readUTF16String(rva + uint32(stringTable.Size()))
	if err != nil {
		return StringTable{}, "", fmt.Errorf("couldn't read StringTable unicode string at RVA 0x%x: %v", rva+uint32(stringTable.Size()), err)
	}
	return stringTable, langID, nil
}
