package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatfatcocofat/rosamsoe/app/models"
	"github.com/fatfatcocofat/rosamsoe/app/response"
	"github.com/fatfatcocofat/rosamsoe/app/utils"
	"github.com/fatfatcocofat/rosamsoe/pkg/config"
	"github.com/fatfatcocofat/rosamsoe/pkg/validator"
	"github.com/fatfatcocofat/rosamsoe/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type WalletController struct {
	Config *config.Config
	Logger *zerolog.Logger
}

func NewWalletController(config *config.Config, logger *zerolog.Logger) *WalletController {
	return &WalletController{
		Config: config,
		Logger: logger,
	}
}

// ListController retrieves wallet information.
//
// @Summary List user's wallets
// @Description Retrieve a list of wallets belonging to the authenticated user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Success
// @Failure 401 {object} response.Unauthorized
// @Failure 502 {object} response.BadGateway
// @Router /wallet [get]
func (c *WalletController) ListController(ctx *fiber.Ctx) error {
	user := utils.ParseUserFromCtx(ctx)

	var wallets []models.Wallet
	results := database.DB.Preload("User").Where("user_id = ?", fmt.Sprint(user.ID)).Find(&wallets)
	if results.Error != nil {
		c.Logger.Error().Err(results.Error).Send()
		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadGateway{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	var walletRes []models.WalletResponse
	for _, w := range wallets {
		walletRes = append(walletRes, models.WalletFilterRecord(&w))
	}

	return ctx.JSON(response.Success{
		Success: true,
		Data: fiber.Map{
			"wallets": walletRes,
		},
	})
}

// ShowController get wallet by address
//
// @Summary Get details of a specific wallet
// @Description Retrieves details of a wallet with the specified address for the authenticated user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param address path string true "Wallet address" format("string")
// @Success 200 {object} response.Success
// @Failure 400 {object} response.BadRequest
// @Failure 401 {object} response.Unauthorized
// @Failure 404 {object} response.NotFound
// @Failure 500 {object} response.ServerError
// @Failure 502 {object} response.BadGateway
// @Router /wallet/{address} [get]
func (c *WalletController) ShowController(ctx *fiber.Ctx) error {
	user := utils.ParseUserFromCtx(ctx)
	address := ctx.Params("address")

	if address == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "Invalid wallet address",
		})
	}

	var wallet models.Wallet
	result := database.DB.Preload("User").Where("user_id = ? and address = ?", fmt.Sprint(user.ID), address).First(&wallet)
	if result.Error != nil {
		if database.IsRecordNotFoundError(result.Error) {
			return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound{
				Success: false,
				Message: "Wallet data with this address was not found",
			})
		}

		c.Logger.Error().Err(result.Error).Send()

		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadGateway{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	return ctx.JSON(response.Success{
		Success: true,
		Data: fiber.Map{
			"wallet": models.WalletFilterRecord(&wallet),
		},
	})
}

// CreateController create new wallet.
//
// @Summary Create a new wallet for the authenticated user
// @Description Creates a new wallet for the user. The user can have a maximum of 3 wallets.
// @Tags Wallet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body models.WalletCreateRequest true "Wallet creation request payload"
// @Success 201 {object} response.Success
// @Failure 400 {object} response.BadRequest
// @Failure 401 {object} response.Unauthorized
// @Failure 502 {object} response.BadGateway
// @Failure 500 {object} response.ServerError
// @Router /wallet [post]
func (c *WalletController) CreateController(ctx *fiber.Ctx) error {
	user := utils.ParseUserFromCtx(ctx)

	var payload *models.WalletCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: err.Error(),
		})
	}

	errors := validator.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Errors:  errors,
		})
	}

	if !utils.ValidWalletCurrency(payload.Currency) {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "Currency not supported, allowed currencies: " + strings.Join(models.WalletCurrencies, ", "),
		})
	}

	var wallets []models.Wallet
	results := database.DB.Where("user_id = ?", fmt.Sprint(user.ID)).Find(&wallets)
	if results.Error != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadGateway{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	if len(wallets) >= 4 {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "You currently have the maximum number of wallets",
		})
	}

	newWallet := models.Wallet{
		UserID:   user.ID,
		Currency: strings.ToUpper(payload.Currency),
	}

	for {
		address, err := utils.GenerateWalletAddress()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(response.ServerError{
				Success: false,
				Message: response.SERVER_ERROR_MSG,
			})
		}

		var wallet models.Wallet
		result := database.DB.First(&wallet, "address = ?", address)
		if result.Error != nil { // FIX: It is currently assumed that the error is record not found.
			newWallet.Address = address
			break
		}
	}

	result := database.DB.Create(&newWallet)
	if result.Error != nil {
		c.Logger.Error().Err(result.Error).Send()

		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadGateway{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	database.DB.Preload("User").First(&newWallet, "address = ?", newWallet.Address)

	return ctx.Status(fiber.StatusCreated).JSON(response.Success{
		Success: true,
		Data: fiber.Map{
			"wallet": models.WalletFilterRecord(&newWallet),
		},
	})
}

// DeleteController delete wallet by address.
//
// @Summary Delete a specific wallet
// @Description Deletes a wallet with the specified address for the authenticated user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param address path string true "Wallet address" format("string")
// @Success 204 "No Content"
// @Failure 400 {object} response.BadRequest
// @Failure 401 {object} response.Unauthorized
// @Failure 404 {object} response.NotFound
// @Failure 502 {object} response.BadGateway
// @Router /wallet/{address} [delete]
func (c *WalletController) DeleteController(ctx *fiber.Ctx) error {
	user := utils.ParseUserFromCtx(ctx)
	address := ctx.Params("address")

	if address == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "Invalid wallet address",
		})
	}

	var wallet models.Wallet
	result := database.DB.Where("user_id = ? and address = ?", fmt.Sprint(user.ID), address).First(&wallet)
	if result.Error != nil {
		if database.IsRecordNotFoundError(result.Error) {
			return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound{
				Success: false,
				Message: "Wallet data with this address was not found",
			})
		}

		c.Logger.Error().Err(result.Error).Send()

		return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound{
			Success: false,
			Message: "Wallet data with this address was not found",
		})
	}

	result = database.DB.Delete(&models.Wallet{}, "id = ?", fmt.Sprint(wallet.ID))

	if result.Error != nil || result.RowsAffected == 0 {
		c.Logger.Error().Err(result.Error).Send()

		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadGateway{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// UpdateController update wallet by address.
//
// @Summary Update details of a specific wallet
// @Description Updates details of a wallet with the specified address for the authenticated user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param address path string true "Wallet address" format("string")
// @Param payload body models.WalletUpdateRequest true "Wallet update request payload"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.BadRequest
// @Failure 401 {object} response.Unauthorized
// @Failure 404 {object} response.NotFound
// @Failure 502 {object} response.BadGateway
// @Router /wallets/{address} [put]
func (c *WalletController) UpdateController(ctx *fiber.Ctx) error {
	user := utils.ParseUserFromCtx(ctx)
	address := ctx.Params("address")

	if address == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: "Invalid wallet address",
		})
	}

	var wallet models.Wallet
	result := database.DB.Where("user_id = ? and address = ?", fmt.Sprint(user.ID), address).First(&wallet)
	if result.Error != nil {
		if database.IsRecordNotFoundError(result.Error) {
			return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound{
				Success: false,
				Message: "Wallet data with this address was not found",
			})
		}

		c.Logger.Error().Err(result.Error).Send()

		return ctx.Status(fiber.StatusNotFound).JSON(response.NotFound{
			Success: false,
			Message: "Wallet data with this address was not found",
		})
	}

	var payload *models.WalletUpdateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Message: err.Error(),
		})
	}

	errors := validator.ValidateStruct(payload)
	if errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
			Success: false,
			Errors:  errors,
		})
	}

	updates := make(map[string]interface{})

	if payload.Currency != "" {
		if !utils.ValidWalletCurrency(payload.Currency) {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.BadRequest{
				Success: false,
				Message: "Currency not supported, allowed currencies: " + strings.Join(models.WalletCurrencies, ", "),
			})
		}

		updates["currency"] = strings.ToUpper(payload.Currency)
	}

	updates["updated_at"] = time.Now()

	result = database.DB.Model(&wallet).Updates(updates)

	if result.Error != nil {
		return ctx.Status(fiber.StatusBadGateway).JSON(response.BadGateway{
			Success: false,
			Message: response.BAD_GATEWAY_MSG,
		})
	}

	database.DB.Preload("User").Where("user_id = ? and address = ?", fmt.Sprint(user.ID), address).First(&wallet)

	return ctx.Status(fiber.StatusOK).JSON(response.Success{
		Success: true,
		Data: fiber.Map{
			"wallet": wallet,
		},
	})
}
