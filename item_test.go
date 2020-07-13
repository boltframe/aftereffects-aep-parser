package aep

import (
	"testing"
)

func TestItemMetadata(t *testing.T) {
	project, err := Open("data/Item-01.aep")
	if err != nil {
		t.Fatal(err)
	}

	expect(t, project.RootFolder.ID, uint32(0), "Root folder should have 0th ID")

	folder01 := project.RootFolder.FolderContents[0]
	expect(t, folder01.Name, "Folder 01")
	expect(t, folder01.ID, uint32(46))
	expect(t, folder01.TypeName, ItemTypeFolder)

	folder02 := project.RootFolder.FolderContents[1]
	expect(t, folder02.Name, "Folder 02")
	expect(t, folder02.ID, uint32(47))
	expect(t, folder02.TypeName, ItemTypeFolder)

	comp01 := folder01.FolderContents[0]
	expect(t, comp01.Name, "Comp 01")
	expect(t, comp01.ID, uint32(48))
	expect(t, comp01.TypeName, ItemTypeComposition)
	expect(t, comp01.FootageDimensions, [2]uint16{351, 856})
	expect(t, comp01.FootageFramerate, float64(21))
	expect(t, comp01.FootageSeconds, float64(31))
	expect(t, comp01.BackgroundColor, [3]byte{15, 75, 82})

	comp02 := folder02.FolderContents[0]
	expect(t, comp02.Name, "Comp 02")
	expect(t, comp02.ID, uint32(59))
	expect(t, comp02.TypeName, ItemTypeComposition)
	expect(t, comp02.FootageDimensions, [2]uint16{452, 639})
	expect(t, comp02.FootageFramerate, float64(29.97))
	expect(t, comp02.FootageSeconds, float64(71.338004671338))
	expect(t, comp02.BackgroundColor, [3]byte{145, 206, 85})

	footageFolder := project.RootFolder.FolderContents[2]
	expect(t, footageFolder.Name, "Footage")
	expect(t, footageFolder.ID, uint32(70))

	placeholderFootage := footageFolder.FolderContents[2]
	expect(t, placeholderFootage.Name, "Missing Footage")
	expect(t, placeholderFootage.ID, uint32(71))
	expect(t, placeholderFootage.TypeName, ItemTypeFootage)
	expect(t, placeholderFootage.FootageSeconds, float64(127))
	expect(t, placeholderFootage.FootageFramerate, float64(123.45669555664062))
	expect(t, placeholderFootage.FootageDimensions, [2]uint16{1234, 5678})

	redSolid := footageFolder.FolderContents[3]
	expect(t, redSolid.FootageType, FootageTypeSolid)
	expect(t, redSolid.Name, "Red Solid 1")
}
