package store

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"github.com/safecility/microservices/go/process/hotdrop/messages"
)

// TODO adjust locationId when changed on local db
const (
	getDeviceStmt = `SELECT name, uid, locationId as groupId, companyId, power_factor, line_voltage
		FROM device
		JOIN safecility.power_device p on device.id = p.deviceId
		WHERE type='power' AND device.uid = ?`
)

// DeviceSql is accessed both directly and by the device Cache, direct access is only for uplinks which show Compliance events
type DeviceSql struct {
	sqlDB          *sql.DB
	getDeviceByUID *sql.Stmt
	getDaliStatus  *sql.Stmt
}

func NewDeviceSql(db *sql.DB) (*DeviceSql, error) {
	sqlDB := &DeviceSql{
		sqlDB: db,
	}
	var err error

	if sqlDB.getDeviceByUID, err = db.Prepare(getDeviceStmt); err != nil {
		return nil, err
	}

	return sqlDB, nil
}

func (db DeviceSql) GetDevice(uid string) (*messages.PowerDevice, error) {
	log.Debug().Str("uid", uid).Msg("getting device from sql")
	row := db.getDeviceByUID.QueryRow(uid)

	serverDevice, err := scanDevice(row)
	if err != nil {
		return nil, err
	}

	return serverDevice, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanDevice(s rowScanner) (*messages.PowerDevice, error) {
	var (
		name        sql.NullString
		uid         sql.NullString
		groupID     sql.NullInt64
		companyID   sql.NullInt64
		voltage     sql.NullFloat64
		powerFactor sql.NullFloat64
	)

	err := s.Scan(&name, &uid, &groupID, &companyID, &powerFactor, &voltage)
	if err != nil {
		return nil, err
	}

	deviceInfo := messages.PowerDevice{
		Device: &lib.Device{
			DeviceUID:  uid.String,
			DeviceName: name.String,
			CompanyID:  companyID.Int64,
			Group: &lib.Group{
				GroupID: groupID.Int64,
			},
		},
		PowerFactor: powerFactor.Float64,
		Voltage:     voltage.Float64,
	}

	return &deviceInfo, nil
}
