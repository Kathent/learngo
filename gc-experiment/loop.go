package gc_experiment

import (
	"math/rand"
	"fmt"
	"os"
	"runtime/pprof"
	"os/signal"
	"syscall"
)

func recrusiveLoop(val int, valArr []int) {
	// fmt.Println(fmt.Sprintf("enter loop, val:%d", val))
	if val > 0 {
		tmpArr := make([]int, 0)
		for i := 0; i < 100; i++ {
			ttmpArr := make([]int, 0)
			for j := 0; j < 100; j++ {
				ttmpArr = append(ttmpArr, rand.Intn(10))
			}

			tmpArr = append(tmpArr, ttmpArr...)
		}

		valArr = append(valArr, tmpArr...)
		recrusiveLoop(val - 1, valArr)
	}else  {
		sum := 0
		for i := 0; i < len(valArr); i++ {
			sum += valArr[i]
		}

		valArr = append(valArr, sum)
	}
}

func GcLoop() {
	go sinalProcess()
	f, _ := os.Create("profile_file")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	arr := make([]int, 0)
	recrusiveLoop(10000, arr)

	m, _ := os.Create("profile.mprof")
	pprof.WriteHeapProfile(m)

	if len(arr) > 0 {
		fmt.Println(fmt.Sprintf("after loop sum is:%d", arr[len(arr) - 1]))
	}else {
		fmt.Println("arr len is 0...")
	}
}

func sinalProcess() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT,
		syscall.SIGTERM)
	for s := range signalChan {
		fmt.Println(s.String())
	}
}
