package main

import (
	"log"
	"nexus-wallet/factory"
)

func main() {
	f, err := factory.NewServiceFactory()
	if err != nil {
		println(err.Error())
		log.Fatalln(err)
	}

	cronKernel, onShutdown, err := f.CreateCronKernel()
	if err != nil {
		println(err.Error())
		log.Fatalln(err)
	}
	defer func() {
		err := onShutdown()
		if err != nil {
			println(err.Error())
			log.Fatalln(err)
		}
	}()

	cronKernel.Run()

	httpKernel, onShutdown, err := f.CreateHttpKernel()
	if err != nil {
		println(err.Error())
		log.Fatalln(err)
	}
	defer func() {
		err := onShutdown()
		if err != nil {
			println(err.Error())
			log.Fatalln(err)
		}
	}()

	err = httpKernel.Run()
	if err != nil {
		println(err.Error())
		log.Fatalln(err)
	}
}
