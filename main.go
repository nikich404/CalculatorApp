package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Calculation struct {
	Id         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

var calculations = []Calculation{}

func calculationExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), err

}
func getCalculation(c echo.Context) error {
	return c.JSON(http.StatusOK, calculations)
}

func postCalculation(c echo.Context) error {
	var req = CalculationRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no valid request"})
	}

	result, err := calculationExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no valid request"})
	}

	calc := Calculation{
		Id:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}

	calculations = append(calculations, calc)
	return c.JSON(http.StatusCreated, calc)

}

func patchCalculations(c echo.Context) error {
	id := c.Param("id")
	var req = CalculationRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no valid request"})
	}

	result, err := calculationExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no valid request"})
	}

	for i, calculation := range calculations {
		if calculation.Id == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = result
			return c.JSON(http.StatusOK, calculations[i])

		}

	}
	return c.JSON(http.StatusBadRequest, map[string]string{"error": "Calculation not found"})
}
func deletColculation(c echo.Context) error {
	id := c.Param("id")

	for i, calculation := range calculations {
		if calculation.Id == id {
			calculations = append(calculations[:i], calculations[i+1:]...)
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/calculations", getCalculation)
	e.POST("/calculations", postCalculation)
	e.PATCH("calculations/:id", patchCalculations)
	e.DELETE("calculations/:id", deletColculation)

	e.Start("localhost:8080")

}
