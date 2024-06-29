package proto_conversion

import (
	"agentske/preprocessing"
	messages "agentske/proto"

	"gonum.org/v1/gonum/mat"
)

func ConvertToProtoData(data preprocessing.Data) (*messages.Data, error) {
	protoData := &messages.Data{}

	// Assign Labels
	protoData.Labels = ConvertToInt32Slice(data.Labels)

	// Convert Histograms
	for _, data := range data.Input {
		protoInput := &messages.Input{}
		protoInput.Values = data

		protoData.Input = append(protoData.Input, protoInput)
	}

	return protoData, nil
}

func ConvertToInt32Slice(labels []float64) []float64 {
	result := make([]float64, len(labels))
	for i, val := range labels {
		result[i] = float64(val)
	}
	return result
}

func GetDataSetsFromProto(data *messages.Data) (*mat.Dense, *mat.Dense, error) {
	rows, cols := len(data.Input), len(data.Input[0].Values)
	trainingData := make([]float64, rows*cols)
	for i, input := range data.Input {
		copy(trainingData[i*cols:(i+1)*cols], input.Values)
	}
	X := mat.NewDense(rows, cols, trainingData)
	Y := mat.NewDense(rows, 1, data.Labels)
	return X, Y, nil
}
