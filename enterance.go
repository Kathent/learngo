package main

//import "learngo/learn_etcd"

func main() {
	//q := queue.New()
	//fmt.Println(q)

	//go rpc.AcceptRpc()
	//go rpc.SendRpcClient()
	//go func() {
	//learn_etcd.LearnEtcd()
		//time.Sleep(time.Second * 5)
		//go_micro_learn.StartClient()
		//go_micro_learn.StartGinClient()
	//}()
	//go_micro_learn.StartGinServer()
	//go_micro_learn.StartServer()
	//go go_micro_learn.StartGrpcServer()
	//time.Sleep(time.Second * 4)
	//go_micro_learn.StartGrcpClient()
	//file_read_analysis.LoadFiles("load_file_dir", "2017-10-02 10:00:00", "2017-10-03 14:00:00")
	//mgo.LearnMgo()
	//mgo.TryMongoDial()
	//yml.LeanYml()
	//mgo.LearnDate()

	arr := make([]int, 3, 5)
	println(&arr)
	modify(arr)

	//print(arr[2])
	println(&arr)
}

func modify(ints []int) {
	ints[0] = 1
	ints[1] = 2
	ints[2] = 3
	println(&ints)

	ints = append(ints, 6)
	println(&ints)

	ints = append(ints, 6)
	ints = append(ints, 6)
	ints = append(ints, 6)
	ints = append(ints, 6)
	println(&ints)
}
