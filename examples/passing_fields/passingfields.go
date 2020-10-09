package main

import (
	"context"
	"errors"
	"fmt"
	goworker "github.com/catmullet/go-workers"
	"math/rand"
)

type WorkerTwoConfig struct {
	AmountToMultiply int
}

func main() {
	ctx := context.Background()
	workerOne := goworker.NewWorker(ctx, workerFunctionOne, 10).
		AddField("amountToMultiply", 2).
		Work()
	workerTwo := goworker.NewWorker(ctx, workerFunctionTwo, 10).
		AddField("amountToMultiply", &WorkerTwoConfig{AmountToMultiply: 4}).
		InFrom(workerOne).
		Work()

	for i := 0; i < 10; i++ {
		workerOne.Send(rand.Intn(100))
	}

	workerOne.Close()
	if err := workerOne.Wait(); err != nil {
		fmt.Println(err)
	}

	workerTwo.Close()
	if err := workerTwo.Wait(); err != nil {
		fmt.Println(err)
	}
}

func workerFunctionOne(w *goworker.Worker) error {
	amountToMultiply, ok := w.GetFieldInt("amountToMultiply")
	if !ok {
		return errors.New("No amount to multiply supplied")
	}

	for in := range w.In() {
		total := in.(int) * amountToMultiply
		fmt.Println(fmt.Sprintf("%d * 2 = %d", in.(int), total))
		w.Out(total)
	}
	return nil
}

func workerFunctionTwo(w *goworker.Worker) error {
	var workerConfig WorkerTwoConfig
	w.GetFieldObject("amountToMultiply", &workerConfig)

	for in := range w.In() {
		totalFromWorkerOne := in.(int)
		fmt.Println(fmt.Sprintf("%d * 4 = %d", totalFromWorkerOne, totalFromWorkerOne*workerConfig.AmountToMultiply))
	}
	return nil
}