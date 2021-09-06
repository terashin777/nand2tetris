// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.
@color
M=0

(WATCH)
  @KBD
  D=M
  @PRESS
  D;JGT
  @color
  D=M+1
  @EMPTY
  D;JEQ
@WATCH
0;JMP
(PRESS)
  @color
  D=M
  @BLACK
  D;JEQ
@WATCH
0;JMP
(EMPTY)
  @color
  M=0
@FILL
0;JMP
(BLACK)
  @color
  M=-1
(FILL)
  @SCREEN
  D=A
  @i
  M=D
  (FILL_LOOP)
    @24575
    D=A
    @i
    D=M-D
    @WATCH
    D;JGT
    @color
    D=M
    @i
    A=M
    M=D
    @i
    M=M+1
  @FILL_LOOP
  0;JMP