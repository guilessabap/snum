package snum

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
)

type INumberRange interface {
	CreateInterval(name string, numberFrom uint32, numberTo uint32, isRolling bool) (INumberRangeInterval, error)
	GetInterval(name string) (INumberRangeInterval, error)
	GetIntervals() []INumberRangeInterval
	GetName() string
}

type INumberRangeInterval interface {
	GetNext() (uint32, error)
	GetName() string
}

func CreateNumberRange(name string) (INumberRange, error) {
	nr := &numberRange{&name, &sync.Mutex{}, make(map[string]*interval), make([]*string, 0), make([]*interval, 0)}
	if err := nr.Validate(); err != nil {
		return nil, err
	}
	numberRanges[name] = nr
	keys = append(keys, &name)
	return numberRanges[name], nil
}

func GetNumberRange(name string) (INumberRange, error) {
	nr, ok := numberRanges[name]
	if !ok {
		return nil, fmt.Errorf("Number Range \"%v\" does not exist", name)
	}
	return nr, nil
}

func ToTable(output io.Writer) {
	writer := tabwriter.NewWriter(output, 1, 1, 1, ' ', 0)
	var underlined []string

	ns := []string{
		"Number Range",
		"Interval",
		"From",
		"To",
		"Actual",
		"Rolling",
	}

	for _, s := range ns {
		underlined = append(underlined, strings.Repeat("-", len(s)))
	}

	fmt.Fprintf(writer, strings.Join(ns, "\t")+"\n"+strings.Join(underlined, "\t")+"\n")

	sort.Slice(keys, func(i, j int) bool {
		return *keys[i] < *keys[j]
	})
	for _, k_nr := range keys {
		sort.Slice(numberRanges[*k_nr].intvKeys, func(i, j int) bool {
			return *numberRanges[*k_nr].intvKeys[i] < *numberRanges[*k_nr].intvKeys[j]
		})
		for _, k_intv := range numberRanges[*k_nr].intvKeys {
			intv := numberRanges[*k_nr].intervals[*k_intv]
			fmt.Fprintf(
				writer,
				"%v\t%v\t%v\t%v\t%v\t%v\n",
				*numberRanges[*k_nr].name, *intv.name, *intv.numberFrom, *intv.numberTo, *intv.numberActual, *intv.isRolling)
		}
	}
	writer.Flush()
}
