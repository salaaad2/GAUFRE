package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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

    //show_bar := flag.Bool("bar_chart", false, "show a bar chart")
    //show_plot := flag.Bool("line_plot", false, "show a line plot")
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
    var inspection map[CustomDate][]SetDetails
    inspection = make(map[CustomDate][]SetDetails)

    for _, session := range workout_data.Workout.Sessions {
        fmt.Printf("%s\n", session.Date.String())
        for _, sets := range session.Sets {
            for _, set := range sets {
                if len(exercise_to_inspect) > 0 && set.ExerciseName == exercise_to_inspect {
                    inspection[session.Date] = append(inspection[session.Date], set)
                }
                fmt.Printf("name: %s, reps: %s, weight: %f, dropset: %t\n",
                    set.ExerciseName, set.Reps, set.Weight, set.DropSet)
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
        // fmt.Printf("%f, %s\n", inspection[i].Weight, inspection[i].Reps)
    }
    lineFromInspectionMap(inspection)
}

func lineFromInspectionMap(inspection map[CustomDate][]SetDetails) *charts.Line {
    line := charts.NewLine()

    line.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "Line chart for chest press",
        }),
   		charts.WithXAxisOpts(opts.XAxis{
			Name: "Date",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "weight(kilos)",
		}),
    )

	var keys []CustomDate
	for key := range inspection {
		keys = append(keys, key)
	}

	// Sort the map preemptively
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j].Time)
	})
	sortedMap := make(map[CustomDate][]SetDetails)
    var xAxisData []string
    var yAxisData []opts.LineData
    var repsData []opts.LineData
	for _, key := range keys {
		sortedMap[key] = inspection[key]
        xAxisData = append(xAxisData, key.String())
        for _, v := range sortedMap[key] {
            yAxisData = append(yAxisData, opts.LineData{Value: v.Weight})
            
            repsData = append(repsData, opts.LineData{Value: getTotalReps(v.Reps)})
        }
	}

    line.AddSeries("Chest weight", yAxisData).SetXAxis(xAxisData).SetSeriesOptions(
			charts.WithMarkPointNameTypeItemOpts(
				opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
				opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
				opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
			),
			charts.WithMarkPointStyleOpts(
				opts.MarkPointStyle{Label: &opts.Label{Show: opts.Bool(true)}}),
    )

    line.AddSeries("Chest reps", repsData).SetXAxis(xAxisData)

    // Render the chart to an HTML file
	f, err := os.Create("line_chart.html")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	if err := line.Render(f); err != nil {
		log.Fatalf("Failed to render chart: %v", err)
	}

	log.Println("Line chart rendered to line_chart.html")
    return line
}

func getTotalReps(reps_str string) int {
    reps_nowhitespace := strings.ReplaceAll(reps_str, " ", "")
    minus_index := strings.Index(reps_nowhitespace, "-")
    if minus_index != -1 {
        reps_nowhitespace = reps_nowhitespace[:minus_index]
    }
    split := strings.Split(reps_nowhitespace, "*")
    if len(split) == 2 {
        res1,_ := strconv.Atoi(split[0])
        res2,_ := strconv.Atoi(split[1])
        return res1 * res2
    }
    return 3*8
}

// func barFromSetDetails(set_details []SetDetails) *widgets.BarChart {
//    	bc := widgets.NewBarChart()
// 	bc.Title = "Bar Chart"
// 	bc.Labels = []string{}
//     t_width, t_height := ui.TerminalDimensions()
//     bc.SetRect(0, 0, t_width, t_height)
//     bc.BarWidth = 5
//     for i := 0 ; i < len(set_details); i++ {
//         bc.Labels = append(bc.Labels, set_details[i].Reps)
//         bc.Data = append(bc.Data, set_details[i].Weight)
//     }
//     bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorYellow}
// 	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
// 	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorWhite)}
//     return bc
// }
// 
// func lineFromSetDetails(set_details []SetDetails) *widgets.Plot {
//    	lc := widgets.NewPlot()
// 	lc.Title = "Line"
//     t_width, t_height := ui.TerminalDimensions()
//     lc.SetRect(0, 0, t_width, t_height)
//     lc.Data = make([][]float64, 2)
//     lc.Data[0] = make([]float64, 0)
//     lc.DataLabels = append(lc.DataLabels, "test")
//     lc.DataLabels = append(lc.DataLabels, "test1")
//    for i := 0 ; i < len(set_details); i++ {
//         lc.Data[0] = append(lc.Data[0], set_details[i].Weight)
//     }
//    	lc.AxesColor = ui.ColorWhite
// 	lc.LineColors[0] = ui.ColorGreen
// 	lc.Marker = widgets.MarkerDot
//     lc.DotMarkerRune = '+'
//     return lc
// }
