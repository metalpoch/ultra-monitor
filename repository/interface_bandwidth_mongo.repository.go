package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/metalpoch/ultra-monitor/entity"
)

type InterfaceBandwidthMongoRepository struct {
	collection *mongo.Collection
}

func NewInterfaceBandwidthMongoRepository(mongoDB *mongo.Database) *InterfaceBandwidthMongoRepository {
	return &InterfaceBandwidthMongoRepository{
		collection: mongoDB.Collection("ehealth2"),
	}
}

func (r *InterfaceBandwidthMongoRepository) GetInterfaceBandwidthFromMongoDB(ctx context.Context, startDate, endDate time.Time) ([]entity.InterfaceBandwidth, error) {
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"trends.group", "RP-ENLACES-GPON-OLT-HUWEI"},
				{"trends.elements.times", bson.D{
					{"$gte", startDate},
					{"$lte", endDate},
				}},
			}},
		},
		bson.D{{"$unwind", "$trends"}},
		bson.D{
			{"$match", bson.D{
				{"trends.group", "RP-ENLACES-GPON-OLT-HUWEI"},
				}},
		},
		bson.D{{"$unwind", "$trends.elements"}},
		bson.D{
			{"$project", bson.D{
				{"interface", "$trends.elements.interface"},
				{"lastBandwidth", bson.D{
					{"$arrayElemAt", bson.A{"$trends.elements.bandwidth", -1}},
				}},
			}},
		},
		bson.D{
			{"$group", bson.D{
				{"_id", "$interface"},
				{"bandwidth", bson.D{
					{"$first", "$lastBandwidth"},
				}},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"interface", "$_id"},
				{"bandwidth", 1},
			}},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Interface string  `bson:"interface"`
		Bandwidth float64 `bson:"bandwidth"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var bandwidthData []entity.InterfaceBandwidth
	for _, result := range results {
		// Extract OLT name from interface string
		oltName := extractOLTFromInterface(result.Interface)

		bandwidthData = append(bandwidthData, entity.InterfaceBandwidth{
			Interface: result.Interface,
			Olt:       oltName,
			Bandwidth: result.Bandwidth,
			CreatedAt: time.Now(),
		})
	}

	return bandwidthData, nil
}

// extractOLTFromInterface extracts the OLT name from the interface string
// Examples:
// - "ncc-sar-00-0/0/1/2-To_OLT-HW-NUEVA-CARACAS-01" -> "OLT-HW-NUEVA-CARACAS-01"
// - "brs-hwsar-00-Eth7-To_OLT-HW-BARINAS-00" -> "OLT-HW-BARINAS-00"
// - "pde-sar-00-0/4/1/4-To_OLT-HW-Prados-del-este-03" -> "OLT-HW-Prados-del-este-03"
// - "cri-dsw-00-1/1/5-To_OLT-CAMURI" -> "OLT-CAMURI"
func extractOLTFromInterface(interfaceStr string) string {
	// Find the position of "-To_" in the string
	toIndex := -1
	for i := 0; i < len(interfaceStr)-3; i++ {
		if interfaceStr[i:i+4] == "-To_" {
			toIndex = i + 4 // Skip "-To_" itself
			break
		}
	}

	// If "-To_" is found, return everything after it
	if toIndex != -1 && toIndex < len(interfaceStr) {
		return interfaceStr[toIndex:]
	}

	// If "-To_" is not found, return the original interface
	return interfaceStr
}

