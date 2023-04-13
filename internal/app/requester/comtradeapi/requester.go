package comtradeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"proj/internal/app/server"
	"time"
)

type comtradeRequester struct {
	cfg Config

	client *http.Client
	logger *zap.Logger
}

func NewComtradeRequester(cfg Config) (*comtradeRequester, error) {
	return &comtradeRequester{
		cfg:    cfg,
		client: &http.Client{},
		logger: zap.NewNop(),
	}, nil
}

func (r *comtradeRequester) UseLogger(logger *zap.Logger) {
	r.logger = logger
}

func (r *comtradeRequester) trySplit(request ApiRequest, splitSize, depth int) (*Response, error) {
	r.logger.Info("start splitting request")
	defer r.logger.Info("stop splitting request")

	var requests []ApiRequest
	depth++
	if splitSize <= 512 {
		r.logger.Info("splitting by reporter", zap.Int("splitSize", splitSize))
		requests = request.SplitByReporter(splitSize)
	}
	if splitSize > 512 || len(requests) == 1 {
		splitSize = r.cfg.SplitFactor
		depth = 1
		r.logger.Info("splitting by partner2", zap.Int("splitSize", splitSize))
		requests = request.SplitByPartner2(splitSize)
	}
	if len(requests) == 0 {
		return nil, fmt.Errorf("nothing to split")
	}

	response := &Response{
		Response: server.Response{
			FlowCode: request.FlowCode,
			Data:     make([]server.Report, 0),
		},
	}
	for _, req := range requests {
		time.Sleep(1 * time.Second)
		resp, err := r.do(req)
		if err != nil {
			return nil, fmt.Errorf("failed request after splitting: %w", err)
		}
		if resp.Count >= 250000 {
			r.logger.Info("bad split")
			resp, err = r.trySplit(req, depth*2, depth)
			if err != nil {
				return nil, fmt.Errorf("error splitting request: %w\n", err)
			}
		}
		response.Append(resp)
	}
	return response, nil
}

func (r *comtradeRequester) do(request ApiRequest) (*Response, error) {
	url := request.Url(r.cfg.ApiUrl)
	r.logger.Info("start requesting data", zap.String("url", url))
	defer r.logger.Info("finish requesting data")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Ocp-Apim-Subscription-Key", r.cfg.Token)
	apiResponse, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	if code := apiResponse.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("api returned %d", code)
	}
	body, err := io.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	_ = apiResponse.Body.Close()

	response := &Response{
		Response: server.Response{
			FlowCode: request.FlowCode,
		},
	}
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return response, nil
}

func (r *comtradeRequester) get(request ApiRequest) (*server.Response, error) {
	response, err := r.do(request)
	if err != nil {
		return nil, err
	}
	if response.Count >= 250000 {
		r.logger.Warn("some data might be lost", zap.Int("responseCount", int(response.Count)))
		response, err = r.trySplit(request, r.cfg.SplitFactor, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to split request: %s", err)
		}
	}
	return &response.Response, nil
}

func (r *comtradeRequester) GetBatch(ctx context.Context, out chan<- *server.Response, batch server.RequestBatch) int {
	r.logger.Info("start get batch")
	ticker := time.NewTicker(time.Duration(r.cfg.RequestTimeoutSec) * time.Second)

	defer func() {
		r.logger.Info("finish get batch")
		ticker.Stop()
		close(out)
	}()

	requests := batch.SplitBatch()
	r.logger.Info("split batch", zap.Any("requests", requests))

	counter := 0
	for _, request := range requests {
		request := request
		select {
		case <-ticker.C:
			resp, err := r.get(ApiRequest(request))
			if err != nil {
				r.logger.Error("error getting data from api", zap.Error(err))
			} else {
				out <- resp
				if len(resp.Data) > 0 {
					counter++
				}
			}
		case <-ctx.Done():
			break
		}
	}
	return counter
}
