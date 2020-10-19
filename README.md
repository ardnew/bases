# bases
#### Evaluate and print expressions in various bases

## Usage

`bases` is a simple utility to evaluate arithmetic and bitwise expressions and print the result in the most commonly used number bases: `BIN` (2), `OCT` (8), `DEC` (10), and `HEX` (16).

Invoked without any arguments, `bases` enters a run-eval-print-loop (REPL) that repeatedly prompts for user input and prints its evaluation, like a simple calculator (recognizing all valid Perl expressions):

```
$ bases
	LN:    HEX          OCT             BIN                                  DEC
> 1 + 0b10 + 03 + 0x4
	001:   0x0000000a   0o00000000012   0b00000000000000000000000000001010   10
> (5 << 3) | (~5 & 0xF)
	002:   0x0000002a   0o00000000052   0b00000000000000000000000000101010   42
> 0775 & ~0022
	003:   0x000001ed   0o00000000755   0b00000000000000000000000111101101   493
> ord('A')
	004:   0x00000041   0o00000000101   0b00000000000000000000000001000001   65
```

Alternatively, input can be provided via stdin, and the header line can be suppressed with the "quiet" command-line flag (`-q`):

```
$ echo 20*2+2 | bases -q
	001:   0x0000002a   0o00000000052   0b00000000000000000000000000101010   42

$ echo $'0xAA\n0xBB\n0xCC' | bases -q
	001:   0x000000aa   0o00000000252   0b00000000000000000000000010101010   170
	002:   0x000000bb   0o00000000273   0b00000000000000000000000010111011   187
	003:   0x000000cc   0o00000000314   0b00000000000000000000000011001100   204
```

By default, the value is output as a 32-bit unsigned integer. If the "byte" command-line flag (`-b`) is added, a second output section is also printed containing the significant 8-bit fields:

```
$ bases -b
> 0xDEADBEEF
[32-bit]
	LN:    HEX          OCT             BIN                                  DEC
	001:   0xdeadbeef   0o33653337357   0b11011110101011011011111011101111   3735928559

[8-bit]
	OFF:   HEX    OCT      BIN          DEC
	 +0:   0xef   0o0357   0b11101111   239
	 +8:   0xbe   0o0276   0b10111110   190
	+16:   0xad   0o0255   0b10101101   173
	+24:   0xde   0o0336   0b11011110   222
```

## Installation

Clone this repo or copy the `bases` Perl script to any directory in your `PATH` environment variable, and enable execute permissions. For example:

```
$ git clone git@github.com:ardnew/bases.git
 Cloning into 'bases'...
 remote: Enumerating objects: 4, done.
 remote: Counting objects: 100% (4/4), done.
 remote: Compressing objects: 100% (4/4), done.
 remote: Total 4 (delta 0), reused 0 (delta 0), pack-reused 0
 Receiving objects: 100% (4/4), done.

$ sudo cp bases/bases /usr/local/bin

$ sudo chmod +x /usr/local/bin/bases

```

Also ensure the path to the Perl executable at the top of the script (e.g., `#!/usr/bin/perl`) is correct.

