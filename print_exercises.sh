#!/bin/bash

# prints all exercises found in log.json
grep "exercise_name" log.json | sort | uniq | cut -d ':' -f 2 | sed -r 's/(^\s*)|(")|(,)//g'
