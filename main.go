package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type CustomDate struct {
	time.Time
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	date_str := string(data[1 : len(data)-1])

	parsed_time, err := time.Parse("02/01", date_str)
	if err != nil {
		return err
	}

	cd.Time = parsed_time
	return nil
}

type SetDetails struct {
	ExerciseName string  `json:"exercise_name"`
	Weight       float64 `json:"weight"`
	Reps         string  `json:"reps"`
	DropSet      bool    `json:"drop_set"`
}

type Session struct {
    Sets []map[string]SetDetails `json:"sets"`
    Date CustomDate `json:"date"`
}

type Workout struct {
    Sessions []Session `json:"sessions"`
}

type WorkoutData struct {
	Workout Workout `json:"workout"`    
}

type Science struct {
    MusclesWorked []string
    Name string
    CaloriesBurnedBySet uint32
}

func main() {
    log_path := "./log.json"
    log_contents, err := os.ReadFile(log_path)
    if err != nil {
        log.Fatal("could not read log file.", err)
    }

    var workout_data WorkoutData
    err = json.Unmarshal(log_contents, &workout_data)

    if err != nil {
        log.Fatal("json.Unmarshal() failed: ", err)
    }

    for _, session := range workout_data.Workout.Sessions {
        fmt.Printf("%s\n", session.Date.String())
        for _, sets := range session.Sets {
            for _, set := range sets {
                fmt.Printf("name: %s, dropset: %d, reps: %s,  weight: %f\n", set.ExerciseName, set.DropSet, set.Reps, set.Weight)
            }
        }
    }
}
    //science_path := "./science.json"
    //science_content, _ := os.ReadFile(science_path)

    //var science Science
    //err := json.Unmarshal(science_content, &science)
    //if err != nil {
    //    log.Fatal("error: could not read science content.", err)
    //}
