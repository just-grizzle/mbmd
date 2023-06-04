package rs485

import . "github.com/volkszaehler/mbmd/meters"

func init() {
	Register("PQPLUS", NewPQPlusProducer)
}

type PQPlusProducer struct {
	Opcodes
}

func NewPQPlusProducer() Producer {
	/**
	 * Opcodes for PQ Plus Modbus devices (Meters and power quality analyzers).
	 * See https://www.pq-plus.de/site/assets/files/2796/modbus_manual_v4_0_de_rev1_2.pdf
	 */
	ops := Opcodes{
		VoltageL1:       0x1100, // 32b float
		VoltageL2:       0x1102, // 32b float
		VoltageL3:       0x1104, // 32b float
		CurrentL1:       0x1200, // 32b float
		CurrentL2:       0x1202, // 32b float
		CurrentL3:       0x1204, // 32b float
		Power:           0x1314, // 32b float
		PowerL1:         0x1320, // 32b float
		PowerL2:         0x1322, // 32b float
		PowerL3:         0x1324, // 32b float
		ReactivePower:   0x1316, // 32b float
		ReactivePowerL1: 0x1328, // 32b float
		ReactivePowerL2: 0x132A, // 32b float
		ReactivePowerL3: 0x132C, // 32b float
		ApparentPower:   0x1318, // 32b float
		ApparentPowerL1: 0x1330, // 32b float
		ApparentPowerL2: 0x1332, // 32b float
		ApparentPowerL3: 0x1334, // 32b float
		Import:          0x2000, // 64b double
		ImportL1:        0x2010, // 64b double
		ImportL2:        0x2014, // 64b double
		ImportL3:        0x2018, // 64b double
		Export:          0x2004, // 64b double
		ExportL1:        0x2020, // 64b double
		ExportL2:        0x2024, // 64b double
		ExportL3:        0x2028, // 64b double
		CosphiL1:        0x130C, // 32b
		CosphiL2:        0x130E, // 32b
		CosphiL3:        0x1310, // 32b
		Frequency:       0x1004, // 32b float
	}
	return &PQPlusProducer{Opcodes: ops}
}

// Description implements Producer interface
func (p *PQPlusProducer) Description() string {
	return "PQ Plus meters and power quality analyzers"
}

// snip creates modbus operation
func (p *PQPlusProducer) snip16(iec Measurement, readlen uint16) Operation {
	return Operation{
		FuncCode: ReadHoldingReg,
		OpCode:   p.Opcode(iec), // adjust according to docs
		ReadLen:  readlen,
		IEC61850: iec,
	}
}

// snip32 creates modbus operation for double register
func (p *PQPlusProducer) snip32(iec Measurement, scaler ...float64) Operation {
	snip := p.snip16(iec, 2)

	snip.Transform = RTUIeee754ToFloat64 // default conversion
	if len(scaler) > 0 {
		snip.Transform = MakeScaledTransform(snip.Transform, scaler[0])
	}

	return snip
}

// Probe implements Producer interface
func (p *PQPlusProducer) Probe() Operation {
	return p.snip32(VoltageL1)
}

// Produce implements Producer interface
func (p *PQPlusProducer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
		CurrentL1, CurrentL2, CurrentL3,
		Power, PowerL1, PowerL2, PowerL3,
		ReactivePower, ReactivePowerL1, ReactivePowerL2, ReactivePowerL3,
		ApparentPower, ApparentPowerL1, ApparentPowerL2, ApparentPowerL3,
	} {
		res = append(res, p.snip32(op))
	}

	// for _, op := range []Measurement{
	// 	CurrentL1, CurrentL2, CurrentL3,
	// } {
	// 	res = append(res, p.snip32(op, 1000))
	// }
	return res
}
