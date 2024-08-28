package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
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
	Date CustomDate              `json:"date"`
}

type Workout struct {
	Sessions []Session `json:"sessions"`
}

type WorkoutData struct {
	Workout Workout `json:"workout"`
}

type Science struct {
	MusclesWorked       []string
	Name                string
	CaloriesBurnedBySet uint32
}

func main() {
	month_flag := flag.String("month", "01", "The month you want to see data for")
	flag.Parse()

	single_month := false
	if month_flag != nil {
		single_month = true
	}

	log_path := "./log.json"
	log_contents, err := os.ReadFile(log_path)
	if err != nil {
		log.Fatal("could not read ./log.json. Are you sure it exists ?", err)
		return
	}

	var workout_data WorkoutData
	err = json.Unmarshal(log_contents, &workout_data)
	if err != nil {
		log.Fatal("json.Unmarshal() failed: ", err)
		return
	}

	var exercises_to_inspect []string
	if len(os.Args) == 2 {
		exercises_to_inspect = strings.Split(os.Args[1], ",")
	} else {
		log.Fatal("Missing exercise name to inspect.\nUse ./print_exercises.sh to see which can be used.")
	}
	charts := makeCharts(workout_data, exercises_to_inspect)

	page := components.NewPage()
	page.AddCharts(charts...)
	f, err := os.Create("line_chart.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

func makeCharts(workout_data WorkoutData, exercises_to_inspect []string) []components.Charter {
	var charts []components.Charter
	inspection := make(map[string]map[CustomDate][]SetDetails)
	for _, exercise_name := range exercises_to_inspect {
		inspection[exercise_name] = make(map[CustomDate][]SetDetails)
		for _, session := range workout_data.Workout.Sessions {
			for _, sets := range session.Sets {
				for _, set := range sets {
					if set.ExerciseName == exercise_name {
						inspection[exercise_name][session.Date] =
							append(inspection[exercise_name][session.Date], set)
					}
				}
			}
		}
		charts = append(
			charts,
			lineFromInspectionMap(inspection[exercise_name], exercise_name))
	}
	return charts
}

func lineFromInspectionMap(inspection map[CustomDate][]SetDetails, exercise_name string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Line chart for " + exercise_name,
			Subtitle: "total sets: " + strconv.Itoa(len(inspection)),
			Right:    "40%",
		}),
		charts.WithLegendOpts(opts.Legend{Left: "60%"}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Date",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "\nweight(kilos)",
		}),
	)

	// Sort the map preemptively
	var keys []CustomDate
	for key := range inspection {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j].Time)
	})
	var date_data_x []string
	var weight_data_y []opts.LineData
	var reps_data_y []opts.LineData
	var bar_data_y []opts.BarData
	for _, key := range keys {
		date_data_x = append(date_data_x, key.String())
		for _, v := range inspection[key] {
			weight_data_y = append(weight_data_y, opts.LineData{Value: v.Weight})
			total_weight := getTotalReps(v.Reps)
			reps_data_y = append(reps_data_y, opts.LineData{Value: total_weight})
			bar_data_y = append(bar_data_y, opts.BarData{Value: v.Weight * float64(total_weight)})
		}
	}

	// weight line plot
	line.AddSeries("weight: "+exercise_name, weight_data_y).SetXAxis(date_data_x).SetSeriesOptions(
		charts.WithMarkPointNameTypeItemOpts(
			opts.MarkPointNameTypeItem{Name: "Maximum", Type: "max"},
			opts.MarkPointNameTypeItem{Name: "Average", Type: "average"},
			opts.MarkPointNameTypeItem{Name: "Minimum", Type: "min"},
		),
		charts.WithMarkPointStyleOpts(
			opts.MarkPointStyle{Label: &opts.Label{Show: opts.Bool(true)}},
		),
		charts.WithLineChartOpts(
			opts.LineChart{
				Smooth: opts.Bool(true),
			},
		),
	)

	// reps line plot
	line.AddSeries("reps: "+exercise_name, reps_data_y).SetXAxis(date_data_x).SetSeriesOptions(
		charts.WithLineChartOpts(
			opts.LineChart{
				Smooth: opts.Bool(true),
			},
		),
	)

	// total weight bar chart (disabled for now)
	bar := charts.NewBar()
	bar.AddSeries("total_weight", bar_data_y).SetXAxis(date_data_x)
	var selected = map[string]bool{}
	selected["total_weight"] = false
	bar.Selected = selected
	// line.Overlap(bar)
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
		res1, _ := strconv.Atoi(split[0])
		res2, _ := strconv.Atoi(split[1])
		return res1 * res2
	}
	// default: 3 sets of 8 reps
	return 3 * 8
}
