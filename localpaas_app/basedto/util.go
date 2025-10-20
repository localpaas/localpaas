package basedto

import "github.com/localpaas/localpaas/localpaas_app/apperrors"

// TransformObjectSlice transforms a slice of objects
func TransformObjectSlice[T, U any, TS ~[]T](objects TS, singleTransformFn func(T) (U, error)) ([]U, error) {
	resp := make([]U, 0, len(objects))
	for _, obj := range objects {
		itemResp, err := singleTransformFn(obj)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, itemResp)
	}
	return resp, nil
}

// TransformObjectSlice1P transforms a slice of objects with one additional parameter
func TransformObjectSlice1P[T, U, P1 any, TS ~[]T](singleTransformFn func(T, P1) (U, error), objects TS,
	p1 P1) ([]U, error) {
	resp := make([]U, 0, len(objects))
	for _, obj := range objects {
		itemResp, err := singleTransformFn(obj, p1)
		if err != nil {
			return nil, apperrors.Wrap(err)
		}
		resp = append(resp, itemResp)
	}
	return resp, nil
}
