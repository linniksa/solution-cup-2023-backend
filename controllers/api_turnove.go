package controllers

import (
	"github.com/gofiber/fiber/v2"
	"solution/ent"
	"strconv"
	"time"
)

func init() {
	type Response struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	}

	controllers = append(controllers, func(api fiber.Router, client *ent.Client) {
		api.Get("/account-turnover/:n", func(ctx *fiber.Ctx) error {
			accountNumber, err := ctx.ParamsInt("n")
			if err != nil {
				return ctx.Status(400).JSON("Wrong account number")
			}

			var startDate time.Time
			var endDate time.Time

			{
				if ctx.Query("startDate", "") != "" {
					if t, err := time.Parse(time.RFC3339, ctx.Query("startDate", "")); err != nil {
						return ctx.Status(400).JSON("Wrong startDate")
					} else {
						startDate = t
					}
				}

				if ctx.Query("endDate", "") != "" {
					if t, err := time.Parse(time.RFC3339, ctx.Query("endDate", "")); err != nil {
						return ctx.Status(400).JSON("Wrong startDate")
					} else {
						endDate = t
					}
				}
			}

			if !endDate.IsZero() && startDate.After(endDate) {
				return ctx.Status(400).JSON("StartDate bigger EndDate")
			}

			account, err := client.Account.Get(ctx.Context(), accountNumber)
			if err != nil {
				return err
			}

			sql := `select sum(amount) from transactions where account_id = $1`

			args := make([]interface{}, 0, 5)
			args = append(args, account.ID)

			if !startDate.IsZero() {
				sql += " and date >= $" + strconv.Itoa(len(args)+1)
				args = append(args, startDate)
			}
			if !endDate.IsZero() {
				sql += " and date <= $" + strconv.Itoa(len(args)+1)
				args = append(args, endDate)
			}

			result, err := client.DB().Query(sql, args...)
			if err != nil {
				return err
			}
			var amount float64
			for result.Next() {
				var d *float64
				if err := result.Scan(&d); err != nil {
					return err
				}
				if d != nil {
					amount = *d
				}
			}

			return ctx.Status(fiber.StatusOK).JSON(Response{
				Amount:   amount,
				Currency: account.Currency.String(),
			})
		})

	})
}
