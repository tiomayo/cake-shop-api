package cakes

import (
	"cake-store/internal/helpers"
	"context"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"strconv"
)

type SvcInterface interface {
	List(ctx echo.Context) error
	Get(ctx echo.Context) error
	Create(ctx echo.Context) error
	Update(ctx echo.Context) error
	Delete(ctx echo.Context) error
}

type svcImplementation struct {
	repo RepoInterface
}

func NewHandler(repo RepoInterface) SvcInterface {
	return svcImplementation{repo}
}

// List godoc
// @Summary List all cakes
// @Description This endpoint for get list of cakes
// @Tags Cakes
// @Accept  json
// @Produce  json
// @Param services query ListRequestDto true "Find query"
// @Success 200 {array} Cake
// @Failure 422 {object} helpers.JSONResponse
// @Failure 500 {object} helpers.JSONResponse
// @Router /cakes [get]
func (s svcImplementation) List(ctx echo.Context) error {
	request := ListRequestDto{}
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ctx.Validate(&request); err != nil {
		return err
	}

	if request.Limit == 0 {
		request.Limit = 10
	}

	res, total, err := s.repo.List(context.TODO(), request)
	if err != nil {
		return err
	}

	page := math.Ceil(float64(total) / float64(request.Limit))
	ctx.Response().Header().Add("Pagination-Rows", strconv.Itoa(int(total)))
	ctx.Response().Header().Add("Pagination-Page", strconv.Itoa(int(page)))
	ctx.Response().Header().Add("Pagination-Limit", strconv.Itoa(request.Limit))
	return ctx.JSON(http.StatusOK, res)
}

// Get godoc
// @Summary Get detail of cake
// @Description This endpoint for get detail of cake
// @Tags Cakes
// @Accept  json
// @Produce  json
// @Param id path string true "cake id"
// @Success 200 {object} Cake
// @Failure 422 {object} helpers.JSONResponse
// @Failure 204 {object} helpers.JSONResponse
// @Failure 500 {object} helpers.JSONResponse
// @Router /cakes/{id} [get]
func (s svcImplementation) Get(ctx echo.Context) error {
	ID, errConv := strconv.Atoi(ctx.Param("id"))
	if errConv != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid id")
	}

	data, errGet := s.repo.Get(context.TODO(), ID)
	if data == nil {
		return echo.NewHTTPError(http.StatusNoContent, "Data Not Found")
	}
	if errGet != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errGet.Error())
	}

	return ctx.JSON(http.StatusOK, data)
}

// Create godoc
// @Summary Create cake
// @Description This endpoint for creating cake
// @Tags Cakes
// @Accept  json
// @Produce  json
// @Param Request body RequestDto true "Create cakes"
// @Success 200 {object} Cake
// @Failure 422 {object} helpers.JSONResponse
// @Failure 500 {object} helpers.JSONResponse
// @Router /cakes [post]
func (s svcImplementation) Create(ctx echo.Context) error {
	request := RequestDto{}
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ctx.Validate(&request); err != nil {
		return err
	}

	errCreate := s.repo.Create(context.TODO(), request)
	if errCreate != nil {
		return errCreate
	}
	return ctx.JSON(http.StatusCreated, helpers.JSONResponse{Message: "Cake Created"})
}

// Update godoc
// @Summary Update cake
// @Description This endpoint for updating cake
// @Tags Cakes
// @Accept  json
// @Produce  json
// @Param id path string true "cake id"
// @Param Request body UpdateRequestDto true "Update cakes"
// @Success 200 {object} Cake
// @Failure 422 {object} helpers.JSONResponse
// @Failure 204 {object} helpers.JSONResponse
// @Failure 500 {object} helpers.JSONResponse
// @Router /cakes/{id} [patch]
func (s svcImplementation) Update(ctx echo.Context) error {
	request := UpdateRequestDto{}
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := ctx.Validate(&request); err != nil {
		return err
	}

	exist, errGet := s.repo.Get(context.TODO(), request.ID)
	if exist == nil {
		return echo.NewHTTPError(http.StatusNoContent, "Data Not Found")
	}
	if errGet != nil {
		return errGet
	}

	errCreate := s.repo.Update(context.TODO(), request)
	if errCreate != nil {
		return errCreate
	}
	return ctx.JSON(http.StatusOK, helpers.JSONResponse{Message: "Cake Updated"})
}

// Delete godoc
// @Summary Delete cake
// @Description This endpoint for deleting cake
// @Tags Cakes
// @Accept  json
// @Produce  json
// @Param id path string true "cake id"
// @Success 200 {object} Cake
// @Failure 422 {object} helpers.JSONResponse
// @Failure 204 {object} helpers.JSONResponse
// @Failure 500 {object} helpers.JSONResponse
// @Router /cakes/{id} [delete]
func (s svcImplementation) Delete(ctx echo.Context) error {
	ID, errConv := strconv.Atoi(ctx.Param("id"))
	if errConv != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Invalid id")
	}
	exist, errGet := s.repo.Get(context.TODO(), ID)
	if exist == nil {
		return echo.NewHTTPError(http.StatusNoContent, "Data Not Found")
	}
	if errGet != nil {
		return errGet
	}

	err := s.repo.Delete(context.TODO(), ID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, "Success")
}
