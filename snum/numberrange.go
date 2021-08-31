package snum

import (
	"fmt"
	"sync"
)

var numberRanges map[string]*numberRange
var keys []*string

type numberRange struct {
	name       *string
	mutex      *sync.Mutex
	intervals  map[string]*interval
	intvKeys   []*string
	intvValues []*interval
}

func (nr *numberRange) CreateInterval(name string, numberFrom uint32, numberTo uint32, isRolling bool) (INumberRangeInterval, error) {
	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	if _, found := nr.intervals[name]; found {
		return nil, fmt.Errorf("Interval \"%v\" does already exist", name)
	} else if numberFrom == 0 {
		return nil, fmt.Errorf("The begin of the Interval must be greater then 0")
	} else if numberFrom > numberTo {
		return nil, fmt.Errorf("Interval begin (%v) is greater then end (%v)", numberFrom, numberTo)
	} else { //check intervals
		for _, intv := range nr.intvValues {
			if !(numberFrom < *intv.numberFrom && numberTo < *intv.numberFrom ||
				numberFrom > *intv.numberTo) {
				return nil, fmt.Errorf("Interval \"%v\" (From: %v, To: %v) is overlapping", *intv.name, *intv.numberFrom, *intv.numberTo)
			}
		}
	}

	nr.intervals[name] = &interval{&name, &numberFrom, &numberTo, new(uint32), &isRolling, nr.name}
	nr.intvKeys = append(nr.intvKeys, &name)
	nr.intvValues = append(nr.intvValues, nr.intervals[name])
	return nr.intervals[name], nil
}

func (nr *numberRange) GetInterval(name string) (INumberRangeInterval, error) {
	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	intv, ok := nr.intervals[name]
	if !ok {
		return nil, fmt.Errorf("Interval \"%v\" does not exist", name)
	}
	return intv, nil
}

func (nr *numberRange) GetIntervals() (retIntervals []INumberRangeInterval) {
	for _, intv := range nr.intervals {
		retIntervals = append(retIntervals, intv)
	}
	return
}

func (nr *numberRange) GetName() *string {
	return nr.name
}

func (nr *numberRange) Validate() error {
	_, found := numberRanges[*nr.GetName()]
	if found {
		return fmt.Errorf("Number Range \"%v\" does already exist", *nr.GetName())
	}

	err := checkName(*nr.GetName())
	return err
}
