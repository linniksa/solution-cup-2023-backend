package controllers

import (
	"github.com/gofiber/fiber/v2"
	"solution/ent"
	"solution/ent/account"
	"solution/ent/schema"
	"time"
)

func init() {
	type TopUpRequest struct {
		Amount       float64   `json:"amount"`
		TopUpDateStr string    `json:"topUpDate"`
		TopUpDate    time.Time `json:"-"`
	}

	controllers = append(controllers, func(api fiber.Router, client *ent.Client) {
		api.Post("/accounts/:n/top-up", func(ctx *fiber.Ctx) error {
			var payload TopUpRequest
			if err := ctx.BodyParser(&payload); err != nil {
				return err
			}

			if payload.Amount <= 0 {
				return ctx.Status(400).JSON("Wrong amount")
			} else if payload.TopUpDateStr == "" {
				return ctx.Status(400).JSON("Wrong top up date")
			} else if t, err := time.Parse(time.RFC3339, payload.TopUpDateStr); err != nil {
				return ctx.Status(400).JSON("Wrong top up date")
			} else {
				payload.TopUpDate = t
			}

			accountNumber, err := ctx.ParamsInt("n")
			if err != nil {
				return ctx.Status(400).JSON("Wrong account number")
			}

			tx, err := client.Tx(ctx.Context())
			if err != nil {
				return err
			}
			defer tx.Rollback()

			// get from db
			act, err := tx.Account.Query().Where(account.ID(accountNumber)).ForUpdate().Only(ctx.Context())
			if err != nil {
				return err
			}

			_, err = tx.Transaction.Create().
				SetOperation(schema.Deposit).
				SetAmount(payload.Amount).
				SetAccount(act).
				SetDate(payload.TopUpDate).
				Save(ctx.Context())
			if err != nil {
				return err
			}

			err = tx.Account.UpdateOneID(act.ID).SetAmount(act.Amount + payload.Amount).Exec(ctx.Context())
			if err != nil {
				return err
			}

			if err := tx.Commit(); err != nil {
				return err
			}

			return ctx.Status(fiber.StatusOK).JSON(struct{}{})
		})

	})
}
