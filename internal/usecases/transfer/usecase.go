package transfer

import (
	"brick/internal/adapters/dto/request"
	"brick/internal/adapters/dto/response"
	"brick/internal/entities"
	"brick/internal/pkg/serror"
	"brick/internal/pkg/utils"
	"brick/internal/usecases/recipientAccount"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	TransferRepository         TransferRepository
	RecipientAccountRepository recipientAccount.RecipientAccountRepository
	TransferPresenter          TransferPresenter
	DB                         *sql.DB
	Log                        *logrus.Logger
}

func NewTransferUsecase(
	userRepository TransferRepository,
	recipientAccountRepository recipientAccount.RecipientAccountRepository,
	presenter TransferPresenter,
	db *sql.DB,
	log *logrus.Logger) TransferUseCase {
	return &usecase{
		TransferRepository:         userRepository,
		RecipientAccountRepository: recipientAccountRepository,
		TransferPresenter:          presenter,
		DB:                         db,
		Log:                        log,
	}
}

func (u *usecase) ValidateAccount(ctx context.Context, validateAccountRequest *request.ValidateAccount) (*response.ValidateAccount, *serror.Error) {
	errChan := make(chan *serror.Error, 2)
	recipientAccountIdChan := make(chan int, 1)

	recipientAccount := new(entities.RecipientAccount)
	recipientAccount.BankCode = validateAccountRequest.BankCode
	recipientAccount.AccountNumber = validateAccountRequest.AccountNumber

	go func() {
		recipientAccount, serr := u.RecipientAccountRepository.GetRecipientAccountByBankCodeAndAccountNumber(ctx, recipientAccount)
		if serr != nil {
			u.Log.Errorf("[Transfer][Usecase] while RecipientAccountRepository.GetRecipientAccountByBankCodeAndAccountNumber: %s", serr.Error())
			errChan <- serr
			return
		}

		if recipientAccount.IsNotExist() {
			recipientAccount.VerificationStatus = "unverified"
			recipientAccount, serr = u.RecipientAccountRepository.CreateRecipientAccount(ctx, recipientAccount)
		} else if recipientAccount.IsExist() {
			recipientAccount, serr = u.RecipientAccountRepository.UpdateRecipientAccountByID(ctx, recipientAccount)
		}

		if serr != nil {
			u.Log.Errorf("[Transfer][Usecase] while Upsert RecipientAccountRepository: %s", serr.Error())
			errChan <- serr
			return
		}

		recipientAccountIdChan <- recipientAccount.ID
		errChan <- nil
	}()

	validationResponse := new(response.ValidateBank)
	go func() {
		// Simulate url call
		// Change second parameter argument to false, to see error
		bankValidationResult, err := u.MockSendValidationResponse(validateAccountRequest, true)
		if err != nil {
			errChan <- err
			return
		}

		validationResponse = bankValidationResult
		accountId := <-recipientAccountIdChan

		recipientAccount.ID = accountId
		recipientAccount.AccountName = validationResponse.AccountName
		recipientAccount.AccountNumber = validationResponse.AccountNumber
		recipientAccount.AccountName = validationResponse.AccountName
		recipientAccount.BankCode = validationResponse.BankCode
		recipientAccount.BankName = validationResponse.BankName
		recipientAccount.VerificationStatus = "verified"

		if validationResponse.Success {
			recipientAccount, err = u.RecipientAccountRepository.UpdateRecipientAccountByID(ctx, recipientAccount)
			if err != nil {
				errChan <- err
				return
			}
		}
		errChan <- nil
	}()

	for i := 0; i < 2; i++ {
		serr := <-errChan
		if serr != nil {
			u.Log.Errorf("[Transfer][Usecase] while check error chan: %s", serr.Error())
			close(recipientAccountIdChan)
			return nil, serr
		}
	}
	close(recipientAccountIdChan)

	return u.TransferPresenter.PresentValidateAccount(recipientAccount), nil
}

func (u *usecase) CreateTransfer(ctx context.Context, disburseRequest *request.Disburse) (*response.Disburse, *serror.Error) {
	transfer := new(entities.Transfer)
	recipientAccount := new(entities.RecipientAccount)

	transfer.Amount = disburseRequest.Amount
	transfer.SenderAccountID = disburseRequest.UserID
	transfer.Status = "processed"

	recipientAccount.BankCode = disburseRequest.BankCode
	recipientAccount.AccountNumber = disburseRequest.AccountNumber

	recipientAccount, serr := u.RecipientAccountRepository.GetRecipientAccountByBankCodeAndAccountNumber(ctx, recipientAccount)
	if serr != nil {
		u.Log.Errorf("[Transfer][Usecase] while RecipientAccountRepository.GetRecipientAccountByBankCodeAndAccountNumber: %s", serr.Error())
		return nil, serr
	}

	transfer.RecipientAccountID = recipientAccount.ID
	transfer, serr = u.TransferRepository.CreateTransfer(ctx, transfer)
	if serr != nil {
		u.Log.Errorf("[Transfer][Usecase] while TransferRepository.CreateTransfer: %s", serr.Error())
		return nil, serr
	}

	return u.TransferPresenter.PresentDisburse(transfer, recipientAccount), nil
}

func (u *usecase) UpdateTransferStatus(ctx context.Context, transferCallbackRequest *request.TransferCallback) *serror.Error {
	transfer := new(entities.Transfer)
	transfer.ID = transferCallbackRequest.TransferID
	transfer.Status = transferCallbackRequest.Status

	serr := u.TransferRepository.UpdateTransferStatus(ctx, transfer)
	if serr != nil {
		u.Log.Errorf("[Transfer][Usecase] while RecipientAccountRepository.GetRecipientAccountByBankCodeAndAccountNumber: %s", serr.Error())
		return serr
	}

	return nil
}

func (u *usecase) MockSendValidationResponse(request *request.ValidateAccount, result bool) (*response.ValidateBank, *serror.Error) {
	// simulate operation
	// basically call http.Get
	time.Sleep(2 * time.Second)
	if !result {
		u.Log.Infof("Bank %s sending error response, Account Number %s is not a valid account", request.BankCode, request.AccountNumber)
		return nil, serror.NewError(http.StatusBadRequest, 0, "validation error", "bank account is not validated")
	}
	u.Log.Infof("Bank %s sending success response, Account Number %s is a valid account", request.BankCode, request.AccountNumber)
	return &response.ValidateBank{
		Success:       result,
		BankCode:      request.BankCode,
		BankName:      uuid.New().String(),
		AccountNumber: request.AccountNumber,
		AccountName:   "random dev",
	}, nil
}

func (u *usecase) MockDisburse(ctx context.Context, disburseRequest *request.Disburse, transferCallbackRequest *request.TransferCallback) {
	// Simulate POST request to send money to bank
	time.Sleep(3 * time.Second)
	u.Log.Infof("Disburse request sent to Bank %s, with amount %2.f, to Account Number: %s - Account Name: %s", disburseRequest.BankCode, disburseRequest.Amount, disburseRequest.AccountNumber, disburseRequest.AccountName)

	time.Sleep(5 * time.Second)
	u.Log.Infof("Bank %s is calling our system callback URL with status: `%s`", disburseRequest.BankCode, transferCallbackRequest.Status)

	jsonRequestBody, err := json.Marshal(transferCallbackRequest)
	if err != nil {
		u.Log.Errorf("Bank found error on our system, case when: json.Marshal")
		return
	}

	appPort, err := utils.ReadStringEnvKey("APP_PORT", true)
	if err != nil {
		u.Log.Errorf("Bank found error on our system, case when: utils.ReadStringEnvKey")
		return
	}

	result, err := http.Post(fmt.Sprintf("http://127.0.0.1:%s/api/v1/transfer/callback", appPort), "application/json", bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		u.Log.Errorf("Bank found error on our system, case when: failed when make call on system")
		return
	}
	defer result.Body.Close()

	if result.StatusCode != 200 {
		u.Log.Errorf("Bank failed to get 'OK' status from our system, update status failed")
		return
	}

	_, err = io.ReadAll(result.Body)
	if err != nil {
		u.Log.Errorf("Bank found error on our system, case when: read response from Callback Url")
		return
	}

	u.Log.Info("Bank successfully calling our callback URL, transfer status is updated")
}
