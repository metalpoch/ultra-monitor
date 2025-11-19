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

	// Group interfaces by OLT
	oltGroups := make(map[string][]struct {
		Interface string
		Bandwidth float64
	})

	for _, result := range results {
		oltName := extractOLTFromInterface(result.Interface)
		if oltName == "" {
			continue
		}

		oltGroups[oltName] = append(oltGroups[oltName], struct {
			Interface string
			Bandwidth float64
		}{
			Interface: result.Interface,
			Bandwidth: result.Bandwidth,
		})
	}

	var bandwidthData []entity.InterfaceBandwidth

	// Process each OLT group
	for oltName, interfaces := range oltGroups {
		var lagInterfaces []struct {
			Interface string
			Bandwidth float64
		}
		var portInterfaces []struct {
			Interface string
			Bandwidth float64
		}

		// Separate LAG interfaces from port interfaces
		for _, iface := range interfaces {
			if isLAGInterface(iface.Interface) {
				lagInterfaces = append(lagInterfaces, iface)
			} else {
				portInterfaces = append(portInterfaces, iface)
			}
		}

		// Apply deduplication logic
		if len(lagInterfaces) > 0 {
			// If there are LAG interfaces, keep only the LAG interface with highest bandwidth
			var bestLAG struct {
				Interface string
				Bandwidth float64
			}

			for _, lag := range lagInterfaces {
				if lag.Bandwidth > bestLAG.Bandwidth {
					bestLAG = lag
				}
			}

			bandwidthData = append(bandwidthData, entity.InterfaceBandwidth{
				OltVerbose: oltName,
				Interface:  bestLAG.Interface,
				Bandwidth:  bestLAG.Bandwidth,
			})
		} else {
			// If no LAG interfaces, keep all port interfaces
			for _, portIface := range portInterfaces {
				bandwidthData = append(bandwidthData, entity.InterfaceBandwidth{
					OltVerbose: oltName,
					Interface:  portIface.Interface,
					Bandwidth:  portIface.Bandwidth,
				})
			}
		}
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

// isLAGInterface determines if an interface is a LAG (Link Aggregation Group) interface
// LAG interfaces typically don't have port numbers with slashes
// Examples of LAG interfaces:
// - "scr-hwsr-01-Eth-6-To_OLT-HW-SAN-CRISTOBAL-03" -> true (no port numbers with slashes)
// - "cco-sar-00-smg-4-To_OLT-HW-CARICUAO-03" -> true (no port numbers with slashes)
// Examples of port interfaces:
// - "scr-hwsr-01-7/1/2-To_OLT-HW-SAN-CRISTOBAL-03" -> false (has port numbers with slashes)
// - "cco-sar-00-0/2/0/3-To_OLT-HW-CARICUAO-03" -> false (has port numbers with slashes)
func isLAGInterface(interfaceStr string) bool {
	// Extract the part between the switch name and "-To_"
	upperInterface := strings.ToUpper(interfaceStr)

	// Find the position of "-TO_" in the string
	toIndex := -1
	for i := 0; i < len(upperInterface)-3; i++ {
		if upperInterface[i:i+4] == "-TO_" || upperInterface[i:i+4] == "-TO-" || upperInterface[i:i+4] == "_TO_" {
			toIndex = i // Position of "-TO_"
			break
		}
	}

	if toIndex == -1 {
		return false
	}

	// Extract the port/interface part (everything between last dash before "-TO_" and "-TO_")
	portPart := interfaceStr[:toIndex]

	// Find the last dash in the port part
	lastDashIndex := -1
	for i := len(portPart) - 1; i >= 0; i-- {
		if portPart[i] == '-' {
			lastDashIndex = i
			break
		}
	}

	if lastDashIndex == -1 {
		return false
	}

	// Extract the actual interface identifier (after the last dash)
	interfaceIdentifier := portPart[lastDashIndex+1:]

	// Check if the interface identifier contains port numbers with slashes
	// If it contains numbers and slashes, it's a port interface, not a LAG
	hasNumbers := false
	hasSlashes := false
	for _, char := range interfaceIdentifier {
		if char >= '0' && char <= '9' {
			hasNumbers = true
		}
		if char == '/' {
			hasSlashes = true
		}
	}

	// If it has both numbers and slashes, it's a port interface (not LAG)
	// Otherwise, it's a LAG interface
	return !(hasNumbers && hasSlashes)
}

// extractSwitchFromInterface extracts the switch prefix from the interface string
// Examples:
// - "scr-hwsr-01-7/1/2-To_OLT-HW-SAN-CRISTOBAL-03" -> "scr-hwsr-01"
// - "merii-hwsar-01-1/1/22-To_OLT-HW-MERIDA-01" -> "merii-hwsar-01"
// - "cch-sar-00-0/4/1/5-To_OLT-HW-CACHAMAY-01" -> "cch-sar-00"
func extractSwitchFromInterface(interfaceStr string) string {
	// Convert the interface string to uppercase for case-insensitive search
	upperInterface := strings.ToUpper(interfaceStr)

	// Find the position of "-TO_" in the string
	toIndex := -1
	for i := 0; i < len(upperInterface)-3; i++ {
		if upperInterface[i:i+4] == "-TO_" || upperInterface[i:i+4] == "-TO-" || upperInterface[i:i+4] == "_TO_" {
			toIndex = i // Position of "-TO_"
			break
		}
	}

	// If "-TO_" is found, extract everything before it
	if toIndex != -1 && toIndex > 0 {
		switchWithPorts := interfaceStr[:toIndex]

		// Now we need to extract just the switch name without port numbers
		// The switch name typically ends before the port pattern (numbers and slashes)
		// Look for the last dash that separates switch name from port numbers
		lastDashIndex := -1
		for i := len(switchWithPorts) - 1; i >= 0; i-- {
			if switchWithPorts[i] == '-' {
				// Check if the part after this dash contains numbers and slashes (port pattern)
				if i+1 < len(switchWithPorts) {
					portPart := switchWithPorts[i+1:]
					// Check if port part contains numbers and slashes
					hasNumbers := false
					hasSlashes := false
					for _, char := range portPart {
						if char >= '0' && char <= '9' {
							hasNumbers = true
						}
						if char == '/' {
							hasSlashes = true
						}
					}
					// If it has both numbers and slashes, this is likely the port part
					if hasNumbers && hasSlashes {
						lastDashIndex = i
						break
					}
				}
			}
		}

		// If we found a dash that separates switch name from ports, return only the switch name
		if lastDashIndex != -1 && lastDashIndex > 0 {
			return switchWithPorts[:lastDashIndex]
		}

		// Otherwise return the original string before "-TO_"
		return switchWithPorts
	}

	// If "-TO_" is not found, return empty string
	return ""
}
