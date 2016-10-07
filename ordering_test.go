package colog

import (
	"fmt"
)

func checkadd(e *Entry, err error) {
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ExampleItemsOrdering() {
	l1 := New("l1", ipfsdb)
	l2 := New("l2", ipfsdb)
	l3 := New("l3", ipfsdb)

	checkadd(l1.Add("init"))
	check(l2.Join(l1))

	checkadd(l1.Add("simple fix"))
	checkadd(l2.Add("fix simple bug"))
	checkadd(l2.Add("oh i missed something"))

	check(l1.Join(l2))
	check(l2.Join(l1))

	checkadd(l1.Add("bump version 0.2"))

	check(l2.Join(l1))

	checkadd(l2.Add("add feature"))

	check(l1.Join(l2))
	check(l3.Join(l2))

	checkadd(l3.Add("add chinese translation"))
	checkadd(l1.Add("add spanish translation"))
	checkadd(l2.Add("add afrikaans translation"))

	check(l1.Join(l3))
	check(l3.Join(l1))

	checkadd(l3.Add("fix i18n api"))
	checkadd(l3.Add("add i18n feature"))

	check(l1.Join(l2))
	check(l1.Join(l3))

	checkadd(l1.Add("version bump 0.3"))

	for _, e := range l1.Items() {
		fmt.Println(e.Hash, e.GetString(), e.Prev)
	}

	//Output:
	//QmT1yirz1BgNjWQJgefjz37UoZsc4aa5DDAdVFKcWQek3L init { null }
	//QmYhvpaYk2MGwVGpHwM6gGS5GQu1esCBeueQBbpSXNPRRa fix simple bug { QmT1yirz1BgNjWQJgefjz37UoZsc4aa5DDAdVFKcWQek3L }
	//QmcoQ6WHtjrFi6x1JB6VBLDEYzfgxaM9mYwXRJ953M2SJ8 oh i missed something { QmYhvpaYk2MGwVGpHwM6gGS5GQu1esCBeueQBbpSXNPRRa }
	//QmYWcu4VkcQ7vvugSqJuwsVLwLijfhDj7UGyae2a8V1xY3 simple fix { QmT1yirz1BgNjWQJgefjz37UoZsc4aa5DDAdVFKcWQek3L }
	//QmZtEaT7KTd1w2uzL46xMscjXMgwxupgy1chpPKHvdSbzC bump version 0.2 { QmYWcu4VkcQ7vvugSqJuwsVLwLijfhDj7UGyae2a8V1xY3, QmcoQ6WHtjrFi6x1JB6VBLDEYzfgxaM9mYwXRJ953M2SJ8 }
	//QmYJStDEEPBWRU85G5N55585LroCaf1axK1ynKLR9gcaFh add feature { QmZtEaT7KTd1w2uzL46xMscjXMgwxupgy1chpPKHvdSbzC }
	//Qme9Dqbu98VsdvmJvbpnbcKeUkh4d332BFxiCMJFFEU1Dj add spanish translation { QmYJStDEEPBWRU85G5N55585LroCaf1axK1ynKLR9gcaFh }
	//Qmcuv6bvnVbahziieBb28oR8E6TpgRWyX3vzWu4MgWyVHa add chinese translation { QmYJStDEEPBWRU85G5N55585LroCaf1axK1ynKLR9gcaFh }
	//QmRrX4qyaZX8wMTkjfHLbT8erD6Hz6suzZvDA4LbTcL646 fix i18n api { Qmcuv6bvnVbahziieBb28oR8E6TpgRWyX3vzWu4MgWyVHa, Qme9Dqbu98VsdvmJvbpnbcKeUkh4d332BFxiCMJFFEU1Dj }
	//QmSRozfJqeznzV7j1tL8LbTyPGHRAhFoqaBGsmzbRyReqh add i18n feature { QmRrX4qyaZX8wMTkjfHLbT8erD6Hz6suzZvDA4LbTcL646 }
	//QmRBzEJQ3eDaKVh9dM6XhrQ1GgMSRm6uge8M7itkHUhorE add afrikaans translation { QmYJStDEEPBWRU85G5N55585LroCaf1axK1ynKLR9gcaFh }
	//QmPse7WVT7ftSpvcMPXCfTPjehwigDaSMsaePpgtjANzzy version bump 0.3 { QmRBzEJQ3eDaKVh9dM6XhrQ1GgMSRm6uge8M7itkHUhorE, QmSRozfJqeznzV7j1tL8LbTyPGHRAhFoqaBGsmzbRyReqh }
}
