package atlas

import "machine"

const VoltAvgLength = 1000

type Device struct {
	adc machine.ADC

	calLow  float64
	calMid  float64
	calHigh float64
}

type DeviceConfig struct {
	Calibration struct {
		Low  float64
		Mid  float64
		High float64
	}
}

func New(adc machine.ADC) Device {
	return Device{
		adc: adc,
	}
}

func (d *Device) Configure(cfg DeviceConfig) {
	if cfg.Calibration.Low == 0 {
		cfg.Calibration.Low = 2030
	}
	if cfg.Calibration.Mid == 0 {
		cfg.Calibration.Mid = 1500
	}
	if cfg.Calibration.High == 0 {
		cfg.Calibration.High = 975
	}

	d.calLow = cfg.Calibration.Low
	d.calMid = cfg.Calibration.Mid
	d.calHigh = cfg.Calibration.High
}

func (d *Device) readVoltage() float64 {
	mv := 0.0
	for i := 1; i < VoltAvgLength; i++ {
		mv += float64(d.adc.Get()) / 1024.0 * 5000.0
	}
	mv /= VoltAvgLength
	return mv
}

func (d *Device) ReadPH() float64 {
	mv := d.readVoltage()
	if mv > d.calMid {
		return 7.0 - 3.0/(d.calLow-d.calMid)*(mv-d.calMid)
	} else {
		return 7.0 - 3.0/(d.calMid-d.calHigh)*(mv-d.calMid)
	}
}
