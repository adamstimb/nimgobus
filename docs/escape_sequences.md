# Escape Sequences

The original RM Nimbus documentation recommended using escape sequences to control the text output (e.g. setting screen mode, setting colours and so on).  It supported ANSI escape sequences augmented by some platform-specific commands.  These have - for the most part - been reimplemented in nimgobus with some modifications, particularly if the original implementation was awkward or confusing (e.g. setting forground and background colours).

Escape sequences can be sent to the console using [Printf](#printfs-string).  A sequence must begin with the escape `\x1b` char followed by an open square-bracket `[` char.  A sequence may contain one or more optional parameters separated by a semicolon and ending with either a single letter indicating an ANSI sequence or a tilde followed by a letter indicating a Nimbus-specific sequence.

Currently implemented features are marked with a ✅ while never-to-be-implemented features are ~~crossed-out~~.  None of the key reassignment features were implemented since that kind of scenario is better handled by ebiten.

## Cursor Movement

✅ **CUU Cursor Up**

[ _n_ A

Move cursor up n rows.

✅ **CUD Cursor Down** 		

[ _n_ B

Move cursor down n rows.

✅ **CUF Cursor Forward**

[ _n_ C

Move cursor forwards n rows.

✅ **CUB Cursor Backward**

[ _n_ D

Move cursor backwards n rows.

✅ **CUP Cursor Position** 				

[ _r_ ; _c_ H

- If no parameters given, cursor goes home
- If outside scrolling area, cursor goes home
- If r is within but c is outside:
	- Cursor goes to r and left-most column
- If c is within but r is outside:
	- Cursor goes to c and top row

HVP Horizontal/Vertical Position is equivalent so has not been implemented.

## Cursor Positioning

✅ **SCP Save Cursor Position** 			

[ s

✅ **RCP Restore Cursor Position** 			

[ u

**~~CPR Cursor Position Report~~**	

~~[ R~~

**~~DSR Device Status Report~~**			

~~[ 6 n~~

_Instead of CPR and DSR use Stdio.GetCurpos() which simply returns the current row, column._

## Deleting and Inserting

✅ **ED Erase in Display** 				

[ _n_ J

- n=0 erase chars from the cursor position to the end of the row (including the char at the cursor position), and all rows below. Cursor does not move.
- n=1 erase chars from the beginning of the row to the cursor (including the one at the cursor position), and all rows above. Cursor does not move.
- n=2 erase whole scrolling area and move cursor home.

✅ **EL Erase in Line**					

[ _n_ K

- n=0 erase chars from the cursor position to the end of the row. Cursor does not move. Cursor does not move.
- n=1 erase chars from start of row to cursor (including the char at cursor position). Cursor does not move.
- n=2 erase whole row. Cursor does not move.

**~~IL Insert Line~~**						

~~[ _n_ L~~

- ~~Causes the row containing the cursor and following rows to be moved down n rows, with the insertion of n blank lines.~~

**~~DL Delete lines~~**

[ _n_ M

- ~~Causes the row containing the cursor and the following rows to be moved up n rows, with the deletion of n rows.~~

## Scrolling

✅ **SU Scroll Up** 						

[ _n_ S

- The current scrolling area scrolls up n rows. Cursor does not move.

✅ **SD Scroll Down** 					

[ _n_ T

- The current scrolling area scrolls down n rows. Cursor does not move.

## Modes of Operation

✅ **DSA Define Scrolling Area**

[ _t_ ; _l_ ; _b_ ; _r_ ~B

- Defines scrolling area in which the cursor will operate and the chars will be displayed:
- t : top row, l : left column, b : bottom row, r : right column

✅ **SCA Set Character Attribute**

[ _m_ ; _n_ ; _p_ ~E

- m=0 not underlined
- m=1 underlined
- n=0 standard charset
- n=1 alternative charset
- p=0 XOR writing off
- p=1 XOR writing on

✅ **SCM Set Cursor Mode** 				

[ _n_ ; _q_ ; _r_ ; _m_ ; _p_ ~A

- Defines cursor attributes and character
- n=0 not underlined
- n=1 underlined
- q=0 standard charset
- q=1 alternative charset
- r=0 not flashing (see also SCLT)
- r=1 flashing
- p=0 cursor is displayed
- p=1 cursor not displayed

✅ **CV  Cursor Visible**

[ ~G

✅ **CNV Cursor Not Visible**

[~F

- CV and CNV cursor on or off without affecting attributes

✅ **PCC Print Control Characters**

[ _n_ ~D

- Enable printing of chars 0-31
- n=0 on, n=1 off

✅ **SGR Set Graphics Rendition** 			

[ _n_ ; _n_ ; _n_ ... m

- See also SCLT
- As shown in the tables below, this is implemented slightly differently in nimgobus.  Originally there was some confusing stuff about bold versus faint which used dark or lighter tones of the same colour as a kind-of bold effect, blinking and "reversing" video which might have made sense 40 years ago but I think it's just really obfuscated, so all these attributes have been removed.  Instead, foreground colours simply begin at 30 and background (paper) colours begin at 50, and correspond to logical colours 0-3 (80 column mode) or 0-15 (40 column mode).

_80 column mode_

| Parameter | Action | Effect on logical colour |
| --------- | ------ | ------------------------ |
| 0 	| All attributes off | f=1,b=0 |
| ~~1~~ 	| ~~Bold on~~ | ~~f=3~~ |
| ~~2~~ 	| ~~Faint~~ | ~~f=1~~ |
| 4 	| Underline on| |
| ~~5+6~~ 	| ~~Blink on (use SCLT to set rate)~~ | ~~f=2~~ |
| ~~7~~ 	| ~~Reverse video~~ | ~~F and B swapped~~ |
| ~~8~~ 	| ~~Concealed on (same f and g)~~ | ~~f=0,b=0~~ |
| 10 	| Select standard charset | |
| 11 	| Select alternative charset | |
| 24 	| Underline off | |
| ~~25~~ 	| ~~Blink off~~ | ~~f=1~~ |
| ~~27~~ 	| ~~Normal video~~ | ~~F and B swapped back if currently in reverse~~ |
| | **As documented:**  | |
| 30/34 | Black fg | f=0 |
| 31/35 | Light grey fg | f=1 |
| 32/36 | White/black flashing fg | f=2 |
| 33/37 | White fg | f=3 |
| 40/44 | Black bg | b=0 |
| 41/45 | Light grey bg | b=1 |
| 42/46 | White/black flashing bg | b=2 |
| 43/47 | White bg | b=3 |
|  | **As implemented in nimgobus:** | |
| 30 | Black fg | f=0 |
| 31 | Light grey fg | f=1 |
| 32 | White/black flashing fg | f=2 |
| 33 | White fg | f=3 |
| 50 | Black bg | b=0 |
| 51 | Light grey bg | b=1 |
| 52 | White/black flashing bg | b=2 |
| 53 | White bg | b=3 |

_40 column mode_

| Parameter | Action | Effect on logical colour |
| --------- | ------ | ------------------------ |
| 0 	| All attributes off | f=7,b=0 |
| ~~1~~ 	| ~~Bold on~~ | ~~no change if f < 8, otherwise f+=8~~ |
| ~~2~~ 	| ~~Faint~~ | ~~no change if f > 8, otherwise f-=8~~ |
| 4 	| Underline on | |
| ~~5+6~~ 	| ~~Blink on (use SCLT to set rate)~~ | ~~f=2~~ |
| ~~7~~ 	| ~~Reverse video~~ | ~~F and B swapped~~ |
| ~~8~~ 	| ~~Concealed on (same f and g)~~ | ~~f=0,b=0~~ |
| 10 	| Select standard charset | |
| 11 	| Select alternative charset | |
| 24 	| Underline off | |
| ~~25~~ 	| ~~Blink off~~ | ~~f=1~~ |
| ~~27~~ 	| ~~Normal video~~ | ~~F and B swapped back if currently in reverse~~ |
| | **As documented:**  | |
| 30 | Black fg | f=0/8 |
| 31 | Red fg | f=1/9 |
| 32 | Green fg | f=2/10 |
| 33 | Yellow fg | f=3/11 |
| 34 | Blue fg | f=4/12 |
| 35 | Magenta fg | f=5/13 |
| 36 | Cyan fg| f=6/14 |
| 37 | White fg | f=7/15 |
| 40 | Black bg | b=0 |
| 41 | Red bg | b=1 |
| 42 | Green bg | b=2 |
| 43 | Yellow bg | b=3 |
| 44 | Blue bg | b=4 |
| 45 | Magenta bg | b=5 |
| 46 | Cyan bg| b=6 |
| 47 | White bg | b=7 |
|  | **As implemented in nimgobus:** | |
| 30 | Black fg | f=0 |
| 31 | Dark Red fg | f=1 |
| 32 | Dark Green fg | f=2 |
| 33 | Brown fg | f=3 |
| 34 | Dark Blue fg | f=4 |
| 35 | Purple fg | f=5 |
| 36 | Dark Cyan fg| f=6 |
| 37 | Light Grey fg | f=7 |
| 38 | Dark Grey fg | f=8 |
| 39 | Light Red fg| f=9 |
| 40 | Light Green fg | f=10 |
| 41 | Yellow fg | f=11 |
| 42 | Light Blue fg | f=12 |
| 43 | Light Purple fg | f=13 |
| 44 | Light Cyan fg | f=14 |
| 45 | White fg | f=15 |
| 50 | Black bg | b=0 |
| 51 | Dark Red bg | b=1 |
| 52 | Dark Green bg | b=2 |
| 53 | Brown bg | b=3 |
| 54 | Dark Blue bg | b=4 |
| 55 | Purple bg | b=5 |
| 56 | Dark Cyan bg| b=6 |
| 57 | Light Grey bg | b=7 |
| 58 | Dark Grey bg | b=8 |
| 59 | Light Red bg| b=9 |
| 60 | Light Green bg | b=10 |
| 61 | Yellow bg | b=11 |
| 62 | Light Blue bg | b=12 |
| 63 | Light Purple bg | b=13 |
| 64 | Light Cyan bg | b=14 |
| 65 | White bg | b=15 |

✅ **SCLT Set Colour Lookup Table**

[ _q_ ; _n_ ; _m_ ; _f_ ; _p_ ~C

- Update the colour lookup table.  As in the original Nimbus we have two lookup tables: the Mode 80 high-resolution table with 4 colours, and the Mode 40 low-resolution table with 16 colours.  Modes can be switched without losing the values set in the lookup tables.
- When working with graphics, the graphics colour lookup table will always be synchronized to the current console colour lookup table.  Furthermore if the graphics colour lookup table is updated, these changes will be reflected back in to the console colour lookup table.  I'm not sure if this is how the original Nimbus behaved but it's probably not important.  It's certainly not important to me anyway!
- Unlike the original Nimbus the order of the console logical colours is the same as the graphical logical colours - this just seemed unnecessarily confusing otherwise.
- q=80 update the 80 column mode table.
- q=40            40 column
- n=0-15 defines the logical colour number
- m=0-15 defines the physical colour (from the table below) assigned to n
- f=0 flashing off, f=1 slow flashing between m and p, f=2 fast flashing
- p=0-15 defines the physical flashing colour

| Number | Physical Colour |
| ------ | --------------- |
| 0 | black |
| 1 | ~~dark red~~ dark blue |
| 2 | ~~dark green~~ dark red |
| 3 | ~~brown~~ purple |
| 4 | ~~dark blue~~ dark green |
| 5 | ~~purple~~ dark cyan |
| 6 | ~~dark cyan~~ brown |
| 7 | light grey |
| 8 | dark grey |
| 9 | ~~light red~~ light blue |
| 10 | light green |
| 11 | ~~yellow~~ light purple |
| 12 | ~~light blue~~ light green |
| 13 | ~~light purple~~ light cyan |
| 14 | ~~cyan~~ yellow |
| 15 | white |

✅ **SM Set Mode**

[ _n_ h

- Set screen column mode.
- n=2 ~~,3,6~~ sets 80 column mode, clears screen, sets previously defined CLT
- n=0 ~~,1,4,5~~ sets 40 column mode, clears screen, sets previously defined CLT
- n=7 sets wrap on: cursor moves to start of next row when it reaches the end of a row.

✅ **RM Reset Mode** 						

[ _n_ l

- Same effect as SM except n=7 sets wrap off: cursor stays at end of row and chars are printed on top of each other

✅ **RIS Reset to Initial State**

[ c

- Resets the following:
- Screen cleared, and default column mode (80)
- Cursor reset and sent home.
- Auto repeat reset.
- Scrolling area set to full screen.
- Default CLT set.
- Buffers flushed.
- Wrap-on mode set.
- Control chars not printable.