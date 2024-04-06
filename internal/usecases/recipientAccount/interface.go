package recipientAccount

import (
	"brick/internal/entities"
	"brick/internal/pkg/serror"
	"context"
)

type RecipientAccountRepository interface {
	GetRecipientAccountByBankCodeAndAccountNumber(ctx context.Context, entityRecipientAccount *entities.RecipientAccount) (*entities.RecipientAccount, *serror.Error)
	UpdateRecipientAccountByID(ctx context.Context, entityRecipientAccount *entities.RecipientAccount) (*entities.RecipientAccount, *serror.Error)
	CreateRecipientAccount(ctx context.Context, entityRecipientAccount *entities.RecipientAccount) (*entities.RecipientAccount, *serror.Error)
}
