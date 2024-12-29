package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math"
	"os"
)

type consts struct {
	eCharge          float64 // Заряд электрона
	eMass            float64 // Масса электрона
	startSpeed       float64 // Начальная скорость
	length           float64 // Длина конденсатора
	innerRadius      float64 // Радиус внутренней стенки
	outerRadius      float64 // Радиус внешней стенки
	startPosition    float64 // Стартовая позиция по y
	finishedPosition float64 // Позиция по y, при которой элекрон не пролетит конденсатор
	dt               float64 // Шаг по времени для вычисления промежуточных состояний
}

var num consts

// CountAccelerationAtDistanceFromAxis Функция расчёта ускорения электрона на расстоянии R от оси цилиндров
func CountAccelerationAtDistanceFromAxis(u, r float64) float64 {
	a := (num.eCharge * u) / (num.eMass * r * math.Log(num.outerRadius/num.innerRadius))
	return a
}

// Simulation Функция, которая моделирует движение электрона и возвращает конечную позицию, конечную скорость и время
func Simulation(u float64) float64 {
	a := CountAccelerationAtDistanceFromAxis(u, num.startPosition+num.innerRadius)
	speedOnYAxis := 0.0
	time := 0.0
	pos := num.startPosition

	for ; pos > 0 && (time*num.startSpeed) <= num.length; time += num.dt {
		pos += -speedOnYAxis*num.dt - a*math.Pow(num.dt, 2)/2
		if pos > 0 {
			speedOnYAxis += a * num.dt
		} else {
			speedOnYAxis = 0
		}

		a = CountAccelerationAtDistanceFromAxis(u, pos+num.innerRadius)
	}

	return pos
}

// FindMinimumPotentialDifference Функция для поиска минимального значения разности потенциалов
// при котором электрон не вылетает из конденсатора
func FindMinimumPotentialDifference() float64 {
	leftBoard := 0.0
	rightBoard := 1000.0
	for (rightBoard - leftBoard) > 0.00001 {
		mid := (leftBoard + rightBoard) / 2
		pos := Simulation(mid)
		if pos <= 0 {
			rightBoard = mid
		} else {
			leftBoard = mid
		}
	}
	return rightBoard
}

func SimulationForGraphics(u float64) ([]float64, []float64,
	[]opts.LineData, []opts.LineData, []opts.LineData) {

	var xVals []float64
	var tVals []float64
	var yVals []opts.LineData
	var aVals []opts.LineData
	var vVals []opts.LineData

	a := CountAccelerationAtDistanceFromAxis(u, num.startPosition+num.innerRadius)
	speedOnYAxis := 0.0
	time := 0.0
	pos := num.startPosition

	for ; pos > 0 && (time*num.startSpeed) <= num.length; time += num.dt {

		xVals = append(xVals, speedOnYAxis*time)
		tVals = append(tVals, time)
		yVals = append(yVals, opts.LineData{Value: pos})
		aVals = append(aVals, opts.LineData{Value: a})
		vVals = append(vVals, opts.LineData{Value: speedOnYAxis})

		pos += -speedOnYAxis*num.dt - a*math.Pow(num.dt, 2)/2
		if pos > 0 {
			speedOnYAxis += a * num.dt
		} else {
			speedOnYAxis = 0
		}

		a = CountAccelerationAtDistanceFromAxis(u, pos+num.innerRadius)
	}

	return xVals, tVals, yVals, aVals, vVals
}

func main() {
	num.eCharge = 1.6 * math.Pow(10, -19)
	num.eMass = 9.1 * math.Pow(10, -31)

	var r1, r2, v, l float64
	fmt.Scanln(&r1, &r2, &v, &l)
	num.startSpeed = v
	num.length = l
	num.innerRadius = r1
	num.outerRadius = r2
	num.startPosition = (num.outerRadius - num.innerRadius) / 2
	num.finishedPosition = 0
	num.dt = num.length / (num.startSpeed * 1000)

	fmt.Println("Минимальная разность потенциалов:", FindMinimumPotentialDifference())

	var potDiff float64
	fmt.Scanln(&potDiff)

	xVals, tVals, yVals, aVals, vVals := SimulationForGraphics(potDiff)

	// Создание графика зависимости y(x)
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "График функции y(x)"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "x, м"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "y, м"}),
	)

	line.SetXAxis(tVals).
		AddSeries("y(x)", yVals).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{}))

	f, _ := os.Create("graph_y_x.html")
	defer f.Close()
	line.Render(f)

	// Создание графика зависимости y(t)
	line = charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "График функции y(t)"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "t, с"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "y, м"}),
	)

	line.SetXAxis(xVals).
		AddSeries("y(t)", yVals).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{}))

	f, _ = os.Create("graph_y_t.html")
	defer f.Close()
	line.Render(f)

	// Создание графика зависимости V_y(t)
	line = charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "График функции проекции скорости Vy(t)"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "t, с"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Vy, м/c"}),
	)

	line.SetXAxis(tVals).
		AddSeries("Vy(t)", vVals).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{}))

	f, _ = os.Create("graph_vy_t.html")
	defer f.Close()
	line.Render(f)

	// Создание графика зависимости a(t)
	line = charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "График функции ускорения a(t)"}),
		charts.WithXAxisOpts(opts.XAxis{Name: "t, с"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "a, м/с²"}),
	)

	line.SetXAxis(tVals).
		AddSeries("a(t)", aVals).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{}))

	f, _ = os.Create("graph_a_t.html")
	defer f.Close()
	line.Render(f)
}
