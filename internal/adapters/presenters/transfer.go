package presenters

import (
	"brick/internal/adapters/dto/response"
	"brick/internal/entities"
	"brick/internal/usecases/transfer"
	"time"

	"github.com/sirupsen/logrus"
)

type transferPresenter struct {
	Log *logrus.Logger
}

func NewTransferPresenter(log *logrus.Logger) transfer.TransferPresenter {
	return &transferPresenter{
		Log: log,
	}
}

func (p *transferPresenter) PresentValidateAccount(entityRecipientAccount *entities.RecipientAccount) *response.ValidateAccount {
	return &response.ValidateAccount{
		BankCode:      entityRecipientAccount.BankCode,
		BankName:      entityRecipientAccount.BankName,
		AccountNumber: entityRecipientAccount.AccountNumber,
		AccountName:   entityRecipientAccount.AccountName,
	}
}

func (p *transferPresenter) PresentDisburse(entityTransfer *entities.Transfer, entityRecipientAccount *entities.RecipientAccount) *response.Disburse {
	return &response.Disburse{
		TransferID:    entityTransfer.ID,
		Status:        entityTransfer.Status,
		Amount:        entityTransfer.Amount,
		BankCode:      entityRecipientAccount.BankCode,
		BankName:      entityRecipientAccount.BankName,
		AccountNumber: entityRecipientAccount.AccountNumber,
		AccountName:   entityRecipientAccount.AccountName,
		CreatedAt:     entityTransfer.CreatedAt.UTC().Format(time.RFC1123Z),
	}
}
