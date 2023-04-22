package controllers

import (
	"github.com/gofiber/fiber/v2"
	"solution/ent"
	"solution/ent/schema"
	"time"
)

func init() {
	type CreateAccountRequest struct {
		FirstName string          `json:"firstName"`
		LastName  string          `json:"lastName"`
		Country   string          `json:"country"`
		BirthDay  string          `json:"birthDay"`
		Currency  schema.Currency `json:"currency"`
	}

	type CreateAccountResponse struct {
		AccountNumber int `json:"accountNumber"`
	}

	type AccountInfoResponse struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}

	controllers = append(controllers, func(api fiber.Router, client *ent.Client) {
		repo := client.Account

		api.Post("/accounts", func(ctx *fiber.Ctx) error {
			var payload CreateAccountRequest

			if err := ctx.BodyParser(&payload); err != nil {
				return err
			}

			if !payload.Currency.Valid() {
				return ctx.Status(400).JSON("Invalid currency")
			} else if payload.Country == "" {
				return ctx.Status(400).JSON("Country is empty")
			} else if payload.BirthDay == "" {
				return ctx.Status(400).JSON("Birth is empty")
			} else if payload.FirstName == "" {
				return ctx.Status(400).JSON("FirstName is empty")
			} else if payload.LastName == "" {
				return ctx.Status(400).JSON("LastName is empty")
			}

			op := repo.Create().
				SetFirstName(payload.FirstName).
				SetLastName(payload.LastName).
				SetCountry(payload.Country).
				SetCurrency(payload.Currency).
				SetAmount(0)

			now := time.Now()

			if birthday, err := time.Parse(time.DateOnly, payload.BirthDay); err != nil {
				return ctx.Status(400).JSON(err.Error())
			} else if birthday.After(now) {
				return ctx.Status(400).JSON("Birthday is in future")
			} else if age := now.Year() - birthday.Year(); age < 14 || age > 120 {
				return ctx.Status(400).JSON("age should be in [14,120]")
			} else {
				op.SetBirthDay(birthday)
			}

			account, err := op.Save(ctx.Context())
			if err != nil {
				return err
			}

			return ctx.Status(200).JSON(CreateAccountResponse{
				AccountNumber: account.ID,
			})
		})

		api.Get("/accounts/:n", func(ctx *fiber.Ctx) error {
			accountNumber, err := ctx.ParamsInt("n")
			if err != nil {
				return ctx.Status(400).JSON("Wrong account number")
			}

			account, err := repo.Get(ctx.Context(), accountNumber)
			if err != nil {
				return err
			}

			return ctx.Status(fiber.StatusOK).JSON(AccountInfoResponse{
				Amount:   account.Amount,
				Currency: account.Currency.String(),
			})
		})

	})
}
