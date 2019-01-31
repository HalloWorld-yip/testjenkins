package main

import (
	"math/rand"
	"time"
	"fmt"
)
//
//const (
//	REQ_PREPROPOSE = iota
//	REQ_PROPOSE
//	REQ_GETPOWER
//	REQ_COMMIT
//
//	MSG_AGREE = iota + 10
//	MSG_REJECT
//	MSG_ACCPEPT
//	MSG_UNACCPET
//	MSG_SETTLE
//
//	STATE_NULL = iota + 10
//	STATE_ACCEPT
//	STATE_PROPOSE
//	STATE_REJECT
//	STATE_PREPROPOSE
//	STATE_SETTLE
//)
//
//func min(a, b uint32) uint32 {
//	if a > b {
//		return b
//	}
//	return a
//}
//
//type TcpBuffer struct {
//	head uint32
//	back uint32
//	cap  uint32
//	free uint32
//	buf  []byte
//	lock sync.Mutex
//}
//
//func (c *TcpBuffer) Write(data []byte, length uint32) error {
//	c.lock.Lock()
//	defer c.lock.Unlock()
//
//	if length > c.free {
//		return errors.New("full")
//	}
//
//	if c.head != 0 {
//		copy(c.buf[0:], c.buf[c.head:])
//		c.back -= c.head
//		c.head = 0
//	}
//	copy(c.buf[c.back:], data[:length])
//	c.back += length
//	c.free -= length
//	return nil
//}
//
//func (c *TcpBuffer) read(data []byte, length uint32) error {
//	c.lock.Lock()
//	defer c.lock.Unlock()
//
//	len := length
//	if len > c.back-c.head {
//		return errors.New("not enough")
//	}
//
//	copy(data, c.buf[c.head:c.head+len])
//	c.head += len
//	c.free += len
//	return nil
//}
//
//type MSG struct {
//	Msg       uint32 `json:"msg"`
//	Power     uint32 `json:"power"`
//	Key       string `json:"key"`
//	Value     uint32 `json:"value"`
//	Port      uint32 `json:"port"`
//	TraceID   uint64 `json:"traceid"`
//	Timestamp uint64 `json:"timestamp"`
//}
//
//type NetWorker struct {
//	l       net.Listener
//	port    uint32
//	handler func(msg *MSG)
//	buf     TcpBuffer
//}
//
//func (s *NetWorker) Send(port uint32, v interface{}) {
//	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
//	if err != nil {
//		log.Error("addr:", err.Error())
//		return
//	}
//	conn, err := net.DialTCP("tcp", nil, addr)
//	if err != nil {
//		log.Error("DialTCP:", err.Error())
//		return
//	}
//	defer conn.Close()
//
//	data, err := json.Marshal(v)
//	if err != nil {
//		log.Error("Marshal:", err.Error())
//		return
//	}
//
//	m := make([]byte, 1)
//	m[0] = byte(len(data))
//	buf := bytes.NewBuffer(m)
//	buf.Write(data)
//
//	conn.Write(buf.Bytes())
//}
//
//func (s *NetWorker) recv(conn net.Conn) {
//	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
//	defer conn.Close()                                    // close connection before exit
//	request := make([]byte, 512)                          // set maxium request length to 128B to prevent flood attack
//	for {
//		read_len, err := conn.Read(request)
//
//		if err != nil {
//			//			log.Info(err)
//			break
//		}
//
//		if read_len == 0 {
//			break // connection already closed by client
//		}
//		//str := string(request[5:read_len])
//
//		if nil != s.buf.Write(request, uint32(read_len)) {
//			break
//		}
//		s.buf.read(request, 1)
//		if nil == s.buf.read(request[1:], uint32(request[0])) {
//			var msg MSG
//
//			err := json.Unmarshal(request[1:uint32(request[0]+1)], &msg)
//			if err != nil {
//				log.Error("Marshal:", err.Error())
//				break
//			}
//			s.handler(&msg)
//		}
//		request = make([]byte, 256) // clear last read content
//	}
//}
//
//func (s *NetWorker) Start() {
//	s.buf.cap = 512
//	s.buf.free = 512
//	s.buf.buf = make([]byte, s.buf.cap)
//	var e error
//	s.l, e = net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", s.port))
//	if e != nil {
//		log.Error("Error listening:", e)
//	}
//	if s.handler == nil {
//		log.Error("nil handler")
//	}
//
//	for {
//		conn, err := s.l.Accept()
//		if err != nil {
//			continue
//		}
//		//go s.recv(conn)
//		s.recv(conn)
//	}
//}
//
//func (s *NetWorker) End() {
//	s.l.Close()
//}
//
//type Leader struct {
//	NetWorker
//	power         uint32
//	changePower   uint32
//	key           string
//	value         uint32
//	state         uint32
//	acceptors     []uint32
//	Learner       uint32
//	acceptCount   uint32
//	unacceptCount uint32
//	agreeCount    uint32
//	rejectCount   uint32
//	traceid       uint64
//}
//
//func (s *Leader) IsMajority(num uint32) bool {
//	return num > uint32(len(s.acceptors))/2
//}
//
//func (s *Leader) IsHalf(num uint32) bool {
//	if len(s.acceptors) == 0 {
//		return false
//	}
//	return num > uint32(len(s.acceptors)-1)/2+1
//}
//
//func (s *Leader) Init(power, value, port uint32, acceptors []uint32, Learner uint32) {
//	s.NetWorker.handler = s.Handle
//	s.NetWorker.port = port
//	s.power = power
//	s.value = value
//	s.state = uint32(STATE_PREPROPOSE)
//	s.acceptors = acceptors
//	s.Learner = Learner
//	go s.NetWorker.Start()
//}
//
//func (s *Leader) Sync() {
//	traceid, e := uuid.Get(1)
//	if e != nil {
//		log.Info(e.Error())
//		return
//	}
//	s.traceid = traceid
//	s.acceptCount = 0
//	s.rejectCount = 0
//	s.agreeCount = 0
//	s.unacceptCount = 0
//	var msg = &MSG{
//		Msg:     uint32(REQ_PREPROPOSE),
//		Power:   s.power,
//		Value:   s.value,
//		Port:    s.port,
//		TraceID: traceid,
//	}
//	for _, u := range s.acceptors {
//		time.Sleep(time.Duration(200) * time.Millisecond)
//		s.Send(u, msg)
//		runtime.Gosched()
//	}
//}
//
//func (s *Leader) GetPower() {
//	s.Send(s.Learner, &MSG{
//		Msg:  uint32(REQ_GETPOWER),
//		Port: s.port,
//	})
//}
//
//func (s *Leader) Handle(msg *MSG) {
//	if msg.TraceID != s.traceid && msg.Msg != REQ_GETPOWER && msg.Msg != MSG_SETTLE {
//		return
//	}
//
//	switch msg.Msg {
//	case MSG_AGREE:
//		if s.state != uint32(STATE_PREPROPOSE) {
//			return
//		}
//		s.agreeCount += 1
//		//log.Info(s.port, "agreeby:", msg.Port, msg.Power, "count:", s.agreeCount)
//		if msg.Value != 0 && s.changePower < msg.Power {
//			s.value = msg.Value
//			s.changePower = msg.Power
//			log.Infof("value change: v:%d,p:%d", s.value, s.changePower)
//		}
//		if s.IsMajority(s.agreeCount) {
//			s.state = uint32(STATE_PROPOSE)
//			log.Info(s.port, "majority argree:", s.port, s.power, s.value)
//
//			//发送PROPOSE
//			msgAccept := &MSG{uint32(REQ_PROPOSE), s.power, s.key, s.value, s.port, msg.TraceID, 0}
//			go func() {
//				for _, u := range s.acceptors {
//
//					//--------------------------------------------------------------------------------------------------------------------------
//					time.Sleep(time.Duration(100) * time.Millisecond)
//					//--------------------------------------------------------------------------------------------------------------------------
//					s.Send(u, msgAccept)
//					runtime.Gosched()
//				}
//			}()
//		}
//	case MSG_REJECT:
//		if s.state != uint32(STATE_PREPROPOSE) {
//			return
//		}
//		s.rejectCount += 1
//		//log.Info(s.port, "rejectby:", msg.Port, msg.Power, "count:", s.rejectCount)
//
//		if s.IsMajority(s.rejectCount) {
//			//提高power，新一轮paxos
//			log.Info(s.port, "proposefailed,raise power:", s.power)
//			s.state = uint32(STATE_NULL)
//			go func() {
//				s.GetPower()
//			}()
//		}
//
//	case MSG_ACCPEPT:
//		if s.state != uint32(STATE_PROPOSE) {
//			return
//		}
//		s.acceptCount += 1
//		//log.Info(s.port, "acceptby:", msg.Port, msg.Power, "count:", s.acceptCount)
//		if s.IsMajority(s.acceptCount) {
//			s.state = uint32(STATE_ACCEPT)
//			log.Info(s.port, "majority accept,now value is set as", s.value)
//			//accepted , send to all acceptor
//			s.Send(s.Learner, &MSG{
//				Msg:     uint32(REQ_COMMIT),
//				Power:   s.power,
//				Value:   s.value,
//				Port:    s.port,
//				TraceID: msg.TraceID,
//			})
//		}
//
//	case MSG_UNACCPET:
//		if s.state != uint32(STATE_PROPOSE) {
//			return
//		}
//		s.unacceptCount += 1
//		//log.Info(s.port, "unacceptby:", msg.Port, msg.Power, "count:", s.unacceptCount)
//		if s.IsMajority(s.unacceptCount) {
//
//			log.Info(s.port, "proposefailed,raise power:", s.power)
//			s.state = uint32(STATE_NULL)
//			go func() {
//				s.GetPower()
//			}()
//		}
//	case MSG_SETTLE:
//		if s.state != STATE_ACCEPT {
//			s.power = msg.Power
//			s.value = msg.Value
//			s.state = uint32(STATE_SETTLE)
//		}
//	case REQ_GETPOWER:
//		s.power = msg.Power
//		s.state = uint32(STATE_PREPROPOSE)
//		s.Sync()
//	}
//
//}
//
//type Acceptor struct {
//	NetWorker
//	power       uint32
//	value       uint32
//	settlePower uint32
//	settleValue uint32
//}
//
//func (s *Acceptor) Init(port uint32) {
//	s.NetWorker.handler = s.Handle
//	s.NetWorker.port = port
//	go s.NetWorker.Start()
//}
//
//func (s *Acceptor) Handle(msg *MSG) {
//	switch msg.Msg {
//	case REQ_PREPROPOSE:
//		if s.settleValue != 0 {
//			s.Send(msg.Port, &MSG{
//				Msg:     uint32(MSG_SETTLE),
//				Power:   s.settlePower,
//				Value:   s.settleValue,
//				Port:    s.port,
//				TraceID: msg.TraceID,
//			})
//			break
//		}
//		if msg.Power <= s.power {
//			log.Info(s.port, "recject:", msg.Port, msg.Power)
//			s.Send(msg.Port, &MSG{
//				Msg:     uint32(MSG_REJECT),
//				Power:   s.power,
//				Port:    s.port,
//				TraceID: msg.TraceID,
//			})
//			return
//		}
//		s.power = msg.Power
//		log.Info(s.port, "agree:", msg.Port, msg.Power)
//		s.Send(msg.Port, &MSG{
//			Msg:     uint32(MSG_AGREE),
//			Value:   s.value,
//			Power:   s.power,
//			Port:    s.port,
//			TraceID: msg.TraceID,
//		})
//	case REQ_PROPOSE:
//		if s.settleValue != 0 {
//			s.Send(msg.Port, &MSG{
//				Msg:     uint32(MSG_SETTLE),
//				Power:   s.settlePower,
//				Value:   s.settleValue,
//				Port:    s.port,
//				TraceID: msg.TraceID,
//			})
//			break
//		}
//		if msg.Power < s.power {
//			log.Info(s.port, "unaccept:", msg.Port, msg.Power)
//			s.Send(msg.Port, &MSG{
//				Msg:     uint32(MSG_UNACCPET),
//				Power:   s.power,
//				Port:    s.port,
//				TraceID: msg.TraceID,
//			})
//			return
//		}
//		s.power = msg.Power
//		s.value = msg.Value
//		log.Info(s.port, "accept:", msg.Port, msg.Power, msg.Value)
//
//		s.Send(msg.Port, &MSG{
//			Msg:     uint32(MSG_ACCPEPT),
//			Value:   s.value,
//			Port:    s.port,
//			TraceID: msg.TraceID,
//		})
//	case MSG_SETTLE:
//		s.power = msg.Power
//		s.value = msg.Value
//		s.settlePower = msg.Power
//		s.settleValue = msg.Value
//	}
//}
//
//type Learner struct {
//	NetWorker
//	power   uint32
//	key     string
//	value   uint32
//	members []uint32
//}
//
//func (s *Learner) Init(port uint32, members []uint32) {
//	s.NetWorker.handler = s.Handle
//	s.NetWorker.port = port
//	s.members = members
//	go s.NetWorker.Start()
//}
//
//func (s *Learner) Handle(msg *MSG) {
//	switch msg.Msg {
//	case REQ_GETPOWER:
//		s.power++
//		rsp := &MSG{
//			Msg:   uint32(REQ_GETPOWER),
//			Key:   msg.Key,
//			Power: s.power,
//		}
//
//		s.Send(msg.Port, rsp)
//	case REQ_COMMIT:
//		rsp := &MSG{
//			Msg:     uint32(MSG_SETTLE),
//			TraceID: msg.TraceID,
//			Value:   msg.Value,
//			Power:   msg.Power,
//		}
//
//		go func() {
//			for _, u := range s.members {
//				time.Sleep(time.Duration(100) * time.Millisecond)
//				s.Send(u, rsp)
//				runtime.Gosched()
//			}
//		}()
//	}
//}
//
//func Const2String(s uint32) string {
//	switch s {
//	case STATE_NULL:
//		return "STATE_NULL		"
//	case STATE_ACCEPT:
//		return "STATE_ACCEPT	"
//	case STATE_PROPOSE:
//		return "STATE_PROPOSE	"
//	case STATE_REJECT:
//		return "STATE_REJECT	"
//	case STATE_PREPROPOSE:
//		return "STATE_PREPROPOSE"
//	case STATE_SETTLE:
//		return "STATE_SETTLE	"
//	}
//
//	return ""
//}
//
//var wg sync.WaitGroup
//
//func genToken() {
//	m := md5.New()
//	timestamp := time.Now().Unix()
//	log.Info(timestamp)
//	origin := fmt.Sprintf("%s%d", "c452d977d8f1513088f09689c6a8b0f7", timestamp)
//	m.Write([]byte(origin))
//	token := m.Sum(nil)
//	log.Info(hex.EncodeToString(token))
//}
//
//type CardInfoList struct {
//	ComponentVerifyTicket string `protobuf:"bytes,1,rep,name=component_verify_ticket,json=componentVerifyTicket,proto3" json:"component_verify_ticket,omitempty"`
//}
//type CardInfoList1 struct {
//	ComponentVerifyTicket int `protobuf:"bytes,1,rep,name=component_verify_ticket,json=componentVerifyTicket,proto3" json:"component_verify_ticket,omitempty"`
//}


func main() {
	//c := `<p style="text-align: center;"><img src="http://www.gd.gov.cn/img/0/30/30805/1408506.png" border="0" align="center" srcset="" class="nfw-cms-img" img-id="30805" data-catchresult="img_catchSuccess" style="text-align: center;"/></p><p style="text-align: center;"><img src="http://www.gd.gov.cn/img/0/30/30804/1408506.png" border="0" align="center" srcset="" class="nfw-cms-img" img-id="30804" data-catchresult="img_catchSuccess"/></p><p style="text-align: center;"><img src="http://www.gd.gov.cn/img/0/30/30803/1408506.png" border="0" align="center" srcset="" class="nfw-cms-img" img-id="30803" data-catchresult="img_catchSuccess"/></p><p style="text-align: center;"><img src="http://www.gd.gov.cn/img/0/30/30802/1408506.png" border="0" align="center" srcset="" class="nfw-cms-img" img-id="30802" data-catchresult="img_catchSuccess"/></p><p style="text-align: justify; text-indent: 2em; margin-top: 10px; line-height: 2em;"><span style="font-family: 微软雅黑, &quot;Microsoft YaHei&quot;;">17日上午，习近平来到天津考察调研。在南开大学，他参观了百年校史主题展览，与部分院士、专家和中青年师生代表互动交流，察看了化学学院和元素有机化学国家重点实验室，详细了解南开大学历史沿革、学科建设、人才队伍、科研创新等情况。</span></p>`
	//
	//parseImg(c)
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Uint64()%300)
	return
	//b, err := ioutil.ReadFile("C:\\Users\\Administrator\\Documents\\测试文件3.pdf")
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//dst := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	//base64.StdEncoding.Encode(dst, b)

	// id, _ := uuid.Get(1)
	// uploadReply, errMsg := soap.CertFileUploadByStream(
	// 	fmt.Sprintf("%d", id),
	// 	dst,
	// 	"pdf",
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info(*uploadReply)

	// temReply, errMsg := soap.CertQueryTemplateInfo(
	// 	"0",
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info(temReply)

	// temReply, errMsg = soap.CertQueryTemplateInfo(
	// 	"1",
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info(temReply)

	// createEnvReply, errMsg := soap.CertCreateEnvelop(
	// 	uploadReply.PdfId,
	// 	"0",
	// 	"env11071614",
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info(createEnvReply.Rlt, createEnvReply.Data)

	// queryEnvReply, errMsg := soap.CertQueryEnvelopStatus(
	// 	createEnvReply.Data.Eid,
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info("query:", queryEnvReply.Data)

	// startEnvReply, errMsg := soap.CertStartEnvelop(
	// 	createEnvReply.Data.Eid,
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info(startEnvReply.Rlt)

	// jsonrule := "{\"callbacktype\":\"0\",\"callbackurl\":\"\"," +
	// 	"\"expireDate\":\"2018-12-31 16:55:56.620\"," +
	// 	"\"signedType\":\"1\",\"signSide\":2}"
	// signEnvReply, errMsg := soap.CertCreateEnvelopSignRequest(
	// 	createEnvReply.Data.Eid,
	// 	jsonrule,
	// 	"S",
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info("create signed:", signEnvReply)

	// // echoReply, errMsg := soap.CertVerifyPdf(
	// // 	dst,
	// // )
	// // if engin.IsErr(errMsg) {
	// // 	log.Error(errMsg)
	// // 	return
	// // }
	// // log.Info(echoReply)

	// queryEnvReply, errMsg = soap.CertQueryEnvelopStatus(
	// 	createEnvReply.Data.Eid,
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info("query:", queryEnvReply.Data)

	// body := "<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:ns=\"http://certbaosignservice.server.ws.gov4biz.com/\" xmlns:tns=\"http://certbaosignservice.server.ws.gov4biz.com/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"><SOAP-ENV:Body><ns:downloadPdfByEnvelop><arg0><appId>gwjhapp</appId></arg0><arg1><userId>gwjh</userId></arg1><arg2>84861671</arg2></ns:downloadPdfByEnvelop></SOAP-ENV:Body></SOAP-ENV:Envelope>"
	// req, err := http.NewRequest("POST", "http://112.74.77.238:8080/certbaoServer/services/certbaoSignService", bytes.NewBuffer([]byte(body)))

	// client := &http.Client{}
	// httpResp, err := client.Do(req)
	// if httpResp != nil {
	// 	defer httpResp.Body.Close()
	// 	defer io.Copy(ioutil.Discard, httpResp.Body)
	// }
	// if err != nil {
	// 	log.Error("http err:", err)
	// 	return
	// }
	// mediaType, params, err := mime.ParseMediaType(httpResp.Header.Get("Content-Type"))
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// log.Info(mediaType, params)
	// r := multipart.NewReader(httpResp.Body, params["boundary"])
	// for {
	// 	p, err := r.NextPart()
	// 	defer p.Close()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	log.Info(p.Header.Get("Content-ID"))
	// 	var b bytes.Buffer
	// 	if _, err := io.Copy(&b, p); err != nil {
	// 		return
	// 	}
	// 	log.Info(string(b.Bytes()))
	// }
	// bodyStr, err := ioutil.ReadAll(httpResp.Body)
	// if err != nil {
	// 	log.Error("http err:", err)
	// 	return
	// }
	// log.Info("len:", len(bodyStr))
	// return

	// downloadbytes, errMsg := soap.CertDownloadPdfByEnvelop(
	// 	"84861671",
	// )
	// if engin.IsErr(errMsg) {
	// 	log.Error(errMsg)
	// 	return
	// }
	// log.Info("len:", len(downloadbytes))
	// return

	//verifyReply, errMsg := soap.CertVerVerifyPdf(
	//	dst,
	//)
	//if engin.IsErr(errMsg) {
	//	log.Error(errMsg)
	//	return
	//}
	//log.Info(*verifyReply)
	//
	//b, err := ioutil.ReadFile("project.config.json")
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//key := "\"appid\":"
	//beforappid := strings.Index(string(b), key)
	//tmpb := b[beforappid+len(key):]
	//
	//afterappid := strings.Index(string(tmpb), ",")
	//
	//newb := string(b[:beforappid+len(key)]) + " \"wx8a48f63c76f0d13f\"" + string(tmpb[afterappid:])
	//
	//ioutil.WriteFile("project.config.json", []byte(newb), 0666)
	//return
	//
	//s, _ := mongoclient.GetSession()
	//defer mongoclient.ReleaseSesion(s)
	//c := s.DB("test").C("tc")
	//runner := txn.NewRunner(c)
	//if runner == nil {
	//	log.Error("runner nil")
	//	return
	//}
	//
	//type testobj struct {
	//	Id string
	//}
	//ops := []txn.Op{{
	//	C:      "test1",
	//	Id:     0,
	//	Assert: txn.DocMissing,
	//	Insert: bson.M{"hha": 1},
	//}, {
	//	C:      "test1",
	//	Id:     1,
	//	Assert: txn.DocMissing,
	//	Insert: bson.M{"hha": 2},
	//}, {
	//	C:      "test1",
	//	Id:     2,
	//	Assert: txn.DocMissing,
	//	Insert: bson.M{"hha": 3},
	//}}
	//e := runner.Run(ops, "", nil)
	//if e != nil {
	//	log.Error(e)
	//}
	//return

	// ldrs := make([]*Leader, 0)
	// acpts := make([]*Acceptor, 0)
	// learner := &Learner{}
	// mems := make([]uint32, 0)
	// for i := 55555; i < 55570; i++ {
	// 	mems = append(mems, uint32(i))
	// }
	// for i := 54321; i < 54325; i++ {
	// 	mems = append(mems, uint32(i))
	// }
	// learner.Init(55554, mems)

	// acptports := make([]uint32, 0)
	// for i := 55555; i < 55570; i++ {
	// 	ltmp := &Acceptor{}
	// 	ltmp.Init(uint32(i))
	// 	acpts = append(acpts, ltmp)
	// 	acptports = append(acptports, uint32(i))
	// }

	// for i := 54321; i < 54325; i++ {
	// 	ltmp := &Leader{}
	// 	rand.Seed(time.Now().UnixNano())
	// 	power := rand.Uint32() % 100
	// 	value := rand.Uint32() % 10000
	// 	ltmp.Init(power, value, uint32(i), acptports, 55554)
	// 	log.Infof("leader:%d,%d,%d", i, power, value)
	// 	ldrs = append(ldrs, ltmp)

	// 	time.Sleep(time.Duration(100) * time.Millisecond)
	// }

	// for _, l := range ldrs {
	// 	l.GetPower()
	// }

	// for {
	// 	time.Sleep(time.Duration(100) * time.Millisecond)

	// 	log.Info("------------------------------------------")
	// 	for _, acpt := range acpts {
	// 		log.Infof("acceptor:%d,v:%d,p:%d", acpt.port, acpt.value, acpt.power)
	// 	}
	// 	for _, ldr := range ldrs {
	// 		log.Infof("proposer:%d,%s,v:%d,p:%d,pac:%d,ac:%d,prj:%d,rj:%d",
	// 			ldr.port, Const2String(ldr.state), ldr.value, ldr.power, ldr.agreeCount, ldr.acceptCount, ldr.rejectCount, ldr.unacceptCount)
	// 	}
	// 	log.Info("------------------------------------------")
	// }
}
