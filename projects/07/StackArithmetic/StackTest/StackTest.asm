@17
D=A
@SP
AM=M+1
A=A-1
M=D
@17
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JEQ_TRUE_1
D;JEQ
@SP
A=M-1
M=0
(JEQ_TRUE_1)
@17
D=A
@SP
AM=M+1
A=A-1
M=D
@16
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JEQ_TRUE_2
D;JEQ
@SP
A=M-1
M=0
(JEQ_TRUE_2)
@16
D=A
@SP
AM=M+1
A=A-1
M=D
@17
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JEQ_TRUE_3
D;JEQ
@SP
A=M-1
M=0
(JEQ_TRUE_3)
@892
D=A
@SP
AM=M+1
A=A-1
M=D
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JLT_TRUE_4
D;JLT
@SP
A=M-1
M=0
(JLT_TRUE_4)
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@892
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JLT_TRUE_5
D;JLT
@SP
A=M-1
M=0
(JLT_TRUE_5)
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@891
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JLT_TRUE_6
D;JLT
@SP
A=M-1
M=0
(JLT_TRUE_6)
@32767
D=A
@SP
AM=M+1
A=A-1
M=D
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JGT_TRUE_7
D;JGT
@SP
A=M-1
M=0
(JGT_TRUE_7)
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@32767
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JGT_TRUE_8
D;JGT
@SP
A=M-1
M=0
(JGT_TRUE_8)
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@32766
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@JGT_TRUE_9
D;JGT
@SP
A=M-1
M=0
(JGT_TRUE_9)
@57
D=A
@SP
AM=M+1
A=A-1
M=D
@31
D=A
@SP
AM=M+1
A=A-1
M=D
@53
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=D+M
@112
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=M-D
@SP
A=M-1
M=-M
@SP
AM=M-1
D=M
A=A-1
M=D&M
@82
D=A
@SP
AM=M+1
A=A-1
M=D
@SP
AM=M-1
D=M
A=A-1
M=D|M
@SP
A=M-1
M=!M
