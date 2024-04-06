package transfer

import (
	"brick/internal/adapters/dto/request"
	"brick/internal/adapters/dto/response"
	"brick/internal/entities"
	"brick/internal/pkg/serror"
	"context"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, entityTransfer *entities.Transfer) (*entities.Transfer, *serror.Error)
	UpdateTransferStatus(ctx context.Context, entityTransfer *entities.Transfer) *serror.Error
}

type TransferPresenter interface {
	PresentValidateAccount(entityRecipientAccount *entities.RecipientAccount) *response.ValidateAccount
	PresentDisburse(entityTransfer *entities.Transfer, entityRecipientAccount *entities.RecipientAccount) *response.Disburse
}

type TransferUseCase interface {
	ValidateAccount(ctx context.Context, validateAccountRequest *request.ValidateAccount) (*response.ValidateAccount, *serror.Error)
	CreateTransfer(ctx context.Context, disburseRequest *request.Disburse) (*response.Disburse, *serror.Error)
	UpdateTransferStatus(ctx context.Context, transferCallbackRequest *request.TransferCallback) *serror.Error
	MockDisburse(ctx context.Context, disburseRequest *request.Disburse, transferCallbackRequest *request.TransferCallback)
	MockSendValidationResponse(validateAccountRequest *request.ValidateAccount, result bool) (*response.ValidateBank, *serror.Error)
}
