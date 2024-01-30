package storage

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestDb_Search(t *testing.T) {
	filename := "celestia.json"
	n := 2000000
	err := generateHugeFile(filename, n)
	require.NoError(t, err)
	defer os.Remove(filename)

	db, err := Connect()
	require.NoError(t, err)
	defer db.Close()

	db.LoadFromJsonFile(filename)
	t.Run("simple search", func(t *testing.T) {
		res, err := db.Search("orbit_type")
		require.NoError(t, err)
		require.Equal(t, n, len(res))
		res, err = db.Search("orbit_type   ")
		require.NoError(t, err)
		require.Equal(t, n, len(res))
	})

	t.Run("diapasone search", func(t *testing.T) {
		res, err := db.Search("120 120,2")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res), 1)
		res, err = db.Search("    120    120.2   ")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res), 1)
	})

	t.Run("field search", func(t *testing.T) {
		res, err := db.Search("ref req")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res), 1)
		res, err = db.Search("   Epoch_Year   5000")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res), 1)
	})

	t.Run("field diapasone search", func(t *testing.T) {
		res, err := db.Search("E 0.8 0.7  ")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res), 1)
		res, err = db.Search("   E 0.8 0.9  ")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res), 1)
		res, err = db.Search("   E 0.8 0.8  ")
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(res), 1)
	})

	t.Run("quoted search", func(t *testing.T) {
		res, err := db.Search("\"ref res\"")
		require.NoError(t, err)
		require.Equal(t, 1, len(res))
		res, err = db.Search("\"ref  res \"")
		require.NoError(t, err)
		require.Equal(t, 0, len(res))
	})

	t.Run("field quoted search", func(t *testing.T) {
		res, err := db.Search("ref \"ref res\"")
		require.NoError(t, err)
		require.Equal(t, 1, len(res))
		res, err = db.Search("E \"ref  res \"")
		require.NoError(t, err)
		require.Equal(t, 0, len(res))
	})

	t.Run("wrong search", func(t *testing.T) {
		res, err := db.Search("")
		require.Error(t, err)
		require.Equal(t, 0, len(res))
		res, err = db.Search("\"")
		require.Error(t, err)
		require.Equal(t, 0, len(res))
		res, err = db.Search("a b c d")
		require.Error(t, err)
		require.Equal(t, 0, len(res))
		res, err = db.Search("a b c")
		require.Error(t, err)
		require.Equal(t, 0, len(res))
		res, err = db.Search("a b")
		require.Error(t, err)
		require.Equal(t, 0, len(res))
	})

}

func generateHugeFile(filename string, count int) error {
	fmt.Println("Celesties JSON file generation started.")

	rand.Seed(time.Now().UnixNano())

	celesties := generateRandomCelesties(count)
	err := writeCelestiesToFile(celesties, filename)
	if err != nil {
		return err
	}

	fmt.Println("Celesties JSON file has been generated.")
	return nil
}

func generateRandomCelesties(n int) []Celesty {
	celesties := make([]Celesty, n)
	for i := 0; i < n; i++ {
		celesties[i] = Celesty{
			OrbitType:              "orbit_type",
			ProvisionalPackedDesig: "provisional_packed_desig",
			YearOfPerihelion:       rand.Intn(100),
			MonthOfPerihelion:      rand.Intn(12) + 1,
			DayOfPerihelion:        rand.Float64(),
			PerihelionDist:         rand.Float64(),
			E:                      rand.Float64(),
			Peri:                   rand.Float64(),
			Node:                   rand.Float64(),
			I:                      rand.Float64(),
			EpochYear:              rand.Intn(3000) - 1000,
			EpochMonth:             rand.Intn(12) + 1,
			EpochDay:               rand.Intn(31) + 1,
			H:                      rand.Float64(),
			G:                      rand.Float64(),
			DesignationAndName:     "designation_and_name",
			Ref:                    "ref",
		}
		if i == 1 {
			celesties[i].E = 0.8
		}
		if i == 2 {
			celesties[i].Peri = 120.1
		}
		if i == 3 {
			celesties[i].Ref = "req ref"
		}
		if i == 4 {
			celesties[i].Ref = "ref res"
		}
		if i == 5 {
			celesties[i].EpochYear = 5000
		}
	}
	return celesties
}

func writeCelestiesToFile(celesties []Celesty, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(celesties)
	if err != nil {
		return err
	}

	return nil
}
