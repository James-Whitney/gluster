#!/bin/bash
echo TESTING SIZE 1;
for i in `seq 1 24`;
do
./wordCountGluster 1 $i 14 0
done
echo TESTING SIZE 100;
for i in `seq 1 24`;
do
./wordCountGluster 100 $i 14 0
done