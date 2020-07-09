package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"compress/gzip"
	"strconv"
	
	"sync"
)

/*
功能说明：爬取url数据，并且分析http body的内容，
提取HBP_GO_DEBUG_PORT字段内容，并且展示到TableView中
*/
func grabWebData(url interface{}, model *FooModel){
	myurl := url.(string)

	client := &http.Client{}
	fullurl := myurl + "/api/v1/workload/default?filterBy=&itemsPerPage=52&name=&page=1&sortBy=d,creationTimestamp"
    req, _ := http.NewRequest("GET", fullurl, nil)
	// 自定义Header

	req.Header.Set("Host","10.161.42.53:32000")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	req.Header.Set("Referer", myurl)
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("http get error", err)
        return
	}
	defer resp.Body.Close()
	
	body := resp.Body
  	if resp.Header.Get("Content-Encoding") == "gzip" {
	    body, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("http resp unzip is failed,err: ", err)
	    }
    }

	data, err := ioutil.ReadAll(body)
    if err != nil {
        fmt.Println("read error", err)
        return
    }


	var httpBody HttpBody
	err = json.Unmarshal(data, &httpBody)
	if err != nil{
		fmt.Println("json unmarshal fail", err)
		return 
	}

	var items []*Foo
	items = make([]*Foo, 0)
	var index int

	for _,v := range httpBody.Deploy.Data{
		var kute KUBE
		
		err = json.Unmarshal([]byte(v.Meta.ANNO.KUBE), &kute)
		if err != nil{
			fmt.Println("json unmarshal fail", err)
			return 
		}
		//fmt.Println(kute)
		for _, v2 := range kute.SPEC1.Template.SPEC2.Container{
			for _, v3:= range v2.Env{
				if v3.Name == "HBP_GO_DEBUG_PORT"{	
					port, _ := strconv.Atoi(v3.Value)
					fool := &Foo{Index: index, Name:v.Meta.Name, Port: port}
					items = append(items,fool)
					index++
				}
			}
		}
	}

	model.items = items
	model.PublishRowsReset()
	model.Sort(model.sortColumn, model.sortOrder)
}

/*抓取pprof*/
func grapPProf(wg *sync.WaitGroup, url, newdir, name string, port int, seconds string){
	cmdline := assemblePProfCmdLine(url, newdir, name,seconds, port)
	fmt.Println(cmdline)
	excuteCmd(cmdline)
	wg.Done()
}

/*抓取trace*/
func grapTrace(wg *sync.WaitGroup, url, newdir, name string, port int, seconds string){
	cmdline := assembleTraceCmdLine(url, newdir, name,seconds, port)
	fmt.Println(cmdline)
	excuteCmd(cmdline)
	wg.Done()
}

/*抓取heap*/
func grapHeap(wg *sync.WaitGroup, url, newdir, name string, port int, seconds string){
	cmdline := assembleHeapCmdLine(url, newdir, name,seconds, port)
	fmt.Println(cmdline)
	excuteCmd(cmdline)
	wg.Done()
}

/*功能说明：遍历列表，抓取数据*/
func startGrap(url , dir , seconds string, m *FooModel){
	var wg sync.WaitGroup

	newdir := makeDir(dir, url)
	fmt.Println(newdir)
	
	for i := range m.items {
		if m.items[i].checked{
			wg.Add(3)
			go grapPProf(&wg, url, newdir,  m.items[i].Name, m.items[i].Port, seconds)
			go grapTrace(&wg, url, newdir,  m.items[i].Name, m.items[i].Port, seconds)
			go grapHeap(&wg, url, newdir,  m.items[i].Name, m.items[i].Port, seconds)		
		}
	}	
	
	wg.Wait()
	fmt.Println("finish grap")
	var tmp walk.Form
	walk.MsgBox(tmp, "提示", "抓取完成", walk.MsgBoxIconInformation)

	return
}



func main() {
	var te *walk.LineEdit
	var ne *walk.NumberEdit
	var fe *walk.LineEdit
	var widgets []Widget
	var tv *walk.TableView
	model := NewFooModel()
	widgets = []Widget{ 
		TableView{
			AssignTo:         &tv,
			AlternatingRowBG: true,
			CheckBoxes:       true,
			ColumnsOrderable: true,
			MultiSelection:   true,
			Columns: []TableViewColumn{
				{Title: "index"},
				{Title: "Name"},
				{Title: "Port", Alignment: AlignFar},
			},
			Model: model,
		},
		Composite{
			//Layout: VBox{MarginsZero: true},
			Layout:  Grid{Columns: 2},
			Children: []Widget{
				Label{
					Text: "k8s的IP地址:",
				},
				LineEdit{
				//	MaxSize: Size{20, 20},
					Text: "http://10.161.42.53:32000",
					AssignTo: &te,
				},
		
				PushButton{
					Text:"爬取Go Debug端口",
					OnClicked: func() {
						grabWebData(te.Text(), model)//开始爬取数据
					},
				},
				PushButton{
					Text:"选择全部",
					OnClicked: func() {
						selectAll( model)
					},
				},
				PushButton{
					Text:"全部不选",
					OnClicked: func() {
						unSelectAll( model)
					},
				},
				PushButton{
					Text:"选择呼叫类",
					OnClicked: func() {
						selectCall( model)
					},
				},
				Label{
					Text: "设置抓取秒数:",
				},
				NumberEdit{
					MinValue: 5.0,
					MaxValue:10000.0,
					Suffix:   " s",
					Decimals: 0,
					AssignTo: &ne,
				},
				PushButton{
					Text: "选择保存目录:",
					OnClicked: func(){
						dlg := new(walk.FileDialog)
						dlg.Title = "选择保存目录"
						var tmp walk.Form
						if _, err := dlg.ShowBrowseFolder(tmp); err != nil{
							fe.SetText("Error")
							return
						}
						s := fmt.Sprintf("%s", dlg.FilePath)
						fe.SetText(s)
					},
				},
				LineEdit{
					Text: "D://",
					AssignTo: &fe,
				},
		
				PushButton{
					Text:"开始抓取pprof和trace数据",
					OnClicked: func() {
						seconds := strconv.Itoa(int(ne.Value()))
						fmt.Println(seconds)
						dir := string(fe.Text())
						url := string(te.Text())
						startGrap(url, dir, seconds,  model)
					},
				},
			},
		},
		
		
	}

	MainWindow{
		Title:   "抓取性能数据",
		MinSize: Size{500, 400},
		Size:Size{600, 400},
		Layout:  Grid{Columns: 2},
		Children: widgets,
	}.Run()
}