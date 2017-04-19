package main

//<!DOCTYPE mame [
//<!ELEMENT mame (machine+)>
//	<!ATTLIST mame build CDATA #IMPLIED>
//	<!ATTLIST mame debug (yes|no) "no">
//	<!ATTLIST mame mameconfig CDATA #REQUIRED>
type Mame struct {
	Build      string `xml:"build,attr"`
	Debug      bool   `xml:"debug,attr"`
	Mameconfig string `xml:"mameconfig,attr"`

	Machines []Machine `xml:"machine"`
}

//	<!ELEMENT machine (description, year?, manufacturer?, biosset*, rom*, disk*, device_ref*, sample*, chip*, display*, sound?, input?, dipswitch*, configuration*, port*, adjuster*, driver?, device*, slot*, softwarelist*, ramoption*)>
//		<!ATTLIST machine name CDATA #REQUIRED>
//		<!ATTLIST machine sourcefile CDATA #IMPLIED>
//		<!ATTLIST machine isbios (yes|no) "no">
//		<!ATTLIST machine isdevice (yes|no) "no">
//		<!ATTLIST machine ismechanical (yes|no) "no">
//		<!ATTLIST machine runnable (yes|no) "yes">
//		<!ATTLIST machine cloneof CDATA #IMPLIED>
//		<!ATTLIST machine romof CDATA #IMPLIED>
//		<!ATTLIST machine sampleof CDATA #IMPLIED>
//		<!ELEMENT description (#PCDATA)>
//		<!ELEMENT year (#PCDATA)>
//		<!ELEMENT manufacturer (#PCDATA)>
type Machine struct {
	Name         string `xml:"name,attr"`
	Sourcefile   string `xml:"sourcefile,attr"`
	Isbios       bool   `xml:"isbios,attr"`
	Isdevice     bool   `xml:"isdevice,attr"`
	Ismechanical bool   `xml:"ismechanical,attr"`
	Runnable     bool   `xml:"runnable,attr"`
	Cloneof      string `xml:"cloneof,attr"`
	Romof        string `xml:"romof,attr"`
	Sampleof     string `xml:"sampleof,attr"`

	Description  string `xml:"description"`
	Year         string `xml:"year"`
	Manufacturer string `xml:"manufacturer"`

	Biossets    []Biosset    `xml:"biosset"`
	Roms        []Rom        `xml:"rom"`
	Disks       []Disk       `xml:"disk"`
	Device_refs []Device_ref `xml:"device_ref"`
	Samples     []Sample     `xml:"sample"`
	Chips       []Chip       `xml:"chip"`
	Displays    []Display    `xml:"display"`
	Sound       Sound        `xml:"sound"`
	Input       Input        `xml:"input"`
	//Dipswitchs     []Dipswitch     `xml:"dipswitch"`
	Configurations []Configuration `xml:"configuration"`
	//Ports          []Port          `xml:"port"`
	Adjusters []Adjuster `xml:"adjuster"`
	Driver    Driver     `xml:"driver"`
	//Devices        []Device        `xml:"device"`
	//Slots          []Slot          `xml:"slot"`
	//Softwarelists  []Softwarelist  `xml:"softwarelist"`
	//Ramoptions     []Ramoption     `xml:"ramoption"`
}

//		<!ELEMENT biosset EMPTY>
//			<!ATTLIST biosset name CDATA #REQUIRED>
//			<!ATTLIST biosset description CDATA #REQUIRED>
//			<!ATTLIST biosset default (yes|no) "no">
type Biosset struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr"`
	Default     bool   `xml:"default,attr"`
}

//		<!ELEMENT rom EMPTY>
//			<!ATTLIST rom name CDATA #REQUIRED>
//			<!ATTLIST rom bios CDATA #IMPLIED>
//			<!ATTLIST rom size CDATA #REQUIRED>
//			<!ATTLIST rom crc CDATA #IMPLIED>
//			<!ATTLIST rom sha1 CDATA #IMPLIED>
//			<!ATTLIST rom merge CDATA #IMPLIED>
//			<!ATTLIST rom region CDATA #IMPLIED>
//			<!ATTLIST rom offset CDATA #IMPLIED>
//			<!ATTLIST rom status (baddump|nodump|good) "good">
//			<!ATTLIST rom optional (yes|no) "no">
type Rom struct {
	Name     string `xml:"name,attr"`
	Bios     string `xml:"bios,attr"`
	Size     string `xml:"size,attr"`
	Crc      string `xml:"crc,attr"`
	Sha1     string `xml:"sha1,attr"`
	Merge    string `xml:"merge,attr"`
	Region   string `xml:"region,attr"`
	Offset   string `xml:"offset,attr"`
	Status   string `xml:"status,attr"`
	Optional bool   `xml:"optional,attr"`

	Availabl bool
}

//		<!ELEMENT disk EMPTY>
//			<!ATTLIST disk name CDATA #REQUIRED>
//			<!ATTLIST disk sha1 CDATA #IMPLIED>
//			<!ATTLIST disk merge CDATA #IMPLIED>
//			<!ATTLIST disk region CDATA #IMPLIED>
//			<!ATTLIST disk index CDATA #IMPLIED>
//			<!ATTLIST disk writable (yes|no) "no">
//			<!ATTLIST disk status (baddump|nodump|good) "good">
//			<!ATTLIST disk optional (yes|no) "no">
type Disk struct {
	Name     string `xml:"name,attr"`
	Sha1     string `xml:"sha1,attr"`
	Merge    string `xml:"merge,attr"`
	Region   string `xml:"region,attr"`
	Index    string `xml:"index,attr"`
	Writable bool   `xml:"writable,attr"`
	Status   string `xml:"status,attr"`
	Optional bool   `xml:"optional,attr"`

	Availabl bool
}

//		<!ELEMENT device_ref EMPTY>
//			<!ATTLIST device_ref name CDATA #REQUIRED>
type Device_ref struct {
	Name string `xml:"name,attr"`
}

//		<!ELEMENT sample EMPTY>
//			<!ATTLIST sample name CDATA #REQUIRED>
type Sample struct {
	Name string `xml:"name,attr"`
}

//		<!ELEMENT chip EMPTY>
//			<!ATTLIST chip name CDATA #REQUIRED>
//			<!ATTLIST chip tag CDATA #IMPLIED>
//			<!ATTLIST chip type (cpu|audio) #REQUIRED>
//			<!ATTLIST chip clock CDATA #IMPLIED>
type Chip struct {
	Name  string `xml:"name,attr"`
	Tag   string `xml:"tag,attr"`
	Type  string `xml:"type,attr"`
	Clock string `xml:"clock,attr"`
}

//		<!ELEMENT display EMPTY>
//			<!ATTLIST display tag CDATA #IMPLIED>
//			<!ATTLIST display type (raster|vector|lcd|unknown) #REQUIRED>
//			<!ATTLIST display rotate (0|90|180|270) #REQUIRED>
//			<!ATTLIST display flipx (yes|no) "no">
//			<!ATTLIST display width CDATA #IMPLIED>
//			<!ATTLIST display height CDATA #IMPLIED>
//			<!ATTLIST display refresh CDATA #REQUIRED>
//			<!ATTLIST display pixclock CDATA #IMPLIED>
//			<!ATTLIST display htotal CDATA #IMPLIED>
//			<!ATTLIST display hbend CDATA #IMPLIED>
//			<!ATTLIST display hbstart CDATA #IMPLIED>
//			<!ATTLIST display vtotal CDATA #IMPLIED>
//			<!ATTLIST display vbend CDATA #IMPLIED>
//			<!ATTLIST display vbstart CDATA #IMPLIED>
type Display struct {
	Tag      string `xml:"tag,attr"`
	Type     string `xml:"type,attr"`
	Rotate   string `xml:"rotate,attr"`
	Flipx    bool   `xml:"flipx,attr"`
	Width    string `xml:"width,attr"`
	Height   string `xml:"height,attr"`
	Refresh  string `xml:"refresh,attr"`
	Pixclock string `xml:"pixclock,attr"`
	Htotal   string `xml:"htotal,attr"`
	Hbend    string `xml:"hbend,attr"`
	Hbstart  string `xml:"hbstart,attr"`
	Vtotal   string `xml:"vtotal,attr"`
	Vbend    string `xml:"vbend,attr"`
	Vbstart  string `xml:"vbstart,attr"`
}

//		<!ELEMENT sound EMPTY>
//			<!ATTLIST sound channels CDATA #REQUIRED>
type Sound struct {
	Channels string `xml:"channels,attr"`
}

//		<!ELEMENT input (control*)>
//			<!ATTLIST input service (yes|no) "no">
//			<!ATTLIST input tilt (yes|no) "no">
//			<!ATTLIST input players CDATA #REQUIRED>
//			<!ATTLIST input buttons CDATA #IMPLIED>
//			<!ATTLIST input coins CDATA #IMPLIED>
type Input struct {
	Service bool   `xml:"service,attr"`
	Tilt    bool   `xml:"tilt,attr"`
	Players string `xml:"players,attr"`
	Buttons string `xml:"buttons,attr"`
	Coins   string `xml:"coins,attr"`

	Controls []Control `xml:"control"`
}

//			<!ELEMENT control EMPTY>
//				<!ATTLIST control type CDATA #REQUIRED>
//				<!ATTLIST control minimum CDATA #IMPLIED>
//				<!ATTLIST control maximum CDATA #IMPLIED>
//				<!ATTLIST control sensitivity CDATA #IMPLIED>
//				<!ATTLIST control keydelta CDATA #IMPLIED>
//				<!ATTLIST control reverse (yes|no) "no">
//				<!ATTLIST control ways CDATA #IMPLIED>
//				<!ATTLIST control ways2 CDATA #IMPLIED>
//				<!ATTLIST control ways3 CDATA #IMPLIED>
type Control struct {
	Type        string `xml:"type,attr"`
	Minimum     string `xml:"minimum,attr"`
	Maximum     string `xml:"maximum,attr"`
	Sensitivity string `xml:"sensitivity,attr"`
	Keydelta    string `xml:"keydelta,attr"`
	Reverse     bool   `xml:"reverse,attr"`
	Ways        string `xml:"ways,attr"`
	Ways2       string `xml:"ways2,attr"`
	Ways3       string `xml:"ways3,attr"`
}

//		<!ELEMENT dipswitch (dipvalue*)>
//			<!ATTLIST dipswitch name CDATA #REQUIRED>
//			<!ATTLIST dipswitch tag CDATA #REQUIRED>
//			<!ATTLIST dipswitch mask CDATA #REQUIRED>
type Dipswitch struct {
	Name string `xml:"name,attr"`
	Tag  string `xml:"tag,attr"`
	Mask string `xml:"mask,attr"`

	Dipvalues []Dipvalue `xml:"dipvalue"`
}

//			<!ELEMENT dipvalue EMPTY>
//				<!ATTLIST dipvalue name CDATA #REQUIRED>
//				<!ATTLIST dipvalue value CDATA #REQUIRED>
//				<!ATTLIST dipvalue default (yes|no) "no">
type Dipvalue struct {
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
	Default bool   `xml:"default,attr"`
}

//		<!ELEMENT configuration (confsetting*)>
//			<!ATTLIST configuration name CDATA #REQUIRED>
//			<!ATTLIST configuration tag CDATA #REQUIRED>
//			<!ATTLIST configuration mask CDATA #REQUIRED>
type Configuration struct {
	Name string `xml:"name,attr"`
	Tag  string `xml:"tag,attr"`
	Mask string `xml:"merge,attr"`

	Confsettings []Confsetting `xml:"confsetting"`
}

//			<!ELEMENT confsetting EMPTY>
//				<!ATTLIST confsetting name CDATA #REQUIRED>
//				<!ATTLIST confsetting value CDATA #REQUIRED>
//				<!ATTLIST confsetting default (yes|no) "no">
type Confsetting struct {
	Name    string `xml:"name,attr"`
	Value   string `xml:"value,attr"`
	Default bool   `xml:"default,attr"`
}

//		<!ELEMENT port (analog*)>
//			<!ATTLIST port tag CDATA #REQUIRED>
type Port struct {
	Tag string `xml:"tag,attr"`

	Analogs []Analog `xml:"analog"`
}

//			<!ELEMENT analog EMPTY>
//				<!ATTLIST analog mask CDATA #REQUIRED>
type Analog struct {
	Mask string `xml:"mask,attr"`
}

//		<!ELEMENT adjuster EMPTY>
//			<!ATTLIST adjuster name CDATA #REQUIRED>
//			<!ATTLIST adjuster default CDATA #REQUIRED>
type Adjuster struct {
	Name    string `xml:"name,attr"`
	Default bool   `xml:"default,attr"`
}

//		<!ELEMENT driver EMPTY>
//			<!ATTLIST driver status (good|imperfect|preliminary) #REQUIRED>
//			<!ATTLIST driver emulation (good|imperfect|preliminary) #REQUIRED>
//			<!ATTLIST driver color (good|imperfect|preliminary) #REQUIRED>
//			<!ATTLIST driver sound (good|imperfect|preliminary) #REQUIRED>
//			<!ATTLIST driver graphic (good|imperfect|preliminary) #REQUIRED>
//			<!ATTLIST driver cocktail (good|imperfect|preliminary) #IMPLIED>
//			<!ATTLIST driver protection (good|imperfect|preliminary) #IMPLIED>
//			<!ATTLIST driver savestate (supported|unsupported) #REQUIRED>
type Driver struct {
	Status     string `xml:"status,attr"`
	Emulation  string `xml:"emulation,attr"`
	Color      string `xml:"color,attr"`
	Sound      string `xml:"sound,attr"`
	Graphic    string `xml:"graphic,attr"`
	Cocktail   string `xml:"cocktail,attr"`
	Protection string `xml:"protection,attr"`
	Savestate  string `xml:"savestate,attr"`
}

//		<!ELEMENT device (instance*, extension*)>
//			<!ATTLIST device type CDATA #REQUIRED>
//			<!ATTLIST device tag CDATA #IMPLIED>
//			<!ATTLIST device mandatory CDATA #IMPLIED>
//			<!ATTLIST device interface CDATA #IMPLIED>
type Device struct {
	Type      string `xml:"type,attr"`
	Tag       string `xml:"tag,attr"`
	Mandatory string `xml:"mandatory,attr"`
	Interface string `xml:"interface,attr"`

	Instances  []Instance  `xml:"instance"`
	Extensions []Extension `xml:"extension"`
}

//			<!ELEMENT instance EMPTY>
//				<!ATTLIST instance name CDATA #REQUIRED>
//				<!ATTLIST instance briefname CDATA #REQUIRED>
type Instance struct {
	Name      string `xml:"name,attr"`
	Briefname string `xml:"briefname,attr"`
}

//			<!ELEMENT extension EMPTY>
//				<!ATTLIST extension name CDATA #REQUIRED>
type Extension struct {
	Name string `xml:"name,attr"`
}

//		<!ELEMENT slot (slotoption*)>
//			<!ATTLIST slot name CDATA #REQUIRED>
type Slot struct {
	Name string `xml:"name,attr"`

	Slotoptions []Slotoption `xml:"slotoption"`
}

//			<!ELEMENT slotoption EMPTY>
//				<!ATTLIST slotoption name CDATA #REQUIRED>
//				<!ATTLIST slotoption devname CDATA #REQUIRED>
//				<!ATTLIST slotoption default (yes|no) "no">
type Slotoption struct {
	Name    string `xml:"name,attr"`
	Devname string `xml:"devname,attr"`
	Default bool   `xml:"default,attr"`
}

//		<!ELEMENT softwarelist EMPTY>
//			<!ATTLIST softwarelist name CDATA #REQUIRED>
//			<!ATTLIST softwarelist status (original|compatible) #REQUIRED>
//			<!ATTLIST softwarelist filter CDATA #IMPLIED>
type Softwarelist struct {
	Name   string `xml:"name,attr"`
	Status string `xml:"status,attr"`
	Filter bool   `xml:"filter,attr"`
}

//		<!ELEMENT ramoption (#PCDATA)>
//			<!ATTLIST ramoption default CDATA #IMPLIED>
type Ramoption struct {
	Default string `xml:"default,attr"`
}

//]>
