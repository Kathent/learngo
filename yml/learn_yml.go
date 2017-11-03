package yml

import (
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v1"
	"fmt"
)

type AppConfig struct {
	//Redis相关配置
	Db struct {
		Addr     string
		Password string
		User     string
		Database int
	}

	//排队系统相关配置
	Server struct {
		NetWork string
		//本地rpc server地址
		Addr string
		//排队系统调用地址
		RemoteAddr string
	}

	//GRpc server端配置
	RpcServer struct {
		Protocol      string
		Port          string
		InterfaceName string
	}

	//上传文件配置
	OSS struct {
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
		Action          string
	}

	//获取信息配置
	Fetch struct {
		FetchUrl       string
		FetchWeChatUrl string
	}

	//rabbitMQ配置信息
	RabbitMQ struct {
		Addr string
	}

	//调试配置
	DebugConfig struct {
		Debug bool
	}

	Etcd struct {
		Addr []string
	}

	Logger struct {
		Level string
	}

	App struct{
		AppKey string
	}

	Http struct {
		Addr string
	}

	Switch struct{
		AgentLimt string
	}
}

func LeanYml(){
	file, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	tmp := AppConfig{}
	err = yaml.Unmarshal(bytes, &tmp)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", tmp))



	file, err = os.Open("glide.yaml")
	if err != nil {
		panic(err)
	}

	bytes, err = ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	type Package struct {
		Package string
		Version string
		SubPackages []string
	}
	type GlideConfig struct {
		Package string
		Import []Package
	}

	ttmp := GlideConfig{}
	err = yaml.Unmarshal(bytes, &ttmp)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v", ttmp))
}
