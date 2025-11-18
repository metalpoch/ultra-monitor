package repository

import (
	"context"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InterfaceBandwidthRepository interface {
	GetInterfaceBandwidthFromMongoDB(ctx context.Context, startDate, endDate time.Time) ([]entity.InterfaceBandwidth, error)
	GetAll(ctx context.Context) ([]entity.InterfaceBandwidth, error)
	CleanInterfacesBandwidth(ctx context.Context) error
	InsertInterfaceBandwidth(ctx context.Context, bandwidthData []entity.InterfaceBandwidth) error
}

type interfaceBandwidthRepository struct {
	db         *sqlx.DB
	collection *mongo.Collection
}

func NewInterfaceBandwidthRepository(db *sqlx.DB, mongoDB *mongo.Database) *interfaceBandwidthRepository {
	return &interfaceBandwidthRepository{
		db:         db,
		collection: mongoDB.Collection("ehealth2"),
	}
}

func (r *interfaceBandwidthRepository) GetInterfaceBandwidthFromMongoDB(ctx context.Context, startDate, endDate time.Time) ([]entity.InterfaceBandwidth, error) {
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
		oltName := extractOLTFromInterface(result.Interface)

		bandwidthData = append(bandwidthData, entity.InterfaceBandwidth{
			OltVerbose: oltName,
			Interface:  result.Interface,
			Bandwidth:  result.Bandwidth,
		})
	}

	return bandwidthData, nil
}

func (r *interfaceBandwidthRepository) CleanInterfacesBandwidth(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM interfaces_bandwidth;")
	return err
}

func (r *interfaceBandwidthRepository) InsertInterfaceBandwidth(ctx context.Context, bandwidthData []entity.InterfaceBandwidth) error {
	query := `INSERT INTO interfaces_bandwidth (olt_verbose, interface, bandwidth) VALUES (:olt_verbose, :interface, :bandwidth);`
	_, err := r.db.NamedExecContext(ctx, query, bandwidthData)
	return err
}

func (r *interfaceBandwidthRepository) GetAll(ctx context.Context) ([]entity.InterfaceBandwidth, error) {
	var res []entity.InterfaceBandwidth
	err := r.db.SelectContext(ctx, &res, `SELECT * FROM interfaces_bandwidth;`)
	return res, err
}

// extractOLTFromInterface extracts the OLT name from the interface string
// Examples:
// - "pde-sar-00-0/4/1/4-To_OLT-HW-Prados-del-este-03" -> "OLT-HW-Prados-del-este-03"
func extractOLTFromInterface(interfaceStr string) string {
	// Convert the interface string to uppercase for case-insensitive search
	upperInterface := strings.ToUpper(interfaceStr)

	// Find the position of "-TO_" in the string
	toIndex := -1
	for i := 0; i < len(upperInterface)-3; i++ {
		if upperInterface[i:i+4] == "-TO_" || upperInterface[i:i+4] == "-TO-" || upperInterface[i:i+4] == "_TO_" {
			toIndex = i + 4 // Skip "-TO_" itself
			break
		}
	}

	// If "-TO_" is found, return everything after it
	if toIndex != -1 && toIndex < len(interfaceStr) {
		return interfaceStr[toIndex:]
	}

	// If "-TO_" is not found, return the original interface
	return ""
}
