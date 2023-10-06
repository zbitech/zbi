package http

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zbitech/controller/app/service-api/request"
	"github.com/zbitech/controller/app/service-api/response"
	"github.com/zbitech/controller/internal/helper"
	"github.com/zbitech/controller/pkg/logger"
	"github.com/zbitech/controller/pkg/model"
)

func SetPolicy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	var input struct {
		Policy *model.PolicyInfo `json:"policy"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}

	helper.SetPolicy((input.Policy))
}

func setBlockchains(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	var input struct {
		Blockchains []model.BlockchainInfo `json:"blockchains"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}

	helper.SetBlockchains(ctx, input.Blockchains, false)
}

func setBlockchain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx)

	var input struct {
		Blockchain *model.BlockchainInfo `json:"blockchain"`
	}

	if err := request.ReadJSON(w, r, &input); err != nil {
		log.WithFields(logrus.Fields{"error": err}).Errorf("failed to read input")
		response.BadRequestResponse(w, r, err)
		return
	}

	helper.SetBlockchain(ctx, input.Blockchain, false)
}
