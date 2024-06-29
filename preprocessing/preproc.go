package preprocessing

import (
	"encoding/csv"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Data struct {
	Labels []float64
	Input  [][]float64
}

func readData(fileLocation string) (Data, error) {
	// Open the CSV file
	file, err := os.Open(fileLocation)
	if err != nil {
		fmt.Println("Error:", err)
		return Data{}, nil
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return Data{}, nil
	}

	var data Data

	// Iterate through the rows
	for i, row := range records {
		if i == 0 { // Skip the header row
			continue
		}

		var row_data []float64

		// Iterate through the cells
		for j, cell := range row {
			switch j {
			case 1: // gender
				gender := cell
				//fmt.Println(gender)
				if gender == "Male" {
					row_data = append(row_data, 0)
				} else if gender == "Female" {
					row_data = append(row_data, 1)
				} else {
					continue // Drop rows with "Other"
				}
			case 2: // age
				age, _ := strconv.ParseFloat(cell, 64)
				//fmt.Println(age)
				row_data = append(row_data, age)
			case 3: // hypertension
				hypertension, _ := strconv.Atoi(cell)
				//fmt.Println(hypertension)
				row_data = append(row_data, float64(hypertension))
			case 4: // heart_disease
				heartDisease, _ := strconv.Atoi(cell)
				row_data = append(row_data, float64(heartDisease))
			case 5: // ever_married
				everMarried := cell
				//fmt.Println(everMarried)
				if everMarried == "No" {
					row_data = append(row_data, 0)
				} else {
					row_data = append(row_data, 1)
				}
			case 6: // work_type
				workType := cell
				//fmt.Println(workType)
				switch workType {
				case "children":
					row_data = append(row_data, 1, 0, 0, 0, 0)
				case "Govt_job":
					row_data = append(row_data, 0, 1, 0, 0, 0)
				case "Never_worked":
					row_data = append(row_data, 0, 0, 1, 0, 0)
				case "Private":
					row_data = append(row_data, 0, 0, 0, 1, 0)
				case "Self-employed":
					row_data = append(row_data, 0, 0, 0, 0, 1)
				}
			case 7: // Residence_type
				residenceType := cell
				//fmt.Println(residenceType)
				if residenceType == "Rural" {
					row_data = append(row_data, 0)
				} else {
					row_data = append(row_data, 1)
				}
			case 8: // avg_glucose_level
				avgGlucoseLevel, _ := strconv.ParseFloat(cell, 64)
				row_data = append(row_data, avgGlucoseLevel)
			case 9: // bmi
				bmi := cell
				//fmt.Println(bmi)
				if bmi == "N/A" {
					continue // Drop rows with "N/A"
				}
				bmiFloat, _ := strconv.ParseFloat(bmi, 64)
				row_data = append(row_data, bmiFloat)
			case 10: // smoking_status
				smokingStatus := cell
				//fmt.Println(smokingStatus)
				switch smokingStatus {
				case "formerly smoked":
					row_data = append(row_data, 1, 0, 0)
				case "never smoked":
					row_data = append(row_data, 0, 1, 0)
				case "smokes":
					row_data = append(row_data, 0, 0, 1)
				case "Unknown":
					row_data = append(row_data, 0, 0, 0)
				}
			case 11: // stroke
				stroke, _ := strconv.Atoi(cell)
				if stroke == 1 {
					// Oversample minority class (stroke == 1)
					for j := 0; j < 20; j++ {
						data.Input = append(data.Input, make([]float64, len(row_data)))
						copy(data.Input[len(data.Input)-1], row_data)
						data.Labels = append(data.Labels, float64(stroke))
					}
				} else {
					data.Input = append(data.Input, make([]float64, len(row_data)))
					copy(data.Input[len(data.Input)-1], row_data)
					data.Labels = append(data.Labels, float64(stroke))
				}
			}
		}
	}
	return data, nil
}

func deepCopy(input []float64) []float64 {
	dup := make([]float64, len(input))
	copy(dup, input)
	return dup
}

func oversample(data Data) Data {
	// Seed the random number generator to ensure different results on each run
	rand.Seed(time.Now().UnixNano())

	// Find the indices of minority class samples
	minorityIndices := make([]int, 0)
	for i, label := range data.Labels {
		if label == 1 {
			minorityIndices = append(minorityIndices, i)
		}
	}

	// Count the number of samples in each class
	majorityCount := len(data.Labels) - len(minorityIndices)
	minorityCount := len(minorityIndices)

	// Calculate the number of samples needed to balance the classes
	samplesNeeded := majorityCount - minorityCount

	// Perform random oversampling
	for i := 0; i < samplesNeeded; i++ {
		// Choose a random index from the minority class indices
		randomIndex := rand.Intn(len(minorityIndices))
		minorityIndex := minorityIndices[randomIndex]

		// Append the duplicate minority class sample with a deep copy
		data.Input = append(data.Input, deepCopy(data.Input[minorityIndex]))
		data.Labels = append(data.Labels, data.Labels[minorityIndex])
	}

	return data
}

func shuffleData(data Data, seed int) Data {
	rand.Seed(int64(seed))
	shuffledData := Data{
		Labels: make([]float64, len(data.Labels)),
		Input:  make([][]float64, len(data.Input)),
	}

	perm := rand.Perm(len(data.Input))
	for i, j := range perm {
		shuffledData.Labels[i] = data.Labels[j]
		shuffledData.Input[i] = data.Input[j]
	}

	return shuffledData
}

func splitData(data Data, splitRatio float64, seed int) (Data, Data) {
	// Shuffle the data first
	shuffledData := shuffleData(data, seed)

	numSamples := len(shuffledData.Input)
	numTrain := int(float64(numSamples) * splitRatio)

	trainData := Data{
		Labels: make([]float64, numTrain),
		Input:  make([][]float64, numTrain),
	}

	valData := Data{
		Labels: make([]float64, numSamples-numTrain),
		Input:  make([][]float64, numSamples-numTrain),
	}

	for i := 0; i < numTrain; i++ {
		trainData.Labels[i] = shuffledData.Labels[i]
		trainData.Input[i] = shuffledData.Input[i]
	}

	for i := numTrain; i < numSamples; i++ {
		valData.Labels[i-numTrain] = shuffledData.Labels[i]
		valData.Input[i-numTrain] = shuffledData.Input[i]
	}

	return trainData, valData
}

func PreprocessImagesForTraining() (Data, Data) {
	data, err := readData(".//..//data//healthcare-dataset-stroke-data.csv")
	if err != nil {
		fmt.Println("Error loading data:", err)
	}
	oversampledData := oversample(data)
	trainData, validationData := splitData(oversampledData, 0.8, 42)
	return trainData, validationData
}

func PreprocessImagesForEvaluation() Data {
	data, err := readData(".//..//data//test_data.csv")
	if err != nil {
		fmt.Println("Error loading data:", err)
	}

	return data
}
