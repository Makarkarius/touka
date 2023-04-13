package mongostorage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"proj/internal/app/server"
)

func (d *database) selectCollection(flowCode string) (*mongo.Collection, error) {
	var collection *mongo.Collection
	switch flowCode {
	case "M":
		collection = d.importCollection
	case "X":
		collection = d.exportCollection
	case "":
		return nil, fmt.Errorf("response does not have a flow code")
	default:
		return nil, fmt.Errorf("response flow code is invalid")
	}
	return collection, nil
}

func (d *database) InsertChan(ctx context.Context, in <-chan *server.Response) {
	d.logger.Info("start insert chan")
	defer d.logger.Info("finish insert chan")

	opts := options.InsertMany().SetOrdered(false)
	for resp := range in {
		if resp.Count == 0 {
			d.logger.Warn("empty response")
			continue
		}

		collection, err := d.selectCollection(resp.FlowCode)
		if err != nil {
			d.logger.Error("error choosing collection", zap.Error(err))
			continue
		}

		docs := make([]interface{}, 0, int(resp.Count))
		for _, record := range resp.Data {
			document := bson.M{
				"period":       record.Period,
				"reporterCode": int(record.ReporterCode),
				"partnerCode":  int(record.PartnerCode),
				"partner2Code": int(record.Partner2Code),
				"cmdCode":      record.CmdCode,
				"primaryValue": record.PrimaryValue,
			}
			docs = append(docs, document)
		}

		d.logger.Info("start inserting data", zap.String("flow_code", resp.FlowCode))
		_, err = collection.InsertMany(ctx, docs, opts)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				d.logger.Warn("duplicate key error")
			} else {
				d.logger.Error("error inserting data", zap.String("flow_code", resp.FlowCode), zap.Error(err))
			}
		} else {
			d.logger.Info("successfully inserted data", zap.Int("documentCount", len(docs)))
		}
		d.logger.Info("finish inserting data", zap.String("flow_code", resp.FlowCode))
	}
}
