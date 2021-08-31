package snum

import "fmt"

type interval struct {
	name         *string
	numberFrom   *uint32
	numberTo     *uint32
	numberActual *uint32
	isRolling    *bool
	rangeName    *string
}

func (intv *interval) GetNext() (*uint32, error) {
	numberRanges[*intv.rangeName].mutex.Lock()
	defer numberRanges[*intv.rangeName].mutex.Unlock()

	if *intv.numberActual == 0 { //first use
		*intv.numberActual = *intv.numberFrom
	} else if *intv.numberActual == *intv.numberTo {
		if *intv.isRolling {
			*intv.numberActual = *intv.numberFrom
		} else {
			return nil, fmt.Errorf("Interval \"%v\" has reached the end (%d)", *intv.name, *intv.numberTo)
		}
	} else {
		*intv.numberActual++
	}
	return intv.numberActual, nil

}

func (intv *interval) GetName() *string {
	return intv.name
}

func (intv *interval) Validate() error {
	var err error

	nr, found := numberRanges[*intv.GetName()]
	if found {
		return fmt.Errorf("Number Range \"%v\" does already exist", *intv.GetName())
	}

	err = checkName(*nr.GetName())
	return err
}
