package eventbus

import (
	"testing"
	"time"
)

const (
	eventName = "test"
)

func benchmarkBusIO(listeners int, b *testing.B) {

	bus := New()

	for i := 0; i < listeners; i++ {
		bus.On(eventName, func() {
			time.Sleep(time.Second * 1)
		})
	}

	for n := 0; n < b.N; n++ {
		bus.Emit(eventName)
	}

}

func BenchmarkBusIO_1(b *testing.B)     { benchmarkBusIO(1, b) }
func BenchmarkBusIO_2(b *testing.B)     { benchmarkBusIO(2, b) }
func BenchmarkBusIO_3(b *testing.B)     { benchmarkBusIO(3, b) }
func BenchmarkBusIO_10(b *testing.B)    { benchmarkBusIO(10, b) }
func BenchmarkBusIO_20(b *testing.B)    { benchmarkBusIO(20, b) }
func BenchmarkBusIO_80(b *testing.B)    { benchmarkBusIO(80, b) }
func BenchmarkBusIO_160(b *testing.B)   { benchmarkBusIO(160, b) }
func BenchmarkBusIO_250(b *testing.B)   { benchmarkBusIO(320, b) }
func BenchmarkBusIO_500(b *testing.B)   { benchmarkBusIO(640, b) }
func BenchmarkBusIO_1000(b *testing.B)  { benchmarkBusIO(1280, b) }
func BenchmarkBusIO_2500(b *testing.B)  { benchmarkBusIO(2560, b) }
func BenchmarkBusIO_10000(b *testing.B) { benchmarkBusIO(1000, b) }

func benchmarkBusCPU(listeners int, b *testing.B) {

	bus := New()

	for i := 0; i < listeners; i++ {
		bus.On(eventName, func(a, b int) {
			a = a * b
		})
	}
	waitAll := make([]WaitCallback, b.N)

	for n := 0; n < b.N; n++ {
		waitAll[n] = bus.Emit(eventName, 1, 2)
	}

	for i := 0; i < b.N; i++ {
		waitAll[i]()
	}

}

func BenchmarkBusCPU_1(b *testing.B)     { benchmarkBusCPU(1, b) }
func BenchmarkBusCPU_2(b *testing.B)     { benchmarkBusCPU(2, b) }
func BenchmarkBusCPU_3(b *testing.B)     { benchmarkBusCPU(3, b) }
func BenchmarkBusCPU_10(b *testing.B)    { benchmarkBusCPU(10, b) }
func BenchmarkBusCPU_20(b *testing.B)    { benchmarkBusCPU(20, b) }
func BenchmarkBusCPU_80(b *testing.B)    { benchmarkBusCPU(80, b) }
func BenchmarkBusCPU_160(b *testing.B)   { benchmarkBusCPU(160, b) }
func BenchmarkBusCPU_250(b *testing.B)   { benchmarkBusCPU(320, b) }
func BenchmarkBusCPU_500(b *testing.B)   { benchmarkBusCPU(640, b) }
func BenchmarkBusCPU_1000(b *testing.B)  { benchmarkBusCPU(1280, b) }
func BenchmarkBusCPU_2500(b *testing.B)  { benchmarkBusCPU(2560, b) }
func BenchmarkBusCPU_10000(b *testing.B) { benchmarkBusCPU(1000, b) }
