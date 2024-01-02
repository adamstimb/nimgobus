package sprite

// Sprite represents (as much as is practical) Nimbus sprite data.
type Sprite struct {
	HighResolution bool      // Set to true if the sprite is intended for high-resolution mode.
	Hotspot        [2]int    // The x, y vector from the bottom-left of the sprite to the hotspot.
	Poses          [][][]int // The sprite images stored as 2D arrays in individual poses.  Up to 2 poses allowed in low-resolution mode, up to 4 in high-resolution mode.
}

// saveTableRow describes a row in the SaveTable type
type SaveTableRow struct {
	X int // x position
	Y int // y position
	C int // colour
}

// SaveTable describes a save table used in sprites
type SaveTable struct {
	Rows []SaveTableRow
}
