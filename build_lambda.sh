#!/bin/bash
set -e # panic if can't execute instruction

rm -rf build
mkdir build
cp -r tweeter build/tmp
cd build/tmp
pip install tweepy -t .
zip -r ../awsFunc.zip *
rm -rf build/tmp
