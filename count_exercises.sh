#!/bin/zsh

# should be 30 as of aug 25
grep "exercise_name" log.json | sort | uniq | wc -l
