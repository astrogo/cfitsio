package cfitsio

// Column represents a column in a FITS table
type Column struct {
	Name    string // column name, corresponding to ``TTYPE`` keyword
	Format  string // column format, corresponding to ``TFORM`` keyword
	Unit    string // column unit, corresponding to ``TUNIT`` keyword
	Null    string // null value, corresponding to ``TNULL`` keyword
	Bscale  int    // bscale value, corresponding to ``TSCAL`` keyword
	Bzero   int    // bzero value, corresponding to ``TZERO`` keyword
	Display string // display format, corresponding to ``TDISP`` keyword
	Dim     int    // column dimension corresponding to ``TDIM`` keyword
	Start   int    // column starting position, corresponding to ``TBCOL`` keyword
	Ascii   bool   // whether this describes a column for an ASCII table
}

// ColDefs is a list of Column definitions
type ColDefs []Column

// EOF
