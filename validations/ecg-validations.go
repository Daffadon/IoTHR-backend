package validations

type ResampleECGDataInput struct {
	ECG_Plot []float64 `json:"ecg_plot" bson:"ecg_plot"`
}
