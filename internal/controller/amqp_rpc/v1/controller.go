package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/leabago/share-radio/adder/pkg/jwt"
	"github.com/leabago/share-radio/adder/pkg/logger"
)

// V1 -.
type V1 struct {
	j *jwt.Manager
	l logger.Interface
	v *validator.Validate
}
