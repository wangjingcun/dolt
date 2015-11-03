package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"sync"

	"github.com/attic-labs/noms/clients/util"
	"github.com/attic-labs/noms/d"
	"github.com/attic-labs/noms/datas"
	"github.com/attic-labs/noms/dataset"
	"github.com/attic-labs/noms/types"
)

var (
	inputID  = flag.String("input-ds", "", "dataset to pull data from.")
	outputID = flag.String("output-ds", "", "dataset to store data in.")
)

func processPitcher(m MapOfStringToString) (id, name string) {
	id = m.Get("-id")
	name = fmt.Sprintf("%s, %s", m.Get("-last_name"), m.Get("-first_name"))
	return
}

func checkPitch(v MapOfStringToString) bool {
	return v.Has("-px") && v.Has("-pz")
}

func getPitch(v MapOfStringToString) Pitch {
	x, _ := strconv.ParseFloat(v.Get("-px"), 64)
	z, _ := strconv.ParseFloat(v.Get("-pz"), 64)
	return NewPitch().SetX(x).SetZ(z)
}

func processPitches(v types.Value) (pitches []Pitch) {
	switch v := v.(type) {
	case types.List:
		for i := uint64(0); i < v.Len(); i++ {
			pitches = append(pitches, processPitches(v.Get(i))...)
		}
	case MapOfStringToString:
		if checkPitch(v) {
			pitches = append(pitches, getPitch(v))
		}
	case nil:
		return // Yes, an at-bat can end with no pitches thrown.
	default:
		d.Chk.Fail("No pitch should be %+v, which is of type %s!\n", v, reflect.TypeOf(v).String())
	}
	return
}

func processInning(m MapOfStringToValue) map[string][]Pitch {
	pitchCounts := map[string][]Pitch{}

	// This is brittle, figure out how to do it without being super verbose.
	top := m.Get("top")
	if _, ok := top.(types.Map); !ok {
		// If "top" is anything other than a map, give up
		return pitchCounts
	}

	halves := []MapOfStringToValue{
		top.(MapOfStringToValue),
	}

	if bot := m.Get("bottom"); bot != nil {
		halves = append(halves, bot.(MapOfStringToValue))
	}

	addPitch := func(ab MapOfStringToValue) {
		pitchData := ab.Get("pitch")
		pitcher := ab.Get("-pitcher").(types.String).String()
		pitchCounts[pitcher] = append(pitchCounts[pitcher], processPitches(pitchData)...)
	}

	for _, half := range halves {
		atbat := half.Get("atbat")
		switch atbat := atbat.(type) {
		case ListOfMapOfStringToValue:
			for i := uint64(0); i < atbat.Len(); i++ {
				ab := atbat.Get(i)
				addPitch(ab)
			}
		case MapOfStringToValue:
			// Apparently, if there's only one, it's encoded directly as a singleton. Yay, data!
			addPitch(atbat)
		default:
		}
	}
	return pitchCounts
}

func getIndex(input types.List) MapOfStringToListOfPitch {
	mu := sync.Mutex{}
	pitchers := NewMapOfStringToString()

	// Walk through the list in inputDataset and basically switch
	// on the top-level key to know if it's an inning or a pitcher.
	innings := input.MapP(512, func(item types.Value, i uint64) interface{} {
		m := item.(MapOfStringToValue)

		if key := "inning"; m.Has(key) {
			return processInning(m.Get(key).(MapOfStringToValue))
		}

		if key := "Player"; m.Has(key) {
			id, name := processPitcher(m.Get(key).(MapOfStringToString))
			if id != "" && name != "" {
				mu.Lock()
				pitchers = pitchers.Set(id, name)
				mu.Unlock()
			}
		}

		return nil
	})

	pitchCounts := NewMapOfStringToListOfPitch()
	for _, inning := range innings {
		if inning == nil {
			continue
		}

		for id, p := range inning.(map[string][]Pitch) {
			pitches := NewListOfPitch()
			if pitchCounts.Has(id) {
				pitches = pitchCounts.Get(id)
			}
			pitchCounts = pitchCounts.Set(id, pitches.Append(p...))
		}
	}

	namedPitchCounts := NewMapOfStringToListOfPitch()
	pitchCounts.Iter(func(id string, p ListOfPitch) (stop bool) {
		if pitchers.Has(id) {
			namedPitchCounts = namedPitchCounts.Set(pitchers.Get(id), p)
		} else {
			d.Chk.Fail("Unknown pitcher!", id)
		}
		return
	})

	return namedPitchCounts
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flags := datas.NewFlags()
	flag.Parse()

	ds, ok := flags.CreateDataStore()
	if !ok || *inputID == "" || *outputID == "" {
		flag.Usage()
		return
	}
	defer ds.Close()

	err := d.Try(func() {
		if util.MaybeStartCPUProfile() {
			defer util.StopCPUProfile()
		}
		inputDataset := dataset.NewDataset(ds, *inputID)
		outputDataset := dataset.NewDataset(ds, *outputID)

		input := inputDataset.Head().Value().(types.List)
		output := getIndex(input)

		_, ok := outputDataset.Commit(output)
		d.Exp.True(ok, "Could not commit due to conflicting edit")

		util.MaybeWriteMemProfile()
	})
	if err != nil {
		log.Fatal(err)
	}
}
