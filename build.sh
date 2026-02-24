#!/bin/bash
v -o rn-resource-checker-v.exe -cflags "-march=native" -prod ./src
# v -o rn-resource-checker-v.exe -prod -cflags "-O3 -mavx2 -march=native -s -static" .