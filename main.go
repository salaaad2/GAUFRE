package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
    "flag"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

    show_bar := flag.Bool("bar_chart", false, "show a bar chart")
    show_plot := flag.Bool("line_plot", false, "show a line plot")
    flag.Parse()
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

    var exercise_to_inspect string 
    if len(os.Args) == 2 {
        exercise_to_inspect = os.Args[1]
    } else {
        exercise_to_inspect = os.Args[2]
    }
    var inspection []SetDetails

    for _, session := range workout_data.Workout.Sessions {
        fmt.Printf("%s\n", session.Date.String())
        for _, sets := range session.Sets {
            for _, set := range sets {
                if len(exercise_to_inspect) > 0 && set.ExerciseName == exercise_to_inspect {
                    inspection = append(inspection, set)
                }
                fmt.Printf("name: %s, reps: %s, weight: %f, dropset: %t\n", set.ExerciseName, set.Reps, set.Weight, set.DropSet)
            }
        }
    }

    if len(exercise_to_inspect) == 0 {
        return ;
    }
    fmt.Printf("Inspecting: %s\n", exercise_to_inspect)
    for i := 0; i < len(inspection); i++{
        if i == 0 {
            fmt.Printf("First session :")
        } else if i == len(inspection) - 1{ 
            fmt.Printf("Latest session:")
        }
        fmt.Printf("%f, %s\n", inspection[i].Weight, inspection[i].Reps)
    }
    
    if err := ui.Init(); err != nil {
        log.Fatal("failed to Init() termui.", err)
    }
    defer ui.Close()

    if *show_plot {
        line_plot := lineFromSetDetails(inspection)
        ui.Render(line_plot)
    } else if *show_bar {
        bar_chart := barFromSetDetails(inspection)
        ui.Render(bar_chart)
    }
    ui_events := ui.PollEvents()
    for {
        e := <-ui_events
		switch e.ID {
		case "q", "<C-c>":
			return
		}
    }
}

func barFromSetDetails(set_details []SetDetails) *widgets.BarChart {
   	bc := widgets.NewBarChart()
	bc.Title = "Bar Chart"
	bc.Labels = []string{}
    t_width, t_height := ui.TerminalDimensions()
    bc.SetRect(0, 0, t_width, t_height)
    bc.BarWidth = 3
    for i := 0 ; i < len(set_details); i++ {
        bc.Labels = append(bc.Labels, set_details[i].Reps)
        bc.Data = append(bc.Data, set_details[i].Weight)
    }
    bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorYellow}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
    return bc
}

func lineFromSetDetails(set_details []SetDetails) *widgets.Plot {
   	lc := widgets.NewPlot()
	lc.Title = "Default plot"
    t_width, t_height := ui.TerminalDimensions()
    lc.SetRect(0, 0, t_width, t_height)
    lc.Data = make([][]float64, 2)
    lc.Data[0] = make([]float64, 0)

   for i := 0 ; i < len(set_details); i++ {
        lc.Data[0] = append(lc.Data[0], set_details[i].Weight)
    }
   	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorGreen
	lc.Marker = widgets.MarkerDot
    lc.DotMarkerRune = '+'
    return lc
}
