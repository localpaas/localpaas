package basedto

import (
	"fmt"
	"math"
	"time"

	vld "github.com/tiendc/go-validator"
	vldbase "github.com/tiendc/go-validator/base"
	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/pkg/timeutil"
)

func ValidateID(id *string, required bool, field string) (res []vld.Validator) {
	if required {
		res = append(res, vld.Required(id).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_ID_REQUIRED"), // use default error key
		))
	}
	if id != nil && *id != "" {
		res = append(res, vld.StrIsULID(id).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_ID_INVALID"), // use default error key
		))
	}
	return res
}

func ValidateIDSlice(ids []string, unique bool, minLen int, field string) (result []vld.Validator) {
	if unique {
		result = append(result, vld.SliceUnique(ids).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_NON_UNIQUE"), // use default error key
		))
	}
	if minLen > 0 {
		result = append(result, vld.SliceLen(ids, minLen, math.MaxInt).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_REQUIRED"), // use default error key
		))
	}
	result = append(result,
		vld.Slice(ids).ForEach(func(element string, index int, elemValidator vld.ItemValidator) {
			elemValidator.Validate(
				vld.StrIsULID(&element).OnError(
					vld.SetField(fmt.Sprintf("%s[%d]", field, index), nil),
					vld.SetCustomKey("ERR_VLD_OBJECT_ID_INVALID"), // use default error key
				),
			)
		}),
	)
	return result
}

func ValidateObjectIDReq(objID *ObjectIDReq, required bool, field string) (res []vld.Validator) {
	var id *string
	if objID != nil {
		id = &objID.ID
	}
	if required {
		res = append(res, vld.Required(id).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_ID_REQUIRED"), // use default error key
		))
	}
	if id != nil && *id != "" {
		res = append(res, vld.StrIsULID(id).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_ID_INVALID"), // use default error key
		))
	}
	return res
}

func ValidateObjectIDSliceReq(ids []*ObjectIDReq, unique bool, minLen int, field string) (result []vld.Validator) {
	if unique {
		result = append(result, vld.SliceUniqueBy(ids, func(item *ObjectIDReq) string {
			return item.ID
		}).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_NON_UNIQUE"), // use default error key
		))
	}
	if minLen > 0 {
		result = append(result, vld.SliceLen(ids, minLen, math.MaxInt).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_REQUIRED"), // use default error key
		))
	}
	result = append(result,
		vld.Slice(ids).ForEach(func(element *ObjectIDReq, index int, elemValidator vld.ItemValidator) {
			elemValidator.Validate(
				vld.StrIsULID(&element.ID).OnError(
					vld.SetField(fmt.Sprintf("%s[%d].id", field, index), nil),
					vld.SetCustomKey("ERR_VLD_OBJECT_ID_INVALID"), // use default error key
				),
			)
		}),
	)
	return result
}

func ValidateRequiredField(field any, fieldName string) (res []vld.Validator) {
	res = append(res, vld.Required(field).OnError(
		vld.SetField(fieldName, nil),
		vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
	))
	return res
}

func ValidateObjectAccessReq(access *ObjectAccessReq, required bool, field string) (res []vld.Validator) {
	var id *string
	if access != nil {
		id = &access.ID
	}
	if required {
		res = append(res, vld.Required(id).OnError(
			vld.SetField(field+".id", nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_ID_REQUIRED"), // use default error key
		))
	}
	if id != nil && *id != "" {
		res = append(res, vld.StrIsULID(id).OnError(
			vld.SetField(field+".id", nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_ID_INVALID"), // use default error key
		))
	}
	return res
}

func ValidateObjectAccessSliceReq(access ObjectAccessSliceReq, unique bool, minLen int, field string) (
	res []vld.Validator) {
	if unique {
		res = append(res, vld.SliceUniqueBy(access, func(item *ObjectAccessReq) string {
			return item.ID
		}).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_NON_UNIQUE"),
		))
	}
	if minLen > 0 {
		res = append(res, vld.SliceLen(access, minLen, math.MaxInt).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_REQUIRED"),
		))
	}
	res = append(res,
		vld.Slice(access).ForEach(func(item *ObjectAccessReq, index int, itemValidator vld.ItemValidator) {
			itemValidator.Validate(
				vld.StrIsULID(&item.ID).OnError(
					vld.SetField(fmt.Sprintf("%s[%d].id", field, index), nil),
					vld.SetCustomKey("ERR_VLD_OBJECT_ID_INVALID"),
				),
			)
		}),
	)
	return res
}

func ValidateModuleAccessReq(access *ModuleAccessReq, required bool, field string) (res []vld.Validator) {
	var id *string
	if access != nil {
		id = &access.ID
	}
	if required {
		res = append(res, vld.Required(id).OnError(
			vld.SetField(field+".id", nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_ID_REQUIRED"), // use default error key
		))
	}
	return res
}

func ValidateModuleAccessSliceReq(access ModuleAccessSliceReq, unique bool, minLen int, field string) (
	res []vld.Validator) {
	if unique {
		res = append(res, vld.SliceUniqueBy(access, func(item *ModuleAccessReq) string {
			return item.ID
		}).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_NON_UNIQUE"),
		))
	}
	if minLen > 0 {
		res = append(res, vld.SliceLen(access, minLen, math.MaxInt).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_OBJECT_IDS_REQUIRED"),
		))
	}
	return res
}

func ValidateEmail(email *string, required bool, field string) (res []vld.Validator) {
	if required {
		res = append(res, vld.Required(email).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_EMAIL_ADDR_REQUIRED"), // use default error key
		))
	}
	if email != nil && *email != "" {
		res = append(res, vld.StrIsEmail(email).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_EMAIL_ADDR_INVALID"), // use default error key
		))
	}
	return res
}

func ValidateEmailSlice(emails []string, unique bool, minLen int, field string) (result []vld.Validator) {
	if unique {
		result = append(result, vld.SliceUnique(emails).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_EMAIL_ADDRS_NON_UNIQUE"), // use default error key
		))
	}
	if minLen > 0 {
		result = append(result, vld.SliceLen(emails, minLen, math.MaxInt).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_EMAIL_ADDRS_REQUIRED"), // use default error key
		))
	}
	result = append(result,
		vld.Slice(emails).ForEach(func(element string, index int, elemValidator vld.ItemValidator) {
			elemValidator.Validate(
				vld.StrIsEmail(&element).OnError(
					vld.SetField(fmt.Sprintf("%s[%d]", field, index), nil),
					vld.SetCustomKey("ERR_VLD_EMAIL_ADDR_INVALID"), // use default error key
				),
			)
		}),
	)
	return result
}

func ValidateSlice[T comparable, S ~[]T](s S, unique bool, minLen int, allowedValues []T, field string) (
	result []vld.Validator) {
	return ValidateSliceEx(s, unique, minLen, math.MaxInt, allowedValues, field)
}

func ValidateSliceEx[T comparable, S ~[]T](s S, unique bool, minLen, maxLen int, allowedValues []T, field string) (
	result []vld.Validator) {
	if unique {
		result = append(result, vld.SliceUnique(s).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUES_NON_UNIQUE"),
		))
	}

	// Check minimum length
	if minLen > 0 && len(s) < minLen {
		result = append(result, vld.SliceLen(s, minLen, math.MaxInt).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
		))
	}

	// Check maximum length separately with specific error message
	if maxLen < math.MaxInt && len(s) > maxLen {
		result = append(result, vld.Must(len(s) <= maxLen).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_TOO_MANY"),
			vld.SetParam("Max", maxLen),
			vld.SetParam("Actual", len(s)),
		))
	}

	if len(allowedValues) > 0 {
		result = append(result,
			vld.SliceElemIn(s, allowedValues...).OnError(
				vld.SetField(field, nil),
				vld.SetCustomKey("ERR_VLD_VALUE_NOT_IN_LIST"),
			))
	}
	return result
}

func ValidateStr(s *string, required bool, minLen, maxLen int, field string) (result []vld.Validator) {
	if required {
		result = append(result, vld.Required(s).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
		))
	}
	if s != nil && *s != "" {
		result = append(result, vld.StrLen(s, minLen, maxLen).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_FIELD_LENGTH_INVALID"),
		))
	}
	return result
}

func ValidateStrIn[T ~string](s *T, required bool, allowedValues []T, field string) (
	result []vld.Validator) {
	return ValidateStrInEx(s, required, allowedValues, field, "")
}

func ValidateStrInEx[T ~string](s *T, required bool, allowedValues []T, field string, allowedValuesKey string) (
	result []vld.Validator) {
	if required {
		result = append(result, vld.Required(s).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
		))
	}
	if s != nil && len(*s) > 0 && len(allowedValues) > 0 {
		customKey := "ERR_VLD_VALUE_NOT_IN_LIST"
		if allowedValuesKey != "" {
			customKey = allowedValuesKey
		}

		result = append(result,
			vld.StrIn(s, allowedValues...).OnError(
				vld.SetField(field, nil),
				vld.SetCustomKey(customKey),
			))
	}
	return result
}

func ValidateTime(t *time.Time, required bool, from, to time.Time, field string) (result []vld.Validator) {
	if required {
		result = append(result, vld.Required(t).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
		))
	}
	if t != nil && !t.IsZero() {
		dt := *t
		if !from.IsZero() {
			result = append(result, vld.TimeGTE(dt, from).OnError(
				vld.SetField(field, nil),
				vld.SetCustomKey("ERR_VLD_VALUE_MUST_GREATER_THAN"),
			))
		}
		if !to.IsZero() {
			result = append(result, vld.TimeLT(dt, to).OnError(
				vld.SetField(field, nil),
				vld.SetCustomKey("ERR_VLD_VALUE_MUST_LESS_THAN"),
			))
		}
	}
	return result
}

func ValidateDate(date *timeutil.Date, required bool, from, to timeutil.Date, field string) (result []vld.Validator) {
	var t *time.Time
	if date != nil {
		t = gofn.ToPtr(date.ToTime())
	}
	return ValidateTime(t, required, from.ToTime(), to.ToTime(), field)
}

func ValidateObjectSliceBy[T, A comparable, S ~[]T](s S, required, unique bool, allowedValues []A,
	mapBy func(T) A, field string) (result []vld.Validator) {
	if required && len(s) == 0 {
		result = append(result, vld.Required(s).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
		))
	}

	values := gofn.MapSlice(s, func(s T) A { return mapBy(s) })
	if unique {
		result = append(result, vld.SliceUnique(values).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUES_NON_UNIQUE"),
		))
	}
	if len(s) > 0 && len(allowedValues) > 0 {
		result = append(result,
			vld.SliceElemIn(values, allowedValues...).OnError(
				vld.SetField(field, nil),
				vld.SetCustomKey("ERR_VLD_VALUE_NOT_IN_LIST"),
			))
	}
	return result
}

func ValidateNumber[T int | uint](v *T, required bool, min, max T, field string) (result []vld.Validator) {
	if required {
		result = append(result, vld.Must(v != nil).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_REQUIRED"),
		))
	}
	if v != nil {
		result = append(result, vld.NumRange(v, min, max).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_NOT_IN_RANGE"),
		))
	}
	return result
}

func ValidateIntSliceRange[T vldbase.Number](s []T, min, max T, field string) (result []vld.Validator) {
	for _, value := range s {
		v := value
		result = append(result, vld.NumRange(&v, min, max).OnError(
			vld.SetField(field, nil),
			vld.SetCustomKey("ERR_VLD_VALUE_NOT_IN_RANGE"),
		))
	}
	return result
}

func ValidateMutualExclusiveFields(shouldValidate bool, fields ...string) (result []vld.Validator) {
	for _, field := range fields {
		result = append(result,
			vld.Must(!shouldValidate).OnError(
				vld.SetField(field, nil),
				vld.SetCustomKey("ERR_VLD_MUTUALLY_EXCLUSIVE_FIELDS_INVALID"),
			),
		)
	}

	return result
}
