# Disk transfer and picture recovery

The supplied DMS has an RLE banner followed by eighty HEAVY2 cylinders. Local
xDMS verifies the archive and recovers a 901,120-byte image. The boot identifies
this supplied release as a 1 MB chip-memory version; that identity is retained.
It has a custom raw-track loader rather than a usable filesystem directory.

The boot copies disk bytes `0x34..0x2c6` into memory `0xc0`, enters supervisor
mode, reads tracks 0–1 to `0x3ec00`, and jumps to `0x3f000` (disk `0x400`).
Following the copied code avoids confusing original disk and runtime addresses.
Each DMS cylinder restores both sides. The resident addresses 160 physical
side tracks, with eleven sectors per side and a 5,632-byte transfer stride.
[`assets/raw/transfers.json`](../assets/raw/transfers.json) records the verified
initial transfers and the later picture/credit stages. Subsequent effect-loader
transfers still need to be mapped before claiming a complete runtime image.

Two standard four-channel modules are present: the loader at disk `0x6e00`,
35,398 bytes, and the main tune at `0x22226`, 230,652 bytes. Their titles are
`?????` and `condom_corruption`; sample-header credits identify Vinnie and
Travolta. Both replay for 270 seconds through DCK's common music facade without
decoder errors. Envelope matching identifies their entrances in the recording
near 5.85 and 44.11 seconds; final synchronization remains to be verified.

## Packed pictures

The unpacker at `0x5273c` uses a twelve-byte header: four DBF offset counters,
a big-endian output size and a packed size. It reads reverse 32-bit words with
a sentinel and backward byte copies. Offset entries consume one additional
bit because the corresponding instruction loops through DBF counter zero.

The Go reader's complete 37,356-byte output for the first packed picture matches
338,608 instructions of the original CPU routine exactly (SHA-256
`376babf81457c628169a8500dc30d6fcff0db56d3d38d5ca6fa4c150cc526b7f`).
That picture is the large **STATE** title. It was initially provisionally
identified as the dragon before the planar preview was decoded.
The four 37,356-byte title/credit pictures begin at source-bank offsets
`0`, `0x14c8`, `0x1fc8` and `0x220c`. Their three 352 by 283 planes have
44-byte rows and 12,452-byte plane spans. The director stores the source
RGB12 palette in its copper register pairs.

The dragon uses a separate 81,920-byte packed region: four 640 by 256 hires
planes. Its copper pointers begin four cleared bytes before the decoded data.
The Go preview retains that offset and its sixteen RGB12 colors, displaying
hires pixels at half width in the low-resolution canvas. Its illustration and
Spaceballs heading are now recovered from the original bank.

An inventory finds 32 plausible packed regions across the disk. Thirty-one
currently decode through the verified reader. One region beginning at disk
`0x1ca80` does not; its original call path and packed-tail handling require
independent CPU checks. The first-title checksum establishes only that
specific original unpacking path, not blanket fidelity of every packed bank.
