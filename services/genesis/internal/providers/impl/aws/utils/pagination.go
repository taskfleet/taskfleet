package awsutils

import "fmt"

// RunPaginatedRequest runs the provided request multiple times, extracting the results and the
// "next token" from each successfull API call.
func RunPaginatedRequest[T any, R any](
	request func(*string) (T, error), getOutputs func(T) ([]R, *string),
) ([]R, error) {
	result := make([]R, 0)
	var nextToken *string
	var outputs []R
	for {
		response, err := request(nextToken)
		if err != nil {
			return nil, fmt.Errorf("paginated request failed: %s", err)
		}
		outputs, nextToken = getOutputs(response)
		result = append(result, outputs...)
		if nextToken == nil {
			return result, nil
		}
	}
}
