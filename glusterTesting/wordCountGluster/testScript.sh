#!/bin/bash
for i in `seq 1 24`;
do
./wordCountGluster 1 $i 14 0
done
