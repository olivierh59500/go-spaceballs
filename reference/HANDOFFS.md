# Loading passages belong to the scene sequence

The first spirograph's 480 ticks are followed by a distinct program entry
at `$401fa`, loaded from side track 79. It selects the bootstrap copper
list, installs its own VBL handler and reads the following contour-effect
block. Its wait compares a new counter at `$40cb6` with 150. This is a
separate minimum three-second interval, not the first effect's counter.
It explains the extra interval before the retained-mask trails.

After the seven-second textured passage, the side-track 109 program enters
at `$501fa`. It is **another animated composition**, not a blank hold.
It combines original monochrome tiles into three display planes with a
moving pointer program and alternating palettes. Its VBL counter at
`$50cd2` must reach 143 before it enters the next spirograph at `$45012`.
That minimum interval is 2.86 seconds. The native tile composition is now
[reconstructed and independently verified](TILES.md); it remains to be
connected to the complete director.

The native preview's main-module starting offsets now include these source
waits. They are local validation offsets, not proof that all disk reads and
decompression handoffs match the recording. The final director must retain
the actual composition of every loading passage and verify absolute music
alignment through the whole production.
