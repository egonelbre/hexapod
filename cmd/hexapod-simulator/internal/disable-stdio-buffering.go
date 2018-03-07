package internal

// #include <stdio.h>
// #include <errno.h>
import "C"

func init() { C.setvbuf(C.stdout, nil, C._IONBF, 0) }
