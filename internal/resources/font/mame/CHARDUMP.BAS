10 REM Basic program to plot all chars from both charsets
20 SET MODE 80 : SET COLOUR 3 TO 0 : SET COLOUR 0 TO 15
30 X% := 0 : Y% := 200
40 FOR F% := 0 TO 1
50   FOR C% := 0 TO 255
60     PLOT CHR$(C%), X%, Y% FONT F%
70     X% := X% + 8
80     IF X% > 639 THEN X% := 0 : Y% := Y% - 10
90   NEXT C%
100 NEXT F%