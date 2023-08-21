package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationErrors map[string]string

type User struct {
	Name     string `validate:"required"`
	Age      int    `validate:"required,min=18,max=99"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,minLen=8,maxLen=20"`
}

func validateStruct(obj interface{}) ValidationErrors {
	errs := make(ValidationErrors)

	val := reflect.ValueOf(obj)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag

		// Check for "required" tag
		if strings.Contains(tag.Get("validate"), "required") {
			if isEmpty(field) {
				errs[typ.Field(i).Name] = "Field is required"
			}
		}

		// Check for "max" and "min" tags
		if strings.Contains(tag.Get("validate"), "max") || strings.Contains(tag.Get("validate"), "min") {
			if field.Kind() == reflect.Int {
				minStr := tag.Get("min")
				maxStr := tag.Get("max")
				val, _ := field.Interface().(int)

				if minStr != "" {
					min, _ := strconv.Atoi(minStr)
					if val < min {
						errs[typ.Field(i).Name] = fmt.Sprintf("Value should be greater than or equal to %s", minStr)
					}
				}

				if maxStr != "" {
					max, _ := strconv.Atoi(maxStr)
					if val > max {
						errs[typ.Field(i).Name] = fmt.Sprintf("Value should be less than or equal to %s", maxStr)
					}
				}
			}
		}

		// Check for "maxLen" and "minLen" tags
		if strings.Contains(tag.Get("validate"), "maxLen") || strings.Contains(tag.Get("validate"), "minLen") {
			if field.Kind() == reflect.String {
				minLenStr := tag.Get("minLen")
				maxLenStr := tag.Get("maxLen")
				val, _ := field.Interface().(string)
				length := len(val)

				if minLenStr != "" {
					minLen, _ := strconv.Atoi(minLenStr)
					if length < minLen {
						errs[typ.Field(i).Name] = fmt.Sprintf("Length should be at least %s characters", minLenStr)
					}
				}

				if maxLenStr != "" {
					maxLen, _ := strconv.Atoi(maxLenStr)
					if length > maxLen {
						errs[typ.Field(i).Name] = fmt.Sprintf("Length should be at most %s characters", maxLenStr)
					}
				}
			}
		}

		// Check for "email" tag
		if strings.Contains(tag.Get("validate"), "email") {
			if field.Kind() == reflect.String {
				email := field.Interface().(string)
				emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
				match, _ := regexp.MatchString(emailPattern, email)
				if !match {
					errs[typ.Field(i).Name] = "Invalid email format"
				}
			}
		}
	}

	return errs
}

func isEmpty(field reflect.Value) bool {
	zero := reflect.Zero(field.Type()).Interface()
	return reflect.DeepEqual(field.Interface(), zero)
}

func main() {
	user := User{Name: "Handy", Age: 24, Email: "handy@gmail.com", Password: "123456789"}
	errors := validateStruct(user)

	if len(errors) > 0 {
		fmt.Println("Validation errors:")
		for field, errMsg := range errors {
			fmt.Printf("%s: %s\n", field, errMsg)
		}
	} else {
		fmt.Println("Validation successful")
	}
}
