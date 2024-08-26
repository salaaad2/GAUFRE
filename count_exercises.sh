#!/bin/bash
# count exercises
grep "exercise_name" log.json | sort | uniq | wc -l
