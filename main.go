package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"code.aliyun.com/tongueplus_true_service/dlock/etcd"
	"github.com/coreos/etcd/clientv3"
)

type countStruct struct {
	a  int    //第一个操作数
	op string //运算符
	b  int    //第二个操作数
	r  int    //运算结果
}

func main() {
	sysc := make(chan os.Signal, 0)
	signal.Notify(sysc)
	go func() {
		sig := <-sysc
		os.Exit(0)
		fmt.Println(sig)
	}()
	// bar := pb.StartNew(int(0))
	// bar.Prefix("计算次数:")
	// bar.ShowBar = true
	// bar.ShowSpeed = true
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{"127.0.0.1:2379"}})
	if err != nil {
		log.Fatalln(err)
	}
	for {
		myMutex, mutexErr := etcd.GetMutex(cli, "testkey", 1)
		if mutexErr != nil {
			log.Fatalln(mutexErr)
		}
		if myMutex.Lock() == false {
			log.Fatalln(" lock failed")
		}
		// inputinterface, _ := cli.Get(cli.Ctx(), "countNumber")
		// input := 0
		// if len(inputinterface.Kvs) > 0 {
		// 	inputbytes := inputinterface.Kvs[0].Value
		// 	input, _ = strconv.Atoi(string(inputbytes))
		// }
		// if result, err := mathCount(input); err != nil {
		// 	log.Fatalln(err)
		// } else {
		// 	cli.Put(cli.Ctx(), "countNumber", strconv.Itoa(result))
		// }
		myMutex.UnLock()
		//	bar.Add64(1)
	}

}
func mathCount(input int) (int, error) {
	cs := new(countStruct)
	cs.a = input
	cs.op = "+"
	cs.b = 1
	// fmt.Scanln(&cs.op)
	// fmt.Scanln(&cs.b)
	switch cs.op {
	case "+":
		cs.r = cs.a + cs.b
	case "-":
		cs.r = cs.a - cs.b
	case "*":
		cs.r = cs.a * cs.b
	case "/":
		cs.r = cs.a / cs.b
	default:
		return input, fmt.Errorf("运算符只支持加减乘除")

	}
	return cs.r, nil
}
