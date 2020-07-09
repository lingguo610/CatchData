package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"	
	"os/exec"
	"log"
	"bytes"
)

func pathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
		fmt.Println("111")
        return true, nil
    }
    if os.IsNotExist(err) {
		fmt.Println("222")
        return false, nil
	}
	fmt.Println("333")
    return false, err
}

/*功能说明：在dir目录下，创建一个以"ip地址/当前时间" 为名字的文件夹，并返回这个文件夹的全路径*/
func makeDir(dir, url string) string{

	ip := getUrlIp(url)
	newdir := dir + "/"+ip+"/"
	if _, err := pathExists(newdir); err == nil{ //目录不存在，创建目录
		os.Mkdir(newdir, os.ModePerm)
	}

	tmp := time.Now().Format("2006-01-02-15-04-05")
	newdir = newdir+"/"+tmp+"/"
	fmt.Println(newdir)
	
	exist := false
	if !exist{
		err := os.Mkdir(newdir, os.ModePerm)
        if err != nil {
            fmt.Printf("mkdir failed![%v]\n", err)
        } else {
			fmt.Printf("mkdir success!\n")
			return newdir
        }
	}
	return dir
} 

/*功能说明：修改下url，其实就是换下端口*/
func modifyPort(url ,port string) string{
	index := strings.LastIndex(url, ":")
	myurl := url[0 : index+1]
	return myurl + port
}

/*功能说明：获取url的ip部分，并且写成“10-161-42-33-” 这样的格式，是为了给文件命令*/
func getUrlIp(url  string)string{
	index := strings.LastIndex(url, ":")
	tmp := url[0 : index]
	index = strings.Index(tmp, "//")
	tmp = tmp[index +2:]
	tmp = strings.ReplaceAll(tmp, ".", "-")
	
	return tmp
}

/*功能说明：组装pprof的cmdline*/
func assemblePProfCmdLine(url, newdir, name, seconds string, port int)string{
	newurl := modifyPort(url, strconv.Itoa(port))
	cmdline := "curl " + newurl+ "/debug/pprof/profile?seconds="+seconds+" -p > " + newdir +"/profile_cpu_"+ name + "_"+seconds+"s.out\n"
	return cmdline
}

/*功能说明：组装pprof的cmdline*/
func assembleTraceCmdLine(url, newdir, name, seconds string, port int)string{
	newurl := modifyPort(url, strconv.Itoa(port))
	cmdline := "curl " + newurl+ "/debug/pprof/trace?seconds="+seconds+" -p > " + newdir +"/trace_"+ name+ "_"+seconds+"s.out\n"
	return cmdline
}

func assembleHeapCmdLine(url, newdir, name, seconds string, port int)string{
	newurl := modifyPort(url, strconv.Itoa(port))
	cmdline := "curl " + newurl+ "/debug/pprof/heap?seconds="+seconds+" -p > " + newdir +"/profile_mem_"+ name+ "_"+seconds+"s.out\n"
	return cmdline
}


func excuteCmd(cmdline string){
	cmd := exec.Command("cmd")
	in := bytes.NewBuffer(nil)
	cmd.Stdin = in //绑定输入
	var out bytes.Buffer
	cmd.Stdout = &out //绑定输出
	
	in.WriteString(cmdline) //写入你的命令，可以有多行，"\n"表示回车
	
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(cmd.Args)
	err = cmd.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
	rt := out.String() 
	fmt.Println(rt)
	return
}