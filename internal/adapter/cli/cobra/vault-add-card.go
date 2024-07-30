package cobra

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fishus/go-advanced-gophkeeper/internal/core/domain"
)

func (cli *cliAdapter) vaultAddCardCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "card",
		Short: "Add a new bank card to the vault",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cardNumber, _ := cmd.Flags().GetString("number")
			cardHolder, _ := cmd.Flags().GetString("holder")
			cardExp, _ := cmd.Flags().GetString("exp")
			cardCvc, _ := cmd.Flags().GetString("cvc")
			info, _ := cmd.Flags().GetString("info")

			cardExpDate := new(domain.CardExpDate)
			if len(cardExp) == 5 && strings.Contains(cardExp, "/") {
				month, err := strconv.ParseUint(cardExp[:2], 10, 0)
				if err != nil {
					return domain.ErrInvalidCardExpDate
				}
				cardExpDate.Month = uint(month)

				year, err := strconv.ParseUint(cardExp[3:], 10, 0)
				if err != nil {
					return domain.ErrInvalidCardExpDate
				}
				cardExpDate.Year = uint(year)
			} else if len(cardExp) > 0 {
				return domain.ErrInvalidCardExpDate
			}

			data := domain.VaultDataCard{
				Number:     cardNumber,
				HolderName: cardHolder,
				ExpDate:    *cardExpDate,
				CvcCode:    cardCvc,
				Info:       info,
			}

			_, err := cli.clientService.VaultAddCard(cmd.Context(), data)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			fmt.Println("New card added successfully")

			return nil
		},
	}

	c.Flags().String("number", "", "Card number")
	c.Flags().String("holder", "", "Card holder name")
	c.Flags().String("exp", "", "Card exp. date (MM/YY)")
	c.Flags().String("cvc", "", "Card CVC code")
	c.Flags().StringP("info", "i", "", "Card metadata")
	c.MarkFlagRequired("number")

	return c
}
