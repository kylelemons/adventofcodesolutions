package aocday

import (
	"sort"
	"strings"
	"testing"
	"time"
)

func part1(t *testing.T, in []string) int {

	type event struct {
		ts    time.Time
		id    int
		event int // 0 - start, 1 - falls asleep, 2 - wakes up
	}
	const (
		ShiftStart  = 0
		FallsAsleep = 1
		WakesUp     = 2
	)

	parseStamp := func(s string) time.Time {
		ts, err := time.Parse("2006-01-02 15:04", s)
		if err != nil {
			t.Fatalf("failed to parse time: %s", err)
		}
		return ts
	}

	guardSleepTotal := make(map[int]int)    // [guard] = total minutes sleeping
	guardSleepingAt := make(map[[3]int]int) // [{guard,hour,minute}] = times asleep

	var events []event
	sort.Strings(in)
	for _, input := range in {
		if input == "" {
			continue
		}
		if stamp, id := "", 0; scanner(input).scan(t, `\[(....-..-.. ..:..)\] Guard #(\d+) begins shift`, &stamp, &id) {
			events = append(events, event{
				ts:    parseStamp(stamp),
				id:    id,
				event: ShiftStart,
			})
		} else if scanner(input).scan(t, `\[(....-..-.. ..:..)\] falls asleep`, &stamp) {
			events = append(events, event{
				ts:    parseStamp(stamp),
				id:    events[len(events)-1].id,
				event: FallsAsleep,
			})
		} else if scanner(input).scan(t, `\[(....-..-.. ..:..)\] wakes up`, &stamp) {
			last := events[len(events)-1]
			events = append(events, event{
				ts:    parseStamp(stamp),
				id:    events[len(events)-1].id,
				event: WakesUp,
			})
			curr := events[len(events)-1]

			id := last.id
			for tt := last.ts; tt.Before(curr.ts); tt = tt.Add(1 * time.Minute) {
				guardSleepTotal[id]++
				guardSleepingAt[[3]int{id, tt.Hour(), tt.Minute()}]++
			}
		} else {
			t.Errorf("Failed to parse %q", input)
		}
	}

	maxGuard, maxMinutes := -1, -1
	for id, mins := range guardSleepTotal {
		if mins <= maxMinutes {
			continue
		}
		maxGuard, maxMinutes = id, mins
	}
	maxMinute, maxCount := -1, -1
	for info, count := range guardSleepingAt {
		if info[0] != maxGuard || count <= maxCount {
			continue
		}
		maxMinute, maxCount = info[2], count
	}

	return maxGuard * maxMinute
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want int
	}{
		{"example1", []string{
			"[1518-11-01 00:00] Guard #10 begins shift",
			"[1518-11-01 00:05] falls asleep",
			"[1518-11-01 00:25] wakes up",
			"[1518-11-01 00:30] falls asleep",
			"[1518-11-01 00:55] wakes up",
			"[1518-11-01 23:58] Guard #99 begins shift",
			"[1518-11-02 00:40] falls asleep",
			"[1518-11-02 00:50] wakes up",
			"[1518-11-03 00:05] Guard #10 begins shift",
			"[1518-11-03 00:24] falls asleep",
			"[1518-11-03 00:29] wakes up",
			"[1518-11-04 00:02] Guard #99 begins shift",
			"[1518-11-04 00:36] falls asleep",
			"[1518-11-04 00:46] wakes up",
			"[1518-11-05 00:03] Guard #99 begins shift",
			"[1518-11-05 00:45] falls asleep",
			"[1518-11-05 00:55] wakes up",
		}, 10 * 24},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, want := part1(t, test.in), test.want; got != want {
				t.Errorf("part1(%#v) = %#v, want %#v", test.in, got, want)
			}
		})
	}
	if t.Failed() {
		return
	}

	t.Run("part1", func(t *testing.T) {
		t.Logf("part1: %v", part1(t, strings.Split(read(t, "input.txt"), "\n")))
	})
}
