package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"math"
	"solution/ent"
	"solution/ent/account"
	"solution/ent/schema"
	"solution/rates"
	"time"
)

func init() {
	type TransferRequest struct {
		ReceiverAccount        int     `json:"receiverAccount"`
		SenderAccount          int     `json:"senderAccount"`
		AmountInSenderCurrency float64 `json:"amountInSenderCurrency"`

		TransferDateStr string    `json:"transferDate"`
		TransferDate    time.Time `json:"-"`
	}

	controllers = append(controllers, func(api fiber.Router, client *ent.Client, logger *zap.Logger, cs *rates.Manager) {
		api.Post("/transfers", func(ctx *fiber.Ctx) error {
			var payload TransferRequest
			if err := ctx.BodyParser(&payload); err != nil {
				return err
			}

			if payload.AmountInSenderCurrency <= 0 {
				return ctx.Status(400).JSON("Wrong amount")
			} else if payload.TransferDateStr == "" {
				return ctx.Status(400).JSON("Wrong top up date")
			} else if t, err := time.Parse(time.RFC3339, payload.TransferDateStr); err != nil {
				return ctx.Status(400).JSON("Wrong top up date")
			} else {
				payload.TransferDate = t
			}

			tx, err := client.Tx(ctx.Context())
			if err != nil {
				return err
			}
			defer tx.Rollback()

			// get from db

			sender, err := tx.Account.Query().Where(account.ID(payload.SenderAccount)).ForUpdate().Only(ctx.Context())
			if err != nil {
				return err
			}

			if sender.Amount < payload.AmountInSenderCurrency {
				return ctx.Status(400).JSON("Not enough money")
			}

			receiver, err := tx.Account.Query().Where(account.ID(payload.ReceiverAccount)).ForUpdate().Only(ctx.Context())
			if err != nil {
				return err
			}

			amountSend := payload.AmountInSenderCurrency
			amountReceiver := payload.AmountInSenderCurrency

			if sender.Currency != receiver.Currency {
				rt := cs.Rate()
				if rt == nil {
					logger.Error("no rate")

					return ctx.Status(500).JSON("no rate")
				}
				rate := rt.GetRate(sender.Currency, receiver.Currency)
				if rate <= 0 {
					logger.Error("rate is zero")

					return ctx.Status(500).JSON("no rate")
				}
				amountReceiver = math.RoundToEven(amountReceiver*rate*100) / 100
			}

			if err := tx.Account.UpdateOne(sender).SetAmount(sender.Amount - amountSend).Exec(ctx.Context()); err != nil {
				return err
			}

			if err := tx.Account.UpdateOne(receiver).SetAmount(receiver.Amount + amountReceiver).Exec(ctx.Context()); err != nil {
				return err
			}

			err = tx.Transaction.Create().
				SetOperation(schema.Transfer).
				SetDate(payload.TransferDate).
				SetAccount(receiver).
				SetAmount(amountReceiver).
				Exec(ctx.Context())
			if err != nil {
				return err
			}

			err = tx.Transaction.Create().
				SetOperation(schema.Transfer).
				SetDate(payload.TransferDate).
				SetAccount(sender).
				SetAmount(amountSend).
				Exec(ctx.Context())

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
