package main

import "testing"

func TestGetNextGenerationState(t *testing.T) {
	testCases := []struct {
		desc            string
		c               cellCharacter
		aliveNeighbours int
		want            cellCharacter
	}{
		{
			desc:            "alive - dies by underpopulation",
			c:               aliveCell,
			aliveNeighbours: 0,
			want:            deadCell,
		},
		{
			desc:            "alive - dies by underpopulation",
			c:               aliveCell,
			aliveNeighbours: 1,
			want:            deadCell,
		},
		{
			desc:            "alive - lives",
			c:               aliveCell,
			aliveNeighbours: 2,
			want:            aliveCell,
		},
		{
			desc:            "alive - lives",
			c:               aliveCell,
			aliveNeighbours: 3,
			want:            aliveCell,
		},
		{
			desc:            "alive - dies by overpopulation",
			c:               aliveCell,
			aliveNeighbours: 4,
			want:            deadCell,
		},
		{
			desc:            "alive - dies by overpopulation",
			c:               aliveCell,
			aliveNeighbours: 5,
			want:            deadCell,
		},
		{
			desc:            "alive - dies by overpopulation",
			c:               aliveCell,
			aliveNeighbours: 6,
			want:            deadCell,
		},
		{
			desc:            "alive - dies by overpopulation",
			c:               aliveCell,
			aliveNeighbours: 7,
			want:            deadCell,
		},
		{
			desc:            "alive - dies by overpopulation",
			c:               aliveCell,
			aliveNeighbours: 8,
			want:            deadCell,
		},

		// dead
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 0,
			want:            deadCell,
		},
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 1,
			want:            deadCell,
		},
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 2,
			want:            deadCell,
		},
		{
			desc:            "dead - becomes alive by reproduction",
			c:               deadCell,
			aliveNeighbours: 3,
			want:            aliveCell,
		},
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 4,
			want:            deadCell,
		},
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 5,
			want:            deadCell,
		},
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 6,
			want:            deadCell,
		},
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 7,
			want:            deadCell,
		},
		{
			desc:            "dead - stays dead",
			c:               deadCell,
			aliveNeighbours: 8,
			want:            deadCell,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := getNextGenerationState(tC.aliveNeighbours, tC.c)

			if got != tC.want {
				t.Errorf("want: %s, got: %s", string(tC.want), string(got))
			}
		})
	}
}
